package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"math"
)

func EncryptResult(data []map[string]interface{}, compress bool) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(0x7F) // magic number

	for _, row := range data {
		buf.WriteByte(0x01) // row start
		for k, v := range row {
			// Key obfuscation
			ks := []byte(k)
			for i := range ks {
				ks[i] ^= 0xAA
			}
			buf.Write(ks)
			buf.WriteByte(0x1F) // key-value separator

			// Value encoding
			switch val := v.(type) {
			case int64:
				buf.WriteByte(0x02) // INT type
				var b [8]byte
				binary.LittleEndian.PutUint64(b[:], uint64(val))
				buf.Write(b[:])
			case float64:
				buf.WriteByte(0x03) // FLOAT type
				var b [8]byte
				binary.LittleEndian.PutUint64(b[:], math.Float64bits(val))
				buf.Write(b[:])
			case string:
				buf.WriteByte(0x04) // STRING type
				buf.Write([]byte(val))
			}
			buf.WriteByte(0x1E) // field end
		}
		buf.WriteByte(0x00) // row end
	}

	// Add checksum
	checksum := sha256.Sum256(buf.Bytes())
	buf.Write(checksum[:4])

	if compress && buf.Len() > 1024 {
		return compressData(buf.Bytes())
	}

	return buf.Bytes()
}

func compressData(data []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data)
	w.Close()
	return buf.Bytes()
}
