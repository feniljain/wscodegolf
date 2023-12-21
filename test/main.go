package main

import (
	"syscall"
	"unsafe"
)

var (
	httpInitMsg = []byte("GET /echo HTTP/1.1\r\nHost: localhost.com:8080\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\nConnection: keep-alive, Upgrade\r\nSec-Fetch-Mode: websocket\r\n\r\n")

	// httpInitMsg = [...]byte{0b1000111, 0b1000101, 0b1010100, 0b100000, 0b101111, 0b100000, 0b1001000, 0b1010100, 0b1010100, 0b1010000, 0b101111, 0b110001, 0b101110, 0b110001, 0b1101, 0b1010, 0b1001000, 0b1101111, 0b1110011, 0b1110100, 0b111010, 0b1101, 0b1010, 0b1010101, 0b1110000, 0b1100111, 0b1110010, 0b1100001, 0b1100100, 0b1100101, 0b111010, 0b1110111, 0b1100101, 0b1100010, 0b1110011, 0b1101111, 0b1100011, 0b1101011, 0b1100101, 0b1110100, 0b1101, 0b1010, 0b1000011, 0b1101111, 0b1101110, 0b1101110, 0b1100101, 0b1100011, 0b1110100, 0b1101001, 0b1101111, 0b1101110, 0b111010, 0b1010101, 0b1110000, 0b1100111, 0b1110010, 0b1100001, 0b1100100, 0b1100101, 0b1101, 0b1010, 0b1010011, 0b1100101, 0b1100011, 0b101101, 0b1010111, 0b1100101, 0b1100010, 0b1010011, 0b1101111, 0b1100011, 0b1101011, 0b1100101, 0b1110100, 0b101101, 0b1001011, 0b1100101, 0b1111001, 0b111010, 0b1100100, 0b1000111, 0b1101000, 0b1101100, 0b1001001, 0b1001000, 0b1001110, 0b1101000, 0b1100010, 0b1011000, 0b1000010, 0b1110011, 0b1011010, 0b1010011, 0b1000010, 0b1110101, 0b1100010, 0b110010, 0b110101, 0b1101010, 0b1011010, 0b1010001, 0b111101, 0b111101, 0b1101, 0b1010, 0b1010011, 0b1100101, 0b1100011, 0b101101, 0b1010111, 0b1100101, 0b1100010, 0b1010011, 0b1101111, 0b1100011, 0b1101011, 0b1100101, 0b1110100, 0b101101, 0b1010110, 0b1100101, 0b1110010, 0b1110011, 0b1101001, 0b1101111, 0b1101110, 0b111010, 0b110001, 0b110011, 0b1101, 0b1010, 0b1000011, 0b1101111, 0b1101110, 0b1101110, 0b1100101, 0b1100011, 0b1110100, 0b1101001, 0b1101111, 0b1101110, 0b111010, 0b1010101, 0b1110000, 0b1100111, 0b1110010, 0b1100001, 0b1100100, 0b1100101, 0b1101, 0b1010, 0b1101, 0b1010}

	packet = []byte{
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

	sockaddr = [16]byte{
		// family - AF_INET (0x2), padded to 16 bits
		0b00000010,
		0b00000000,
		// port - 8080, padded to 16 bits
		0b00011111,
		0b10010000,
		// addr - 127.0.0.1, 32 bits
		// 127 << 0 | 0 << 8 | 0 << 16 | 1 << 24
		0b01111111,
		0b00000000,
		0b00000000,
		0b00000001,
		// 64 bits of padding
		0b00000000, 0b00000000, 0b00000000, 0b00000000,
		0b00000000, 0b00000000, 0b00000000, 0b00000000,
	}

	response [1024]byte
)

//go:linkname alloc runtime.alloc
func alloc(size uintptr, layout unsafe.Pointer) unsafe.Pointer {
	return nil
}

//go:section .text_noalign
//export actual_main
func main() {
	// __NR_socket, AF_INET, SOCK_STREAM
	var sock, _, _ = syscall.Syscall(359, 0x2, 0x1, 0)

	// __NR_connect, fd, sockaddr_in, len(sockaddr_in)
	syscall.Syscall6(362, sock, uintptr(unsafe.Pointer(&sockaddr[0])), uintptr(len(sockaddr)), 0, 0, 0)

	// __NR_sendto, fd, buf, len(buf), flags, addr, addr_len
	syscall.Syscall6(369, sock, uintptr(unsafe.Pointer(&httpInitMsg[0])), uintptr(len(httpInitMsg)), 0, 0, 0)

	// __NR_recvfrom, fd, buf, len(buf), flags, addr, addr_len
	var n, _, _ = syscall.Syscall6(371, sock, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)), 0, 0, 0)

	// __NR_sendto
	syscall.Syscall6(369, sock, uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)), 0, 0, 0)

	// __NR_recvfrom
	syscall.Syscall6(371, sock, uintptr(unsafe.Pointer(&response[n])), uintptr(len(response))-n, 0, 0, 0)

	// __NR_close
	syscall.Syscall(6, sock, 0, 0)

	// __NR_write, STDOUT_FILENO
	syscall.Syscall(4, 1, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)))

	// __NR_exit
	syscall.Syscall(1, 1, 0, 0)
}
