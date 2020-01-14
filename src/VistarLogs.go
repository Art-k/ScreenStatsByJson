package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type WorkStatus struct {
	CurrentPage    int64
	TotalPageCount int64
}

var workStatus WorkStatus

// 'https://maxtvmedia.com/cms/.cron/vistar_new/logs.php?from=2019-12-13&to=2019-12-14&event=get_assets&per-page=10&page=2'
// get_assets | get_proofs | send_proofs

func GetVistarLogs(SD, ED, Type string) (int, int) {

	var per_page int = 100

	var Downloaded int = 0
	var Skipped int = 0

	VistarLogsUrl := os.Getenv("VISTAR_LOGS_URL")

	var currentPage int = 1
	var DownloadedRec int64 = 0

	for {

		req, _ := http.NewRequest("GET", VistarLogsUrl, nil)
		q := req.URL.Query()
		q.Add("event", Type)
		q.Add("from", SD)
		q.Add("to", ED)
		q.Add("per-page", strconv.Itoa(per_page))
		q.Add("page", strconv.Itoa(currentPage))
		req.URL.RawQuery = q.Encode()

		resp, _ := http.Get(req.URL.String())
		body, _ := ioutil.ReadAll(resp.Body)

		var response IncomingVLPage
		jsonErr := json.Unmarshal(body, &response)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		fmt.Println(currentPage, response.Total/int64(per_page))
		workStatus.CurrentPage = int64(currentPage)
		workStatus.TotalPageCount = response.Total / int64(per_page)

		for _, rec := range response.Entities {
			var vlrec VLRec
			n, _ := strconv.ParseInt(rec.MaxTvId, 10, 64)
			Db.Where("max_tv_id = ?", n).Find(&vlrec)
			if vlrec.ID == 0 {
				vlrec.MaxTvId = n
				vlrec.Link = rec.Link
				vlrec.Type = rec.Type
				vlrec.Event = rec.Event
				vlrec.Processed = false
				Db.Create(&vlrec)
				Downloaded++
			} else {
				Skipped++
			}
		}

		currentPage++
		DownloadedRec = DownloadedRec + int64(len(response.Entities))

		if DownloadedRec >= response.Total {
			break
		}
	}

	return Downloaded, Skipped
}

func VistarLogsStat(w http.ResponseWriter, r *http.Request) {

	addedRecordString, _ := json.Marshal(workStatus)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, string(addedRecordString))
	fmt.Println(n)
	return

}

func ProceedVistarGetAssetsRequest() (VLRec, VLRec, VistarGetAssetsRequest, int) {
	var record_request VLRec
	Db.Where("processed <> ?", 1).
		Where("type = ?", "request").
		Where("event = ?", "get_assets").
		First(&record_request)

	//record_request := records[0]
	resp, _ := http.Get(record_request.Link)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(body)

	var response VistarGetAssetsRequestIncoming
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var v_getassets VistarGetAssetsRequest
	v_getassets.LogFileId = record_request.ID
	v_getassets.DeviceId = response.DeviceId
	v_getassets.VenueId = response.VenueId
	v_getassets.DisplayTime = response.DisplayTime
	v_getassets.DirectConnection = response.DirectConnection
	Db.Create(&v_getassets)

	newLink := strings.Replace(record_request.Link, "request", "response", -1)

	assetsCount := 0
	var record_response VLRec
	Db.Where("link = ?", newLink).Find(&record_response)
	if record_response.Link != "" {
		resp, _ := http.Get(record_response.Link)
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "[]" {
			Db.Model(VistarGetAssetsRequest{}).Where("id = ?", v_getassets.ID).Update(VistarGetAssetsRequest{HasResponse: false})
			Db.Model(VLRec{}).Where("id = ?", record_request.ID).Update(VLRec{Processed: true})
			Db.Model(VLRec{}).Where("id = ?", record_response.ID).Update(VLRec{Processed: true})
			return record_request, record_response, v_getassets, 0
		}

		var incVistarAssetResponse IncVistarAssetResponse
		jsonErr := json.Unmarshal(body, &incVistarAssetResponse)
		if jsonErr != nil {
			Db.Model(VistarGetAssetsRequest{}).Where("id = ?", v_getassets.ID).Update(VistarGetAssetsRequest{HasResponse: false})
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

		Db.Model(VistarGetAssetsRequest{}).Where("id = ?", v_getassets.ID).Update(VistarGetAssetsRequest{HasResponse: false})
	}
	Db.Model(VLRec{}).Where("id = ?", record_request.ID).Update(VLRec{Processed: true})

	return record_request, record_response, v_getassets, assetsCount
}

func VistarFileDBBulk(w http.ResponseWriter, r *http.Request) {
	Count := 0
	for {
		request_line, _, _, _ := ProceedVistarGetAssetsRequest()
		Count++
		if request_line.ID == 0 {
			break
		}
		fmt.Println(Count)
	}
}

func VistarFileDB(w http.ResponseWriter, r *http.Request) {

	request_line, response_line, request, assetsCount := ProceedVistarGetAssetsRequest()

	type Response struct {
		RequestLine  VLRec                  `json:"request_line"`
		ResponseLine VLRec                  `json:"response_line"`
		Request      VistarGetAssetsRequest `json:"request"`
		AssetsCount  int                    `json:"assetsCount"`
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
	fmt.Println(GetVistarLogs(startDate, endDate, Type))
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, "")
	fmt.Println(n)
}
