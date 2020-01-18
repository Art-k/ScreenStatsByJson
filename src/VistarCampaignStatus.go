package src

func GetAppearanceOfCampaign(CampaignId uint) (int64, int64) {
	var vistarAsset VistarAsset
	DbVistar.Where("campaign_id = ?", CampaignId).First(&vistarAsset)

	var vistarAssetResponseFirst VistarAssetResponse
	DbVistar.Where("vistar_asset_id = ?", vistarAsset.ID).First(&vistarAssetResponseFirst)
	var vistarAssetRequestFirst VistarAssetsRequest
	DbVistar.Where("id = ?", vistarAssetResponseFirst.RequestID).First(&vistarAssetRequestFirst)

	var vistarAssetResponseLast VistarAssetResponse
	DbVistar.Where("vistar_asset_id = ?", vistarAsset.ID).First(&vistarAssetResponseLast)
	var vistarAssetRequestLast VistarAssetsRequest
	DbVistar.Where("id = ?", vistarAssetResponseLast.RequestID).First(&vistarAssetRequestLast)

	return vistarAssetRequestFirst.DisplayTime, vistarAssetRequestLast.DisplayTime
}

func GetListOfAllScreens() []Screen {
	var screens []Screen
	Db.Order("sys_id desc").Find(&screens)
	return screens
}
