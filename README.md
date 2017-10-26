# WhoAmI

Tiny Go WebServer that prints OS information (Host, IP, ENV) and HTTP request to output.

```sh
$ docker run -d -p 8080:80 --name iamwho z0r1k/whoami

$ curl "http://0.0.0.0:8080"
Hostname: Jakku
IP: 127.0.0.1
IP: ::1
IP: fe80::1
IP: 192.168.99.1
ENV: TERM=xterm-256color
ENV: LANG=en_US
ENV: SHELL=/bin/bash
ENV: TMPDIR=/var/folders/8f/w83bx8k16s5dp9lfb2s2nl000000gn/T/
ENV: HOME=/Users/z0r1k
GET / HTTP/1.1
Host: localhost:8080
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.8,ru;q=0.6,de;q=0.4,uk;q=0.2
Cache-Control: max-age=0
Connection: keep-alive
Cookie: JSESSIONID=zZSrpsQPH9VWoKQsZI-ynkVI8LhOXZyeb0rUm8XR
Dnt: 1
Upgrade-Insecure-Requests: 1
```

`http://0.0.0.0:8080/api` will get you the same response but in the JSON format.