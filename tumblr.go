package main

import (
	"github.com/qbhy/go-utils"
	"path"
	"fmt"
	"github.com/ovear/go-tumblr-crawler/downloader"
	myutils "github.com/ovear/go-tumblr-crawler/utils"
	"runtime"
	config2 "github.com/ovear/go-tumblr-crawler/config"
	"time"
	"github.com/ovear/go-tumblr-crawler/pool"
	"github.com/ovear/go-tumblr-crawler/config/global"
	"github.com/ovear/go-tumblr-crawler/db"
	"github.com/ovear/go-tumblr-crawler/config/loader"
)

func main() {

	configLoader := config2.NewConfig()
	currentPath := myutils.CurrentPath()
	if currentPath == "" {
		fmt.Println("无法获取当前目录")
		return
	}
	var sites []config2.SiteConfig
	var proxies config2.ProxyConfig

	//获取程序运行时设置
	settingPath := path.Join(currentPath, "settings.yaml")
	if exists, _ := utils.PathExists(settingPath); exists {
		fmt.Printf("加载配置文件：%s\n", settingPath)
		configLoader.LoadYaml(settingPath, &global.Settings)
		fmt.Printf("%+v\n", global.Settings)
	} else {
		fmt.Printf("找不到配置文件：%s\n", settingPath)
	}

	if global.Settings.EnableDB {
		result := db.Init()
		if !result {
			fmt.Println("数据库初始化失败")
			return
		}
		result = db.Test()
		if !result {
			fmt.Println("数据库测试失败")
			return
		}
		fmt.Println("数据库测试成功")
	}

	// 获取代理配置
	proxyPath := path.Join(currentPath, "proxies.json")
	if exists, _ := utils.PathExists(proxyPath); exists {
		fmt.Printf("加载代理配置文件：%s\n", proxyPath)
		proxies = config2.ProxyConfig{}
		configLoader.Load(proxyPath, &proxies)
		fmt.Println(proxies)
	}

	// 获取站点配置
	sitesPath := path.Join(currentPath, "sites.json")
	if exists, _ := utils.PathExists(sitesPath); !exists && global.Settings.SiteLoading == "config" {
		fmt.Printf("站点配置文件[%s]不存在\n", sitesPath)
		return
	}
	sites = loader.LoadSites(sitesPath)

	// 设置最大协程数
	maxProcesses := runtime.NumCPU() //获取cpu个数
	runtime.GOMAXPROCS(maxProcesses) //限制同时运行的goroutines数量

	//启动Worker
	fmt.Println("启动Worker:", global.Settings.DownloadThread)
	for i := 0; i < global.Settings.DownloadThread; i++ {
		pool.StartWorker(i)
	}

	fmt.Println("CPU数量:", maxProcesses)

	for {
		// 如果所有线程退出，就会退出阻塞
		if len(sites) > 0 {
			for _, site := range sites {
				siteInstance := downloader.NewSite(site, proxies)
				siteInstance.StartDownload()
			}

			downloader.WaitGroupInstance.Wait()
		} else {
			fmt.Println("没有配置站点")
		}

		//是否启用循环模式
		if global.Settings.EnableLoop {
			fmt.Printf("采集完成，等待%d秒\n", global.Settings.LoopInterval)
			time.Sleep(time.Duration(global.Settings.LoopInterval) * time.Second)
		} else {
			break
		}
		sites = loader.LoadSites(sitesPath)
	}
	fmt.Println("采集完成")



}
