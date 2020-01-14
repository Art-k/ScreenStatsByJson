package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type WorkStatus struct {
	CurrentPage    int64
	TotalPageCount int64
}

var workStatus WorkStatus

// 'https://maxtvmedia.com/cms/.cron/vistar_new/logs.php?from=2019-12-13&to=2019-12-14&event=get_assets&per-page=10&page=2'
// get_assets | get_proofs | send_proofs

func GetVistarLogs(SD, ED, Type string) (int, int) {

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
		q.Add("per-page", strconv.Itoa(10))
		q.Add("page", strconv.Itoa(currentPage))
		req.URL.RawQuery = q.Encode()

		resp, _ := http.Get(req.URL.String())
		body, _ := ioutil.ReadAll(resp.Body)

		var response IncomingVLPage
		jsonErr := json.Unmarshal(body, &response)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		fmt.Println(currentPage, response.Total/10)
		workStatus.CurrentPage = int64(currentPage)
		workStatus.TotalPageCount = response.Total / 10

		for _, rec := range response.Entities {
			var vlrec VLRec
			n, _ := strconv.ParseInt(rec.MaxTvId, 10, 64)
			Db.Where("max_tv_id = ?", n).Find(&vlrec)
			if vlrec.ID == 0 {
				vlrec.MaxTvId = n
				vlrec.Link = rec.Link
				vlrec.Type = rec.Type
				vlrec.Event = rec.Event
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

func VistarLogs(w http.ResponseWriter, r *http.Request) {

	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")
	Type := r.URL.Query().Get("type")
	fmt.Println(GetVistarLogs(startDate, endDate, Type))
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, "")
	fmt.Println(n)
}
