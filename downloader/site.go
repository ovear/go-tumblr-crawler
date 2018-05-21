package downloader

import (
	"github.com/qbhy/go-utils"
	"path"
	"os"
	"strings"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"encoding/xml"
	response2 "github.com/ovear/go-tumblr-crawler/downloader/response"
	"math/big"
	"github.com/ovear/go-tumblr-crawler/config"
	myutils "github.com/ovear/go-tumblr-crawler/utils"
)

func NewSite(site config.SiteConfig, config config.ProxyConfig) *Site {
	return &Site{
		Site:        site,
		ProxyConfig: config,
	}
}

const (
	BaseUrl    = "https://{site}.tumblr.com/api/read?type={mediaType}&num={num}&start={start}"
	PageNumber = 50
)

func GenerateMediaUrl(site string, mediaType string, num int64, start int64) string {
	mediaUrl := strings.Replace(BaseUrl, "{site}", site, -1)
	mediaUrl = strings.Replace(mediaUrl, "{mediaType}", mediaType, -1)
	mediaUrl = strings.Replace(mediaUrl, "{num}", big.NewInt(num).String(), -1)
	mediaUrl = strings.Replace(mediaUrl, "{start}", big.NewInt(start).String(), -1)
	//fmt.Println("site:", site)
	//fmt.Println("mediaType:", mediaType)
	//fmt.Println("num:", num)
	//fmt.Println("start:", start)
	//fmt.Println("mediaUrl:", mediaUrl)
	return mediaUrl
}

type Site struct {
	Site        config.SiteConfig
	ProxyConfig config.ProxyConfig
	currentPath string
	sitePath    string
	videoPath   string
	photoPath   string
}

func (this *Site) StartDownload() {
	this.Init()

	if this.Site.VideoDownload || this.Site.VideoDB {
		WaitGroupInstance.Add(1)
		go this.DownloadVideo()
	}
	if this.Site.PhotoDownload || this.Site.PhotoDB {
		WaitGroupInstance.Add(1)
		go this.DownloadPhoto()
	}
}

func (this *Site) Init() {
	this.currentPath = path.Join(myutils.CurrentPath(), "files")

	if exists, _ := utils.PathExists(this.currentPath); !exists {
		os.Mkdir(this.currentPath, 0755)
	}

	this.sitePath = path.Join(this.currentPath, this.Site.Site)

	if exists, _ := utils.PathExists(this.sitePath); !exists {
		os.Mkdir(this.sitePath, 0755)
	}
}

func (this *Site) DownloadVideo() {
	this.videoPath = path.Join(this.sitePath, "video")
	if exists, _ := utils.PathExists(this.videoPath); !exists {
		os.Mkdir(this.videoPath, 0755)
	}
	this.DownloadMedia("video", 0)
}

func (this *Site) DownloadPhoto() {
	this.photoPath = path.Join(this.sitePath, "photo")
	if exists, _ := utils.PathExists(this.photoPath); !exists {
		os.Mkdir(this.photoPath, 0755)
	}
	this.DownloadMedia("photo", 0)
}

func (this *Site) DownloadMedia(mediaType string, start int64) {

	for {

		mediaUrl := GenerateMediaUrl(this.Site.Site, mediaType, PageNumber, start)

		request := gorequest.New().Proxy(this.ProxyConfig.Https)
		res, responseString, err := request.Get(mediaUrl).End()
		fmt.Printf("开始抓取site[%+v] start[%d] mediaUrl[%s]\n", this.Site, start, mediaUrl)

		if err != nil || res.StatusCode == 404 {
			fmt.Printf("下载site[%+v]时发生错误res[%+v] error[%+v] mediaUrl[%s]\n", this.Site, res, err, mediaUrl)
			break
		}

		if mediaType == "video" {
			video := response2.NewVideo()
			err := xml.Unmarshal([]byte(responseString), &video)
			if err != nil {
				fmt.Printf("error: %v", err)
				break
			} else if len(video.Posts.Post) <= 0 {
				fmt.Printf("抓取结束[%+v] mediaUrl[%s]\n", this.Site.Site, mediaUrl)
				break
			}

			downloadVideos(this, video.Posts)
		} else {
			photo := response2.NewPhoto()
			err := xml.Unmarshal([]byte(responseString), &photo)
			if err != nil {
				fmt.Printf("error: %v", err)
				break
			} else if len(photo.Posts.Post) <= 0 {
				fmt.Println("没有更多内容了 ", this.Site.Site)
				break
			}

			downloadPhotos(this, photo.Posts)
		}
		start += PageNumber
	}

	defer WaitGroupInstance.Done()
}
