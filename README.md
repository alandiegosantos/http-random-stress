# http-random-stress

Run this to generate your certs: 

```
openssl req -new -x509 -config assets/example.org.cnf -keyout assets/server.key -out assets/server.crt -days 365
```