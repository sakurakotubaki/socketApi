# GoでWebsocketを使う

WebSocketについて解説します。

WebSocketとは：
1. **基本的な特徴**
- 双方向通信が可能なプロトコル
- サーバーとクライアント間で常時接続を維持
- HTTPと違い、一度接続すれば継続的に通信可能
- リアルタイムでデータをやり取りできる

2. **使用例**
- チャットアプリケーション
- リアルタイムゲーム
- 株価や為替のリアルタイム更新
- ライブ通知システム
- オンライン協同作業ツール

3. **従来のHTTPとの違い**
```
HTTP:
クライアント → リクエスト → サーバー
クライアント ← レスポンス ← サーバー
（毎回接続と切断を繰り返す）

WebSocket:
クライアント ⟷ 継続的な双方向通信 ⟷ サーバー
（一度接続したら継続的に通信可能）
```

4. **Mac環境でのテストツール**

基本的なツール：
```bash
# WebSocketクライアントツールのインストール
brew install websocat

# または
npm install -g wscat
```

GUI ツール：
1. **Postman**
    - 最新版はWebSocket対応
    - 視覚的に使いやすい
    - リクエスト/レスポンスの履歴が見やすい

2. **Simple WebSocket Client**
    - Chrome拡張機能
    - インストール後すぐに使える
    - シンプルで軽量

テスト用サンプルコード：

```go
package main

import (
    "fmt"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func echo(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()

    log.Println("Client connected")

    for {
        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        
        // エコーバック
        err = c.WriteMessage(mt, message)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
}

func main() {
    http.HandleFunc("/ws", echo)
    fmt.Println("Server is running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

```

5. **テスト方法**

基本的なテスト：
```bash
# サーバー起動
go run main.go

# クライアント1（受信用）
websocat ws://localhost:8080/ws

# クライアント2（送信用）
echo "Hello" | websocat ws://localhost:8080/ws
```

6. **参考になるサイト**

公式ドキュメント：
- [MDN WebSocket API](https://developer.mozilla.org/ja/docs/Web/API/WebSocket)
- [RFC 6455 - WebSocket Protocol](https://tools.ietf.org/html/rfc6455)

チュートリアルと解説：
- [WebSocketの仕組みについて](https://zenn.dev/nameless_sn/articles/websocket_impression)
- [GoのWebSocketサンプル](https://pkg.go.dev/golang.org/x/net/websocket)
- [Flutter WebSocket Guide](https://flutter.dev/docs/cookbook/networking/web-sockets)

7. **使用上の注意点**
- セキュリティ（WSS使用推奨）
- エラーハンドリング
- 接続の再確立処理
- スケーラビリティ考慮
- クライアント数の制限

8. **デバッグツール**：
- Chrome DevTools（Network タブ）
- Safari Web Inspector
- Firefox WebSocket Inspector

これらのツールと情報を使って、WebSocketの動作を理解し、効果的にテストすることができます。実際の開発では、まず小さなテストから始めて、徐々に機能を追加していくことをお勧めします。