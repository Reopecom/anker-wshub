package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":9200", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func status(subscription *Subscription) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(subscription.Status())
		w.Header().Set("content-type", "application/json")
		log.Println(data)
		_, _ = w.Write(data)
	}
}

func main() {
	flag.Parse()
	subscription := NewSubscription()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/status", status(subscription))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(subscription, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
