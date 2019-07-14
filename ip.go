// Package net contains helper function for handling
// e.g. ip addresses or domain names
package net

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"strings"
)

// IsIPAddr return true if string ip contains a valid
// representation of an IPv4 or IPv6 address
func IsIPAddr(ip string) bool {
	ipaddr := net.ParseIP(ip)
	if ipaddr != nil {
		if IsIPv4(ipaddr) || IsIPv6(ipaddr) {
			return true
		}
	}

	return false
}

// IsIPv4 return true if string ip contains a valid
// representation of an IPv4 address
func IsIPv4(ip net.IP) bool {
	return strings.Count(ip.String(), ":") < 2
}

// IsIPv6 return true if string ip contains a valid
// representation of an IPv6 address
func IsIPv6(ip net.IP) bool {
	return strings.Count(ip.String(), ":") >= 2
}

// ReverseIPAddr reverses string ip
// (use e.g. for DNS blacklists)
func ReverseIPAddr(ip string) (string, error) {
	result := ""

	if !IsIPAddr(ip) {
		return result, errors.New("invalid IP address")
	}

	ipaddr := net.ParseIP(ip)
	if IsIPv4(ipaddr) {
		ipaddr = ipaddr.To4()
	} else {
		ipaddr = ipaddr.To16()
	}
	for i := 0; i < len(ipaddr); i++ {
		result = fmt.Sprintf("%v.%s", ipaddr[i], result)
	}

	return result, nil
}

// IsNetwork return true if string network contains a valid
// representation of an ip network
func IsNetwork(network string) bool {
	_, ipn, err := net.ParseCIDR(network)
	if err == nil {
		// attn: comparing ipn.IP.String() to the network passed to this function
		// is important to avoid entries like 1.2.3.4/3 being detected as networks
		// while they are in fact URLs!!
		if strings.Split(network, "/")[0] == ipn.IP.String() {
			return true
		}
	}

	return false
}

// IsIPRange return true if string r contains a valid representation
// of an ip network (e.g. 192.168.10.1-192.168.10.199)
func IsIPRange(r string) bool {
	f := strings.Split(r, "-")
	if len(f) == 2 {
		f[0] = strings.TrimSpace(f[0])
		f[1] = strings.TrimSpace(f[1])

		if IsIPAddr(f[0]) && IsIPAddr(f[1]) {
			return true
		}
	}

	return false
}

// IntToIP return net.IP from the integer representation
// of an ip address (use e.g. for IP2Location databases )
func IntToIP(i string) net.IP {
	var ip net.IP

	ni := big.NewInt(0)
	ni.SetString(i, 10)

	b := ni.Bytes()

	if len(b) == 4 {
		// IPv4
		ip = net.IP(b)

	} else {
		// IPv6
		b2 := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		offset := len(b2) - len(b)
		for i := range b {
			b2[i+offset] = b[i]
		}
		ip = net.IP(b2)
	}

	return ip
}
