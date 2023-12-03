package main

import (
	"net"
	"os"
)

func main() {
	httpInitMsg := "GET /echo HTTP/1.1\r\nHost: localhost.com:8080\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\nConnection: keep-alive, Upgrade\r\nSec-Fetch-Mode: websocket\r\n\r\n"

	conn, _ := net.Dial("tcp", "localhost:8080")

	data := []byte(httpInitMsg)
	conn.Write(data)

	buf := make([]byte, 1024)
	conn.Read(buf)

	data = []byte{
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
	conn.Write(data)

	n, _ := conn.Read(buf)

	os.Stdout.Write(buf[2:n])
}
