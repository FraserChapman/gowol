package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/FraserChapman/gowol/wol"
)

func main() {
	// Define the command-line flags
	macAddr := flag.String("mac", "", "MAC address of the target device")
	broadcastIP := flag.String("broadcast", "192.168.1.255:9", "Broadcast address")
	mode := flag.String("mode", "send", "Mode of operation: \"send\" or \"receive\"")
	port := flag.Int("port", 9, "Port to listen for incoming WOL packets")

	// Parse the command-line flags
	flag.Parse()

	// If no arguments are provided, print the local MAC and broadcast address
	if len(os.Args) == 1 {
		mac, broadcast, err := wol.GetPrimaryMACAndBroadcast()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("--------------------------------------------")
		fmt.Printf("Local MAC Address: %s\n", mac)
		fmt.Printf("Broadcast Address: %s\n", broadcast)
		fmt.Println("--------------------------------------------")
		flag.Usage()
		return
	}

	switch strings.ToLower(*mode) {
	case "send":
		if *macAddr == "" {
			fmt.Println("Error: MAC address is required in send mode")
			flag.Usage()
			os.Exit(1)
		}

		// Send the Wake-on-LAN packet
		err := wol.SendWOL(*broadcastIP, *macAddr)
		if err != nil {
			fmt.Println("Error sending WOL packet:", err)
		} else {
			fmt.Println("WOL packet sent successfully")
		}

	case "receive":
		// Listen for incoming WOL packets
		err := wol.ListenForMagicPackets(*port)
		if err != nil {
			fmt.Println("Error listening for WOL packets:", err)
		}

	default:
		fmt.Println("Error: Unknown mode. Use 'send' or 'receive'.")
		flag.Usage()
		os.Exit(1)
	}
}
