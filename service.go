package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	httpClient   *http.Client = &http.Client{}
	urlFlag      UrlFlag
	emailAddress string
	pollInterval time.Duration
)

type UrlFlag struct {
	urls []*url.URL
}

type Website struct {
	Url           *url.URL
	CachedContent []byte
}

func loadConfig() {
	flag.Var(&urlFlag, "urls", "Website url(s) to check")
	flag.StringVar(&emailAddress, "email", "your@email.com", "Notification email address")
	flag.DurationVar(&pollInterval, "poll", 5*time.Second, "Poll interval period")
	flag.Parse()
}

func (uf *UrlFlag) String() string {
	return fmt.Sprint(uf.urls)
}

func (uf *UrlFlag) Set(value string) error {
	urlSlice := strings.Split(value, ",")
	for _, item := range urlSlice {
		if parsedUrl, err := url.Parse(item); err != nil {
			return err
		} else {
			uf.urls = append(uf.urls, parsedUrl)
		}
	}
	return nil
}

func (w *Website) FetchSiteContent() []byte {
	data, err := httpClient.Get(w.Url.String())
	if err != nil {
		log.Printf("Error fetching content from [%s]: %s", w.Url, err)
		return nil
	}
	defer data.Body.Close()

	content, err := ioutil.ReadAll(data.Body)
	if len(content) == 0 {
		log.Printf("Warning [%s] is empty", w.Url)
	}
	return content
}

func (w *Website) HasChanged(content []byte) bool {
	if bytes.Equal(w.CachedContent, content) {
		return false
	}
	return true
}

func (website *Website) StartNotifier() {
	go func() {
		for _ = range time.NewTicker(pollInterval).C {
			// log here for demo purposes, could put a verbose debug flag to reduce noise
			log.Println("...checking ", website.Url)
			content := website.FetchSiteContent()
			changed := website.HasChanged(content)
			if changed {
				website.CachedContent = content
				SendMail(emailAddress, fmt.Sprintf("'url [%s] has changed'", website.Url))
			}
		}
	}()
}

func SendMail(email string, content string) {
	log.Printf("email sent to %s with content %s", email, content)
	// implement connection with external mail server and send
}

func main() {
	loadConfig()
	for _, aUrl := range urlFlag.urls {
		website := &Website{Url: aUrl}
		website.StartNotifier()
	}

	log.Println("Content checking service started...")
	_ = http.ListenAndServe("localhost:8080", nil)
}
