/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-02-17 14:33:03
 * @LastEditTime: 2023-02-23 14:49:01
 * @LastEditors: yuanshisan
 */
package util

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// The structure of the data package is as follows:
// | (uint32) || (uint32) || (uint32) || (uint32) || (binary)
// |  4-byte  ||  4-byte  ||  4-byte  ||  4-byte  ||  N-byte
// -----------------------------------------------------------...
//    packet     protocol     userid     serverid     data
//     size         ID
//

var (
	HeaderLen = 16
)

type MsgHeader struct {
	Len uint32
	Pid uint32
	Uid uint32
	Sid uint32
}

func PackHeader(buf []byte, header MsgHeader) []byte {
	headerBuf := make([]byte, HeaderLen, HeaderLen)
	pkgLen := HeaderLen + len(buf)
	binary.BigEndian.PutUint32(headerBuf[0:4], uint32(pkgLen))
	binary.BigEndian.PutUint32(headerBuf[4:8], header.Pid)
	binary.BigEndian.PutUint32(headerBuf[8:12], header.Uid)
	binary.BigEndian.PutUint32(headerBuf[12:16], header.Sid)
	return headerBuf
}

func UnPackHeader(buf []byte) (MsgHeader, error) {
	length := len(buf)
	var header MsgHeader
	if length < HeaderLen {
		err := errors.New("Read msg size failed")
		return header, err
	}
	lenBuf := bytes.NewBuffer(buf[0:4])
	pBuf := bytes.NewBuffer(buf[4:8])
	userBuf := bytes.NewBuffer(buf[8:12])
	serverBuf := bytes.NewBuffer(buf[12:16])
	binary.Read(lenBuf, binary.BigEndian, &header.Len)
	binary.Read(pBuf, binary.BigEndian, &header.Pid)
	binary.Read(userBuf, binary.BigEndian, &header.Uid)
	binary.Read(serverBuf, binary.BigEndian, &header.Sid)
	return header, nil
}
