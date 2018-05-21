package config

type Setting struct {
	SiteLoading      string `yaml:"site_loading"`
	SavePath         string `yaml:"save_path"`
	DownloadThread   int    `yaml:"download_thread"`
	DownloadPhoto    bool   `yaml:"download_photo"`
	DownloadVideo    bool   `yaml:"download_video"`
	DownloadRawPhoto bool   `yaml:"download_raw_photo"`
	EnableLoop       bool   `yaml:"enable_loop"`
	LoopInterval     int    `yaml:"loop_interval"`
	EnableDB         bool   `yaml:"enable_db"`
	DBHost           string `yaml:"db_host"`
	DBPort           int    `yaml:"db_port"`
	DBDatabase       string `yaml:"db_database"`
	DBUsername       string `yaml:"db_username"`
	DBPassword       string `yaml:"db_password"`
	DBTable          string `yaml:"db_table"`
	DBSitesTable     string `yaml:"db_sites_table"`
}
