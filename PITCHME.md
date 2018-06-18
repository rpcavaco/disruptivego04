# Disruptive Go

## Building a REST API with fasthttp (GO#04)

---

## Requirements

In order to reproduce code examples, you should have installed on your system:

- Go
- Redis, an in-memory datastructure store to use as our REST API backend

 
 Go can be dowloaded from [here](https://golang.org/dl/), and Redis from [here](https://redis.io/)

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

---?image=assets/fortunes_fasthttp_light.png&size=auto 75%

## Some benchmarking

Techempower.com is my favorite benchmark

Results available at
https://www.techempower.com/benchmarks/

---?image=assets/fortunes_fasthttp.png&position=bottom 50px right 100px&size=auto 45%

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
