package loader

import (
	"fmt"
	"github.com/ovear/go-tumblr-crawler/config"
	"github.com/ovear/go-tumblr-crawler/config/global"
	"github.com/ovear/go-tumblr-crawler/db/sites"
)

func LoadSites(path string) (sites []config.SiteConfig) {
	if global.Settings.SiteLoading == "config" {
		sites = loadSitesFromConfig(path)
	} else if global.Settings.SiteLoading == "db" {
		sites = loadSitesFromDB()
	}
	fmt.Printf("站点加载成功 %+v\n", sites)
	return
}

func loadSitesFromConfig(path string) (sites []config.SiteConfig) {
	fmt.Printf("从配置文件加载站点配置文件：%s\n%+v\n", path)
	configLoader := config.NewConfig()
	configLoader.Load(path, &sites)
	return sites
}

func loadSitesFromDB() (s []config.SiteConfig) {
	s = sites.Sites()
	fmt.Printf("从数据库加载站点配置文件：%+v\n", s)
	return
}