package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ktsakalozos-canonical/fancni/internal/cni"
	"github.com/ktsakalozos-canonical/fancni/internal/config"
	"github.com/ktsakalozos-canonical/fancni/internal/fan"
	"github.com/ktsakalozos-canonical/fancni/internal/ipam"
)

func main() {
	// 1. Set up logging to /var/log/fancni.log (append mode).
	logFile, err := os.OpenFile("/var/log/fancni.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Fall back to stderr if we can't open the log file.
		log.SetOutput(os.Stderr)
		log.Printf("warning: could not open log file: %v", err)
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Printf("fancni invoked: CNI_COMMAND=%s CNI_CONTAINERID=%s",
		os.Getenv("CNI_COMMAND"), os.Getenv("CNI_CONTAINERID"))

	if err := run(); err != nil {
		writeCNIError(err)
		os.Exit(1)
	}
}

func run() error {
	command := os.Getenv("CNI_COMMAND")

	// VERSION does not need config or host IP.
	if command == "VERSION" {
		p := cni.NewPlugin(config.NetConfig{}, nil, nil)
		return p.HandleVersion()
	}

	// 2. Read CNI config from stdin.
	cfg, err := config.Parse(os.Stdin)
	if err != nil {
		return fmt.Errorf("parsing CNI config: %w", err)
	}

	// 3. Detect host IP via UDP dial trick (no traffic sent).
	hostIP, err := detectHostIP()
	if err != nil {
		return fmt.Errorf("detecting host IP: %w", err)
	}

	// 4. Compute pod CIDR.
	podCIDR, err := fan.ComputeSubnet(cfg.OverlayNetwork, hostIP)
	if err != nil {
		return fmt.Errorf("computing pod CIDR: %w", err)
	}

	// 5. Create IPAM.
	fileIPAM := ipam.NewFileIPAM("/var/lib/cni/fancni", podCIDR)

	// 6. Create Plugin.
	plugin := cni.NewPlugin(cfg, fileIPAM, hostIP)

	// 7. Dispatch command.
	switch command {
	case "ADD":
		return plugin.HandleAdd()
	case "DEL":
		return plugin.HandleDel()
	case "CHECK":
		return plugin.HandleCheck()
	default:
		return fmt.Errorf("unknown CNI command: %q", command)
	}
}

// detectHostIP returns the IP of the interface that would route to 8.8.8.8.
// No actual traffic is sent.
func detectHostIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, fmt.Errorf("dialing to detect host IP: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// writeCNIError writes a CNI 1.0.0 error JSON to stdout.
func writeCNIError(err error) {
	cniErr := struct {
		CNIVersion string `json:"cniVersion"`
		Code       int    `json:"code"`
		Msg        string `json:"msg"`
	}{
		CNIVersion: "1.0.0",
		Code:       100,
		Msg:        err.Error(),
	}
	data, jsonErr := json.Marshal(cniErr)
	if jsonErr != nil {
		// Last-resort fallback.
		fmt.Fprintf(os.Stdout, `{"cniVersion":"1.0.0","code":100,"msg":%q}`, err.Error())
		return
	}
	os.Stdout.Write(data)
}
