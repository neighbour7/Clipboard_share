package main

import (
	"encoding/binary"
	"errors"
	"net"
)

const headerLen = 8

func validateIp(ip string) (string, error) {
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("ip error")
}

func validatePort(port int) (int, error) {
	if port < 0 || port > 65535 {
		return -1, errors.New("")
	}
	return port, nil
}

func Int64ToBytes(num int64) []byte {
	byteArray := make([]byte, headerLen)
	binary.LittleEndian.PutUint64(byteArray, uint64(num))

	return byteArray
}

func BytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes[:]))
}
