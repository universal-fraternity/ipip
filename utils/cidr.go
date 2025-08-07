// Package utils provides a unified processing method
package utils

import "net"

// LastIP Calculate the last IP address of a given CIDR
func LastIP(cidr *net.IPNet) net.IP {
	ip := cidr.IP
	mask := cidr.Mask
	last := make(net.IP, len(ip))
	for i := range ip {
		last[i] = ip[i] | ^mask[i]
	}
	return last
}
