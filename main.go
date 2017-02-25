package main

import (
  "github.com/gorilla/websocket"
  "log"
  "net/http"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println("upgrade:", err)
  }

  defer conn.Close()

  for {
    mt, message, err := conn.ReadMessage()
    if err != nil {
      log.Println("read:", err)
      break
    }
    log.Printf("recv: %s", message)
    err = conn.WriteMessage(mt, message)
    if err != nil {
      log.Println("write:", err)
      break
    }
  }
}

func main() {
  http.HandleFunc("/echo", echo)
  http.Handle("/", http.FileServer(http.Dir("./")))
  log.Fatal(http.ListenAndServe(":9999", nil))
}
