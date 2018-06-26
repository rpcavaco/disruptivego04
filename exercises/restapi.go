package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/mediocregopher/radix.v2/redis"
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

func requestMux(ctx *fasthttp.RequestCtx) {
	
	path := string(ctx.Path())
	
	switch {
		case strings.Contains(path, "/docs"):
			docsHandler(ctx)
		case strings.Contains(path, "/test"):
			testHandler(ctx)
		default:
			ctx.Error(fmt.Sprintf("HTTP not found: %s", ctx.Path()), fasthttp.StatusNotFound)
	}
}

func docsHandler(ctx *fasthttp.RequestCtx) {
	
	path := string(ctx.Path())
	
	switch ctx.Method() {
		case "GET":
			client, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				ctx.Error(fmt.Sprintf("Erro acesso Redis: %s", err)), fasthttp.StatusInternalServerError)
			} 
			// List all docs if path ends with /docs
			if strings.HasSuffix(path, "/docs") {
				
			}
			
		case "POST":
			pass
		default:
			ctx.Error(fmt.Sprintf("Not implemented: %s", ctx.Method()), fasthttp.StatusNotImplemented)
			
	}

	ctx.SetContentType("text/plain; charset=utf8")	
}

func testHandler(ctx *fasthttp.RequestCtx) {

	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())

	ctx.SetContentType("text/plain; charset=utf8")	
}
