package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/skip2/go-qrcode"
)

var (
	version = Version()
)

func main() {
	cfg := parseFlags()

	if cfg.doctor {
		runDoctor()
		os.Exit(0)
	}

	if cfg.version {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	if cfg.list {
		if err := checkRootPrivileges(); err != nil {
			log.Fatal(err)
		}

		connections, err := getConnectionFiles()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Previously connected WiFi networks:")
		for _, conn := range connections {
			fmt.Println("📶", strings.TrimSuffix(conn, ".nmconnection"))
		}
		os.Exit(0)
	}

	ssid := cfg.connection

	if cfg.connection == "" {
		id, err := getSSID()
		if err != nil {
			log.Fatal(err)
		}

		ssid = id
	}

	if err := checkRootPrivileges(); err != nil {
		log.Fatal(err)
	}

	connFile, err := findConnectionFile(ssid)
	if err != nil {
		log.Fatal(err)
	}

	password, err := getWifiPassword(connFile)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.qr {
		qrPath := cfg.qrPath

		if cfg.qrPath == "" {
			qrPath = fmt.Sprintf("%s.png", ssid)
		}

		wifiConfig := fmt.Sprintf("WIFI:T:WPA;S:%s;P:%s;;", ssid, password)

		if cfg.saveQr {
			err := qrcode.WriteFile(wifiConfig, qrcode.Medium, 256, qrPath)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("QR code saved as %s\n", qrPath)

			os.Exit(0)
		}

		qrCode, _ := qrcode.New(wifiConfig, qrcode.Medium)
		fmt.Println(qrCode.ToString(true))
	}

	fmt.Printf("📶 WiFi: %s 🔑 Password: %s\n", ssid, password)
}
