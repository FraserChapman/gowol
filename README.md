# gowol

Go-based Wake-on-LAN (WOL) tool with cross-platform support for Linux, Windows, and macOS

Either used as a command-line tool or library for sending and receiving WOL magic packets.

It can be used to wake devices on your network remotely using their MAC address, and you can also listen for WOL packets for testing purposes.

## Features
- Send WOL Packets: Wake a device on the local network using its MAC address.
- Listen for WOL Packets: Test and capture incoming WOL magic packets.
- Can be used as a module: Integrate it into other Go projects.
- CLI Tool: Use it from the command line for convenience.

## Installation

To install the command-line tool or use the module in your project, run

```bash
go get github.com/FraserChapman/gowol
```

## CLI

```bash
gowol [options]
```

### Options

- `-mac` The MAC address of the target device (required in send mode).
- `-broadcast` The broadcast address to send the WOL packet to (default: `192.168.1.255:9`).
- `-mode` "send" to send WOL packets, "receive" to listen for WOL packets (default: `"send"`).
- `-port` The UDP port to listen on for incoming WOL packets (default: `9`).

### Examples

Send a WOL Packet:

```bash
gowol -mac AA:BB:CC:DD:EE:FF -broadcast 192.168.1.255:9 -mode send
```

Listen for WOL Packets:

```bash
gowol -mode receive -port 9
```

### Display Local MAC and Broadcast Address:

If no arguments are provided, the tool will display the local machine’s MAC address and broadcast address and the help text.

```bash
gowol

--------------------------------------------
Local MAC Address: 00:10:20:30:40:50
Broadcast Address: 123.45.67.89
--------------------------------------------
Usage of ./gowol:
  -broadcast string
        Broadcast address (default "192.168.1.255:9")
  -mac string
        MAC address of the target device
  -mode string
        Mode of operation: "send" or "receive" (default "send")
  -port int
        Port to listen for incoming WOL packets (default 9)
```

## Module

To use the gowol module in your Go project, import it as follows:

```go
import "github.com/FraserChapman/gowol/wol"
```

## Functions

```go
SendWOL(broadcastAddr, mac string) error
```
Sends a WOL magic packet to the given broadcast address with the specified MAC address.

---

```go
ListenForMagicPackets(port int) error
```
Listens for incoming WOL packets on the specified UDP port.

---

---

```go
CreateMagicPacket(mac string) ([]byte, error)
```
Creates a WOL magic packet for a given MAC address.

---

```go
IsMagigPacket(packet []byte) bool
```
Checks if the given packet is a valid WOL magic packet.

---

```go
GetPrimaryMACAndBroadcast() (string, string, error)
```
Returns the MAC address and broadcast address of the current machine’s primary network interface.

---

### Example usage

```go
package main

import (
    "fmt"
    "github.com/FraserChapman/gowol/wol"
)

func main() {
    mac := "AA:BB:CC:DD:EE:FF"
    broadcast := "192.168.1.255:9"

    err := wol.SendWOL(broadcast, mac)
    if err != nil {
        fmt.Printf("Error sending WOL packet: %v\n", err)
    } else {
        fmt.Println("WOL packet sent successfully")
    }
}
```
