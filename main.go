package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

type config struct {
	prettyurl string
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method
		log.Printf("%s: %s", method, uri)
	})
}

func (c *config) goget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")

	pretty := path.Join(c.prettyurl, r.URL.Path)
	content := path.Join(r.Host, r.URL.Path)
	fmt.Fprintf(
		w, `<head><meta name="go-import" content="%s git https://%s"></head>`,
		pretty, content,
	)
	return
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.prettyurl, "pretty-url", "", "pretty url for your go module")
	flag.Parse()

	if cfg.prettyurl == "" {
		fmt.Println("goget: required flag --pretty-url")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	log.Println("starting server on 0.0.0.0:6868")
	mux.Handle("/", logger(http.HandlerFunc(cfg.goget)))
	log.Fatal(http.ListenAndServe("0.0.0.0:6868", mux))
}
