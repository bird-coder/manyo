/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-27 22:13:50
 * @LastEditTime: 2023-09-27 22:17:24
 * @LastEditors: yuanshisan
 */
package util

import (
	"fmt"
	"net"
)

func GetLocalHost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}
	for i := 0; i < len(netInterfaces); i++ {
		if netInterfaces[i].Flags&net.FlagUp != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return ""
}
