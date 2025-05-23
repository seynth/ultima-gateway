package auxiliary

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
)

type Header struct {
	Key string
	Val string
}

func ConvertHeaders(headers network.Headers) []*fetch.HeaderEntry {
	var result []*fetch.HeaderEntry
	for k, v := range headers {
		if s, ok := v.(string); ok {
			result = append(result, &fetch.HeaderEntry{
				Name:  k,
				Value: s,
			})
		}
	}
	return result
}

func AddAndOverwriteHeaders(customHeaders []Header, original []*fetch.HeaderEntry) []*fetch.HeaderEntry {
	headerMap := make(map[string]*fetch.HeaderEntry)
	for _, h := range original {
		headerMap[strings.ToLower(h.Name)] = h
	}

	for _, custom := range customHeaders {
		lowerKey := strings.ToLower(custom.Key)
		if existing, found := headerMap[lowerKey]; found {
			existing.Value = custom.Val
		} else {
			newEntry := &fetch.HeaderEntry{
				Name:  custom.Key,
				Value: custom.Val,
			}
			original = append(original, newEntry)
			headerMap[lowerKey] = newEntry
		}
	}

	return original
}

func Sha256Encode(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
