package src

import "github.com/jinzhu/gorm"

var Db *gorm.DB

const DbLogMode = true

var Err error
var Port string
var ScreenStatUrl string

const Version = "0.0.1"

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
	attempt                  string
}

type Video struct {
	gorm.Model
	file        string
	duration_ms int
}

type Ad struct {
	gorm.Model
	maxtv_id      string
	company_id    string
	parent_id     string
	title         string
	status        string
	campaign_date string
	Type          string `json:"type"`
	end_date      string
	name          string
	file          string
	duration_ms   int
	spot          string
	exclusive     string
	category_id   string
	Screen        uint
	attempt       string
}

type IncomingAd struct {
	id            string
	company_id    string
	parent_id     string
	title         string
	status        string
	campaign_date string
	Type          string `json:"type"`
	end_date      string
	name          string
	file          string
	duration_ms   int
	spot          string
	exclusive     string
	category_id   string
}

type Screen struct {
	gorm.Model
	maxtv_id       string
	name           string
	sysid          string
	vistar         bool
	campsite       bool
	hivestack      bool
	dwelTime       int
	trafficPerWeek float64
	impression     float64
	ads            []Ad
	Building       uint
	attempt        string
}

type IncomingScreen struct {
	id             string
	name           string
	sysid          string
	vistar         bool
	campsite       bool
	hivestack      bool
	dwelTime       int
	trafficPerWeek float64
	impression     float64
	ads            []Ad
}

type Building struct {
	gorm.Model
	maxtv_id string
	name     string
	address  string
	displays []Screen
	attempt  string
}

type IncomingBuilding struct {
	id       string
	name     string
	address  string
	displays []Screen
}
