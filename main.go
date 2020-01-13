package main

import (
	"./src"
	"fmt"
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

	go func() {
		src.DoEvery(time.Duration(src.GetMaxTvStatInterval)*time.Second, src.GetScreenStat)
	}()

	handleHTTP()
}

func handleHTTP() {

	http.HandleFunc("/get_statistic", src.GetStatistic)
	http.HandleFunc("/get_maxtv_statistic", src.GetMaxTVStatistic)

	fmt.Printf("Starting Server to HANDLE maxtv.tech back end\nPort : " + src.Port + "\nAPI revision " + src.Version + "\n\n")
	if err := http.ListenAndServe(":"+src.Port, nil); err != nil {
		log.Fatal(err)
	}
}
