/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-03-13 17:31:09
 * @LastEditTime: 2023-09-27 22:06:10
 * @LastEditors: yuanshisan
 */
package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func EncryptByAes(data []byte, key []byte) (string, error) {
	res, err := aesEncrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func DecryptByAes(data string, key []byte) ([]byte, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return aesDecrypt(dataBytes, key)
}

// aes 加密
func aesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

// aes 解密
func aesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

// pkcs7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, paddingText...)
}

// pkcs7 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("encrypt string is invalid")
	}
	padding := int(data[length-1])
	return data[:(length - padding)], nil
}

func Md5Encode(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
