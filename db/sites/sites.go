package sites

import (
	"github.com/ovear/go-tumblr-crawler/config"
	"github.com/ovear/go-tumblr-crawler/db"
	"fmt"
	"github.com/ovear/go-tumblr-crawler/config/global"
	"github.com/ovear/go-tumblr-crawler/utils"
)

func Sites() (sites []config.SiteConfig) {
	conn := db.Connection()
	//site
	//video_download
	//photo_download
	//video_db
	//photo_db
	sql := "select site, video_download, photo_download, video_db, photo_db from %s"
	sql = fmt.Sprintf(sql, global.Settings.DBSitesTable)
	result, err := conn.Query(sql)
	if (err != nil) {
		fmt.Printf("[DB]执行数据库指令失败[%s][%s]\n", sql, err)
		panic(err)
	}
	for result.Next() {
		var site string
		var video_download int
		var photo_download int
		var video_db int
		var photo_db int

		result.Scan(&site, &video_download, &photo_download, &video_db, &photo_db)
		sites = append(sites, config.SiteConfig{Site: site, VideoDownload: utils.IntToBool(video_download), PhotoDownload: utils.IntToBool(photo_download),
			VideoDB: utils.IntToBool(video_db), PhotoDB: utils.IntToBool(photo_db)})
	}
	return

}
