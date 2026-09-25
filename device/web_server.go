package device

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

const port = ":8099"

//go:embed static/*.html static/style.css
var content embed.FS

func RunWebServer(ds State, errCh chan error, httpStarted chan bool) {

	webFS, err := fs.Sub(content, "static")
	if err != nil {
		errCh <- fmt.Errorf("Error creating sub filesystem: %v\n", err)
	}

	fs := http.FileServer(http.FS(webFS))
	http.Handle("/", fs)

	http.Handle("/metrics", ds.Metrics.MetricsHandler())

	//  -- Datastar endpoints
	http.Handle("/indextbl", HandleIndexTbl(ds))
	http.Handle("/messages", HandleLogMsg(ds))

	httpStarted <- true

	err = http.ListenAndServe(port, nil)

	if err != nil {
		errCh <- fmt.Errorf("Error starting server: %v\n", err)
		return
	} else {
		fmt.Println("HTTP server started successfully")
	}
}
