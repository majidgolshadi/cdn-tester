package main

import (
	"strings"
	"fmt"
	"os"
	"net/http"
	"io"
	"time"
)

func downloadFromUrl(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	output, err := os.Create("download/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func main() {
	cdnUrl := "https://fs3.qmery.com/static_file/cdntest/test.jpg"
	originUrl := "http://cdn3.fs.soroush-hamrah.ir/static_file/cdntest/test.jpg"

	for i := 0; i < 300; i++ {
		downloadFromUrl(cdnUrl)
	}

	fmt.Printf("|%6s|%ds|", "cdn", "origin")
	fmt.Printf("|%6s|%ds|", "------", "------")
	for i := 0; i < 10; i++ {
		downloadFromCdnStartedAt := time.Now().UnixNano()
		downloadFromUrl(cdnUrl)
		downloadFromCdnCompletedAt := time.Now().UnixNano()

		downloadFromOriginStartedAt := time.Now().UnixNano()
		downloadFromUrl(originUrl)
		downloadFromOriginCompletedAt := time.Now().UnixNano()

		cdnTime := downloadFromCdnCompletedAt - downloadFromCdnStartedAt
		originTime := downloadFromOriginCompletedAt - downloadFromOriginStartedAt
		fmt.Printf("|%6d|%6d|\n", cdnTime, originTime)
	}
}