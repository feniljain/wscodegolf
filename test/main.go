package main

import (
	"syscall"
	"unsafe"
)

// Will only work on linux arm64 systems
// https://chromium.googlesource.com/chromiumos/docs/+/master/constants/syscalls.md#arm64-64_bit
func main() {
	// 0xc6 = 198, socket
	r0, _, _ := syscall.RawSyscall(0xc6, 0x2, 0x1, 0x0)

	// https://github.com/bminor/musl/blob/f314e133929b6379eccc632bef32eaebb66a7335/include/netinet/in.h#L16-L21
	structaddr := [16]byte{
		// 2 << 0 | 0 << 8
		0b00000010,
		0b00000000, // AF_INET
		0b00011111,
		0b10010000, // 8080
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

	wsPacket := []byte{
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

	// 0xcb = 203, connect
	_, _, err := syscall.RawSyscall(0xcb, r0, uintptr(unsafe.Pointer(&structaddr[0])), uintptr(len(structaddr)))
	if err != 0 {
		panic(err)
	}

	httpInitMsg := []byte("GET /echo HTTP/1.1\r\nHost: localhost.com:8080\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\n\r\n")
	// 0xce = 206, sendto
	syscall.RawSyscall6(0xce, r0, uintptr(unsafe.Pointer(&httpInitMsg[0])), uintptr(len(httpInitMsg)), 0, 0, 0)

	// 0xcf = 207, recvfrom
	var response [135]byte
	syscall.RawSyscall6(0xcf, r0, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)), 0, 0, 0)

	syscall.RawSyscall6(0xce, r0, uintptr(unsafe.Pointer(&wsPacket[0])), uintptr(len(wsPacket)), 0, 0, 0)

	syscall.RawSyscall6(0xcf, r0, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)), 0, 0, 0)

	// 0x40 = 64, write
	syscall.RawSyscall(0x40, 1, uintptr(unsafe.Pointer(&response[0])), uintptr(len(response)))
}
