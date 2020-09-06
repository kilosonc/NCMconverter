package converter

import (
	"crypto/aes"
)

func decryptAes128(key, data []byte) ([]byte, error) {
	data = data[:len(data)/aes.BlockSize*aes.BlockSize]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	dataLen := len(data)
	decryptedData := make([]byte, dataLen)
	bs := block.BlockSize()
	for i := 0; i <= dataLen-bs; i += bs {
		block.Decrypt(decryptedData[i:i+bs], data[i:i+bs])
	}

	length := len(decryptedData)
	unpadding := int(decryptedData[length-1])
	return decryptedData[:(length - unpadding)], nil
}

func buildKeyBox(key []byte) []byte {
	box := make([]byte, 256)
	for i := 0; i < 256; i++ {
		box[i] = byte(i)
	}
	keyLen := byte(len(key))
	var c, lastByte, keyOffset byte
	for i := 0; i < 256; i++ {
		c = (box[i] + lastByte + key[keyOffset]) & 0xff
		keyOffset++
		if keyOffset >= keyLen {
			keyOffset = 0
		}
		box[i], box[c] = box[c], box[i]
		lastByte = c
	}
	return box
}
