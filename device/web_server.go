package device

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
)

const port = ":8099"

//go:embed static/*.html static/style.css
var content embed.FS

func RunWebServer(ds State, errCh chan error, httpStarted chan bool) {

	webFS, err := fs.Sub(content, "static")
	if err != nil {
		errCh <- fmt.Errorf("Error creating sub filesystem: %v\n", err)
		return
	}

	mux := http.NewServeMux()
	staticFS := http.FileServer(http.FS(webFS))
	mux.Handle("/", staticFS)

	mux.Handle("/metrics", ds.Metrics.MetricsHandler())

	//  -- Datastar endpoints
	mux.Handle("/indextbl", HandleIndexTbl(ds))
	mux.Handle("/messages", HandleLogMsg(ds))

	listener, err := net.Listen("tcp", port)
	if err != nil {
		errCh <- fmt.Errorf("Error binding HTTP listener: %v\n", err)
		return
	}

	httpStarted <- true

	fmt.Println("Starting HTTP server on port", port)
	err = http.Serve(listener, mux)

	if err != nil {
		errCh <- fmt.Errorf("Error starting server: %v\n", err)
		return
	} else {
		fmt.Println("HTTP server started successfully")
	}
}
