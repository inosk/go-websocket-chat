# これはなに

go/websocketでchatの実装

# コードの説明

コレの写経
https://github.com/gorilla/websocket/tree/master/examples/chat

理解のためにコメントを追加

- main.go
  - ただのhttpサーバ
- hub.go
  - clientを管理するための構造体
    - clientの参加や離脱を管理している
  - client間はhubを通じてmessageのやり取りをする
  - hubを経由したbroadcastのみをサポート
- client.go
  - チャット参加者を表現する構造体
  - readPumpはユーザからの入力をwebsocket経由で受け付けてbroadcat channelに投げる
  - writePumpはwebsocketからのsendを受け付けて画面に表示するためのイベント発行
- home.html
  - フロントエンド
  - onmessageで受信したmessageのdivを作ってbodyに挿入
