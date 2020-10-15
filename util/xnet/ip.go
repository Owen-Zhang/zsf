package xnet

import "net"

//GetLocalMainIP 获取机器ip地址
func GetLocalMainIP() (string, int, error) {
	// UDP Connect, no handshake
	conn, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), localAddr.Port, nil
}
