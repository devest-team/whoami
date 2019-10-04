package dns

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

// DNS Server Port
const (
	dnsPort = "53"
)

type (
	// Diagnosis 정보
	Info struct {
		Dns          net.IP `json:"clientDns"`
		Ip           string `json:"clientIp"`
		UserAgent    string `json:"userAgent"`
		ResponseTime string `json:"responseTime"`
		ReceiveTime  string `json:"receiveTime"`
	}

	// DNS Server
	Server struct {
		Dns *dns.Server
	}
)

func NewServer() *Server {
	return &Server{
		// DNS 메시지 길이는 512byte 가 넘지 않고 Zone Transfer 요청도 없음으로 udp 만
		Dns:       &dns.Server{Addr: "[::]:" + dnsPort, Net: "udp4", TsigSecret: nil},
	}
}

func (v *Info) String() string {
	bytes, _ := json.Marshal(v)
	str := string(bytes)
	return str
}

func (s *Server) Start(domain string) {
	// DNS Server
	var err error
	go func() {
		s.DnsHandler(domain)
		err = s.Dns.ListenAndServe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to setup the dns server: %s\n", err.Error())
		}
	}()
	fmt.Fprintln(os.Stdout, "server started: ", " dns(", s.Dns.Addr, ")")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	catch := <-sig

	fmt.Printf("Signal (%s) received, stopping\n", catch)
}
