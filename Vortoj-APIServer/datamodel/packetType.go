package datamodel

const (
	PATH               string = "root:root@tcp(127.0.0.1:3306)/"
	DATABASE_NAME      string = "vortojpacket"
	PACKET_TABLENAME   string = "packet"
	DISTANCE_TABLENAME string = "distance"
	DBTYPE             string = "mysql"
)

//DBPacket packetキャプチャしたデータのDBモデル
type DBPacket struct {
	ID        int64  `db:"id"`
	DeviceID  string `db:"deviceid"`
	SrcMAC    string `db:"src_mac"`
	DstMAC    string `db:"dst_mac"`
	SrcIP     string `db:"src_ip"`
	DstIP     string `db:"dst_ip"`
	SrcPort   string `db:"src_port"`
	DstPort   string `db:"dst_port"`
	SYN       bool   `db:"syn"`
	ACK       bool   `db:"ack"`
	Sequence  int64  `db:"sequence"`
	Protocol  string `db:"protocol"`
	Length    int64  `db:"length"`
	DataChank []byte `db:"datachank"`
}

//DBDistance wifi電波強度による距離データのDBモデル
type DBDistance struct {
	ID     int64  `db:"id"`
	Mac    string `db:"macaddr"`
	RpiMac string `db:"rpimac"`
	Dist   int64  `db:"distance"`
	Pwr    int64  `db:"pwr"`
}
