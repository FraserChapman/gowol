package wol

import (
	"bytes"
	"fmt"
	"net"
)

// GetPrimaryMACAndBroadcast fetches the MAC address and broadcast address of the current machine's active network interface.
func GetPrimaryMACAndBroadcast() (string, string, error) {
	// Connect to an external IP (Google's DNS) to determine the primary interface
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", "", fmt.Errorf("unable to determine primary network interface: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP

	// Find the interface that matches the local IP address
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", "", fmt.Errorf("unable to retrieve network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.Equal(localAddr) {
					mac := iface.HardwareAddr.String()
					if mac == "" {
						return "", "", fmt.Errorf("no MAC address found for interface %s", iface.Name)
					}

					// Calculate the broadcast address
					mask := ipNet.Mask
					ip := ipNet.IP.To4()
					broadcast := make(net.IP, 4)
					for i := 0; i < 4; i++ {
						broadcast[i] = ip[i] | ^mask[i]
					}
					return mac, broadcast.String(), nil
				}
			}
		} else {
			return "", "", fmt.Errorf("unable to retrieve addresses for interface %s: %w", iface.Name, err)
		}
	}

	return "", "", fmt.Errorf("no suitable network interface found")
}

// CreateMagicPacket creates a magic packet for WOL using the given MAC address.
func CreateMagicPacket(mac string) ([]byte, error) {
	hwAddr, err := net.ParseMAC(mac)
	if err != nil || len(hwAddr) != 6 {
		return nil, fmt.Errorf("invalid MAC address: %w", err)
	}

	// Pre-allocate the full size of the magic packet (6 bytes of 0xFF + 16 repetitions of MAC address)
	packet := make([]byte, 6+(16*len(hwAddr)))

	// First 6 bytes are 0xFF
	copy(packet[:6], bytes.Repeat([]byte{0xFF}, 6))

	// Followed by 16 repetitions of the MAC address
	for i := 0; i < 16; i++ {
		copy(packet[6+(i*6):6+(i+1)*6], hwAddr)
	}

	return packet, nil
}

// SendWOL sends the magic packet to wake up a device on the network using its broadcast address and MAC address.
func SendWOL(broadcastAddr, mac string) error {
	packet, err := CreateMagicPacket(mac)
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return fmt.Errorf("unable to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("unable to dial UDP: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("unable to send magic packet: %w", err)
	}

	return nil
}

// ListenForMagicPackets listens for incoming WOL packets on a specified UDP port.
func ListenForMagicPackets(port int) error {
	addr := fmt.Sprintf(":%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("unable to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("unable to listen on UDP port: %w", err)
	}
	defer conn.Close()

	fmt.Printf("Listening for WOL packets on port %d...\n", port)

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return fmt.Errorf("error reading from UDP: %w", err)
		}

		fmt.Printf("Received packet from %s (%d bytes)\n", remoteAddr.String(), n)

		if IsMagicPacket(buf[:n]) {
			fmt.Println("Magic packet detected!")
		} else {
			fmt.Println("Received packet, but it is not a magic packet")
		}
	}
}

// IsMagicPacket checks if the received packet is a valid magic packet.
func IsMagicPacket(packet []byte) bool {
	if len(packet) < 102 { // Magic packets are exactly 102 bytes
		return false
	}

	// First 6 bytes must be 0xFF
	if !bytes.Equal(packet[:6], bytes.Repeat([]byte{0xFF}, 6)) {
		return false
	}

	// The remaining bytes should consist of the MAC address repeated 16 times
	macAddr := packet[6:12]
	for i := 0; i < 16; i++ {
		if !bytes.Equal(packet[6+(i*6):12+(i*6)], macAddr) {
			return false
		}
	}

	return true
}
