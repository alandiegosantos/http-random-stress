# http-random-stress

This load generator tool generates requests using a [Poisson Distribution](https://en.wikipedia.org/wiki/Poisson_distribution), i.e, it generates random number of requests per second based on the rate informed. This simulates how users will interact with your system.

## How to compile it:

```
$ git clone https://github.com/alandiegosantos/http-random-stress
$ make -f build/Makefile
```

## How to use it

Use this command to requests from *https://localhost:34412/* using a rate of **10 rps** for **20s**
```
$ ./http-random-stress -rate 10 -url https://localhost:32412/ -timeout 20
```