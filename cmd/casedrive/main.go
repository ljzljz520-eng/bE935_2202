package main

import (
	"flag"
	"lawdrive/internal/httpapi"
	"lawdrive/internal/query"
	"lawdrive/internal/store"
	"lawdrive/internal/workflow"
	"log"
	"net/http"
	"os"
)

func main() {
	path := flag.String("db", "casedrive.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	s, e := store.Open(*path)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	r := store.NewRepository(s)
	w := workflow.New(r)
	q := query.New(r)
	log.Printf("case drive listening on %s", *addr)
	if e = http.ListenAndServe(*addr, httpapi.New(w, q).Handler()); e != nil && !os.IsTimeout(e) {
		log.Fatal(e)
	}
}
