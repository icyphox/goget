package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
)

func goget( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	p := path.Join(r.Host, r.URL.Path)
	fmt.Fprintf(
		w, `<head><meta name="go-import" content="%s git https://%s"></head>`,
		p, p,
	)
	return
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", goget)
	log.Fatal(http.ListenAndServe("0.0.0.0:6868", mux))
}
