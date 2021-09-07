package main

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"net/http"
)

func apiNotFound(ctx *fasthttp.RequestCtx) {
	b, err := json.Marshal(map[string]interface{}{
		"error": http.StatusText(fasthttp.StatusNotFound),
	})

	if err != nil {
		panic(err)
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.SetBodyString(string(b))
}

func apiMethodNotAllowed(ctx *fasthttp.RequestCtx) {
	b, err := json.Marshal(map[string]interface{}{
		"error": http.StatusText(fasthttp.StatusMethodNotAllowed),
	})

	if err != nil {
		panic(err)
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	ctx.SetBodyString(string(b))
}

func apiPanicHandler(ctx *fasthttp.RequestCtx, e interface{}) {
	b, err := json.Marshal(map[string]interface{}{
		"error":   http.StatusText(fasthttp.StatusInternalServerError),
		"message": e,
	})

	if err != nil {
		panic(err)
	}

	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetContentType("application/json")
	ctx.Response.Header.SetServer("webshot/v1")
	ctx.SetBodyString(string(b))
}

func apiRoot(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/v1", fasthttp.StatusMovedPermanently)
}

func apiHelp(ctx *fasthttp.RequestCtx) {
	methods := map[string]interface{}{
		"/v1":       "Help",
		"/v1/add":   "Add a new URL to the queue to process the screenshot",
		"/v1/check": "Check if the URL being processed is complete",
		"/v1/info":  "Get screenshot URL to view / download image and other useful info",
	}

	b, err := json.MarshalIndent(methods, "", "    ")
	if err != nil {
		panic(err)
	}

	ctx.SetStatusCode(200)
	ctx.Response.Header.SetContentType("text/plain")
	ctx.SetBodyString(string(b))
}

func apiAdd(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotImplemented)
	ctx.SetBodyString(http.StatusText(fasthttp.StatusNotImplemented))
}

func apiCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotImplemented)
	ctx.SetBodyString(http.StatusText(fasthttp.StatusNotImplemented))
}

func apiInfo(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotImplemented)
	ctx.SetBodyString(http.StatusText(fasthttp.StatusNotImplemented))
}

func middleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/json")
		ctx.Response.Header.SetServer("webshot/v1")
		h(ctx)
	}
}

func main() {
	app := router.New()

	app.NotFound = middleware(apiNotFound)
	app.MethodNotAllowed = middleware(apiMethodNotAllowed)
	app.PanicHandler = apiPanicHandler

	app.Handle("GET", "/", middleware(apiRoot))
	app.Handle("GET", "/v1", middleware(apiHelp))
	app.Handle("GET", "/v1/add", middleware(apiAdd))
	app.Handle("GET", "/v1/check", middleware(apiCheck))
	app.Handle("GET", "/v1/info", middleware(apiInfo))

	if err := fasthttp.ListenAndServe(":8080", app.Handler); err != nil {
		panic(err)
	}
}
