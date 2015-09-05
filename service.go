package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"
)

var (
	httpClient       *http.Client = &http.Client{}
	urlFlag          UrlFlag
	pollInterval     time.Duration
	recipientAddress string
	senderAddress    string
	smtpUser         string
	smtpPass         string
	smtpHost         string
	smtpPort         string
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
	flag.DurationVar(&pollInterval, "poll", 5*time.Second, "Poll interval period")
	flag.StringVar(&recipientAddress, "emailTo", "", "Notification email recipient address")
	flag.StringVar(&senderAddress, "emailFrom", "", "Notification email sender address")
	flag.StringVar(&smtpUser, "smtpUser", "", "SMTP Username")
	flag.StringVar(&smtpPass, "smtpPass", "", "SMTP Password")
	flag.StringVar(&smtpHost, "smtpHost", "", "SMTP Host")
	flag.StringVar(&smtpPort, "smtpPort", "25", "SMTP Port")
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
				go SendMail(recipientAddress, fmt.Sprintf("[%s] content changed", website.Url))
			}
		}
	}()
}

func SendMail(recipient string, subject string) {
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	to := []string{recipient}
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"This site has recently changed\r\n", recipient, subject))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderAddress, to, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s:%s email sent to %s with subject '%s'", smtpHost, smtpPort, recipient, subject)
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
