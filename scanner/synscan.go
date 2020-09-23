package scanner

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
	"sync"
	"time"
)

func localIPPort(dstIp net.IP) (net.IP, int) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstIp.String()+":12345")
	if err != nil {
		log.Fatal(err)
	}
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpAddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpAddr.IP, udpAddr.Port
		}
	}
	log.Fatal("could not get local ip: " + err.Error())
	return nil, -1
}

func SynScan(hostname string, port int, results *[]ScanResult, wg *sync.WaitGroup, srcIp net.IP, sPort int) {
	defer wg.Done()

	dstAddrs, err := net.LookupIP(hostname)
	if err != nil {
		log.Fatal(err)
	}

	// parse the destination host and port from the command line os.Args
	dstIp := dstAddrs[0].To4()
	dstPort := layers.TCPPort(port)

	//srcIp, sPort := localIPPort(dstIp)
	srcPort := layers.TCPPort(sPort)

	// Our IP header... not used, but necessary for TCP checksumming.
	ip := &layers.IPv4{
		SrcIP:    srcIp,
		DstIP:    dstIp,
		Protocol: layers.IPProtocolTCP,
	}
	// Our TCP header
	tcp := &layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		Seq:     1105024978,
		SYN:     true,
		Window:  14600,
	}
	tcp.SetNetworkLayerForChecksum(ip)

	// Serialize.  Note:  we only serialize the TCP layer, because the
	// socket we get with net.ListenPacket wraps our data in IPv4 packets
	// already.  We do still need the IP layer to compute checksums
	// correctly, though.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	log.Println("writing request")
	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIp}); err != nil {
		log.Fatal(err)
	}

	// Set deadline so we don't wait forever.
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, 4096)
		log.Println("reading from conn")
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			log.Println("error reading packet: ", err)
			return
		} else if addr.String() == dstIp.String() {
			// Decode a packet
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcPort {
					if tcp.SYN && tcp.ACK {
						log.Printf("Port %d is OPEN\n", dstPort)
					} else {
						log.Printf("Port %d is CLOSED\n", dstPort)
					}
					return
				}
			}
		} else {
			log.Printf("Got packet not matching addr")
		}
	}
}
