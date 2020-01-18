package main

import (
	"./src"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	err := godotenv.Load("p.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	src.Port = os.Getenv("PORT")
	src.VistarLogsBaseUrl = os.Getenv("VISTAR_BASE_PATH_TO_LOGS")
	src.AuthorizationHash = os.Getenv("AUTH")
	src.GetMaxTvStatInterval, _ = strconv.ParseInt(os.Getenv("GET_MAXTV_STAT_INTERVAL_SECOND"), 10, 64)

	src.Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))

	src.Db, src.Err = gorm.Open("sqlite3", "screen.db")
	if src.Err != nil {
		panic("failed to connect database")
	}

	src.DbVistar, src.Err = gorm.Open("sqlite3", "vistar.db")
	if src.Err != nil {
		panic("failed to connect database")
	}

	defer src.Db.Close()
	src.Db.LogMode(src.DbLogMode)

	//src.Db.AutoMigrate(&src.Ad{})
	//src.Db.AutoMigrate(&src.Screen{})
	//src.Db.AutoMigrate(&src.Building{})
	//src.Db.AutoMigrate(&src.Building{})
	//src.Db.AutoMigrate(&src.Stat{})
	//src.Db.AutoMigrate(&src.StatGetAttempt{})
	//src.Db.AutoMigrate(&src.Video{})
	//src.Db.AutoMigrate(&src.Spot{})

	// All related to VISTAR statistic
	src.DbVistar.AutoMigrate(&src.VistarEventLog{})
	src.DbVistar.AutoMigrate(&src.VistarAssetsRequest{})
	src.DbVistar.AutoMigrate(&src.VistarAsset{})
	src.DbVistar.AutoMigrate(&src.VistarAssetResponse{})
	src.DbVistar.AutoMigrate(&src.VistarEventType{})
	src.DbVistar.AutoMigrate(&src.VistarLogType{})

	src.Db.AutoMigrate(&src.JobLogs{})
	//src.Db.AutoMigrate(&src.TmpVistarEventLog{})

	//go func() {
	//	src.DoEvery(time.Duration(src.GetMaxTvStatInterval)*time.Second, src.GetScreenStat)
	//}()

	handleHTTP()
}

func handleHTTP() {

	r := mux.NewRouter()
	r.Use(authMiddleware)
	r.Use(headerMiddleware)

	// Jsons end points
	//r.HandleFunc("/coverage", src.GetCoverage)

	// Vistar End Points
	// Grab log files from maxtv server
	// https://maxtvmedia.com/cms/.cron/vistar_new/logs.php?from=2019-12-13&to=2019-12-14&event=get_assets&per-page=10&page=2'
	// GET
	r.HandleFunc("/vistar_logs_request", src.VistarLogs)

	r.HandleFunc("/jobs", src.JobsStatus)
	r.HandleFunc("/job/{hash}", src.JobStatus)
	//r.HandleFunc("/run_scan", src.RunScan)
	//
	//
	//r.HandleFunc("/asset_request", src.AssetsRequestHTTP)
	//r.HandleFunc("/asset_response/{belongs_to_request_id}", src.AssetsResponseHTTP)
	//
	//// Read logs to Database file -> DB
	//// Read One file and put to DB
	//r.HandleFunc("/vistar_proceed_get_assets", src.VistarProceedAssetRequestsLogFile)
	//// Read All files and put to DB (
	//r.HandleFunc("/vistar_proceed_get_assets_bulk", src.VistarProceedAssetRequestsLogFiles)
	//
	//// by post
	//r.HandleFunc("/vistar_assets_request", src.VistarAssetRequestHTTP)
	//r.HandleFunc("/vistar_assets_response", src.VistarAssetResponseHTTP)

	fmt.Printf("Starting Server to HANDLE maxtv.tech back end\nPort : " + src.Port + "\nAPI revision " + src.Version + "\n\n")
	if err := http.ListenAndServe(":"+src.Port, r); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Authorization := r.Header.Get("Authorization")
		if Authorization == "Bearer "+src.AuthorizationHash {
			next.ServeHTTP(w, r)
		}
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		src.FillAnswerHeader(w)
		src.OptionsAnswer(w)
		next.ServeHTTP(w, r)
	})
}
