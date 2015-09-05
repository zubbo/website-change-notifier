package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unit Test: Website change notifier", func() {

	Describe("HasChanged", func() {
		oldContent := []byte(`this is some text`)
		newContent := []byte(`this is some different text`)
		var website *Website

		BeforeEach(func() {
			aUrl, _ := url.Parse("http://google.com")
			website = &Website{Url: aUrl, CachedContent: oldContent}
		})

		It("should return false when inputs are same", func() {
			result := website.HasChanged(oldContent)
			Expect(result).To(Equal(false))
		})

		It("should return true when inputs are different", func() {
			result := website.HasChanged(newContent)
			Expect(result).To(Equal(true))
		})
	})

	Describe("FetchSiteContent", func() {
		var previousHttpClient *http.Client

		BeforeEach(func() {
			previousHttpClient = httpClient
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					w.WriteHeader(200)
					fmt.Fprintln(w, `SOME CONTENT`)
				}
			}))
			transport := &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(server.URL)
				},
			}

			httpClient = &http.Client{Transport: transport}
		})

		AfterEach(func() {
			httpClient = previousHttpClient
		})

		It("should return content", func() {
			aUrl, _ := url.Parse("http://google.com")
			website := &Website{Url: aUrl}
			content := website.FetchSiteContent()
			Expect(len(content)).To(Equal(13))
		})
	})
})
