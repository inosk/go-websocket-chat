package main

import (
  "github.com/gorilla/websocket"
  "log"
  "net/http"
)

// httpのUpgradeヘッダーでwsをつかうようにする
var upgrader = websocket.Upgrader{}

// echoするためのハンドラ
func echo(w http.ResponseWriter, r *http.Request) {
  // upgrader.Upgradeでhttp -> websocketに変換
  // webscocketのconnectionオブジェクトを返す
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println("upgrade:", err)
  }

  // 関数から抜けるときにコネクションを切断する
  defer conn.Close()

  // 無限ループ
  for {
    // バッファーからメッセージを読む
    mt, message, err := conn.ReadMessage()

    // readに失敗した場合は、ループを抜ける
    if err != nil {
      log.Println("read:", err)
      break
    }

    log.Printf("recv: %s", message)

    // バッファーに書き込む
    err = conn.WriteMessage(mt, message)

    // writeに失敗した場合は、ループを抜ける
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
