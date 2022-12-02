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
	giturl    string
	addr      string
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method
		ua := r.UserAgent()
		log.Printf("%s: %s -- %s", method, uri, ua)
	})
}

func (c *config) goget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")

	pretty := path.Join(c.prettyurl, r.URL.Path)
	content := path.Join(c.giturl, r.URL.Path)
	fmt.Fprintf(
		w, `<head><meta name="go-import" content="%s git https://%s"></head>`,
		pretty, content,
	)
	return
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.prettyurl, "pretty-url", "", "pretty url for your go module")
	flag.StringVar(&cfg.giturl, "git-url", "", "actual git url of your go module")
	flag.StringVar(&cfg.addr, "addr", "0.0.0.0:6868", "listen address")
	flag.Parse()

	if cfg.giturl == "" || cfg.prettyurl == "" {
		fmt.Println("goget: required options --pretty-url and --git-url")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	log.Printf("starting server on %s", cfg.addr)
	mux.Handle("/", logger(http.HandlerFunc(cfg.goget)))
	log.Fatal(http.ListenAndServe(cfg.addr, mux))
}
