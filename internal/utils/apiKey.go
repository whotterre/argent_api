package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

func GenerateNewAPIKeyString() string {
	prefix := "sk_live_"

	end := GenString(45)

	return prefix + end
}

func ExpiryStringToTimestamp(expiryStr string) (time.Time, error) {
	now := time.Now()
	switch expiryStr {
	case "1H":
		return now.Add(time.Hour), nil
	case "1D":
		return now.AddDate(0, 0, 1), nil 
	case "1M":
		return now.AddDate(0, 1, 0), nil 
	case "1Y":
		return now.AddDate(1, 0, 0), nil 
	default:
		return time.Time{}, errors.New("invalid expiry string") 
	}
}


func GenString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	result := base64.RawURLEncoding.EncodeToString(bytes)
	return result
}