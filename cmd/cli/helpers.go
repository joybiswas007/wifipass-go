package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"syscall"

	"gopkg.in/ini.v1"
)

const (
	networkConfigDir = "/etc/NetworkManager/system-connections/"
	appName          = "wifipass"
	author           = "Joy Biswas"
	email            = "joybiswas040701@gmail.com"
)

// config holds the command-line options for the WiFi tool
type config struct {
	// - list: Lists previously connected WiFi networks
	list bool

	// - qr: Generates and displays a QR code for the current WiFi password
	qr bool

	// - saveQr: When true, saves the generated QR code to a file instead of displaying it
	saveQr bool

	// - qrPath: Specifies the file path to save the QR code
	qrPath string

	// - connection: Displays the password for a specific WiFi connection
	connection string

	// - doctor: Runs diagnostics to check WiFi settings, configurations, and required packages
	doctor bool

	//version: display version and exit
	version bool
}

// parseFlags parses the command-line flags provided by the user and returns a config struct.
func parseFlags() config {
	var cfg config

	// Override default usage message
	flag.Usage = usage

	flag.BoolVar(&cfg.list, "list", false, "List previously connected WiFi networks (run as sudo)")
	flag.BoolVar(&cfg.qr, "qr", false, "Generate QR code for the current WiFi password and display it")
	flag.BoolVar(&cfg.saveQr, "save-qr", false, "Save the generated QR code to a file instead of displaying it")
	flag.StringVar(&cfg.qrPath, "qr-path", "", "Specify path to save the QR code")
	flag.StringVar(&cfg.connection, "connection", "", "Show password for a specific WiFi connection (run as sudo)")
	flag.BoolVar(&cfg.doctor, "doctor", false, "Run diagnostics to check WiFi settings, configurations, and required packages")
	flag.BoolVar(&cfg.version, "version", false, "Display version and exit")

	flag.Parse()

	return cfg
}

// checkRootPrivileges ensures the user is running the command with superuser privileges
func checkRootPrivileges() error {
	if os.Getgid() != 0 || os.Getuid() != 0 {
		return errors.New("You need to run this command as super user")
	}
	return nil
}

// getSSID retrieves the currently connected WiFi SSID
func getSSID() (string, error) {
	cmd := exec.Command("iwgetid", "-r")

	ssid, err := cmd.CombinedOutput()
	if err != nil {
		// Check if the error is due to "exit status 255"
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Extract the exit status code
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				switch status.ExitStatus() {
				case 255:
					return "", errors.New("You're not connected to any WiFi network.")
				default:
					return "", err
				}
			}
		}
	}

	return strings.TrimSpace(string(ssid)), nil
}

// getConnectionFiles lists all saved WiFi connection files
func getConnectionFiles() ([]string, error) {
	root := os.DirFS(networkConfigDir)
	return fs.Glob(root, "*.nmconnection")
}

// findConnectionFile searches for the configuration file of a given SSID
func findConnectionFile(ssid string) (string, error) {
	connFiles, err := getConnectionFiles()
	if err != nil {
		return "", err
	}

	ssidFile := fmt.Sprintf("%s.nmconnection", ssid)
	for _, conn := range connFiles {
		if conn == ssidFile {
			return conn, nil
		}
	}

	return "", fmt.Errorf("SSID config for %s not found", ssid)
}

// getWifiPassword retrieves the WiFi password from a connection configuration file
func getWifiPassword(configFile string) (string, error) {
	cfg, err := ini.Load(networkConfigDir + configFile)
	if err != nil {
		return "", err
	}
	return cfg.Section("wifi-security").Key("psk").String(), nil
}

// checkDependencies verifies if required commands are installed
func checkDependencies() {
	commands := []string{"NetworkManager", "iwgetid"}
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err != nil {
			fmt.Printf("❌ %s is not installed. Please install it to proceed.\n", cmd)
		} else {
			fmt.Printf("✅ %s is installed.\n", cmd)
		}
	}
}

// runDoctor performs system diagnostics
func runDoctor() {
	fmt.Println("🔍 Running system diagnostics...")
	checkDependencies()
	fmt.Println("✅ Diagnostics completed.")
}

// Version retrieves version information from build metadata
// It extracts the build timestamp, VCS revision, and whether the build was modifie
func Version() string {
	var (
		time     string
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.time":
				time = s.Value
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}

	return fmt.Sprintf("%s-%s", time, revision)
}

// Custom usage function
func usage() {
	fmt.Fprintf(os.Stderr, `%s - A simple CLI tool written in Golang to retrieve WiFi credentials and generate QR codes for easy sharing

Author: %s <%s>

Usage:
`, appName, author, email)
	flag.PrintDefaults()
}
