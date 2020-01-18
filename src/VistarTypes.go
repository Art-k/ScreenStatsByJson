package src

import (
	"github.com/jinzhu/gorm"
)

// #####################################################################

var VistarLogsBaseUrl string

type IncVistarEventLog struct {
	MaxTvId string `json:"id"`
	Link    string `json:"link"`
	Type    string `json:"type"`
	Event   string `json:"event"`
}

type IncVistarEventLogPage struct {
	Total    int64               `json:"total"`
	Entities []IncVistarEventLog `json:"entities"`
}

// ===================================================================

var DbVistar *gorm.DB

type Model struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT"`
}

type VistarEventType struct {
	Model
	TypeName string           `json:"event_type"`
	Events   []VistarEventLog `json:"events"`
}

type VistarLogType struct {
	Model
	LogTypeName string           `json:"log_type"`
	Events      []VistarEventLog `json:"events"`
}

type VistarEventLog struct {
	Model
	MaxTvId   int64  `json:"maxtv_id"`
	Link      string `json:"link"`
	TypeID    int64  `json:"type"`
	EventID   int64  `json:"event"`
	Processed bool   `json:"processed";gorm:"default:false"`
}

//##################################################################################
type IncVistarAssetsRequestDisplayArea struct {
	Id             string   `json:"id"`
	Width          int      `json:"width"`
	Height         int      `json:"height"`
	MinDuration    int      `json:"min_duration"`
	SupportedMedia []string `json:"supported_media"`
	AllowAudio     bool     `json:"allow_audio"`
	StaticDuration int      `json:"static_duration"`
}

type IncVistarAssetsRequest struct {
	NetworkId        string                              `json:"network_id"`
	ApiKey           string                              `json:"api_key"`
	Duration         int                                 `json:"duration"`
	Interval         int                                 `json:"interval"`
	DeviceId         string                              `json:"device_id"`
	VenueId          string                              `json:"venue_id"`
	DisplayArea      []IncVistarAssetsRequestDisplayArea `json:"display_area"`
	DisplayTime      int64                               `json:"display_time"`
	DirectConnection bool                                `json:"direct_connection"`
}

type VistarAssetsRequest struct {
	Model
	ScreenID    uint  `json:"device_id"`
	DisplayTime int64 `json:"display_time"`
	HasResponse bool  `json:"has_non_empty_response"`
}

// ====================================================================

type VistarAsset struct {
	Model
	AssetId          string `json:"asset_id"`
	CreativeId       string `json:"creative_id"`
	CampaignId       int64  `json:"campaign_id"`
	AssetURL         string `json:"asset_url"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	MimeType         string `json:"mime_type"`
	LengthSec        int    `json:"length_in_seconds"`
	LengthMSec       int    `json:"length_in_milliseconds"`
	CreativeCategory string `json:"creative_category"`
	Advertiser       string `json:"advertiser"`
	CreativeName     string `json:"creative_name"`
}

type IncVistarAssetResponse struct {
	Asset []IncVistarAsset `json:"asset"`
}

type VistarAssetResponse struct {
	gorm.Model
	RequestID     int64
	VistarAssetID int64
}

type IncVistarAsset struct {
	AssetId          string `json:"asset_id"`
	CreativeId       string `json:"creative_id"`
	CampaignId       int64  `json:"campaign_id"`
	AssetURL         string `json:"asset_url"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	MimeType         string `json:"mime_type"`
	LengthSec        int    `json:"length_in_seconds"`
	LengthMSec       int    `json:"length_in_milliseconds"`
	CreativeCategory string `json:"creative_category"`
	Advertiser       string `json:"advertiser"`
	CreativeName     string `json:"creative_name"`
}

// ====================================================================

type VistarProofs struct {
	Model
	VistarID      int64 `json:"id"`
	AssetsID      int64 `json:"assets_id"`
	PopURL        int64 `json:"proof_of_play_url"`
	ExpURL        int64 `json:"expiration_url"`
	OrderID       int64 `json:"order_id"`
	LeaseExpiry   int64 `json:"lease_expiry"`
	DisplayAreaID int64 `json:"display_area_id"`
	DealID        int64 `json:"deal_id"`
	DisplayTime   int64 `json:"display_time"`
}

type IncVistarProofs struct {
	PopId                  string `json:"id"`
	Proof_of_play_url      string `json:"proof_of_play_url"`
	Expiration_url         string `json:"expiration_url"`
	Order_id               string `json:"order_id"`
	Display_time           int64  `json:"display_time"`
	Lease_expiry           int64  `json:"lease_expiry"`
	Display_area_id        string `json:"display_area_id"`
	Creative_id            string `json:"creative_id"`
	Asset_id               string `json:"asset_id"`
	Asset_url              string `json:"asset_url"`
	Width                  int    `json:"width"`
	Height                 int    `json:"height"`
	Mime_type              string `json:"mime_type"`
	Length_in_seconds      int    `json:"length_in_seconds"`
	Length_in_milliseconds int    `json:"length_in_milliseconds"`
	Campaign_id            int64  `json:"campaign_id"`
	Creative_category      string `json:"creative_category"`
	Advertiser             string `json:"advertiser"`
	Deal_id                string `json:"deal_id"`
}

type IncProofsResponse struct {
	Advertisement []IncVistarProofs `json:"advertisement"`
}
