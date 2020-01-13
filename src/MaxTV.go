package src

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type StatGetAttempt struct {
	gorm.Model
	Hash string
}

type Stat struct {
	gorm.Model
	ScreensCount             int
	TotalAdsSpaceSeconds     int64
	TotalCoveredSpaceSeconds int64
	Attempt                  string

	VistarSpotsS    int
	CampsiteSpotsS  int
	HivestackSpotsS int
	MaxTVSpotsS     int
	VistarSpots     int
	CampsiteSpots   int
	HivestackSpots  int
	MaxTVSpots      int

	TotalSpotsS         int
	TotalCoverageSpotsS int
}

type Video struct {
	gorm.Model
	File           string `json:"file"`
	DurationMs     int    `json:"duration_ms"`
	RealDurationMs int    `json:"real_duration_ms"`
}

type Ad struct {
	gorm.Model
	MaxtvId      string `json:"maxtv_id"`
	CompanyId    string `json:"company_id"`
	ParentId     string `json:"parent_id"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	CampaignDate string `json:"campaign_date"`
	Type         string `json:"type"`
	EndDate      string `json:"end_date"`
	Name         string `json:"name"`
	File         string `json:"file"`
	DurationMs   int    `json:"duration_ms"`
	Spot         string `json:"spot"`
	Exclusive    string `json:"exclusive"`
	CategoryId   string `json:"category_id"`
	Screen       uint   `json:"Screen"`
	Attempt      string `json:"Attempt"`
	VideoID      uint   `json:"VideoID"`
}

type IncomingAd struct {
	ID           string `json:"id"`
	CompanyId    string `json:"company_id"`
	ParentId     string `json:"parent_id"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	CampaignDate string `json:"campaign_date"`
	Type         string `json:"type"`
	EndDate      string `json:"end_date"`
	Name         string `json:"name"`
	File         string `json:"file"`
	DurationMs   int    `json:"duration_ms"`
	Spot         string `json:"spot"`
	Exclusive    string `json:"exclusive"`
	CategoryId   string `json:"category_id"`
}

type Spot struct {
	gorm.Model
	SpotCode string
	ScreenID uint
}

type Screen struct {
	gorm.Model
	MaxtvId        string  `json:"maxtv_id"`
	Name           string  `json:"name"`
	SysId          string  `json:"sysid"`
	Vistar         bool    `json:"vistar"`
	Campsite       bool    `json:"campsite"`
	Hivestack      bool    `json:"hivestack"`
	DwelTime       int     `json:"dwel_time"`
	TrafficPerWeek float64 `json:"traffic_per_week"`
	Impression     float64 `json:"impression"`
	Ads            []Ad    `json:"ads"`
	Building       uint    `json:"Building"`
	Attempt        string  `json:"Attempt"`
	Spots          []Spot  `json:"spots"`
}

type IncomingScreen struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Sysid          string       `json:"sysid"`
	Vistar         bool         `json:"vistar"`
	Campsite       bool         `json:"campsite"`
	Hivestack      bool         `json:"hivestack"`
	DwelTime       int          `json:"dwelTime"`
	TrafficPerWeek float64      `json:"trafficPerWeek"`
	Impression     float64      `json:"impression"`
	Ads            []IncomingAd `json:"ads"`
}

type Building struct {
	gorm.Model
	MaxtvId  string
	Name     string
	Address  string
	Displays []Screen
	Attempt  string
}

type IncomingBuilding struct {
	Id       string           `json:"id"`
	Name     string           `json:"name"`
	Address  string           `json:"address"`
	Displays []IncomingScreen `json:"displays"`
}

func GetScreenStat(t time.Time) {

	fmt.Println("Timer >>>")
	GetMaxTvData()

}

func GetMaxTvData() {
	fmt.Println()
	ScreenStatUrl := os.Getenv("SCREEN_STAT_URL")
	resp, err := http.Get(ScreenStatUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(string(body))

	type incomingJson struct {
		Total    int                `json:"total"`
		Entities []IncomingBuilding `json:"entities"`
	}

	var response incomingJson
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return
	}

	var Attempt StatGetAttempt
	Attempt.Hash = GetHash()

	var stat Stat
	stat.Attempt = Attempt.Hash

	for ind, record := range response.Entities {
		fmt.Println(ind)
		var building Building
		building.MaxtvId = record.Id
		building.Name = record.Name
		building.Address = record.Address
		building.Attempt = Attempt.Hash
		Db.Create(&building)

		for _, display := range record.Displays {

			stat.ScreensCount++

			var screen Screen
			screen.MaxtvId = display.ID
			screen.Name = display.Name
			screen.SysId = display.Sysid
			screen.Vistar = display.Vistar
			screen.Campsite = display.Campsite
			screen.Hivestack = display.Hivestack
			screen.DwelTime = display.DwelTime
			screen.TrafficPerWeek = display.TrafficPerWeek
			screen.Impression = display.Impression
			screen.Building = building.ID
			screen.Attempt = Attempt.Hash
			Db.Create(&screen)

			for _, adv := range display.Ads {
				var ad Ad

				var video Video
				Db.Where("file = ?", adv.File).Find(&video)
				if video.ID == 0 {
					video.File = adv.File
					video.DurationMs = adv.DurationMs
					Db.Create(&video)
				}

				ad.MaxtvId = adv.ID
				ad.CompanyId = adv.Name
				ad.ParentId = adv.ParentId
				ad.Title = adv.Title
				ad.Status = adv.Status
				ad.CampaignDate = adv.CampaignDate
				ad.Type = adv.Type
				ad.EndDate = adv.EndDate
				ad.Name = adv.Name
				ad.File = adv.File
				ad.DurationMs = adv.DurationMs
				ad.Spot = adv.Spot

				sl, _ := strconv.Atoi(adv.Spot)
				stat.TotalCoverageSpotsS = stat.TotalCoverageSpotsS + sl

				var spot Spot
				spot.ScreenID = screen.ID

				switch ad.Name {
				case "Vistar":
					stat.VistarSpots++
					stat.VistarSpotsS = stat.VistarSpotsS + sl
					spot.SpotCode = "V" + strconv.FormatUint(uint64(video.ID), 10)
					break
				case "CAMPSITE":
					stat.CampsiteSpots++
					stat.CampsiteSpotsS = stat.CampsiteSpotsS + sl
					spot.SpotCode = "C" + strconv.FormatUint(uint64(video.ID), 10)
					break
				case "Hivestack":
					stat.HivestackSpots++
					stat.HivestackSpotsS = stat.HivestackSpotsS + sl
					spot.SpotCode = "H" + strconv.FormatUint(uint64(video.ID), 10)
					break
				default:
					stat.MaxTVSpots++
					stat.MaxTVSpotsS = stat.MaxTVSpotsS + sl
					spot.SpotCode = "M" + strconv.FormatUint(uint64(video.ID), 10)
					break
				}

				Db.Create(&spot)

				ad.Exclusive = adv.Exclusive
				ad.CategoryId = adv.CategoryId
				ad.Screen = screen.ID
				ad.Attempt = Attempt.Hash
				ad.VideoID = video.ID

				Db.Create(&ad)
			}
		}
	}
	stat.TotalSpotsS = stat.ScreensCount * 360
	Db.Create(&stat)
	Db.Create(&Attempt)
}

func GetMaxTVStatistic(w http.ResponseWriter, r *http.Request) {
	GetMaxTvData()
	w.WriteHeader(http.StatusOK)
	n, _ := fmt.Fprintf(w, "")
	fmt.Println(n)
}
