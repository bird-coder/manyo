/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-22 23:16:06
 * @LastEditTime: 2024-05-23 17:11:12
 * @LastEditors: yujiajie
 */
package util

import (
	"math/rand"
	"regexp"
	"strings"
)

const (
	symbol = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~"
	letter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func MultiReplace(s string, words []string, replace string) string {
	tmp := make([]string, 0, len(words)*2)
	for _, v := range words {
		tmp = append(tmp, v, replace)
	}
	r := strings.NewReplacer(tmp...)
	return r.Replace(s)
}

func MultiContains(s string, words []string) bool {
	reg, _ := regexp.Compile(strings.Join(words, "|"))
	return reg.Match([]byte(s))
}

func GenerateRandomKey(length int) string {
	return generateRandString(length, letter)
}

func generateRandString(length int, s string) string {
	var chars = []byte(s)
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4))
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
