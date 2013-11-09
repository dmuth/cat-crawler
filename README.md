## Cat Crawler

A webcrawler I'm writing in Google Go that I can use to find and download cat pictures.

### Installation

- Make sure your GOPATH environment variable is set up properly:
   `export GOPATH=$HOME/golib`
- Make sure the bin directory is in your path:
   `PATH=$PATH:$GOPATH/bin`
- Now install the package
   `go get -v github.com/dmuth/cat-crawler`

### Running the crawler
    cat-crawler [--seed-url url[,url[,url[...]]]] [ --num-connections n ] [--allow-urls [url,[url,[...]]]] [--search-string cat]
        --seed-url What URL to start at? More than one URL may be 
            specified in comma-delimited format.
        --num-connections How many concurrent connections?
        --search-string A string we want to search for in ALT and TITLE attributes on images
        --allow-urls If specified, only URLs starting with the URLs listed here are crawled
        --stats Print out stats once a second using my stats package

### Examples
    cat-crawler --seed-url cnn.com --num-connections 1
Get top stories. :-)

    cat-crawler --seed-url (any URL) --num-connections 1000
This will saturate your download bandwidth. Seriously, don't do it.

    cat-crawler --seed-url cnn.com  --num-connections 1 --allow-urls cnn.com
Don't leave CNN's website

    cat-crawler --seed-url cnn.com  --num-connections 1 --allow-urls foobar
After crawling the first page, nothing will happen.  Oops.

### Sequence diagram

![Sequence Diagram](https://raw.github.com/dmuth/cat-crawler/master/docs/sequence-diagram.png "Sequence Diagram")


### Development

    go get -v github.com/dmuth/cat-crawler && cat-crawler [options]


### Running the tests

    go get -v -a github.com/dmuth/procedural-webserver # Dependency
    go test -v github.com/dmuth/cat-crawler

You should see results like this:

    === RUN TestSplitHostnames
    --- PASS: TestSplitHostnames (0.00 seconds)
    === RUN TestHtmlNew
    --- PASS: TestHtmlNew (0.00 seconds)
    === RUN TestHtmlBadImg
    --- PASS: TestHtmlBadImg (0.00 seconds)
    === RUN TestHtmlLinksAndImages
    --- PASS: TestHtmlLinksAndImages (0.00 seconds)
    === RUN TestHtmlNoLinks
    --- PASS: TestHtmlNoLinks (0.00 seconds)
    === RUN TestHtmlNoImages
    --- PASS: TestHtmlNoImages (0.00 seconds)
    === RUN TestHtmlNoLinksNorImages
    --- PASS: TestHtmlNoLinksNorImages (0.00 seconds)
    === RUN TestHtmlPortNumberInBaseUrl
    --- PASS: TestHtmlPortNumberInBaseUrl (0.00 seconds)
    === RUN TestGetFilenameFromUrl
    --- PASS: TestGetFilenameFromUrl (0.00 seconds)
    === RUN Test
    --- PASS: Test (0.00 seconds)
    === RUN TestFilterUrl
    --- PASS: TestFilterUrl (0.00 seconds)
    === RUN TestIsUrlAllowed
    --- PASS: TestIsUrlAllowed (0.00 seconds)
    PASS
    ok      github.com/dmuth/cat-crawler    0.037s


### Depdendencies

This repo uses other packages I wrote:
- [log4go](https://github.com/dmuth/google-go-log4go)
- [golang-stats](https://github.com/dmuth/golang-stats)


### TODO

- Rate limiting by domain in URL crawler
	- I could have an array of key=domain, value=count and a goroutine 
		that decrements count regularly
		- Could get a bit crazy on the memory, though!
- Write instrumentation to detect how many goroutines are active/idle
	- GoStatStart(key)
	- GoStatStop(key)
	- go GoStatDump(interval)


### Contact

Questions? Complaints? Here's my contact info: http://www.dmuth.org/contact



