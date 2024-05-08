package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

func main() {
	// 密钥
	key := "acac3871a3ff9c98e5d2dccc547ea37c"

	// 随机数
	nonce := make([]byte, aes.BlockSize)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	fmt.Println("随机数:", base64.StdEncoding.EncodeToString(nonce))

	// 加密
	content, _ := encrypt([]byte(key), nonce, "你好")
	fmt.Println("密文:", base64.StdEncoding.EncodeToString(content))

	// 解密
	ciphertext, _ := decrypt([]byte(key), nonce, content)
	fmt.Println("解密:", string(ciphertext))
}

func encrypt(key, nonce []byte, content string) ([]byte, error) {
	// 填充
	contentBytes := []byte(content)
	padding := aes.BlockSize - len(contentBytes)%aes.BlockSize
	contentFinal := append(contentBytes, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// 加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil
	}
	cipherByte := make([]byte, len(contentFinal))
	cipher.NewCBCEncrypter(block, nonce).CryptBlocks(cipherByte, contentFinal)

	return cipherByte, nil
}

func decrypt(key, nonce, content []byte) ([]byte, error) {
	// 解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(content))
	cipher.NewCBCDecrypter(block, nonce).CryptBlocks(ciphertext, content)

	// 去填充
	if ciphertext, err = removePadding(ciphertext); err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func removePadding(data []byte) ([]byte, error) {
	length, padLength := len(data), int(data[len(data)-1])
	if length == 0 {
		return nil, errors.New("empty input")
	}
	if length%aes.BlockSize != 0 {
		return nil, errors.New("input is not a multiple of the block size")
	}
	if padLength > aes.BlockSize || padLength == 0 {
		return nil, errors.New("invalid padding")
	}
	for i := 1; i <= padLength; i++ {
		if data[length-i] != byte(padLength) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:length-padLength], nil
}
