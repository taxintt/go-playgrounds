# gzip

## normal request
```
🕙[ 00:19:46 ] ❯ curl localhost:3333 -v                    
*   Trying 127.0.0.1:3333...
* Connected to localhost (127.0.0.1) port 3333 (#0)
> GET / HTTP/1.1
> Host: localhost:3333
> User-Agent: curl/7.87.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: text/html
< Date: Sun, 10 Dec 2023 15:20:49 GMT
< Content-Length: 2000
< 
* Connection #0 to host localhost left intact
ce1e6681ec...
```

## request compression
```
🕙[ 00:20:49 ] ❯ curl -v http://localhost:3333 --compressed
*   Trying 127.0.0.1:3333...
* Connected to localhost (127.0.0.1) port 3333 (#0)
> GET / HTTP/1.1
> Host: localhost:3333
> User-Agent: curl/7.87.0
> Accept: */*
> Accept-Encoding: deflate, gzip
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Encoding: gzip
< Content-Type: text/html
< Vary: Accept-Encoding
< Date: Sun, 10 Dec 2023 15:20:54 GMT
< Content-Length: 1060
< 
* Connection #0 to host localhost left intact
448f1ae56...
```