## Cat Crawler

A webcrawler I'm writing in Google Go that I can use to find and download cat pictures.

### Current status

**In progress.**  This doesn't yet actually work. :-)  I am writing code and a few unit 
tests at a time.

### TODO

- Create image crawler
  - Add search capability to alt and title tags
- Loop detection in URL crawler
	- urls[domain][uri]
- Rate limiting by domain in URL crawler
	- I could have an array of key=domain, value=count and a goroutine that decrements count regularly
		- Could get a bit crazy on the memory, though!
- Write instrumentation to detect how many goroutines are active/idle
	- GoStatStart(key)
	- GoStatStop(key)
	- go GoStatDump(interval)

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




