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

### Running httpr in Docker
To run **httpr** in Docker:

  ```docker run -p 8081:8081 netbucket/httpr log -e -p```
  
### Installing httpr binary
### On macOS:
1. Download the macOS binary:

   ```curl -L0 https://sourceforge.net/projects/httpr/files/release/1.0.0/darwin/httpr -o httpr```
  
  
2. Make it executalbe:

   ```chmod 755 httpr```

### On Linux:
1. Download the Linux binary:

   ```curl -L0 https://sourceforge.net/projects/httpr/files/release/1.0.0/linux/httpr -o httpr```
  
  
2. Make it executalbe:

   ```chmod 755 httpr```

### On Windows:
1. Download the Windows binary:

   ```curl -L0 https://sourceforge.net/projects/httpr/files/release/1.0.0/windows/httpr -o httpr.exe```

### More to come
Work is underway to make **httpr** also available via *brew* for macOS, and *apt-get* for Linux.
 
## Logging HTTP Requests
To log incoming HTTP requests to standard output, use the `httpr log` command. Note that by default, **httpr** will start the HTTP server on port 8081. See `httpr help log` for more options.
 
## Returning a Specific HTTP Status Code
To return a specific HTTP status code back to the HTTP client, use the *-r code* option. For instance, to return the HTTP Service Unavailable code, use `httpr log -r 503`

## Simulating Latency
To simulate a delay in returning the HTTP response to the client, use the *-d millis* option. For instance, to simulate 500 millisecond latency, use `httpr log -d 500`

## Simulating Transient HTTP Failures
**httpr** can make it easy to simulate transient HTTP errors. This is useful when testing HTTP retry logic, or the circuit breaking capabilities in HTTP clients (see https://martinfowler.com/bliki/CircuitBreaker.html). To do that, use the *-f* option. The *-f* option supports additional modifiers to specicfy exactly how the transient failures should be simulated. For instance, to simulate a series of 5 transient failures that return HTTP status 503, followed by 10 successful responses with status code 200, use:

   ```httpr log -f --simulate-failure-count=5  --simulate-failure-code=503 --simulate-success-count=10```
  
 Note that *-f* and *-d* can be used together to simulate latency and transient errors at once.
 
 ## Proxying to Simulate Latency and Transient Failures
 Using **httpr**, it is easy to simulate latency or transient failures in front of an existing HTTP based endpoint. To do that, use the `httpr proxy` command.
 For instance, to log and then proxy HTTP requests to `https://www.google.com`, while simulating a transient failure, use:

   ```httpr proxy https://www.google.com -f```
   
## TLS/HTTPS Support
To start **httpr** server in HTTPS mode, use the *-t* option. By default, **httpr** will generate and use
a self-signed certificate, and print the PEM-encoded certificate to the console. To supply your own
certificate, use the *--tls-cert-file* and *--tls-key-file* options with the *-t* flag to specify
the path/name of the certificate file and the private key file.

To ingore upstream TLS errors when proxying HTTPS requests with *httpr proxy*, use the *-k* flag.


