package src

import "github.com/jinzhu/gorm"

type VistarAsset struct {
	gorm.Model
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
	RequestID     uint
	VistarAssetID uint
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
