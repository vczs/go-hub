package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	createKey("private_key.pem", "public_key.pem")

	cipher, _ := encrypt("public_key.pem", "你好")
	fmt.Println("密文:", base64.StdEncoding.EncodeToString(cipher))

	ciphertext, _ := decrypt("private_key.pem", cipher)
	fmt.Println("密文解密:", string(ciphertext))
}

func encrypt(publicKeyFile, content string) ([]byte, error) {
	// 加载公钥
	pemData, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, err
	}

	// 解码PEM
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// 解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	// 加密
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(content))
	if err != nil {
		return nil, err
	}

	return cipher, nil
}

func decrypt(privateKeyFile string, cipher []byte) ([]byte, error) {
	// 加载私钥
	pemData, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}

	// 解码PEM
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 解密
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipher)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func createKey(privateKeyFile, publicKeyFile string) error {
	// 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 私钥
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err := os.WriteFile(privateKeyFile, privateKeyPEM, 0600); err != nil {
		return err
	}

	// 公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err = os.WriteFile(publicKeyFile, publicKeyPEM, 0644); err != nil {
		return err
	}

	return nil
}
