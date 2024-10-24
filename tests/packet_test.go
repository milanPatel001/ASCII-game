package tests

import (
	"ascii/utils"
	"fmt"
	"net"
	"testing"
)

type TestStruct struct {
	A int
	B string
	C []int
}

func TestPacketSerDer(t *testing.T) {

	ip := net.IP.To4(net.ParseIP("192.168.1.1"))

	messageType := utils.AUTH
	payload := TestStruct{A: 200, B: "ok crodie", C: []int{2, 3, 4, 5}}

	payloadBytes, err := utils.ConvComplexPayloadToBytes(payload)

	if err != nil {
		t.Errorf("Error converting payload to bytes array: %v", err)
	}

	p := utils.NewPacket(ip, byte(messageType), payloadBytes)

	out, err := p.Serialize()

	fmt.Printf("Packet size total: %v\n", len(out))

	if err != nil {
		t.Errorf("Error in serializing: %v", err)
	}

	p2, err := utils.Deserialize(out)

	if p2.MessageType != p.MessageType || p2.SrcIP != p.SrcIP || p2.Timestamp != p.Timestamp {
		t.Errorf("Error in headers: %v", err)
	}

	if len(payloadBytes) != len(p2.Payload) {
		t.Errorf("Error parsing the payload")
	}

	//fmt.Println(p2.Payload)

}
