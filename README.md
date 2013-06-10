## Cat Crawler

A webcrawler I'm writing in Google Go that I can use to find and download cat pictures.

### Current status

**In progress.**  This doesn't yet actually work. :-)  I am writing code and a few unit 
tests at a time.

### TODO

- Write main.go
  - Command line argument parsing
- Loop detection in URL crawler
- Create image crawler
  - Add search capability to alt and title tags
- Write instrumentation to detect how many goroutines are active/idle

### Installation

    git clone git@github.com:dmuth/cat-crawler.git
    
### Running the tests

    go get -v -a github.com/dmuth/procedural-webserver
    go test -i
    go test

You should see results like this:

    PASS
    ok      _/Users/doug/development/google-go/cat-crawler  0.024s

### Contact

Questions? Complaints? Here's my contact info: http://www.dmuth.org/contact




