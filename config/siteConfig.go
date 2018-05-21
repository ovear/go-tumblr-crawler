package config

type SiteConfig struct {
	Site          string `json:"site"`
	VideoDownload bool   `json:"video_download"`
	PhotoDownload bool   `json:"photo_download"`
	VideoDB       bool   `json:"video_db"`
	PhotoDB       bool   `json:"photo_db"`
}