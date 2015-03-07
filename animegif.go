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
	"time"
)

func main() {
	keyword, count := args()
	urls := fetchUrls(keyword, count)
	html := generateHtml(urls)
	openHtml(html)
}

func args() (keyword string, count int) {
	flag.StringVar(&keyword, "k", "yuyushiki", "keyword")
	flag.IntVar(&count, "c", 8, "count")
	flag.Parse()

	return keyword, count
}

func fetchUrls(keyword string, count int) (urls []string) {
	page := 1
	var _urls []string
	for len(urls) <= count {
		_urls = search(page, keyword)
		for _, url := range _urls {
			urls = append(urls, url)
		}
		page += 1
	}

	return urls
}

func generateHtml(urls []string) (html string) {
	html = "<!DOCTYPE HTML><html><body>"
	for _, url := range urls {
		html = html + "<a href='" + url + "' target='_blank'><img src='" + url + "' /></a>"
	}
	html = html + "</body></html>"
	return html
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

	params := url.Values{}
	params.Add("q", keyword)
	params.Add("rsz", fmt.Sprint(perPage))
	params.Add("safe", "off")
	params.Add("v", "1.0")
	params.Add("as_filetype", "gif")
	params.Add("imgsz", "large")
	params.Add("start", fmt.Sprint(start))
	params.Add("as_sitesearch", "tumblr.com")

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
