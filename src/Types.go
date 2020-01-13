package src

import "github.com/jinzhu/gorm"

var Db *gorm.DB

const DbLogMode = false

var Err error
var Port string
var ScreenStatUrl string

const Version = "0.0.2"

var GetMaxTvStatInterval int64

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

type VLRec struct {
	gorm.Model
	MaxTvId int64  `json:"maxtv_id"`
	Link    string `json:"link"`
	Type    string `json:"type"`
	Event   string `json:"event"`
}

type IncomingVLRec struct {
	MaxTvId string `json:"id"`
	Link    string `json:"link"`
	Type    string `json:"type"`
	Event   string `json:"event"`
}

type IncomingVLPage struct {
	Total    int64           `json:"total"`
	Entities []IncomingVLRec `json:"entities"`
}

type LogsPagesDownloded struct {
	gorm.Model
	Link string
}

type VistarGetAssetsRequest struct {
	gorm.Model
	venue_id          string
	device_id         string
	display_time      int64
	direct_connection bool
}
