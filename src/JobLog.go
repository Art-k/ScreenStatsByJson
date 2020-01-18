package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

type JobLogs struct {
	gorm.Model
	JobHash   string
	JobName   string
	JobDone   bool
	JobStatus string
	Action    string
}

func JobsStatus(w http.ResponseWriter, r *http.Request) {
	var jobLogs []JobLogs
	Db.Find(&jobLogs)

	entry, _ := json.Marshal(jobLogs)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, string(entry))
	fmt.Println(n)
}

func JobStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		params := mux.Vars(r)
		hash := params["hash"]

		entry, _ := Cache.Get(hash)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(entry))
		fmt.Println(n)

	case "DELETE":
		params := mux.Vars(r)
		Hash := params["hash"]

		Cache.Set(Hash+"_ACTION", []byte("{\"Action\":\"CANCEL\"}"))
		Db.Model(JobLogs{}).Where("job_hash = ?", Hash).Update(JobLogs{JobStatus: "Terminated", JobDone: false})

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string("{\"Status\":\"Terminated\"}"))
		fmt.Println(n)
	}
}
