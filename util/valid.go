/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-14 20:28:45
 * @LastEditTime: 2023-10-15 23:31:03
 * @LastEditors: yuanshisan
 */
package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ValidateMobile(phone string) bool {
	if _, err := regexp.Match(`^0?(1[3-9][0-9])[0-9]{8}$`, []byte(phone)); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ValidateIp(ip string) bool {
	if _, err := regexp.Match(`^(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`, []byte(ip)); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if _, err := regexp.Match(`^(?!\d+$)(?![a-zA-Z]+$)(?![!@#$%^&*,.]+$)([a-zA-z\d]|[a-zA-z!@#$%^&*,.]|[\d!@#$%^&*,.]|[a-zA-Z\d!@#$%^&*,.]){8,15}$`, []byte(password)); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func ValidateIDCard(idCard string) bool {
	idCard = strings.ToUpper(idCard)
	if _, err := regexp.Match(`^(\d{15})|(\d{17}(\d|X))$`, []byte(idCard)); err != nil {
		fmt.Println(err)
		return false
	}
	if len(idCard) == 15 {
		//校验15位
		reg := `^(\d{6})+(\d{2})+(\d{2})+(\d{2})+(\d{3})$`
		if _, err := regexp.Match(reg, []byte(idCard)); err != nil {
			fmt.Println(err)
			return false
		} else {
			re := regexp.MustCompile(reg)
			matches := re.FindSubmatch([]byte(idCard))
			birthday := fmt.Sprintf("19%s-%s-%s", matches[2], matches[3], matches[4])

			//检查生日日期是否正确
			if _, err := DateToTime(birthday); err != nil {
				return false
			}
		}
	} else {
		//校验18位
		reg := `^(\d{6})+(\d{4})+(\d{2})+(\d{2})+(\d{3})([0-9]|X)$`
		if _, err := regexp.Match(reg, []byte(idCard)); err != nil {
			fmt.Println(err)
			return false
		} else {
			re := regexp.MustCompile(reg)
			matches := re.FindSubmatch([]byte(idCard))
			birthday := fmt.Sprintf("%s-%s-%s", matches[2], matches[3], matches[4])

			//检查生日日期是否正确
			if _, err := DateToTime(birthday); err != nil {
				return false
			} else {
				//检验18位身份证的校验码是否正确。
				//校验位按照ISO 7064:1983.MOD 11-2的规定生成，X可以认为是数字10。
				arr_int := [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
				arr_ch := [11]string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
				sign := 0
				for i := 0; i < 17; i++ {
					b, _ := strconv.Atoi(string(idCard[i]))
					w := arr_int[i]
					sign += b * w
				}
				n := sign % 11
				val := arr_ch[n]
				if val != string(idCard[17]) {
					return false
				}
			}
		}
	}
	return true
}

func ValidateUsername(username string) bool {
	guestExp := `\xA1\xA1|\xAC\xA3|^Guest|^\xD3\xCE\xBF\xCD|\xB9\x43\xAB\xC8`
	reg := `(?is)(\s+|^c:\\con\\con|[%,\|\*\"\s\<\>\&]|` + guestExp + ")"
	if len(username) > 20 || len(username) < 3 {
		return false
	}
	if ok, _ := regexp.Match(reg, []byte(username)); ok {
		return false
	}
	if _, err := regexp.Match(`^[0-9@A-Z_a-z]*$`, []byte(username)); err != nil {
		return false
	}
	if ok, _ := regexp.Match(`^(?!_|\s\')[A-Za-z0-9_\x80-\xff\s\']+$`, []byte(username)); ok {
		return true
	}
	return false
}
