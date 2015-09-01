package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	httpClient           *http.Client = &http.Client{}
	cachedWebsiteContent []byte
	websiteUrl           string
	emailAddress         string
	pollInterval         time.Duration
)

func loadConfig() {
	flag.StringVar(&websiteUrl, "url", "http://google.com", "Website url to check")
	flag.StringVar(&websiteUrl, "email", "your@email.com", "Notification email address")
	flag.DurationVar(&pollInterval, "poll", 5*time.Second, "Poll interval period")
	flag.Parse()
}

func IsTextDifferent(t1, t2 []byte) bool {
	return !bytes.Equal(t1, t2)
}

func FetchSiteContent(url string) ([]byte, error) {
	data, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	content, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func CheckContent(url string) error {
	content, err := FetchSiteContent(url)

	if err != nil {
		log.Printf("Error fetching content from [%s]: %s", url, err)
		return err
	}

	if len(content) == 0 {
		log.Printf("Warning [%s] is empty", url)
	}

	if IsTextDifferent(cachedWebsiteContent, content) {
		cachedWebsiteContent = content
		SendMail(emailAddress, fmt.Sprintf("url [%s] has changed", url))
	}

	return nil
}

func SendMail(email string, content string) {
	log.Printf("email sent to %s with content %s", email, content)
	// implement connection with external mail server and send
}

func main() {
	loadConfig()
	go func() {
		for _ = range time.NewTicker(pollInterval).C {
			log.Println("tick")
			CheckContent(websiteUrl)
		}
	}()

	log.Println("Content checking service started...")
	_ = http.ListenAndServe("localhost:8080", nil)
}
