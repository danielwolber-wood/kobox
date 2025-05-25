# kobox

**kobox** is a set of two related programs for use with Kobo eReaders: 
1. A server that accepts links and webpages, converts them to an Epub, and then uploads them to the user's Dropbox account which syncs them to the user's eReader
2. A browser extension that communicates with the server, allowing one-click "Send to Reader" functionality

## TODO
* Deployment:
    * Dockerization
    * compose.yml for easy distribution to end users
    * CI/CD
    * Automated testing
* Actual deployment onto GCP:
    * Real deployment into the real world for real users
    * Real, actual HTTPS instead of a self-signed cert
    * Sign in with Google flow for user identification
    * Multiuser version
* KoboxMono:
    * Better configuration management
    * Better handling of failed tasks