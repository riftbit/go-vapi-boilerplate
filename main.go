package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	"github.com/riftbit/go-vapi"
	"github.com/valyala/fasthttp"
)

type AppBase struct {
	ListenTo string
	FHServer *fasthttp.Server
	Vapi     *vapi.VAPI
}

func main() {
	app := new(AppBase)
	app.ListenTo = ":8785"
	app.Vapi = vapi.NewServer()

	err := app.Vapi.RegisterService(new(DemoAPI), "demo")
	if err != nil {
		log.Fatalf("error with vapi register service 'demo': %s", err)
	}

	RunAPI(app)
}

// APIHandler handle api request, process it and return result
func (app *AppBase) APIHandler(ctx *fasthttp.RequestCtx) {
	method := ctx.UserValue("method").(string)
	app.Vapi.CallAPI(ctx, method)
	return
}

// RunAPI starts api http server
func RunAPI(app *AppBase) {

	serviceRouter := router.New()

	hndl := LogAfterMW(LogBeforeMW(ServerInfoMW(app.APIHandler)))

	serviceRouter.POST("/api/:method", hndl)

	app.FHServer = &fasthttp.Server{
		Handler:              serviceRouter.Handler,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		MaxConnsPerIP:        500,
		MaxRequestsPerConn:   500,
		MaxKeepaliveDuration: 30 * time.Second,
	}

	// Error handling
	listenErr := make(chan error, 1)
	var shutdownCh = make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Go-VAPI Server started")
		listenErr <- app.FHServer.ListenAndServe(app.ListenTo)
	}()

	for {
		select {
		// If FHServer.ListenAndServe cannot start due to errors such
		// as "port in use" it will return an error.
		case err := <-listenErr:
			if err != nil {
				log.Fatalf("api server error: %s", err)
			}
			os.Exit(0)
		// handle termination signal
		case sig := <-shutdownCh:
			log.Printf("signal received: %s", sig.String())

			// Servers in the process of shutting down should disable KeepAlives
			app.FHServer.DisableKeepalive = true

			// Attempt the graceful shutdown by closing the listener
			// and completing all inflight requests.
			if err := app.FHServer.Shutdown(); err != nil {
				log.Fatalf("error with graceful close: %s", err)
			}
			log.Println("server gracefully stopped")
		}
	}
}
