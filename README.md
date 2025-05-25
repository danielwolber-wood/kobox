# kobox

**kobox** is a set of two related programs for use with Kobo eReaders: 
1. A server that accepts links and webpages, converts them to an Epub, and then uploads them to the user's Dropbox account which syncs them to the user's eReader
2. A browser extension that communicates with the server, allowing one-click "Send to Reader" functionality

## TODO
* Dockerization
* Tests
* CI/CD
* Real, actual HTTPS instead of a self-signed cert
* Real deployment into the real world for real users
* Web portal for auth/configuration rather than CLI (basically needed for dockerization)
* Multi-user GCP-hosted version of the application
* Better configuration management
* Better handling of failed tasks