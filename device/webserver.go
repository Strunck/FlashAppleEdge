package device

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

const port = ":8099"

//go:embed static/*.html static/style.css
var content embed.FS

func RunWebServer(ds State, errCh chan error) {
	var err error

	webFS, err := fs.Sub(content, "static")
	if err != nil {
		fmt.Printf("Error creating sub filesystem: %v\n", err)
		os.Exit(1)
	}

	fs := http.FileServer(http.FS(webFS))
	http.Handle("/", fs)

	// http.Handle("/metrics", ds.Metrics.MetricsHandler())

	//  -- Datastar endpoints
	http.Handle("/indextbl", HandleIndexTbl(ds))

	fmt.Println("Http-Server-Port", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	} else {
		ds.PollClients(errCh)
	}
}
