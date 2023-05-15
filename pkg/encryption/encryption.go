package encryption

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

type encryptedData struct {
	Parts [][]byte
}

func EncryptData(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	keySize := publicKey.Size()
	chunkSize := keySize - 100 // вычитаем размер дополнительных данных RSA

	numChunks := (len(data) + chunkSize - 1) / chunkSize

	encryptedParts := make([][]byte, numChunks)

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize

		if end > len(data) {
			end = len(data)
		}

		chunk := data[start:end]

		encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, chunk, nil)
		if err != nil {
			return nil, fmt.Errorf("encryption error: %v", err)
		}

		encryptedParts[i] = encrypted
	}

	encData := encryptedData{Parts: encryptedParts}
	jsonData, err := json.Marshal(encData)
	if err != nil {
		return nil, fmt.Errorf("JSON encoding error: %v", err)
	}

	return jsonData, nil
}

func DecryptData(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	var encrypted encryptedData
	err := json.Unmarshal(data, &encrypted)
	if err != nil {
		return nil, fmt.Errorf("JSON decoding error: %v", err)
	}

	numChunks := len(encrypted.Parts)
	decryptedParts := make([][]byte, numChunks)

	for i, encryptedPart := range encrypted.Parts {
		decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedPart, nil)
		if err != nil {
			return nil, fmt.Errorf("decryption error: %v", err)
		}

		decryptedParts[i] = decrypted
	}

	decryptedData := bytes.Join(decryptedParts, nil)
	return decryptedData, nil
}

func parseFile(filePath string) (*pem.Block, error) {
	keyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block from private key file %s", filePath)
	}
	return block, nil
}

func GetPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	block, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func GetPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	block, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key file %s contains a non-RSA key", filePath)
	}
	return rsaPubKey, nil
}

func NewRandomKeys() (*rsa.PrivateKey, *rsa.PublicKey) {
	bits := 2048
	privKey, _ := rsa.GenerateKey(rand.Reader, bits)

	return privKey, &privKey.PublicKey
}
