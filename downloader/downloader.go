package downloader

import (
	"github.com/ovear/go-tumblr-crawler/downloader/response"
	"sync"
	"strings"
	"path/filepath"
	"github.com/ovear/go-tumblr-crawler/pool"
	"github.com/ovear/go-tumblr-crawler/downloader/job"
	"github.com/ovear/go-tumblr-crawler/utils"
	"github.com/ovear/go-tumblr-crawler/config/global"
	"github.com/ovear/go-tumblr-crawler/db"
)

var WaitGroupInstance sync.WaitGroup

func init() {
}

func downloadVideos(site *Site, posts response.VideoPosts) {
	for _, post := range posts.Post {
		if ok, videoUrl := post.ParseVideoUrl(); ok {
			//下载原视频
			//TOOD 增加选项
			if utils.IsNumeric(filepath.Base(videoUrl)) {
				pos := strings.LastIndex(videoUrl, "/")
				videoUrl = videoUrl[:pos]
			}
			if site.Site.VideoDB {
				//fmt.Printf("%+v\n", post.BasePost)
				db.Insert(post.BasePost, videoUrl)
			}
			filename := filepath.Base(videoUrl) + ".mp4"
			filename = strings.Replace(filename, " ", "", -1)

			if global.Settings.DownloadVideo && site.Site.VideoDownload {
				WaitGroupInstance.Add(1)
				dp := job.DownloadParameter{site.ProxyConfig, videoUrl, site.videoPath, filename, &WaitGroupInstance}
				pool.NewWork(dp)
			}
		}
	}
}
//http://data.tumblr.com/4a762c9aed5983bb14d11a3455865b22/tumblr_o5i2srwgyY1qaxnseo1_raw.jpg
func downloadPhotos(site *Site, posts response.PhotoPosts) {
	for _, post := range posts.Post {
		for _, url := range post.ParsePhotosUrl() {
			if site.Site.PhotoDB {
				//fmt.Printf("%+v\n", post.BasePost)
				db.Insert(post.BasePost, url)
			}
			if global.Settings.DownloadPhoto && site.Site.PhotoDownload {
				if global.Settings.DownloadRawPhoto {
					url = utils.ConvertToRawTumblrAddress(url)
				}
				WaitGroupInstance.Add(1)
				dp := job.DownloadParameter{site.ProxyConfig, url, site.photoPath, filepath.Base(url), &WaitGroupInstance}
				pool.NewWork(dp)
			}
		}
	}
}