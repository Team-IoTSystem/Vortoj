package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
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

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
const SOCK_ADDRESS = "/tmp/Vortoj-Packet.sock"

var datachannel = make(chan []byte)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	conn, err := net.Dial("unix", SOCK_ADDRESS)
	if err != nil {
		panic(err)
	}

	go func(conn net.Conn) {
		for {
			buf := make([]byte, 500)
			nr, err := conn.Read(buf)
			if err != nil {
				log.Printf("error: %v\n", err)
				break
			}
			buf = buf[:nr]
			datachannel <- buf
			log.Printf("receive: %s\n", buf)
		}
	}(conn)

	http.HandleFunc("/sample", ServerDataAll)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func ServerDataAll(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// mt, message, err := c.ReadMessage()
	// if err != nil {
	// 	log.Println("read:", err)
	// }
	for {
		message := <-datachannel
		err = c.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
