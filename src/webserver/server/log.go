package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/kjk/dailyrotate"
	"github.com/kjk/siser"
)

var (
	httpLogDailyFile *dailyrotate.File
	httpLogSiser     *siser.Writer
)

// HTTPReqInfo represents HTTP request info
type HTTPReqInfo struct {
	method    string
	url       string
	referer   string
	ipaddr    string
	code      int
	size      int64
	duration  time.Duration
	userAgent string
}

// OpenHTTPLog opens a http logger
func OpenHTTPLog() {
	var err error
	dir := getHTTPLogDirMust()
	path := filepath.Join(dir, "2006-01-02.txt")
	httpLogDailyFile, err = dailyrotate.NewFile(path, nil)
	if err != nil {
		panic(err)
	}
	httpLogSiser = siser.NewWriter(httpLogDailyFile)
}

func closeHTTPLog() {
	_ = httpLogDailyFile.Close()
	httpLogDailyFile = nil
	httpLogSiser = nil
}

var (
	muLogHTTP sync.Mutex
)

func logHTTPReq(ri *HTTPReqInfo) {
	var rec siser.Record
	rec.Name = "httplog"
	rec.Append("method", ri.method)
	rec.Append("uri", ri.url)
	if ri.referer != "" {
		rec.Append("referer", ri.referer)
	}
	rec.Append("ipaddr", ri.ipaddr)
	rec.Append("code", strconv.Itoa(ri.code))
	rec.Append("size", strconv.FormatInt(ri.size, 10))
	durMs := ri.duration / time.Millisecond
	rec.Append("duration", strconv.FormatInt(int64(durMs), 10))
	rec.Append("ua", ri.userAgent)

	muLogHTTP.Lock()
	defer muLogHTTP.Unlock()
	_, _ = httpLogSiser.WriteRecord(&rec)
	log.Printf("[%s] HTTP Request { %v }\n", rec.Timestamp.Format("2006-01-02 15:04:05"), rec.Entries)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func logf(format string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Print(format)
		return
	}
	fmt.Printf(format, args...)
}

func makeDirMust(dir string) string {
	err := os.MkdirAll(dir, 0755)
	must(err)
	return dir
}

func getHTTPLogDirMust() string {
	dir := filepath.Join("./log_http")
	return makeDirMust(dir)
}
