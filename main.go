package main

import (
	"strings"
	"fmt"
	"os"
	"net/http"
	"io"
	"time"
	"flag"
)

func downloadFromUrl(url string) (int64, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	output, err := os.Create("download/" + fileName)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func main() {
	var cdnUrl = flag.String("cdn", "", "cdn url")
	var try = flag.Int("t", 300, "download to cache in cdn edges")
	var originUrl = flag.String("origin", "", "origin url")
	var compareTime = flag.Int("c", 10, "compare time")

	flag.Parse()

	if *cdnUrl == "" || *originUrl == "" {
		println("Usage: -cdn <CDN_URL> -t <CDN_DOWNLOAD_TIME> -origin <ORIGIN_URL> -c <COMPARE_TIME>")
		return
	}

	fmt.Printf("\n origin url: %s\n cdn url: %s\n cdn download time: %d\n compare time: %d\n\n", *cdnUrl, *originUrl, *try, *compareTime)

	for i := 0; i < *try; i++ {
		if _, err := downloadFromUrl(*cdnUrl); err != nil {
			fmt.Println("Error while downloading ", cdnUrl, " - ", err.Error())
		} else {
			fmt.Println(i+1, "downloaded")
		}
	}

	fmt.Printf("\n|%12s|%12s|\n", "cdn", "origin")
	fmt.Printf("|%12s|%12s|\n", "------------", "------------")
	var err error
	for i := 0; i < *compareTime; i++ {
		downloadFromCdnStartedAt := time.Now().UnixNano()
		_, err = downloadFromUrl(*cdnUrl)
		downloadFromCdnCompletedAt := time.Now().UnixNano()

		if err != nil{
			println(err.Error())
			continue
		}

		downloadFromOriginStartedAt := time.Now().UnixNano()
		_, err = downloadFromUrl(*originUrl)
		downloadFromOriginCompletedAt := time.Now().UnixNano()

		if err != nil{
			println(err.Error())
			continue
		}

		cdnTime := downloadFromCdnCompletedAt - downloadFromCdnStartedAt
		originTime := downloadFromOriginCompletedAt - downloadFromOriginStartedAt
		fmt.Printf("|%12d|%12d|\n", cdnTime, originTime)
	}
}