// Package cryptox reimplementa o subconjunto do formato de criptografia da
// biblioteca PHP defuse/php-encryption (v2) necessário para decriptar valores
// gerados pelo zeum-admin-api (ex.: Application.apiKey), permitindo validar
// a mesma API Key sem duplicar o segredo em texto plano entre os serviços.
package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
)

const (
	keyHeaderSize        = 4
	keyByteSize          = 32
	checksumByteSize     = 32
	ciphertextHeaderSize = 4
	saltByteSize         = 32
	ivByteSize           = 16
	macByteSize          = 32
	minCiphertextSize    = ciphertextHeaderSize + saltByteSize + ivByteSize + macByteSize

	encryptionInfo     = "DefusePHP|V2|KeyForEncryption"
	authenticationInfo = "DefusePHP|V2|KeyForAuthentication"
)

var (
	keyHeader        = []byte{0xDE, 0xF0, 0x00, 0x00}
	ciphertextHeader = []byte{0xDE, 0xF5, 0x02, 0x00}

	ErrInvalidKey        = errors.New("chave em formato inválido")
	ErrInvalidCiphertext = errors.New("ciphertext em formato inválido")
	ErrIntegrityCheck    = errors.New("falha na verificação de integridade")
)

// ParseKey decodifica uma chave gerada por Defuse\Crypto\Key::saveToAsciiSafeString()
// (o valor usado como DEFUSE_SECRET no zeum-admin-api) e retorna os 32 bytes brutos da chave.
func ParseKey(asciiSafeKey string) ([]byte, error) {

	raw, err := hex.DecodeString(asciiSafeKey)

	if err != nil {
		return nil, ErrInvalidKey
	}

	if len(raw) != keyHeaderSize+keyByteSize+checksumByteSize {
		return nil, ErrInvalidKey
	}

	header := raw[:keyHeaderSize]
	body := raw[:len(raw)-checksumByteSize]
	checksum := raw[len(raw)-checksumByteSize:]

	if !bytes.Equal(header, keyHeader) {
		return nil, ErrInvalidKey
	}

	expectedChecksum := sha256.Sum256(body)

	if subtle.ConstantTimeCompare(checksum, expectedChecksum[:]) != 1 {
		return nil, ErrInvalidKey
	}

	return raw[keyHeaderSize : len(raw)-checksumByteSize], nil
}

// Decrypt reproduz Defuse\Crypto\Crypto::decrypt() (formato v2) a partir da chave bruta,
// permitindo validar ciphertexts gerados pelo zeum-admin-api.
func Decrypt(hexCiphertext string, rawKey []byte) (string, error) {

	data, err := hex.DecodeString(hexCiphertext)

	if err != nil {
		return "", ErrInvalidCiphertext
	}

	if len(data) < minCiphertextSize {
		return "", ErrInvalidCiphertext
	}

	header := data[:ciphertextHeaderSize]

	if !bytes.Equal(header, ciphertextHeader) {
		return "", ErrInvalidCiphertext
	}

	salt := data[ciphertextHeaderSize : ciphertextHeaderSize+saltByteSize]
	iv := data[ciphertextHeaderSize+saltByteSize : ciphertextHeaderSize+saltByteSize+ivByteSize]
	mac := data[len(data)-macByteSize:]
	encrypted := data[ciphertextHeaderSize+saltByteSize+ivByteSize : len(data)-macByteSize]

	authKey, err := hkdf.Key(sha256.New, rawKey, salt, authenticationInfo, keyByteSize)

	if err != nil {
		return "", err
	}

	encKey, err := hkdf.Key(sha256.New, rawKey, salt, encryptionInfo, keyByteSize)

	if err != nil {
		return "", err
	}

	mac2 := hmac.New(sha256.New, authKey)
	mac2.Write(header)
	mac2.Write(salt)
	mac2.Write(iv)
	mac2.Write(encrypted)
	expectedMAC := mac2.Sum(nil)

	if !hmac.Equal(mac, expectedMAC) {
		return "", ErrIntegrityCheck
	}

	block, err := aes.NewCipher(encKey)

	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(encrypted))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, encrypted)

	return string(plaintext), nil
}
