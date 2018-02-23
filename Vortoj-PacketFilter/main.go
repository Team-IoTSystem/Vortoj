package main

import (
	"fmt"
	"os"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket/layers"
)

const NFQUEUE_NUM_ID = 10
const MAX_PACKET_IN_QUEUE = 300

const EXCLUDE_IN_MACAddr = "00:22:cf:f9:1d:03"
const EXCLUDE_IN_Port = "80"
const EXCLUDE_IN_IP = "192.168.0.2"

func isSelectedExcludeInIP(packet *netfilter.NFPacket, target string) {
	if target == EXCLUDE_IN_IP {
		packet.SetVerdict(netfilter.NF_DROP)
		fmt.Println("Drop is IP")
	}
}

func isSelectedExcludeInPort(packet *netfilter.NFPacket, target string) {
	if target == EXCLUDE_IN_Port {
		packet.SetVerdict(netfilter.NF_DROP)
		fmt.Println("Drop is Port")
	}
}
func isSelecteExcludedInMACAddr(packet *netfilter.NFPacket, target string) {
	if target == EXCLUDE_IN_MACAddr {
		packet.SetVerdict(netfilter.NF_DROP)
		fmt.Println("Drop is MACAddress")
	}
}

//start: sudo iptables -A OUTPUT -j NFQUEUE --queue-num 10
//end  : sudo iptables -D OUTPUT -j NFQUEUE --queue-num 10

func main() {
	var err error

	nfq, err := netfilter.NewNFQueue(NFQUEUE_NUM_ID, MAX_PACKET_IN_QUEUE, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer nfq.Close()
	packets := nfq.GetPackets()

	for true {
		select {
		case packet := <-packets:
			ethernetLayer := packet.Packet.Layer(layers.LayerTypeEthernet)
			ipLayer := packet.Packet.Layer(layers.LayerTypeIPv4)
			tcpLayer := packet.Packet.Layer(layers.LayerTypeTCP)
			if ethernetLayer != nil {
				ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
				isSelecteExcludedInMACAddr(&packet, ethernetPacket.SrcMAC.String())
			} else if ipLayer != nil {
				ip, _ := ipLayer.(*layers.IPv4)
				isSelectedExcludeInIP(&packet, ip.SrcIP.String())
			} else if tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				isSelectedExcludeInPort(&packet, tcp.SrcPort.String())
			} else {
				packet.SetVerdict(netfilter.NF_ACCEPT)
			}
		}
	}
}
