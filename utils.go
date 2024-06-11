package main

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

const headerLen = 8

func validateIp(ip string) (string, error) {
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("ip error")
}

func validatePort(port string) (int, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, err
	}
	return p, nil
}

func Int64ToBytes(num int64) []byte {
	byteArray := make([]byte, headerLen)
	binary.LittleEndian.PutUint64(byteArray, uint64(num))

	return byteArray
}

func BytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes[:]))
}
