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


- **emailFrom** *Notification email sender address*
- **emailTo** *Notification email recipient address*
- **poll** *Poll interval period (default 5s)*
- **smtpHost** *SMTP Host*
- **smtpPass** *SMTP Password*
- **smtpPort** *SMTP Port (default "25")*
- **smtpUser** *SMTP Username*
- **urls** *Website url(s) to check*

To run with flags in working directory

Please ensure you enter correct smtp settings, and correct access rights required on some servers

> $ ./website-change-notifier -urls=http://yoururl.com -emailTo=your@email.com -emailFrom=another@email.com -smtpHost=smtp.gmail.com -smtpUser=your@username -smtpPass=yourpassword -smtpPort=yourport -poll 10s

####

#### Check multiple sites

There is now support for multiple sites, just seperate with a comma

> $ ./website-change-notifier -urls=http://yoururl.com,https://www.google.com,http://www.amazon.com

#### TODO

- custom poll period for each site
- implement sending mail
