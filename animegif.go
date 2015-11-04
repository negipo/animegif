package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const defaultCount = 8
const maxCount = 80
const defaultKeyword = "yuyushiki"
const defaultPage = 1

func main() {
	keyword, count, page := args()
	urls := fetchUrls(keyword, count, page)
	if len(urls) == 0 {
		fmt.Println("no image found")
		return
	}

	html := generateHtml(urls)
	openHtml(html)
}

func args() (keyword string, count, page int) {
	flag.StringVar(&keyword, "k", defaultKeyword, "keyword")
	flag.IntVar(&count, "c", defaultCount, "count")
	flag.IntVar(&page, "p", defaultPage, "initial page")
	flag.Parse()

	if count > maxCount {
		count = maxCount
	}

	return keyword, count, page
}

func fetchUrls(keyword string, count, page int) (urls []string) {
	var _urls []string
	for len(urls) < count {
		_urls = search(page, keyword)

		if len(_urls) == 0 {
			return urls
		}

		urls = append(urls, _urls...)
		page += 1
	}

	return urls
}

func generateHtml(urls []string) (html string) {
	htmls := []string{"<!DOCTYPE HTML><html><body>"}
	for _, url := range urls {
		htmls = append(htmls, "<a href='"+url+"' target='_blank'><img src='"+url+"' /></a>")
	}
	htmls = append(htmls, "</body></html>")
	return strings.Join(htmls, "")
}

func openHtml(html string) {
	file, err := ioutil.TempFile(os.TempDir(), "animegif")
	printError(err)
	ioutil.WriteFile(file.Name(), []byte(html), 0644)
	exec.Command("open", file.Name()).Start()
	time.Sleep(time.Second * 1)

	defer os.Remove(file.Name())
}

func printError(err error) {
	if err != nil {
		panic(err)
	}
}

type ResultType struct {
	Url string `json:"url"`
}

type ResponseDataType struct {
	Results []ResultType `json:"results"`
}

type ResponseType struct {
	ResponseData ResponseDataType `json:"responseData"`
}

func search(page int, keyword string) (urls []string) {
	perPage := 8
	base := "http://ajax.googleapis.com/ajax/services/search/images?"
	start := (page-1)*perPage + 1

	params := url.Values{
		"q":             {keyword},
		"rsz":           {fmt.Sprint(perPage)},
		"safe":          {"off"},
		"v":             {"1.0"},
		"as_filetype":   {"gif"},
		"imgsz":         {"large"},
		"start":         {fmt.Sprint(start)},
		"as_sitesearch": {"tumblr.com"},
	}

	body := openUrl(base + params.Encode())

	var response ResponseType
	err := json.Unmarshal(body, &response)
	printError(err)
	for _, value := range response.ResponseData.Results {
		urls = append(urls, value.Url)
	}
	return urls
}

func openUrl(req string) (body []byte) {
	res, err := http.Get(req)
	printError(err)
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	printError(err)
	return body
}
