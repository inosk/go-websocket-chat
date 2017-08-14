package main

import (
  "flag"
  "log"
  "net/http"

  "goji.io"
  "goji.io/pat"
)

var addr = flag.String("addr", ":8080", "http servise address")

func serveHome(w http.ResponseWriter, r *http.Request) {
  log.Println(r.URL)
  if r.URL.Path != "/" {
    http.Error(w, "Not found", 404)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return
  }

  http.ServeFile(w, r, "home.html")
}

func main() {
  flag.Parse()

  hub := newHub()
  go hub.run()

  mux := goji.NewMux()
  mux.HandleFunc(pat.Get("/"), serveHome)
  mux.HandleFunc(pat.Get("/ws"), func(w http.ResponseWriter, r *http.Request) {
    serveWs(hub, w, r)
  })

  err := http.ListenAndServe(*addr, mux)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
