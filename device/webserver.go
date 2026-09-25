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

func RunWebServer(ds State, errCh chan error) (err error) {

	webFS, err := fs.Sub(content, "static")
	if err != nil {
		return fmt.Errorf("Error creating sub filesystem: %v\n", err)
	}

	fs := http.FileServer(http.FS(webFS))
	http.Handle("/", fs)

	http.Handle("/metrics", ds.Metrics.MetricsHandler())

	//  -- Datastar endpoints
	http.Handle("/indextbl", HandleIndexTbl(ds))
	http.Handle("/messages", HandleLogMsg(ds))

	fmt.Println("Http-Server-Port", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return fmt.Errorf("Error starting server: %v\n", err)
	}
	return nil
}
