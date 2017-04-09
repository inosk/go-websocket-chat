package main

// Hubはチャットに参加しているクライアントとブロードキャストするメッセージを管理する
type Hub struct {
  // チャットに参加しているクライアントのmap
  clients map[*Client]bool

  // ブロードキャスト用のチャンネル
  broadcast chan []byte

  // クライアントの参加受付チャンネル
  register chan *Client

  // クライアントの離脱受付チャンネル
  unregister chan *Client
}

// Hubのconstrucor
func newHub() *Hub {
  return &Hub {
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    clients:    make(map[*Client]bool),
  }
}

// Hubの起動
func (h *Hub) run() {
  // 無限ループ
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}
