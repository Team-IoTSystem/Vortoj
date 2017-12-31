package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"./lifecycle"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device      = flag.String("i", "en0", "using device to network interface name") //ifconfigやip aとかで調べて使ってください
	snapshotLen = flag.Int64("len", 1024, "date snapshot size")
	promiscuous = flag.Bool("promise", false, "is promiscuous mode")
	err         error
	timeout     = time.Duration(*(flag.Int("t", 30, "promisc"))) * time.Second
	handle      *pcap.Handle
	// Will reuse these for each packet
	ethLayer       layers.Ethernet
	ipLayer        layers.IPv4
	tcpLayer       layers.TCP
	devicefindFlag = flag.Bool("find", false, "device network interface to find")
)

func networkInterfeceFind() {
	// Find all devices
	devices, err := pcap.FindAllDevs()

	if err != nil {
		log.Fatal(err)
	}

	// Print device information
	log.Println("Devices found:")
	for _, device := range devices {
		fmt.Println("\nName: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	if *devicefindFlag { //デバイスを探すだけなのでこれをしたら落ちる仕様
		networkInterfeceFind()
		os.Exit(1)
	}

	// Open device pcap
	handle, err = pcap.OpenLive(*device, int32(*snapshotLen), *promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	//packet analayzer
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		parser := gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet,
			&ethLayer,
			&ipLayer,
			&tcpLayer,
		)
		foundLayerTypes := []gopacket.LayerType{}
		err := parser.DecodeLayers(packet.Data(), &foundLayerTypes)
		if err != nil {
			fmt.Println("Trouble decoding layers: ", err)
		}

		//パケットの解析をしていく
		for _, layerType := range foundLayerTypes {
			if layerType == layers.LayerTypeTCP {

				applicationLayer := packet.ApplicationLayer()
				if applicationLayer != nil {
					lifecycle.InsertPacketData(&lifecycle.TPacket{
						DeviceID:  "NOT_DEVICEID",
						SrcMAC:    ethLayer.SrcMAC.String(),
						DstMAC:    ethLayer.DstMAC.String(),
						SrcIP:     ipLayer.SrcIP.String(),
						DstIP:     ipLayer.DstIP.String(),
						SrcPort:   tcpLayer.SrcPort.String(),
						DstPort:   tcpLayer.DstPort.String(),
						SYN:       tcpLayer.SYN,
						ACK:       tcpLayer.ACK,
						Sequence:  int64(tcpLayer.Seq),
						Protocol:  string(ipLayer.Protocol.String()),
						Length:    int64(ipLayer.Length),
						DataChank: applicationLayer.Payload(),
					})
				}
			}

		}

	}
}
