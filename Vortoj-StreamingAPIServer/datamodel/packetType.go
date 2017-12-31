package datamodel

const (
	SOCK_ADDRESS = "/tmp/Vortoj-Packet.sock"
)

type Packet struct {
	ID        int16
	DeviceID  string
	SrcMAC    string
	DstMAC    string
	SrcIP     string
	DstIP     string
	SrcPort   string
	DstPort   string
	SYN       bool
	ACK       bool
	Sequence  int64
	Protocol  string
	Length    int64
	DataChank []byte
}
