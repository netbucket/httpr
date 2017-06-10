# HTTP Rake - httpr
HTTP Rake - **httpr** - is a lightweight and flexible HTTP testing tool that is useful for examining  and testing HTTP requests without the need for a full-fledged Web server or a proxy. Implemented in Go, **httpr** is extremely compact, and can be run locally, which is useful in cases when sending HTTP traffic to a third-party hosted solutions is not desireable, or impractical. 

It provides several capabilties that are useful when testing distributed HTTP-based interactions:
 * Log incoming requests in raw or JSON format
 * Simulate latency
 * Simulate a transient HTTP failure
 * Return a specific HTTP status code to the HTTP client
 * Act as a proxy capable of simulating latency and/or transient failures in front of the upstream service
 
 See **httpr help** for more information
 
 ## Installing httpr
 Currently, the quickest way to install **httpr** is to *go get github.com/netbucket/httpr* followed by *go install*. 
 Work is underway to make **httpr** available via *brew* for macOS, *DockerHub*, and *apt-get* for Linux.
 
 
