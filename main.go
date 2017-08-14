package main

import (
  "flag"
  "log"
  "net/http"
  "html/template"
  "strconv"

  "goji.io"
  "goji.io/pat"
)

var addr = flag.String("addr", ":8080", "http servise address")

func serveRooms(hubs []*Hub, w http.ResponseWriter, r *http.Request) {
  log.Println(r.URL)

  tmpl, parseErr := template.ParseFiles("rooms.html.tmpl")
  if parseErr != nil {
    panic(parseErr)
  }

  links := make([]string, len(hubs))
  for i := range hubs {
    links[i] = "/rooms/" + strconv.Itoa(i + 1)
  }

  execErr := tmpl.Execute(w, struct {
    Rooms []string
  }{
    Rooms: links,
  })

  if execErr != nil {
    panic(execErr)
  }
}

/*
 * htmlはtmplateを使う
 * goのtemplateは正直イマイチな感じなのでなんかいい感じのテンプレートが欲しい
 */
func serveRoom(roomId string, w http.ResponseWriter, r *http.Request) {
  log.Println(r.URL)

  tmpl, parseErr := template.ParseFiles("room.html.tmpl")
  if parseErr != nil {
    panic(parseErr)
  }

  execErr := tmpl.Execute(w, struct {
    RoomId string
    Websocket string
  }{
    RoomId: roomId,
    Websocket: "ws/" + roomId,
  })

  if execErr != nil {
    panic(execErr)
  }
}


/*
 * 新しいroomを作ってroomのURLにリダイレクト
 */
func serveCreateRoom (hubs *[]*Hub, w http.ResponseWriter, r *http.Request) {
  hub := newHub()
  go hub.run()
  *hubs = append(*hubs, hub)
  roomId := len(*hubs)
  path := "/rooms/" + strconv.Itoa(roomId)
  http.Redirect(w, r, path, http.StatusFound)
}

func main() {
  flag.Parse()

  // グループごとのhubを管理する配列
  // 部屋の数だけのサイズ配列ができる
  hubs := []*Hub{}

  // routing
  mux := goji.NewMux()
  mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
    serveRooms(hubs, w, r)
  })
  mux.HandleFunc(pat.Get("/rooms/new"), func(w http.ResponseWriter, r *http.Request) {
    // hubsはサーバ全体で共通管理すべきなので参照渡しする
    serveCreateRoom(&hubs, w, r)
  })
  mux.HandleFunc(pat.Get("/ws/:roomId"), func(w http.ResponseWriter, r *http.Request) {
    roomId, _ := strconv.Atoi(pat.Param(r, "roomId"))
    hubIndex := roomId - 1
    hub := hubs[hubIndex]
    serveWs(hub, w, r)
  })
  mux.HandleFunc(pat.Get("/rooms/:roomId"), func(w http.ResponseWriter, r *http.Request) {
    roomId := pat.Param(r, "roomId")
    serveRoom(roomId, w, r)
  })

  err := http.ListenAndServe(*addr, mux)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
