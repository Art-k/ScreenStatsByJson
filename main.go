package main

import (
	Src "./src"
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
	Src.Port = os.Getenv("PORT")
	Src.ScreenStatUrl = os.Getenv("SCREEN_STAT_URL")
	Src.GetMaxTvStatInterval, _ = strconv.ParseInt(os.Getenv("GET_MAXTV_STAT_INTERVAL_SECOND"), 10, 64)

	Src.Db, Src.Err = gorm.Open("sqlite3", "screen_stat.db")
	if Src.Err != nil {
		panic("failed to connect database")
	}
	defer Src.Db.Close()
	Src.Db.LogMode(Src.DbLogMode)

	Src.Db.AutoMigrate(&Src.Ad{})
	Src.Db.AutoMigrate(&Src.Screen{})
	Src.Db.AutoMigrate(&Src.Building{})
	Src.Db.AutoMigrate(&Src.Building{})
	Src.Db.AutoMigrate(&Src.Stat{})
	Src.Db.AutoMigrate(&Src.StatGetAttempt{})
	Src.Db.AutoMigrate(&Src.Video{})
	Src.Db.AutoMigrate(&Src.Spot{})

	go func() {
		Src.DoEvery(time.Duration(Src.GetMaxTvStatInterval)*time.Second, Src.GetScreenStat)
	}()

	handleHTTP()
}

func handleHTTP() {

	http.HandleFunc("/get_statistic", Src.GetStatistic)
	http.HandleFunc("/get_maxtv_statistic", Src.GetMaxTVStatistic)

	fmt.Printf("Starting Server to HANDLE maxtv.tech back end\nPort : " + Src.Port + "\nAPI revision " + Src.Version + "\n\n")
	if err := http.ListenAndServe(":"+Src.Port, nil); err != nil {
		log.Fatal(err)
	}
}
