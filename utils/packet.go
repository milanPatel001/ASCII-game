package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

const (
	AUTH               = 0x01
	MOVE               = 0x02
	REDIRECT_TO_SERVER = 0x03
	BROADCAST          = 0x04
)

type Packet struct {
	SrcIP       [4]byte
	Timestamp   uint32
	MessageType byte
	Payload     []byte
}

func NewPacket(ip net.IP, messageType byte, payload []byte) *Packet {
	timestamp := uint32(time.Now().Unix())
	return &Packet{SrcIP: ConvIpv4ToBytes(ip), Timestamp: timestamp, MessageType: messageType, Payload: payload}
}

func (p *Packet) Serialize() ([]byte, error) {
	var buffer bytes.Buffer

	// Write IPv4 (4 bytes)
	if err := binary.Write(&buffer, binary.BigEndian, p.SrcIP); err != nil {
		return nil, err
	}

	// Write timestamp (8 bytes)
	if err := binary.Write(&buffer, binary.BigEndian, p.Timestamp); err != nil {
		return nil, err
	}

	// Write message type (2 bytes)
	if err := binary.Write(&buffer, binary.BigEndian, p.MessageType); err != nil {
		return nil, err
	}

	n, err := buffer.Write(p.Payload)

	if err != nil {
		return nil, err
	}

	log.Printf("Payload size: %v", n)

	return buffer.Bytes(), nil
}

func Deserialize(packet []byte) (Packet, error) {

	if len(packet) < 9 || len(packet) > 1048576 {
		log.Fatal("packetData is too small or too large")
	}

	// Create a Packet instance
	var pkt Packet

	// Read SrcIP (first 4 bytes)
	copy(pkt.SrcIP[:], packet[0:4])

	// Read Timestamp (next 4 bytes)
	pkt.Timestamp = binary.BigEndian.Uint32(packet[4:8])

	// Read MessageType (next 1 byte)
	pkt.MessageType = packet[8]

	// Read Payload (remaining bytes)
	pkt.Payload = packet[9:]

	return pkt, nil
}

func (self *Packet) IsEqual(other *Packet) bool {
	return self.MessageType == other.MessageType && self.SrcIP == other.SrcIP && self.Timestamp == other.Timestamp && len(self.Payload) == len(other.Payload)
}
