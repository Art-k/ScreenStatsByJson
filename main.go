package main

import (
	"./src"
	"fmt"
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
	//src.ScreenStatUrl = os.Getenv("SCREEN_STAT_URL")
	src.GetMaxTvStatInterval, _ = strconv.ParseInt(os.Getenv("GET_MAXTV_STAT_INTERVAL_SECOND"), 10, 64)

	src.Db, src.Err = gorm.Open("sqlite3", "screen_stat.db")
	if src.Err != nil {
		panic("failed to connect database")
	}

	defer src.Db.Close()
	src.Db.LogMode(src.DbLogMode)

	src.Db.AutoMigrate(&src.Ad{})
	src.Db.AutoMigrate(&src.Screen{})
	src.Db.AutoMigrate(&src.Building{})
	src.Db.AutoMigrate(&src.Building{})
	src.Db.AutoMigrate(&src.Stat{})
	src.Db.AutoMigrate(&src.StatGetAttempt{})
	src.Db.AutoMigrate(&src.Video{})
	src.Db.AutoMigrate(&src.Spot{})
	src.Db.AutoMigrate(&src.LogsPagesDownloded{})
	src.Db.AutoMigrate(&src.VLRec{})
	src.Db.AutoMigrate(&src.VistarGetAssetsRequest{})
	src.Db.AutoMigrate(&src.VistarAsset{})
	src.Db.AutoMigrate(&src.VistarAssetResponse{})

	go func() {
		src.DoEvery(time.Duration(src.GetMaxTvStatInterval)*time.Second, src.GetScreenStat)
	}()

	handleHTTP()
}

func handleHTTP() {

	r := mux.NewRouter()
	r.Use(authMiddleware)
	r.Use(headerMiddleware)

	// Jsons end points
	r.HandleFunc("/coverage", src.GetCoverage)

	// Vistar End Points
	// Grab log files from maxtv server
	// https://maxtvmedia.com/cms/.cron/vistar_new/logs.php?from=2019-12-13&to=2019-12-14&event=get_assets&per-page=10&page=2'
	// GET
	r.HandleFunc("/vistar_logs", src.VistarLogs)

	// Show the progress
	// GET
	r.HandleFunc("/vistar_logs_stat", src.VistarLogsStat)

	// Read logs to Database file -> DB
	r.HandleFunc("/vistar_file_db", src.VistarFileDB)
	r.HandleFunc("/vistar_file_db_bulk", src.VistarFileDBBulk)

	fmt.Printf("Starting Server to HANDLE maxtv.tech back end\nPort : " + src.Port + "\nAPI revision " + src.Version + "\n\n")
	if err := http.ListenAndServe(":"+src.Port, r); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		src.FillAnswerHeader(w)
		src.OptionsAnswer(w)
		next.ServeHTTP(w, r)
	})
}
