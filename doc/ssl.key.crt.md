# 인증서 생성

1. key 생성
2. csr 파일 생성
3. crt 파일 생성
```
# openssl genrsa -out my.key 2048
# openssl req -new -key my.key -out my.csr
# openssl req -new -x509 -days 365 -nodes -keyout my.key -out my.crt
```
### 명령어 
```
# key 파일 생성
# openssl genrsa -out test.key 2048

# csr 파일 생성
# openssl req -new -key test.key -out test.csr

# crt 파일 생성
# openssl x509 -req -days 365 -in test.csr -signkey test.key  -out test.crt

# csr없이 key,crt 파일 한번에 생성
# openssl req -new -x509 -days 365 -nodes -keyout test.key -out test.crt
```

```
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func handleWebSocket(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Read data from the client
		var data string
		err := websocket.Message.Receive(conn, &data)
		if err != nil {
			fmt.Println("Error while reading:", err)
			break
		}

		// Process the received data (you can replace this with your desired logic)
		fmt.Println("Received:", data)

		// Send a response back to the client
		response := "Received: " + data
		err = websocket.Message.Send(conn, response)
		if err != nil {
			fmt.Println("Error while writing:", err)
			break
		}
	}
}

func main() {
	http.Handle("/websocket", websocket.Handler(handleWebSocket))

	// Load the SSL/TLS certificate and key
	certFile := "path_to_your_certificate.crt"
	keyFile := "path_to_your_private_key.key"

	// Start the secure WebSocket server with TLS configuration
	err := http.ListenAndServeTLS(":8080", certFile, keyFile, nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS:", err)
	}
}
```