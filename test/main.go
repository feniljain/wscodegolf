package main

import (
	"syscall"
)

func main() {
	sockfd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)

	syscall.Connect(sockfd, &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	})

	syscall.Sendto(sockfd, []byte("GET /echo HTTP/1.1\r\nHost: localhost.com:8080\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\n\r\n"), 0, nil)

	var response [135]byte
	syscall.Recvfrom(sockfd, response[:], 0)

	syscall.Sendto(sockfd, []byte{
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
	}, 0, nil)

	syscall.Recvfrom(sockfd, response[:], 0)

	syscall.Write(1, response[2:6])
}
