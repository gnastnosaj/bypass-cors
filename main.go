package main

import (
	"flag"
	"fmt"

	"github.com/valyala/fasthttp"
)

var hc = &fasthttp.Client{}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	if method != "OPTIONS" {
		url := string(ctx.RequestURI())[1:]
		ctx.Request.SetRequestURI(url)

		err := hc.Do(&ctx.Request, &ctx.Response)
		if err != nil {
			fmt.Fprintf(ctx, err.Error())
			return
		}
	}

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
}

func main() {
	var port = flag.String("p", "8080", "port")

	flag.Parse()

	fasthttp.ListenAndServe(fmt.Sprintf(":%s", *port), fastHTTPHandler)
}
