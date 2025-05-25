# kobox

**kobox** is a simple program for following RSS feeds, automatically fetching articles, and uploading them to Dropbox
for use with compatible Kobo e-readers. 

This is a demonstration of my ability to program a microservice-architecture program rather than the simplest approach;
for example, for the "Send this Website to Kobo" function, it would be trivial to instead use a simple user-script to
scrapes the current webpage, extract the contents, convert it to an .EPUB file, and point a Dropbox Saver at that .EPUB.

## Architectural Thoughts

The primary functionality is an RSS feed crawler that gets posts, converts them reader files, and uploads them to
Dropbox. In addition, there is a desktop browser extension that allows one-click "Send this Article to Kobo"
functionality.

I think there would likely be a few services:

* **Gateway**: a unified endpoint for end-applications such as the feed manager and the browser extension
* **Crawler**: crawls RSS feeds for new posts
* **Bypasser**: checks whether the post is cached or if the request needs to be redirected
* **Fetcher**: fetches HTML from RSS feeds
* **Processor**: uses readability.js to extract the body and title from the webpage
* **Assembler**: converts the extracted components into a styled .EPUB file
* **Uploader**: uploaders the .EPUB to Dropbox via the HTTP API.

### Version 2

There are really four programs:
1. A Dropbox Connector API that accepts files and uploads them to Dropbox, and can get the list of files currently uploaded to Dropbox
2. A Downloader/Processor that takes URLs, parses the content, and produces an EPUB
3. A feed reader/tracker/crawler that keeps track of feeds, waits for new posts, etc
4. A browser extension which works for processing individual webpages, turning them into epubs, and sending them to Dropbox

#### Downloader/Processor

Basically, the pipeline looks like this:

```txt
url -> full page -> readability object -> minimal html -> epub -> dropbox
```

Where the API accepts any of the steps leading up to Dropbox, and processes them all the way down to uploading them to
dropbox. Since the minimal html -> epub is so minimal, I combine them into a single "Generate" step.

## TODO

* Caching
* Redirection
* Dropbox integration
* Crawling + Crawl Tracking
* End-to-End API
* OPML integration
* Configuration file
* Download tracking
* Feed storage
* There are really two types of RSS feeds: those which include the content of the post, and those which include a link
  to the post. It maye be better to find an external RSS feed parser
 
* Proper handling of secrets/env configuration
* Dockerization
* Test
* CI/CD
* Real, Actual HTTPS instead of a self-signed cert