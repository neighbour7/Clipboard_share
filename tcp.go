package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

type TcpMsg struct {
	Type    string
	Content []byte
}

type Tcp struct {
	conn    net.Conn
	name    string
	watchCh chan *TcpMsg
}

func NewTcp(c net.Conn) *Tcp {
	return &Tcp{
		conn:    c,
		name:    c.RemoteAddr().String(),
		watchCh: make(chan *TcpMsg, 1),
	}
}

func (t *Tcp) Watch() <-chan *TcpMsg {
	go func() {
		for {
			c, err := t.Read()
			if err != nil {
				panic(err)
			}
			t.watchCh <- c
		}
	}()
	return t.watchCh
}

func (t *Tcp) Read() (*TcpMsg, error) {
	msgLenInfo := make([]byte, 8)
	binLen, err := t.conn.Read(msgLenInfo)
	if err != nil {
		return nil, err
	}
	if binLen != 8 {
		return nil, errors.New("msg len not match: get " + fmt.Sprint(binLen))
	}
	msgLen := BytesToInt64(msgLenInfo)

	binContent := make([]byte, msgLen)
	var totalBinContentLen int
	for {
		n, err := t.conn.Read(binContent[totalBinContentLen:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		totalBinContentLen += n
		if totalBinContentLen == len(binContent) {
			break
		}
	}

	if int64(totalBinContentLen) != msgLen {
		return nil, errors.New("content len not match: get" + fmt.Sprint(int64(totalBinContentLen)) + "  " + fmt.Sprint(msgLen))
	}
	msg := &TcpMsg{}
	err = json.Unmarshal(binContent, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (t *Tcp) Send(msgType string, content []byte) error {
	msg := &TcpMsg{
		Type:    msgType,
		Content: content,
	}
	msgBin, err := json.Marshal(msg)
	if err != nil {
		return errors.New("json marshll error")
	}
	msgLen := len(msgBin)
	contentLenBytes := Int64ToBytes(int64(msgLen))
	_, err = t.conn.Write(contentLenBytes)
	if err != nil {
		return errors.New("send headerLen error")
	}
	byteLen, err := t.conn.Write(msgBin)

	if err != nil {
		return err
	}
	if byteLen != msgLen {
		return errors.New("byteLen != msgLen")
	}
	fmt.Println("send successfully: ", msgType)
	return nil
}

func (t *Tcp) Close() {
	if t.watchCh != nil {
		close(t.watchCh)
	}
	t.conn.Close()
}
