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

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		keyword string
		count   int
	)
	flag.StringVar(&keyword, "k", "yuyushiki", "keyword")
	flag.IntVar(&count, "c", 8, "count")
	flag.Parse()

	page := 1
	var (
		urls  []string
		_urls []string
	)
	for len(urls) <= count {
		_urls = search(page, keyword)
		for _, url := range _urls {
			urls = append(urls, url)
		}
		page += 1
	}

	html := "<!DOCTYPE HTML><html><body>"
	for _, url := range urls {
		html = html + "<a href='" + url + "' target='_blank'><img src='" + url + "' /></a>"
	}
	html = html + "</body></html>"

	file, err := ioutil.TempFile(os.TempDir(), "animegif")
	perror(err)
	ioutil.WriteFile(file.Name(), []byte(html), 0644)
	exec.Command("open", file.Name()).Output()
	perror(err)
	time.Sleep(time.Second * 1)

	defer os.Remove(file.Name())
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
	per_page := 8
	base := "http://ajax.googleapis.com/ajax/services/search/images?"
	start := (page-1)*per_page + 1

	params := url.Values{}
	params.Add("q", keyword)
	params.Add("rsz", fmt.Sprint(per_page))
	params.Add("safe", "off")
	params.Add("v", "1.0")
	params.Add("as_filetype", "gif")
	params.Add("imgsz", "large")
	params.Add("start", fmt.Sprint(start))
	params.Add("as_sitesearch", "tumblr.com")

	body := openUrl(base + params.Encode())

	var response ResponseType
	err := json.Unmarshal(body, &response)
	perror(err)
	for _, value := range response.ResponseData.Results {
		urls = append(urls, value.Url)
	}
	return urls
}

func openUrl(req string) (body []byte) {
	res, err := http.Get(req)
	perror(err)
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	perror(err)
	return body
}
