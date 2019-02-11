package main

import (
	"github.com/valyala/fasthttp"
)

func ServerInfoMW(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Delegate request to the given handle
		h(ctx)
		ctx.Response.Header.Set("Server", "VAPI TEST")
		ctx.Response.Header.Set("X-Powered-By", "Riftbit ErgoZ")
		return
	})
}

func LogBeforeMW(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		//log.Println("Log before:", string(ctx.Request.Body()))
		h(ctx)
		return
	})
}

func LogAfterMW(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		//log.Println("Log after:", string(ctx.Response.Body()))
		return
	})
}
