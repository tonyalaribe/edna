package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func (config *Config) RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
	h := strings.Split(r.Host, ".")
	log.Println(h)
	if h[0] == "www" || len(h) > 2 {
		t := template.New("index.html").Delims("<&", "&>")
		t, err := t.ParseFiles("static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		t := template.New("dashboard.html").Delims("<&", "&>")
		t, err := t.ParseFiles("static/dashboard.html")
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, "")
		if err != nil {
			log.Fatal(err)
		}
	}

}
