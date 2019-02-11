package main

import (
	"github.com/riftbit/go-vapi"
	"github.com/valyala/fasthttp"
)

// DemoAPI area
type DemoAPI struct{}

// Test Method to test
func (h *DemoAPI) Test(ctx *fasthttp.RequestCtx, args *struct{ ID string }, reply *struct{ LogID string }) error {
	reply.LogID = args.ID
	return nil
}

// ErrorTest Method to test
func (h *DemoAPI) ErrorTest(ctx *fasthttp.RequestCtx, args *struct{}, reply *struct{}) error {

	errs := &vapi.Error{
		ErrorHTTPCode: 333,
		ErrorCode:     606,
		ErrorMessage:  "Test Wrong answer",
		Data: struct {
			LOL string
		}{
			LOL: "OLOLO",
		},
	}

	return errs
}
