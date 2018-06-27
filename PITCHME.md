# Disruptive Go

## Building a REST API with fasthttp (GO#04)

<small>This presentation and exercise code available in GitHub repo [rpcavaco/disruptivego04](https://github.com/rpcavaco/disruptivego04)</small>
---

## Requirements

In order to reproduce code examples, you should have installed on your system:

- Go
- fasthttp library from Aliaksandr Valialkin
- radix.v2 or other Redis Go Client
- Redis, an in-memory datastructure store to use as our REST API backend

<small>We can use a a Redis cloud server - like the ones available from [Redis Labs](https://redislabs.com/)</small>

---

Go can be dowloaded from [here](https://golang.org/dl/), and Redis from [here](https://redis.io/).

Fasthttp library available as GitHub go repository in [https://github.com/valyala/fasthttp](https://github.com/valyala/fasthttp)

<small>To install it, execute this command line (after Go's installation):</small>

```
go get -u github.com/valyala/fasthttp
```

Redis Go client available as GitHub go repository in [https://github.com/mediocregopher/radix.v2](https://github.com/mediocregopher/radix.v2)

<small>To install it, execute this command line:</small>

```
go get github.com/mediocregopher/radix.v2
```
---

## Why use Golang to program HTTP services?

### Other languages might be more popular and effective

@ul

- Python is very easy to write and maintain
- <code>C++</code> is both robust and efficient
- Java: extremely efficient with a solid knowledge base
- ASP.NET is as good as Java
- PHP is hugely popular and simple

@ulend

---

## Why use Golang to program HTTP services?

We may try to get something out of Python 3 or, we are just in time to give new language a try! ;)

<small>Where should we look for it?</small>

---?image=assets/img/fortunes_fasthttp_light.png&size=auto 75%

## Some benchmarking

Techempower.com is my favorite benchmark

Results available at
https://www.techempower.com/benchmarks/

---?image=assets/img/fortunes_fasthttp.png&position=bottom 50px right 100px&size=auto 45%

Looking at its results (May 2017) I first saw Go as a top performing language.<br><br>

As I’m a huge fan of PostgreSQL, <b>fasthttp-postgresql</b> immediately caught my attention.<br>
<br>
<small>First Python framework comes 46th</small><br>
<small>First PHP at 56th place ... </small><br>

<small>(In latest tests, fasthttp is not so well ...)</small>

---

## Nice features of Go and fasthttp lib for HTTP REST

- efficiency
- simplicity
- good straightforward JSON support out-of-the-box

---

## What is fasthttp?

Fasthttp is an improvement to Go’s net/http package.

Package author states it delivers 3x to 4x better perfomance than standard net/http lib.
See this [Google Groups discussion](https://groups.google.com/forum/m/#!topic/golang-nuts/OaQu4QezAr0/discussion) with the author.

---

### A little smell of **net/http**

```go
http.Handle("/foo", fooHandler)

http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))
```
@[7](ListenAndServe initiates an HTTP server at the given address, this case localhost:8080, and 'multiplexer' function.)
@[7](Passing nil as second parameter, internal multiplexer *DefaultServeMux* (not shown) is used.)
@[1-4](*Handle* and *HandleFunc* add handlers to *DefaultServeMux*.)

---

### The 'multiplexer'

```go
http.Handle("/foo", fooHandler)

http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))
```

<small>The 'multiplexer' routes the request to a given *handler*, matching the request path with a pattern like "/foo" or "/bar".</small>

<small>*Handle* and *HandleFunc* are two different ways of automatically adding handler functions to the *DefaultServeMux* multiplexer</small>

---

#### Building the first net/http example

<small>Let's create an empty directory (name it as you wish). Inside it create our first Go source named *basichttp.go* containing the following code:</small>

```go
package main

import (
	"fmt"
	"net/http"
	"html"
	"log"
)

func main() {
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
---

<small>To run the previous code, justo do:</small>

```
> go build basichttp.go
> ./basichttp
```

---?image=assets/img/output_basic.png&position=bottom 20px right 100px&size=auto 35%

#### Running first example

<small>Posting http://localhost:8080 on a browser's address bar, you should get something as the image shows.</small>

Our first Golang web server is running!

---

### Now, let's go a little faster

---

#### A simple HTTP server on *fasthttp*

```go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
)

func main() {
	flag.Parse()

	h := requestMux

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
```
@[1-9](Define package and imports -- including Aleksandr's library)
@[11-13](Get command line argument containing TCP port to use, default 8080)
@[15-23](*main* function)
@[16](Actual command line argument parsing)
@[18-22](Defining *ListenAndServe*, similarly to net/http's)

---

##### A multiplexer associating paths with handlers

```go
func requestMux(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	switch string(path) {
		case "/hello":
			helloHandler(ctx)
		case "/raw":
			rawHandler(ctx)
		default:
			fmt.Fprintf(ctx, "HTTP not found: %s", ctx.Path())
			ctx.SetContentType("text/plain; charset=utf8")
	}
}
```

---

##### The handlers

```go
func helloHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	ctx.SetContentType("text/plain; charset=utf8")
}

func rawHandler(ctx *fasthttp.RequestCtx) {

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")
}
```

@[1-16](For "/hello")
@[18-23](For "/raw")

---

##### Testing

Full example source available in [rpcavaco/disruptivego04/exercises](https://github.com/rpcavaco/disruptivego04/exercises/basicfhttp.go)

<small>To run the previous code, justo do:</small>

```
> go build basicfhttp.go
> ./basicfhttp
```
---

### A REST API example

Let's expand the previous example to build a very simple REST API.

<small>It would be nicer if we achieve some server persistence for information items bouncing from client to server and back.</small>

We'll use here Redis, an in memory key-value store very easy to use.

---?image=assets/img/redis_downloads.png&position=bottom 20px right 100px&size=auto 35%

#### A fresh REDIS installation

In case you haven't done it previously, now is the time to have a REDIS server installed in your system.

[redis.io/download](https://redis.io/download)


<small>You can use a cloud-based Redis server instead</small>
