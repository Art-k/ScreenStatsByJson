package src

import "github.com/jinzhu/gorm"

type Slices struct {
	gorm.Model
	Hash string
}

type MaxTVBuilding struct {
	gorm.Model
	MaxTvId string
	Name    string
	Address string
}

type SlicedBuilding struct {
	gorm.Model
	BuildingID uint
	Displays   []Screen
	SliceHash  string
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
	MaxTvId      string `json:"maxtv_id"`
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

type SlicedScreenProgrammaticStatus struct {
	gorm.Model
	ScreenID  uint
	Vistar    bool   `json:"vistar"`
	Campsite  bool   `json:"campsite"`
	Hivestack bool   `json:"hivestack"`
	SliceHash string `json:"slice_hash"`
}

type SlicedScreenImpressions struct {
	gorm.Model
	ScreenID       uint
	DwelTime       int     `json:"dwel_time"`
	TrafficPerWeek float64 `json:"traffic_per_week"`
	Impression     float64 `json:"impression"`
	SliceHash      string  `json:"slice_hash"`
}

type Screen struct {
	gorm.Model
	MaxTvId    string `json:"maxtv_id"`
	Name       string `json:"name"`
	SysId      string `json:"sysid"`
	Ads        []Ad   `json:"ads"`
	BuildingID uint   `json:"Building"`
	Spots      []Spot `json:"spots"`
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

type IncomingBuilding struct {
	Id       string           `json:"id"`
	Name     string           `json:"name"`
	Address  string           `json:"address"`
	Displays []IncomingScreen `json:"displays"`
}
