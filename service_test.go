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

	Describe("IsTextDifferent", func() {
		text := []byte(`this is some text`)
		differentText := []byte(`this is some different text`)

		It("should return false when inputs are same", func() {
			result := IsTextDifferent(text, text)
			Expect(result).To(Equal(false))
		})

		It("should return true when inputs are different", func() {
			result := IsTextDifferent(text, differentText)
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
			content, _ := FetchSiteContent("http://www.google.com")
			Expect(len(content)).To(Equal(13))
		})
	})

	Describe("CheckContent", func() {
		var previousHttpClient *http.Client
		originalCache := []byte(`this is the original site cache`)

		BeforeEach(func() {
			cachedWebsiteContent = originalCache

			previousHttpClient = httpClient
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					w.WriteHeader(200)
					fmt.Fprintln(w, `new content`)
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
			cachedWebsiteContent = originalCache
		})

		It("should update cached content when changed", func() {
			_ = CheckContent("http://www.google.com")
			expectedCacheString := "new content"
			Expect(len(cachedWebsiteContent)).To(Equal(len(expectedCacheString) + 1))
		})

		It("should return an error when error with fetching content", func() {
			err := CheckContent("asdf")
			Expect(err).Should(HaveOccurred())

			currentCache := cachedWebsiteContent
			Expect(currentCache).To(Equal(originalCache))
		})
	})
})
