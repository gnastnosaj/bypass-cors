package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

var fasthttpClient = &fasthttp.Client{}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	if method != "OPTIONS" {
		uri := fasthttp.AcquireURI()
		defer fasthttp.ReleaseURI(uri)
		ctx.URI().CopyTo(uri)

		url := string(ctx.RequestURI())[1:]
		ctx.Request.SetRequestURI(url)

		err := fasthttpClient.Do(&ctx.Request, &ctx.Response)
		if err != nil {
			fmt.Fprintf(ctx, err.Error())
			return
		}

		if ctx.Response.Header.StatusCode() == http.StatusFound {
			location := string(ctx.Response.Header.Peek("Location"))
			if strings.HasPrefix(location, "magnet") {
				ctx.Response.Header.Del("Location")
				ctx.Response.SetStatusCode(http.StatusOK)
				ctx.Response.SetBodyString(location)
			} else {
				uri.SetPath(location)
				ctx.Response.Header.Set("Location", uri.String())
			}
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
