# kobox

**kobox** is a set of three related programs for use with Kobo eReaders: 
1. A server that accepts links and webpages, converts them to an Epub, and then uploads them to the user's Dropbox account which syncs them to the user's eReader
2. A browser extension that communicates with the server, allowing one-click "Send to Reader" functionality
3. An OPML-based feed reader and tracker, for automatically sending posts from followed feeds to the server

## TODO
* Caching
* Redirection
* Crawling + Crawl Tracking
* OPML integration
* Configuration file
* Download tracking
* Feed storage
* Proper handling of secrets/env configuration
* Dockerization
* Tests
* CI/CD
* Real, Actual HTTPS instead of a self-signed cert
* Real deployment into the real world for real users