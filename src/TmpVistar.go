package src

//func GetCommandToSaveLogsHTTP(w http.ResponseWriter, r *http.Request){
//	startDate := r.URL.Query().Get("start")
//	endDate := r.URL.Query().Get("end")
//	Type := r.URL.Query().Get("type")
//	PerPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
//
//	var jobLogsStart JobLogs
//	jobLogsStart.JobHash = GetHash()
//	jobLogsStart.JobName = "Get logs from maxtvmedia server " + startDate + " - " + endDate + " : " + Type
//	Db.Create(&jobLogsStart)
//
//	go func() {
//		fmt.Println(GetVistarLogs(startDate, endDate, Type, PerPage, jobLogsStart.JobHash))
//		Db.Model(JobLogs{}).Where("job_hash = ?", jobLogsStart.JobHash).Update(JobLogs{JobStatus: "successfully done", JobDone: true})
//	}()
//
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("content-type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	n, _ := fmt.Fprintf(w, "{\"JobHash\" : \""+jobLogsStart.JobHash+"\"}")
//	fmt.Println(n)
//
//}
//
//func GetAllLogsFromTo(Year, month, startDay, DayCount int) {
//	for c_day = startDay; c_day<startDay+DayCount; c_day++ {
//
//	}
//
//}
