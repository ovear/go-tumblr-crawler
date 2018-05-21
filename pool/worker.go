package pool

import (
	"fmt"
	"github.com/ovear/go-tumblr-crawler/downloader/job"
)

var jobs chan job.DownloadParameter

func init() {
	jobs = make(chan job.DownloadParameter, 200)
}

func StartWorker(id int) {
	go worker(id)
}

func worker(id int) {
	for j := range jobs {
		fmt.Printf("[worker%d] download[%s]\n", id, j.Filename)
		job.DownloadFile(j)
	}
}

func NewWork(p job.DownloadParameter) {
	jobs <- p
}