package main

import (
  "bytes"
  "log"
  "net/http"
  "time"

  "github.com/gorilla/websocket"
)

const (
  // メッセージの記入時間？
  writeWait = 10 * time.Second

  // pongの許容時間
  pongWait = 60 * time.Second

  // pingの送信間隔。pongWaitよりは短くすべき
  pingPeriod = (pongWait * 9) / 10

  // 最大メッセージサイズ
  maxMessageSize = 512
)

var (
  newline = []byte{'\n'}
  space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

type Client struct {
  hub *Hub

  // websocket のコネクション
  conn *websocket.Conn

  // 送信メッセージのバッファーチャンネル
  send chan []byte
}

// websocketコネクションからメッセージを読んでHubに書き込み
func (c *Client) readPump() {
  // 終了時にunregisterしてコネクションを閉じる
  defer func() {
    c.hub.unregister <- c
    c.conn.Close()
  }()

  c.conn.SetReadLimit(maxMessageSize)
  c.conn.SetReadDeadline(time.Now().Add(pongWait))
  c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

  // 無限ループ
  for {
    // コネクションからread
    _, message, err := c.conn.ReadMessage()
    if err != nil{
      if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
        log.Printf("error: %v", err)
      }
      break
    }
    log.Println(message)
    // 末尾の改行とスペースを削除
    message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
    // hubへの書き込み
    c.hub.broadcast <- message
  }
}

// hubからメッセージを読んでwebsocket経由で書き込み
func (c *Client) writePump() {
  ticker := time.NewTicker(pingPeriod)

  defer func() {
    ticker.Stop()
    c.conn.Close()
  }()

  for {
    select {
    case message, ok := <-c.send:
      log.Println(message)

      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }

      w, err := c.conn.NextWriter(websocket.TextMessage)
      if err != nil {
        return
      }

      w.Write(message)

      n := len(c.send)
      for i := 0; i < n; i++ {
        w.Write(newline)
        w.Write(<-c.send)
      }

      if err := w.Close(); err != nil {
        return
      }
    case <-ticker.C:
      // 定期的にpingして応答がなければ自信を終了する
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
  }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }
  client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
  client.hub.register <- client
  go client.writePump()
  client.readPump()
}
