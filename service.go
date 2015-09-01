package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	httpClient           *http.Client = &http.Client{}
	cachedWebsiteContent []byte
	defaultUrl           string = "http://www.google.com"
	defaultEmail         string = "your@email.com"
)

func main() {
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
		SendMail(defaultEmail, fmt.Sprintf("url %s has changed", url))
	}
	return nil
}

func SendMail(email string, content string) {
	log.Printf("email sent to %s with content %s", email, content)
	// implement connection with mail server
}
