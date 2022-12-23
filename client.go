package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Config struct {
	MaxConnectionsPerIP int       `json:"max_connections_per_ip"`
	BanDuration         time.Time `json:"ban_duration"`
}

var (
	// config is the DDoS protection configuration.
	config Config
	// bannedIPs is a map of banned IP addresses and the time when the ban will expire.
	bannedIPs map[string]time.Time
)

func main() {
	// Load the DDoS protection configuration from a JSON file.
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		fmt.Println(err)
		return
	}

	bannedIPs = make(map[string]time.Time)

	// Listen for incoming TCP connections.
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleTCPConnection(conn)
		}
	}()

	// Listen for incoming UDP connections.
	addr, err := net.ResolveUDPAddr("udp", ":80")
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			handleUDPConnection(conn)
		}
	}()

	// Run indefinitely.
	select {}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	if !checkIP(ip) {
		return
	}
	// At this point, the connection is allowed. You can now process the request as needed.
}

func handleUDPConnection(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buf)
	if err
