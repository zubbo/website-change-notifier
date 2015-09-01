# Website change notifier (v0.0.1)

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

- url
- email
- poll

To run with flags in working directory

$ ./website-change-notifier -url "http://yoururl.com" -email "your@email.com" -poll 10s


#### TODO

- to scale design to multiple sites, can either run bash script with multiple run flags *OR* can add text file (CSV) functionality to load as multiple go routines
