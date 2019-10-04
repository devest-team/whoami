package dns

import (
	"log"
	"net"
	"strconv"

	"github.com/miekg/dns"
)

// DNS Server Handle Functions
func (s *Server) DnsHandler(domain string) {
	dns.HandleFunc("whoami." + domain, whoami)
}

// whoami 는 client( local cache dns ) ip 를 return 해줌
func whoami(w dns.ResponseWriter, r *dns.Msg) {
	var (
		v4  bool
		rr  dns.RR
		port string
		a   net.IP
	)
	m := new(dns.Msg)
	m.SetReply(r)

	if ip, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		a = ip.IP
		port = strconv.Itoa(ip.Port)
		v4 = a.To4() != nil
	}

	if a == nil {
		return
	}

	// Query 응답
	// Client Cache DNS IP에 대해 TTL 0
	// 단순 A, TXT 레코드 질의만 처리하기 때문에, Secure 또는 기타 옵션에 대한 예외처리 안함
	// IPv4 or IPv6
	if v4 {
		rr = &dns.A{
			Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
			A:   a.To4(),
		}
	} else {
		rr = &dns.AAAA{
			Hdr:  dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 0},
			AAAA: a,
		}
	}
	t := &dns.TXT{
		Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
		Txt: []string{port},
	}

	// TXT or A (AAAA)
	// A (AAAA) 레코드 질의시 Answer Section 에 Client DNS server IP를 응답.
	// Port 정보를 Additional Section 에 담아 응답. (큰 의미는 없음)
	switch r.Question[0].Qtype {
	case dns.TypeTXT:
		m.Answer = append(m.Answer, t)
		m.Extra = append(m.Extra, rr)
	default:
		fallthrough
	case dns.TypeAAAA, dns.TypeA:
		m.Answer = append(m.Answer, rr)
		m.Extra = append(m.Extra, t)
	}

	log.Println("[whoami]:", m.Question[0].Name, "-> ", m.Answer)
	err := w.WriteMsg(m)
	if err != nil {
		panic(err)
	}
}