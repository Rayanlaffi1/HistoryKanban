package api

import (
	"bytes"
	"io"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

var hautWindows1252 = map[byte]rune{
	0x80: '€', 0x82: '‚', 0x83: 'ƒ', 0x84: '„', 0x85: '…', 0x86: '†', 0x87: '‡',
	0x88: 'ˆ', 0x89: '‰', 0x8A: 'Š', 0x8B: '‹', 0x8C: 'Œ', 0x8E: 'Ž', 0x91: '‘',
	0x92: '’', 0x93: '“', 0x94: '”', 0x95: '•', 0x96: '–', 0x97: '—', 0x98: '˜',
	0x99: '™', 0x9A: 'š', 0x9B: '›', 0x9C: 'œ', 0x9E: 'ž', 0x9F: 'Ÿ',
}

func versUTF8(octets []byte) []byte {
	var tampon bytes.Buffer
	for _, octet := range octets {
		if caractere, present := hautWindows1252[octet]; present {
			tampon.WriteRune(caractere)
			continue
		}
		tampon.WriteRune(rune(octet))
	}
	return tampon.Bytes()
}

func normaliserCorps() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil || c.Request.ContentLength == 0 {
			c.Next()
			return
		}
		octets, erreur := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		if erreur != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"erreur": "lecture du corps impossible"})
			return
		}
		octets = bytes.TrimPrefix(octets, []byte{0xEF, 0xBB, 0xBF})
		octets = bytes.TrimSpace(octets)
		if !utf8.Valid(octets) {
			octets = versUTF8(octets)
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(octets))
		c.Request.ContentLength = int64(len(octets))
		c.Next()
	}
}
