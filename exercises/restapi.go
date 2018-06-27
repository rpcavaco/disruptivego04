package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"strconv"
	"encoding/json"

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

func testHandler(ctx *fasthttp.RequestCtx) {

	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())

	ctx.SetContentType("text/plain; charset=utf8")	
}

func docsHandler(ctx *fasthttp.RequestCtx) {
	
	path := string(ctx.Path())
	elstr := []string{}

	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		ctx.Error(fmt.Sprintf("Error acessing Redis: %s", err), fasthttp.StatusInternalServerError)
		return
	} 
	
	switch string(ctx.Method()) {

		case "GET":
			
			if strings.Contains(path, "/docs") {
				// List all docs if path ends with /docs
				if strings.HasSuffix(path, "/docs") {
					
					// getting documents list from REDIS as a Range object
					rng := client.Cmd("lrange", "docs", 0, -1)
					if rng.Err != nil {
						// handle error
						ctx.Error(fmt.Sprintf("Error listing from Redis: %s", rng.Err), fasthttp.StatusInternalServerError)
						return
					}
					
					// elstr contains a string array of all documents
					l, _ := rng.List()
					for _, elemStr := range l {
						elstr = append(elstr, elemStr)
					}
					
					// marshaling the list to JSON
					json, err2 := json.Marshal(elstr)
					if err2 != nil {
						ctx.Error(fmt.Sprintf("Error -- JSON Marshal: %s", err2), fasthttp.StatusInternalServerError)
						return
					}
					
					ctx.SetContentType("application/json; charset=utf8")	
					fmt.Fprint(ctx, string(json))
					
				} else {
					
					// extract doc index from URL
					splits := strings.Split(path, "/")
					ord, err3 := strconv.Atoi(splits[len(splits)-1])
					if err3 != nil {
						ctx.Error(fmt.Sprintf("Error in  URL: %s", err3), fasthttp.StatusInternalServerError)
						return						
					}
					
					// get length of docs list
					len, err5 := client.Cmd("llen", "docs").Int()

					if len == 0 {
						ctx.Error(fmt.Sprint("Empty docs list"), fasthttp.StatusInternalServerError)
						return						
					}
					
					if err5 != nil || ord < 0 || ord >= len {
						ctx.Error(fmt.Sprintf("Requested doc index out of range: 0-%d", len-1), fasthttp.StatusInternalServerError)
						return						
					}

					// retrieving requested doc
					rng := client.Cmd("lrange", "docs", ord, ord)
					if rng.Err != nil {
						// handle error
						ctx.Error(fmt.Sprintf("Redis listing error: %s", rng.Err), fasthttp.StatusInternalServerError)
						return
					}
						
					// list of one doc	
					l, _ := rng.List()				
					json, err4 := json.Marshal(l[0])
					if err4 != nil {
						ctx.Error(fmt.Sprintf("Error -- JSON Marshal: %s", err4), fasthttp.StatusInternalServerError)
						return
					}
					
					ctx.SetContentType("application/json; charset=utf8")	
					fmt.Fprint(ctx, string(json))					
				}
			}
			
		case "POST":
		
			// setting / saving a doc
		
			if strings.HasSuffix(path, "/docs") {
				
				bdy := string(ctx.PostBody())
				
				// pushing into REDIS list
				client.Cmd("lpush", "docs", bdy)

				// list len to obtain new doc index
				len, err := client.Cmd("llen", "docs").Int()
				
				if err != nil {
					ctx.Error(fmt.Sprintf("Error acessing Redis: %s", err), fasthttp.StatusInternalServerError)
					return						
				}
				
				if len == 0 {
					ctx.Error(fmt.Sprint("Empty docs list"), fasthttp.StatusInternalServerError)
					return						
				}
				
				// new doc index is returned

				ctx.SetContentType("text/plain; charset=utf8")	
				fmt.Fprint(ctx, "%d", len-1)
			}

		default:
			ctx.Error(fmt.Sprintf("Not implemented: %s", ctx.Method()), fasthttp.StatusNotImplemented)
			
	}
}


