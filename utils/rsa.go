/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: bizhiyuan <bizhiyuan@kylinos.cn>
 * Date: Thu Mar 27 09:32:28 2025 +0800
 */

 package utils

 import (
	 "crypto/rand"
	 "crypto/rsa"
	 "crypto/x509"
	 "encoding/base64"
	 "encoding/pem"
	 "errors"
	 "fmt"
	 "os"
 )
 
 // CryptoProvider 加密算法提供者接口
 type CryptoProvider interface {
	 Encrypt(plaintext string) (string, error)
	 Decrypt(ciphertext string) (string, error)
	 Algorithm() string // 返回算法名称
 }
 
 // RSACrypto RSA加密实现
 type RSACrypto struct {
	 publicKeyPath  string
	 privateKeyPath string
 }
 
 // NewRSACrypto 创建RSA加密处理器
 func NewRSACrypto(publicKeyPath, privateKeyPath string) *RSACrypto {
	 return &RSACrypto{
		 publicKeyPath:  publicKeyPath,
		 privateKeyPath: privateKeyPath,
	 }
 }
 
 func (r *RSACrypto) Algorithm() string {
	 return "RSA-OAEP"
 }
 
 // Encrypt 实现RSA加密
 func (r *RSACrypto) Encrypt(plaintext string) (string, error) {
	 publicKey, err := os.ReadFile(r.publicKeyPath)
	 if err != nil {
		 return "", fmt.Errorf("read public key failed: %v", err)
	 }
 
	 block, _ := pem.Decode(publicKey)
	 if block == nil {
		 return "", errors.New("public key decode failed")
	 }
 
	 pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	 if err != nil {
		 return "", fmt.Errorf("public key parse failed: %v", err)
	 }
 
	 rsaPub, ok := pub.(*rsa.PublicKey)
	 if !ok {
		 return "", errors.New("not a valid RSA public key")
	 }
 
	 cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(plaintext))
	 if err != nil {
		 return "", fmt.Errorf("encrypt failed: %v", err)
	 }
	 // 二进制密文 -> Base64字符串
	 return base64.StdEncoding.EncodeToString(cipherText), nil
 }
 
 // 算法名称
 const (
	 AlgorithmRSA string = "RSA"
	 AlgorithmAES string = "AES"
 )
 
 // NewCryptoProvider 创建加密处理器
 func NewCryptoProvider(algorithm string, config map[string]string) (CryptoProvider, error) {
	 switch algorithm {
	 case AlgorithmRSA:
		 return NewRSACrypto(config["publicKeyPath"], config["privateKeyPath"]), nil
	 default:
		 return nil, fmt.Errorf("不支持的加密算法: %s", algorithm)
	 }
 }
 
 func ReadPublicKey(filePath string) ([]byte, error) {
	 data, err := os.ReadFile(filePath)
	 if err != nil {
		 return nil, fmt.Errorf("read public key failed: %v", err)
	 }
	 return data, nil
 }
 