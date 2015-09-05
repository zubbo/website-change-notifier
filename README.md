# Website change notifier (v0.0.2)

A basic notification tool, that polls specified websites and notifies of changes via email


---

#### Add testing suite

Using Ginkgo testing suite for unit tests

- $ go get github.com/onsi/ginkgo/ginkgo
- $ go get github.com/onsi/gomega
- $ go test *this should pass*

#### Install project

- $ go install github.com/zubbo/website-change-notifier
- $ *OR* go build *if GO set up differently*

#### Run project

Website notifier has the following flags

- urls
- email
- poll

To run with flags in working directory

$ ./website-change-notifier -urls=http://yoururl.com -email "your@email.com" -poll 10s

#### Check multiple sites

There is now support for multiple sites, just seperate with a comma

$ ./website-change-notifier -urls=http://yoururl.com,https://www.google.com,http://www.amazon.com

#### TODO

- custom poll period for each site
- implement sending mail
