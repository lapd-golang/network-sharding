package network

import (
	"net"
	"strconv"
	"strings"
)

//Peer ...
type Peer struct {
	IP             *[4]byte
	ListenPortHost uint32
}

//IsIPValid ...
func (p *Peer) IsIPValid() bool {
	var s []string
	for _, i := range *(p).IP {
		s = append(s, strconv.Itoa(int(i)))
	}
	addr := net.ParseIP(strings.Join(s, "."))

	if addr == nil {
		return false
	}

	return true
}
