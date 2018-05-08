package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (config *Config) RootHandler(w http.ResponseWriter, r *http.Request) {
	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	// host := r.Header.Get("Host")
	log.Printf("%#v", r.Header)
	host := r.Host
	log.Println(r.Host)
	log.Println(host)
	h := strings.Split(host, ".")
	log.Println(len(h))
	if h[0] == "www" || len(h) == 2 {
		log.Println("landing page")
		xx := template.New("index.html").Delims("<&", "&>")
		xx, err := xx.ParseFiles("static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/html")
		err = xx.ExecuteTemplate(w, "index.html", "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("dashboard")
		tt := template.New("dashboard.html").Delims("<&", "&>")
		tt, err := tt.ParseFiles("static/dashboard.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/html")
		err = tt.ExecuteTemplate(w, "dashboard.html", "")
		if err != nil {
			log.Fatal(err)
		}
	}

}
