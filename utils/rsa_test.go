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
	 "os"
	 "testing"
 )
 
 // 测试用的临时密钥文件路径
 const (
	 testPublicKeyPath  = "test_public.pem"
	 testPrivateKeyPath = "test_private.pem"
 )
 
 // 生成测试用的RSA密钥对并保存到文件
 func generateTestKeys() error {
	 privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	 if err != nil {
		 return err
	 }
 
	 // 保存私钥
	 privateFile, err := os.Create(testPrivateKeyPath)
	 if err != nil {
		 return err
	 }
	 defer privateFile.Close()
 
	 privateBlock := &pem.Block{
		 Type:  "RSA PRIVATE KEY",
		 Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	 }
	 if err := pem.Encode(privateFile, privateBlock); err != nil {
		 return err
	 }
 
	 // 保存公钥
	 publicFile, err := os.Create(testPublicKeyPath)
	 if err != nil {
		 return err
	 }
	 defer publicFile.Close()
 
	 publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	 if err != nil {
		 return err
	 }
 
	 publicBlock := &pem.Block{
		 Type:  "PUBLIC KEY",
		 Bytes: publicBytes,
	 }
	 if err := pem.Encode(publicFile, publicBlock); err != nil {
		 return err
	 }
 
	 return nil
 }
 
 // 清理测试文件
 func cleanupTestKeys() {
	 _ = os.Remove(testPublicKeyPath)
	 _ = os.Remove(testPrivateKeyPath)
 }
 
 func TestRSACrypto(t *testing.T) {
	 // 生成测试密钥
	 if err := generateTestKeys(); err != nil {
		 t.Fatalf("Failed to generate test keys: %v", err)
	 }
	 defer cleanupTestKeys()
 
	 // 创建测试用的RSACrypto实例
	 crypto := NewRSACrypto(testPublicKeyPath, testPrivateKeyPath)
 
	 t.Run("Test Algorithm", func(t *testing.T) {
		 if algo := crypto.Algorithm(); algo != "RSA-OAEP" {
			 t.Errorf("Expected algorithm 'RSA-OAEP', got '%s'", algo)
		 }
	 })
 
	 t.Run("Test Encrypt and Decrypt", func(t *testing.T) {
		 testCases := []struct {
			 name      string
			 plaintext string
		 }{
			 {"Empty string", ""},
			 {"Short text", "hello world"},
			 {"Long text", "This is a longer text that needs to be encrypted and decrypted properly using RSA algorithm."},
			 {"Special chars", "!@#$%^&*()_+-=[]{};':\",./<>?"},
		 }
 
		 for _, tc := range testCases {
			 t.Run(tc.name, func(t *testing.T) {
				 // 加密
				 cipherText, err := crypto.Encrypt(tc.plaintext)
				 if err != nil {
					 t.Fatalf("Encrypt failed: %v", err)
				 }
 
				 // 确保密文不是明文
				 if cipherText == tc.plaintext {
					 t.Error("Ciphertext is same as plaintext")
				 }
 
				 // 确保是有效的Base64
				 if _, err := base64.StdEncoding.DecodeString(cipherText); err != nil {
					 t.Errorf("Ciphertext is not valid base64: %v", err)
				 }
 
				 // 解密
				 decrypted, err := crypto.Decrypt(cipherText)
				 if err != nil {
					 t.Fatalf("Decrypt failed: %v", err)
				 }
 
				 // 比较解密结果
				 if decrypted != tc.plaintext {
					 t.Errorf("Decrypted text doesn't match original. Expected: '%s', Got: '%s'", tc.plaintext, decrypted)
				 }
			 })
		 }
	 })
 
	 t.Run("Test Invalid Keys", func(t *testing.T) {
		 // 测试无效的公钥路径
		 invalidCrypto := NewRSACrypto("nonexistent_public.pem", testPrivateKeyPath)
		 if _, err := invalidCrypto.Encrypt("test"); err == nil {
			 t.Error("Expected error with invalid public key path")
		 }
 
		 // 测试无效的私钥路径
		 invalidCrypto = NewRSACrypto(testPublicKeyPath, "nonexistent_private.pem")
		 if _, err := invalidCrypto.Decrypt("test"); err == nil {
			 t.Error("Expected error with invalid private key path")
		 }
	 })
 
	 t.Run("Test Invalid Ciphertext", func(t *testing.T) {
		 // 测试无效的Base64
		 if _, err := crypto.Decrypt("not a base64 string"); err == nil {
			 t.Error("Expected error with invalid base64 ciphertext")
		 }
 
		 // 测试无效的密文(随机Base64)
		 randomBytes := make([]byte, 256)
		 rand.Read(randomBytes)
		 randomBase64 := base64.StdEncoding.EncodeToString(randomBytes)
		 if _, err := crypto.Decrypt(randomBase64); err == nil {
			 t.Error("Expected error with random ciphertext")
		 }
	 })
 }
 
func TestNewCryptoProvider(t *testing.T) {
	// 生成测试密钥
	if err := generateTestKeys(); err != nil {
		t.Fatalf("Failed to generate test keys: %v", err)
	}
	defer cleanupTestKeys()

	t.Run("Test RSA Provider", func(t *testing.T) {
		config := map[string]string{
			"publicKeyPath":  testPublicKeyPath,
			"privateKeyPath": testPrivateKeyPath,
		}

		provider, err := NewCryptoProvider(AlgorithmRSA, config)
		if err != nil {
			t.Fatalf("Failed to create RSA provider: %v", err)
		}

		if provider.Algorithm() != "RSA-OAEP" {
			t.Errorf("Expected algorithm 'RSA-OAEP', got '%s'", provider.Algorithm())
		}

		// 测试加密解密
		plainText := "test message"
		cipherText, err := provider.Encrypt(plainText)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		decrypted, err := provider.Decrypt(cipherText)
		if err != nil {
			t.Fatalf("Decrypt failed: %v", err)
		}

		if decrypted != plainText {
			t.Errorf("Decrypted text doesn't match original. Expected: '%s', Got: '%s'", plainText, decrypted)
		}
	})

	t.Run("Test Unsupported Algorithm", func(t *testing.T) {
		_, err := NewCryptoProvider("invalid-algorithm", nil)
		if err == nil {
			t.Error("Expected error for unsupported algorithm")
		}
	})
}