package helper

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func convertToIntIP(ip string) (int, error) {
	ips := strings.Split(ip, ".")
	E := errors.New("Not A IP.")
	if len(ips) != 4 {
		return 0, E
	}
	var intIP int
	for k, v := range ips {
		i, err := strconv.Atoi(v)
		if err != nil || i > 255 {
			return 0, E
		}
		intIP = intIP | i<<uint(8*(3-k))
	}
	return intIP, nil
}

func GetLocalIpToInt() (int, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}
	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return convertToIntIP(ipnet.IP.String())
			}
		}
	}
	return 0, errors.New("can not find the client ip address")
}
