package wol

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

// Test for IsMagicPacket
func TestIsMagicPacket(t *testing.T) {
	mac := "01:23:45:67:89:ab"
	packet, _ := CreateMagicPacket(mac)

	t.Run("Valid Magic Packet", func(t *testing.T) {
		if !IsMagicPacket(packet) {
			t.Error("Valid magic packet was not detected")
		}
	})

	t.Run("Invalid Magic Packet - First byte changed", func(t *testing.T) {
		packet[0] = 0x00
		if IsMagicPacket(packet) {
			t.Error("Invalid magic packet (first byte changed) was incorrectly detected as valid")
		}
	})

	t.Run("Invalid Magic Packet - Truncated", func(t *testing.T) {
		truncatedPacket := packet[:50]
		if IsMagicPacket(truncatedPacket) {
			t.Error("Truncated packet was incorrectly detected as valid")
		}
	})

	t.Run("Invalid Magic Packet - Incorrect repeated MAC", func(t *testing.T) {
		incorrectMACPacket := append(packet[:6], bytes.Repeat([]byte{0x00}, 96)...)
		if IsMagicPacket(incorrectMACPacket) {
			t.Error("Invalid magic packet (incorrect repeated MAC) was incorrectly detected as valid")
		}
	})
}

// Test for SendWOL
func TestSendWOL(t *testing.T) {
	mac := "01:23:45:67:89:ab"
	broadcast := "192.168.1.255:9"

	err := SendWOL(broadcast, mac)
	if err != nil {
		t.Fatalf("Failed to send WOL packet: %v", err)
	}
}

// TestListenForMagicPackets checks if a magic packet is correctly received
func TestListenForMagicPackets(t *testing.T) {
	// Channels for signaling
	done := make(chan int, 1)
	doneTesting := make(chan bool, 1)

	// Start listener in a goroutine
	go func() {
		laddr := net.UDPAddr{Port: 0} // Let the OS choose the port
		conn, err := net.ListenUDP("udp", &laddr)
		if err != nil {
			t.Fatalf("Failed to listen on UDP port: %v", err)
		}
		defer conn.Close()

		done <- conn.LocalAddr().(*net.UDPAddr).Port

		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				t.Errorf("Error reading from UDP: %v", err)
				return
			}

			if IsMagicPacket(buf[:n]) {
				doneTesting <- true
				return
			}
		}
	}()

	// Wait for the listener to be ready
	var listenerPort int
	select {
	case listenerPort = <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Listener did not start in time")
	}

	time.Sleep(500 * time.Millisecond)

	// Create and send a magic packet
	mac := "01:23:45:67:89:ab"
	packet, err := CreateMagicPacket(mac)
	if err != nil {
		t.Fatalf("Failed to create magic packet: %v", err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", listenerPort)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		t.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		t.Fatalf("Failed to send magic packet: %v", err)
	}

	// Set a timeout to ensure the test doesn't run indefinitely
	select {
	case <-doneTesting:
	case <-time.After(5 * time.Second):
		t.Fatal("TestListenForMagicPackets timed out")
	}
}

// Test for GetPrimaryMACAndBroadcast
func TestGetPrimaryMACAndBroadcast(t *testing.T) {
	mac, broadcast, err := GetPrimaryMACAndBroadcast()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mac == "" {
		t.Error("Expected non-empty MAC address")
	}

	if broadcast == "" {
		t.Error("Expected non-empty broadcast address")
	}
}
