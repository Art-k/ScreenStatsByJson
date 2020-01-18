package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var NothingToDo bool = false

type GetLogsStatus struct {
	CurrentPage       int64
	TotalPageCount    int64
	TimePerPageMs     int64
	TimePerLogMs      int64
	EstimateSec       int64
	EstimateHumanTime string
	Status            string
	JobHash           string
}

// 'https://maxtvmedia.com/cms/.cron/vistar_new/logs.php?from=2019-12-13&to=2019-12-14&event=get_assets&per-page=10&page=2'
// get_assets | get_proofs | send_proofs

func SaveAssetsRequestToDatabase(IncRequest IncVistarAssetsRequest) int64 {
	var vistarAssetsRequest VistarAssetsRequest

	vistarAssetsRequest.DisplayTime = IncRequest.DisplayTime
	vistarAssetsRequest.HasResponse = false

	var screen Screen
	screen.SysId = IncRequest.DeviceId
	Db.Where("sys_id = ?", IncRequest.DeviceId).FirstOrCreate(&screen)

	vistarAssetsRequest.ScreenID = screen.ID
	DbVistar.Create(&vistarAssetsRequest)

	return vistarAssetsRequest.ID
}

func AssetsRequestHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var incomingData IncVistarAssetsRequest
		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
			log.Println(n)
			return
		} else {

			var screen Screen
			screen.SysId = incomingData.DeviceId
			Db.Where("sys_id = ?", incomingData.DeviceId).FirstOrCreate(&screen)
			var vistarAssetsRequest VistarAssetsRequest
			DbVistar.Where("display_time = ?", incomingData.DisplayTime).Where("screen_id = ?", screen.ID).Find(&vistarAssetsRequest)

			if vistarAssetsRequest.ID == 0 {
				ID := SaveAssetsRequestToDatabase(incomingData)
				w.WriteHeader(http.StatusCreated)
				n, _ := fmt.Fprintf(w, "{\"id\":"+strconv.Itoa(int(ID))+"}")
				log.Println(n)
				return
			} else {
				fmt.Println(">>> DUPLICATE >>>")
				w.WriteHeader(http.StatusBadRequest)
				n, _ := fmt.Fprintf(w, "{\"message\" : \"POST is allowed method\"}")
				log.Println(n)
				return
			}
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		n, _ := fmt.Fprintf(w, "{\"message\" : \"POST is allowed method\"}")
		log.Println(n)
		return
	}
}

func SaveAssetsResponsesToDatabase(request_id int64, response IncVistarAssetResponse) {
	for _, incAsset := range response.Asset {
		var asset VistarAsset
		DbVistar.Where("asset_id = ?", incAsset.AssetId).
			Where("creative_id = ?", incAsset.CreativeId).
			Where("campaign_id = ?", incAsset.CampaignId).
			Find(&asset)
		if asset.ID == 0 {
			asset.CampaignId = incAsset.CampaignId
			asset.AssetId = incAsset.AssetId
			asset.CreativeId = incAsset.CreativeId
			asset.CreativeName = incAsset.CreativeName
			asset.Advertiser = incAsset.Advertiser
			asset.CreativeCategory = incAsset.CreativeCategory
			asset.LengthMSec = incAsset.LengthMSec
			asset.LengthSec = incAsset.LengthSec
			asset.MimeType = incAsset.MimeType
			asset.Height = incAsset.Height
			asset.Width = incAsset.Width
			asset.AssetURL = incAsset.AssetURL
			DbVistar.Create(&asset)
		}

		var vistartAssetResponse VistarAssetResponse
		vistartAssetResponse.RequestID = request_id
		vistartAssetResponse.VistarAssetID = asset.ID
		DbVistar.Create(&vistartAssetResponse)

	}
	if len(response.Asset) != 0 {
		DbVistar.Model(VistarAssetsRequest{}).Where("id = ?", request_id).Update(VistarAssetsRequest{HasResponse: true})
	}
}

func AssetsResponseHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		params := mux.Vars(r)
		request_id, _ := strconv.Atoi(params["belongs_to_request_id"])

		var incomingData IncVistarAssetResponse
		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			n, _ := fmt.Fprintf(w, "{\"message\" : \"Unexpected Object Received\"}")
			log.Println(n)
			return
		} else {
			SaveAssetsResponsesToDatabase(int64(request_id), incomingData)
			w.WriteHeader(http.StatusCreated)
			n, _ := fmt.Fprintf(w, "")
			log.Println(n)
			return
		}

	default:

	}
}

func GetVistarLogs(SD, ED, Type string, PerPage int, Hash string) (int, int) {

	var getLogsStatus GetLogsStatus
	getLogsStatus.JobHash = Hash

	var Downloaded int = 0
	var Skipped int = 0

	VistarLogsUrl := os.Getenv("VISTAR_LOGS_URL")

	var currentPage int = 1
	var DownloadedRec int64 = 0

	for {

		pageProcessingStart := time.Now().UnixNano() / 1000000
		req, _ := http.NewRequest("GET", VistarLogsUrl, nil)
		q := req.URL.Query()
		q.Add("event", Type)
		q.Add("from", SD)
		q.Add("to", ED)
		q.Add("per-page", strconv.Itoa(PerPage))
		q.Add("page", strconv.Itoa(currentPage))
		req.URL.RawQuery = q.Encode()

		resp, _ := http.Get(req.URL.String())
		body, _ := ioutil.ReadAll(resp.Body)

		var response IncVistarEventLogPage
		jsonErr := json.Unmarshal(body, &response)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		fmt.Println(currentPage, response.Total/int64(PerPage))
		getLogsStatus.CurrentPage = int64(currentPage)
		getLogsStatus.TotalPageCount = response.Total / int64(PerPage)

		for _, rec := range response.Entities {
			oneRecProcessingStart := time.Now().UnixNano() / 1000000
			var vlrec VistarEventLog
			n, _ := strconv.ParseInt(rec.MaxTvId, 10, 64)
			DbVistar.Where("max_tv_id = ?", n).Find(&vlrec)
			if vlrec.ID == 0 {
				vlrec.MaxTvId = n
				//vlrec.Link = rec.Link
				vlrec.Link = strings.Replace(rec.Link, VistarLogsBaseUrl, "", -1)

				var eventType VistarEventType
				DbVistar.Where("type_name = ?", rec.Event).Find(&eventType)
				if eventType.ID == 0 {
					eventType.TypeName = rec.Event
					DbVistar.Create(&eventType)
				}
				vlrec.EventID = eventType.ID

				var logType VistarLogType
				DbVistar.Where("log_type_name = ?", rec.Type).Find(&logType)
				if logType.ID == 0 {
					logType.LogTypeName = rec.Type
					DbVistar.Create(&logType)
				}
				vlrec.TypeID = logType.ID
				vlrec.Processed = false
				DbVistar.Create(&vlrec)
				Downloaded++
			} else {
				Skipped++
			}
			oneRecProcessingEnd := time.Now().UnixNano() / 1000000
			getLogsStatus.TimePerLogMs = oneRecProcessingEnd - oneRecProcessingStart
		}

		currentPage++
		DownloadedRec = DownloadedRec + int64(len(response.Entities))

		if DownloadedRec >= response.Total {
			break
		}

		var tmp JobLogs
		cache, _ := Cache.Get(Hash + "_ACTION")
		json.Unmarshal(cache, &tmp)
		if tmp.Action == "CANCEL" {
			getLogsStatus.Status = "Terminated"
			getLogsStatus.EstimateSec = 0
			stringData, _ := json.Marshal(getLogsStatus)
			Cache.Set(Hash, stringData)
			break
		}

		pageProcessingEnd := time.Now().UnixNano() / 1000000
		getLogsStatus.TimePerPageMs = pageProcessingEnd - pageProcessingStart
		getLogsStatus.EstimateSec = ((getLogsStatus.TotalPageCount - getLogsStatus.CurrentPage) * getLogsStatus.TimePerPageMs) / 1000
		stringData, _ := json.Marshal(getLogsStatus)
		Cache.Set(Hash, stringData)

	}

	getLogsStatus.Status = "Successfully Completed"
	getLogsStatus.EstimateSec = 0
	stringData, _ := json.Marshal(getLogsStatus)
	Cache.Set(Hash, stringData)
	return Downloaded, Skipped
}

func ProceedVistarGetAssetsRequest() (VistarEventLog, VistarEventLog, VistarAssetsRequest, int) {
	var record_request VistarEventLog

	var vistarEventType VistarEventType
	Db.Where("type_name = ?", "get_assets").Find(&vistarEventType)

	var vistarLogType VistarLogType
	Db.Where("log_type_name = ?", "request").Find(&vistarLogType)

	Db.Where("processed <> ?", 1).
		Where("type = ?", vistarLogType.ID).
		Where("event = ?", vistarEventType.ID).
		First(&record_request)

	var record_response VistarEventLog
	var v_getassets VistarAssetsRequest
	assetsCount := 0
	if record_request.Link != "" {

		//record_request := records[0]
		tmpLink := record_request.Link
		Link := strings.Replace(tmpLink, "https://maxtvmedia.com/cms/.cron/vistar/", "http://127.0.0.1:50000/", -1)

		//resp, _ := http.Get(record_request.Link)
		resp, _ := http.Get(Link)
		body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(body)

		var response IncVistarAssetsRequest
		jsonErr := json.Unmarshal(body, &response)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		//var v_getassets VistarAssetsRequest
		//v_getassets.LogFileId = record_request.ID

		var screen Screen
		Db.Where("sys_id = ?", response.DeviceId).Find(&screen)
		if screen.ID == 0 {
			//log.Fatal("Screen not found")
			Db.Model(VistarEventLog{}).Where("id = ?", record_request.ID).Update(VistarEventLog{Processed: true})
			return record_request, record_response, v_getassets, 0
		}
		v_getassets.ScreenID = screen.ID
		v_getassets.DisplayTime = response.DisplayTime
		//v_getassets.DirectConnection = response.DirectConnection
		Db.Create(&v_getassets)

		newLink := strings.Replace(record_request.Link, "request", "response", -1)

		var record_response VistarEventLog
		Db.Where("link = ?", newLink).Find(&record_response)
		if record_response.Link != "" {
			resp, _ := http.Get(record_response.Link)
			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) == "[]" {
				Db.Model(VistarAssetsRequest{}).Where("id = ?", v_getassets.ID).Update(VistarAssetsRequest{HasResponse: false})
				Db.Model(VistarEventLog{}).Where("id = ?", record_request.ID).Update(VistarEventLog{Processed: true})
				Db.Model(VistarEventLog{}).Where("id = ?", record_response.ID).Update(VistarEventLog{Processed: true})
				return record_request, record_response, v_getassets, 0
			}

			var incVistarAssetResponse IncVistarAssetResponse
			jsonErr := json.Unmarshal(body, &incVistarAssetResponse)
			if jsonErr != nil {
				Db.Model(VistarAssetsRequest{}).Where("id = ?", v_getassets.ID).Update(VistarAssetsRequest{HasResponse: false})
				return record_request, record_response, v_getassets, 0
			}

			if len(incVistarAssetResponse.Asset) != 0 {
				for _, inc_asset := range incVistarAssetResponse.Asset {
					var asset VistarAsset
					Db.Where("asset_id = ?", inc_asset.AssetId).Last(&asset)
					if asset.ID == 0 {
						asset.AssetId = inc_asset.AssetId
						asset.CreativeId = inc_asset.CreativeId
						asset.CampaignId = inc_asset.CampaignId
						asset.AssetURL = inc_asset.AssetURL
						asset.Width = inc_asset.Width
						asset.Height = inc_asset.Height
						asset.MimeType = inc_asset.MimeType
						asset.LengthSec = inc_asset.LengthSec
						asset.LengthMSec = inc_asset.LengthMSec
						asset.CreativeCategory = inc_asset.CreativeCategory
						asset.Advertiser = inc_asset.Advertiser
						asset.CreativeName = inc_asset.CreativeName
						//return
						Db.Create(&asset)
					}
					var vistarAssetResponse VistarAssetResponse
					vistarAssetResponse.RequestID = v_getassets.ID
					vistarAssetResponse.VistarAssetID = asset.ID
					Db.Create(&vistarAssetResponse)
					assetsCount++
				}
			}
		} else {

			Db.Model(VistarAssetsRequest{}).Where("id = ?", v_getassets.ID).Update(VistarAssetsRequest{HasResponse: false})
		}
		Db.Model(VistarEventLog{}).Where("id = ?", record_request.ID).Update(VistarEventLog{Processed: true})
	} else {
		NothingToDo = true
	}
	return record_request, record_response, v_getassets, assetsCount
}

func VistarProceedAssetRequestsLogFiles(w http.ResponseWriter, r *http.Request) {
	Count := 0
	TotalCalls := 0
	for {
		if Count < 10 {
			TotalCalls++
			go func() {
				ProceedVistarGetAssetsRequest()
			}()
			Count--
		}
		Count++
		if NothingToDo {
			break
		}
		fmt.Println(strconv.Itoa(TotalCalls) + ": Number of treads " + strconv.Itoa(Count))
	}
	fmt.Println("Number of treads JOB IS DONE")
}

func VistarAssetRequestHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
	default:
	}
}

func VistarAssetResponseHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
	default:
	}
}

func VistarProceedAssetRequestsLogFile(w http.ResponseWriter, r *http.Request) {

	request_line, response_line, request, assetsCount := ProceedVistarGetAssetsRequest()

	type Response struct {
		RequestLine  VistarEventLog      `json:"request_line"`
		ResponseLine VistarEventLog      `json:"response_line"`
		Request      VistarAssetsRequest `json:"request"`
		AssetsCount  int                 `json:"assetsCount"`
	}

	var resp Response
	resp.RequestLine = request_line
	resp.ResponseLine = response_line
	resp.Request = request
	resp.AssetsCount = assetsCount

	addedRecordString, _ := json.Marshal(resp)
	//fmt.Println(err, addedRecordString)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, string(addedRecordString))
	fmt.Println(n)

	return
}

func VistarLogs(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")
	Type := r.URL.Query().Get("type")
	PerPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	var jobLogsStart JobLogs
	jobLogsStart.JobHash = GetHash()
	jobLogsStart.JobName = "Get logs from maxtvmedia server " + startDate + " - " + endDate + " : " + Type
	Db.Create(&jobLogsStart)

	go func() {
		fmt.Println(GetVistarLogs(startDate, endDate, Type, PerPage, jobLogsStart.JobHash))
		Db.Model(JobLogs{}).Where("job_hash = ?", jobLogsStart.JobHash).Update(JobLogs{JobStatus: "Successfully Completed", JobDone: true})
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, "{\"JobHash\" : \""+jobLogsStart.JobHash+"\"}")
	fmt.Println(n)
}
