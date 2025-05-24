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
* Handle Atom syndication feeds http://www.w3.org/2005/Atom https://datatracker.ietf.org/doc/html/rfc4287