package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Print("Ошибка кодирования!")
	}
	return data
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	if len(data) > 0 && passphrase != "" {
		block, _ := aes.NewCipher([]byte(createHash(passphrase)))
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Print("Ошибка шифрования, " + err.Error())
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			log.Print("Ошибка шифрования, " + err.Error())
		}
		ciphertext := gcm.Seal(nonce, nonce, data, nil)
		return []byte(encodeBase64(ciphertext))
	} else {
		return []byte("")
	}
}

func decrypt(data []byte, passphrase string) []byte {
	if len(data) > 0 && passphrase != "" {
		data = decodeBase64(string(data))

		key := []byte(createHash(passphrase))
		block, err := aes.NewCipher(key)
		if err != nil {
			log.Print("Ошибка расшифровки, " + err.Error())
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Print("Ошибка расшифровки, " + err.Error())
		}
		nonceSize := gcm.NonceSize()
		nonce, ciphertext := data[:nonceSize], data[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			log.Print("Ошибка расшифровки, " + err.Error())
		}
		return plaintext
	} else {
		return []byte("")
	}
}

//sendResponse send response
func sendResponse(logger *zap.Logger, w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Error("couldn't send data to connection", zap.Error(err))
	}
}
