package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":9200", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(
		`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>WsPubSubStatus</title>
</head>
<body>
<table>
    <thead>
    <th>Topic</th>
    <th>Subscribers</th>
    </thead>
    <tbody id="subs">
    </tbody>
</table>
<script>
    async function getStatus() {
        const res = await fetch("/status");
        const data = await res.json();
        return data ? Object.entries(data).map(([k, v]) => ({topic: k, subscribers: v})) : [];
    }

    setInterval(() => {
        getStatus().then(x => {
            document.getElementById("subs").innerHTML = ""
            x.map((y) => {
                const row = document.createElement("tr");
                const topic = document.createElement("td");
                topic.innerText = y.topic;
                const count = document.createElement("td");
                count.innerText = y.subscribers;
                row.append(topic, count);
                document.getElementById("subs").appendChild(row);
            })
        }).catch(x => console.log(x))
    }, 1000)
</script>

</body>
</html>`,
	))
}

func status(subscription *Subscription) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(subscription.Status())
		w.Header().Set("content-type", "application/json")
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
	log.Printf("starting server on port %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
