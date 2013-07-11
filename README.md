## Cat Crawler

A webcrawler I'm writing in Google Go that I can use to find and download cat pictures.


### Current status

**Sorta working** 

Images are now downloaded and saved in the downloads/ directory, though rather indiscriminately.

I'll work on alt and title tag searching next.

### TODO

- Create image crawler
  - Add search capability to alt and title tags
- Rate limiting by domain in URL crawler
	- I could have an array of key=domain, value=count and a goroutine 
		that decrements count regularly
		- Could get a bit crazy on the memory, though!
- Write instrumentation to detect how many goroutines are active/idle
	- GoStatStart(key)
	- GoStatStop(key)
	- go GoStatDump(interval)


### Installation

    git clone git@github.com:dmuth/cat-crawler.git

    
### Running the crawler
	go run ./*.go [--seed-url url[,url[,url[...]]]] [ --num-connections n ] [--allow-urls [url,[url,[...]]]]
		--seed-url What URL to start at? More than one URL may be 
			specified in comma-delimited format.
		--num-connections How many concurrent connections?
		--allow-urls If specified, only URLs starting with the URLs listed here are crawled


### Examples
    go build && ./cat-crawler --seed-url cnn.com --num-connections 1
Get top stories. :-)

    go build && ./cat-crawler --seed-url (any URL) --num-connections 1000
This will saturate your download bandwidth. Seriously, don't do it.

    go build && ./cat-crawler --seed-url cnn.com  --num-connections 1 --allow-urls cnn.com
Don't leave CNN's website

    go build && ./cat-crawler --seed-url cnn.com  --num-connections 1 --allow-urls foobar
After crawling the first page, nothing will happen.  Oops.


### Running the tests

    go get -v -a github.com/dmuth/procedural-webserver
    go test -i
    go test

You should see results like this:

    PASS
    ok      _/Users/doug/development/google-go/cat-crawler  0.024s


### Contact

Questions? Complaints? Here's my contact info: http://www.dmuth.org/contact



