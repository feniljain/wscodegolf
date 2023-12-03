package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/websocket"
)

func main() {
	// wsEasiest()
	directSocket()
}

func directSocket() {
	// Most probs only Connection: upgrade, Upgrade: websocket and Sec-WebSocket-Version: 13 headers are needed
	httpInitMsg := "GET /echo HTTP/1.1\r\nHost: localhost.com:8080\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\nConnection: keep-alive, Upgrade\r\nSec-Fetch-Mode: websocket\r\n\r\n"

	//       0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-------+-+-------------+-------------------------------+
	// |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
	// |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
	// |N|V|V|V|       |S|             |   (if payload len==126/127)   |
	// | |1|2|3|       |K|             |                               |
	// +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
	// |     Extended payload length continued, if payload len == 127  |
	// + - - - - - - - - - - - - - - - +-------------------------------+
	// |                               |Masking-key, if MASK set to 1  |
	// +-------------------------------+-------------------------------+
	// | Masking-key (continued)       |          Payload Data         |
	// +-------------------------------- - - - - - - - - - - - - - - - +
	// :                     Payload Data continued ...                :
	// + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
	// |                     Payload Data continued ...                |
	// +---------------------------------------------------------------+

	wsMsg := []byte{
		0b10000001, // FIN, RSV1, RSV2, RSV3, OpCode
		0b10000101, // Mask Bit (Compulsary for client to set) + Payload
		// NOTE: We don't need to set extended payload bits if our
		// msg is less than 126 length
		0b00000001,
		0b00000010,
		0b00000011,
		0b00000100, // Mask
		0b01101001,
		0b01100111,
		0b01101111,
		0b01101000,
		0b01101110, // Payload
	}

	// [104, 101, 108, 108, 111]
	// After mask: [105, 103, 111, 104, 110]

	// Listen for incoming connections
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("connection successful!")

	defer conn.Close()

	fmt.Println("will try to send http msg:", httpInitMsg)

	data := []byte(httpInitMsg)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("starting to wait to read data")

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	_, err = conn.Write(wsMsg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	// for _, byt := range buf {
	// 	if byt&1 != 0 {
	// 		fmt.Printf("%08b\n", byt)
	// 	}
	// }

	fmt.Println("Received data:", string(buf[2:n]))
}

func wsEasiest() {
	url := "ws://localhost:8080/echo"
	ws, err := websocket.Dial(url, "", url)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := ws.Write([]byte("hello")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s\n", msg[:n])
}

// To initiate talk with server send an upgrade request like this
/*
```
 GET /chat HTTP/1.1
 Host: example.com:8000
 Upgrade: websocket
 Connection: Upgrade
 Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
 Sec-WebSocket-Version: 13
```
*/
// > https://github.com/tinygo-org/tinygo/releases/download/v0.30.0/tinygo_0.30.0_arm64.deb
// > each line should end with \r\n
// > The MASK bit tells whether the message is encoded. Messages from the client must be masked, so your server must expect this to be 1.
// > The opcode field defines how to interpret the payload data: 0x0 for continuation, 0x1 for text (which is always encoded in UTF-8), 0x2 for binary,
// > The FIN bit tells whether this is the last message in a series. If it's 0, then the server keeps listening for more parts of the message

/*
Firefox Sent headers:
 ```
 GET /echo HTTP/1.1
 Host: localhost:8080
 User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:121.0) Gecko/20100101 Firefox/121.0
 Accept: *\/*
 Accept-Language: en-US,en;q=0.5
 Accept-Encoding: gzip, deflate, br
 Sec-WebSocket-Version: 13
 Origin: https://okanexe.medium.com
 Sec-WebSocket-Extensions: permessage-deflate
 Sec-WebSocket-Key: Fdpt1jH0r+vSDVMRgln4Ww==
 Connection: keep-alive, Upgrade
 Sec-Fetch-Dest: empty
 Sec-Fetch-Mode: websocket
 Sec-Fetch-Site: cross-site
 Pragma: no-cache
 Cache-Control: no-cache
 Upgrade: websocket
// ```
*/

/*
"hello" message with [1, 2, 3, 4, 5] XOR mask applied:
105
103
111
104
106
*/

// this should be the response from server:
// ```
// HTTP/1.1 101 Switching Protocols
// Upgrade: websocket
// Connection: Upgrade
// Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
// ```

// https://stackoverflow.com/questions/69085092/is-it-possible-to-make-a-go-binary-smaller-by-compiling-it-with-tinygo
// https://gophercoding.com/reduce-go-binary-size/
// https://upx.github.io/
// https://stackoverflow.com/questions/27067112/why-are-binaries-built-with-gccgo-smaller-among-other-differences
// https://pkg.go.dev/github.com/nerzal/tinywebsocket
// https://github.com/gobwas/ws
// https://github.com/tinygo-org/awesome-tinygo
// https://okanexe.medium.com/the-complete-guide-to-tcp-ip-connections-in-golang-1216dae27b5a
// https://www.openmymind.net/WebSocket-Framing-Masking-Fragmentation-and-More/
// https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers
// https://tinygo.org/docs/guides/optimizing-binaries/
// https://totallygamerjet.hashnode.dev/the-smallest-go-binary-5kb
// https://github.com/pusher/websockets-from-scratch-tutorial/blob/master/README.md

// Attempt 1:
// go build  -ldflags="-s -w" main.go
// upx ./main
//
