package job

import (
	"github.com/ovear/go-tumblr-crawler/config"
	"path"
	"github.com/qbhy/go-utils"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"sync"
	"os"
	"io"
	"net/http"
)

type DownloadParameter struct {
	Proxy config.ProxyConfig
	Url string
	FilePath string
	Filename string
	WaitGroupInstance *sync.WaitGroup
}

func DownloadFile(p DownloadParameter) {
	realPath := path.Join(p.FilePath, p.Filename)
	if exists, _ := utils.PathExists(realPath); exists {
		defer p.WaitGroupInstance.Done()
		return
	}
	request := gorequest.New()
	//res, body, err := request.Proxy(p.Proxy.Https).
	//	Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36").
	//	Get(p.Url).
	//	End()
	request = request.Proxy(p.Proxy.Https).
		Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36").
		Get(p.Url)
	req, err := request.MakeRequest()

	if err != nil {
		fmt.Println("下载失败:", p.Url, err, req)
		defer p.WaitGroupInstance.Done()
		return
	}

	client := &http.Client{Transport: request.Transport}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("下载失败:", p.Url, err, resp)
		defer p.WaitGroupInstance.Done()
		return
	}
	defer resp.Body.Close()

	output, err := os.Create(realPath)
	if err != nil {
		fmt.Println("下载失败:", p.Url, err)
		defer p.WaitGroupInstance.Done()
		return
	}
	defer output.Close()

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		fmt.Println("下载失败:", p.Url, err)
		defer p.WaitGroupInstance.Done()
		return
	}
	defer p.WaitGroupInstance.Done()
}

