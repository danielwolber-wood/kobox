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

Based on input/output, there are basically nine pathways:

```txt
input |   output  | pipeline
url   ->  RR        url -> webpage -> RR
url   ->  epub      url -> webpage -> RR -> epub
url   ->  dropbox   url -> webpage -> RR -> epub -> dropbox
page  ->  RR        webpage -> RR
page  ->  epub      webpage -> RR -> epub
page  ->  dropbox   webpage -> RR -> epub -> dropbox
RR    ->  epub      RR -> epub
RR    ->  dropbox   RR -> epub -> dropbox
epub  ->  dropbox   epub -> dropbox
```

So I think my API is going to look something like this:

```txt
# Direct transformations
POST /convert/url-to-rr
POST /convert/url-to-epub  
POST /convert/url-to-dropbox
POST /convert/page-to-rr
POST /convert/page-to-epub
POST /convert/page-to-dropbox
POST /convert/rr-to-epub
POST /convert/rr-to-dropbox
POST /convert/epub-to-dropbox

# Individual steps
POST /fetch     url -> webpage
POST /extract   webpage -> rr 
POST /generate  rr -> epub
POST /upload    epub -> dropbox
```

Rather than having a separate Dropbox Connector, it might make sense to just add a single GET endpoint to this one

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