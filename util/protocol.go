/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-02-17 14:33:03
 * @LastEditTime: 2023-02-17 14:49:09
 * @LastEditors: yuanshisan
 */
package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// The structure of the data package is as follows:
// | (uint32) || (uint32) || (uint32) || (uint32) || (binary)
// |  4-byte  ||  4-byte  ||  4-byte  ||  4-byte  ||  N-byte
// -----------------------------------------------------------...
//    packet     protocol     userid     serverid     data
//     size         ID
//

func PackHeader(buf []byte, pid uint32, uid uint32, sid uint32) []byte {
	headerBuf := make([]byte, 16, 16)
	pkgLen := 16 + uint32(len(buf))
	binary.BigEndian.PutUint32(headerBuf[0:4], pkgLen)
	binary.BigEndian.PutUint32(headerBuf[4:8], pid)
	binary.BigEndian.PutUint32(headerBuf[8:12], uid)
	binary.BigEndian.PutUint32(headerBuf[12:16], sid)
	return headerBuf
}

func UnPackHeader(buf []byte) ([]uint32, error) {
	length := len(buf)
	headers := make([]uint32, 4, 4)
	if length < 16 {
		err := fmt.Errorf("Read msg size failed")
		return headers, err
	}
	pkgBuf := bytes.NewBuffer(buf[0:4])
	pBuf := bytes.NewBuffer(buf[4:8])
	userBuf := bytes.NewBuffer(buf[8:12])
	serverBuf := bytes.NewBuffer(buf[12:16])
	binary.Read(pkgBuf, binary.BigEndian, &headers[0])
	binary.Read(pBuf, binary.BigEndian, &headers[1])
	binary.Read(userBuf, binary.BigEndian, &headers[2])
	binary.Read(serverBuf, binary.BigEndian, &headers[3])
	return headers, nil
}
