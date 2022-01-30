package utils

import (
	// "crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type Hash struct {
	Value string
	Salt  string
}

func AddCookie(writer http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expire,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, &cookie)
}

func GetCookie(cookieString string, key string) string {
	if strings.Contains(cookieString, key) {
		return strings.Split(strings.Split(cookieString, key+"=")[1], ";")[0]

	} else {
		return ""
	}
}

func GenerateCombination(adjectives int, delimiter string, capitalize bool) string {
	result := ""

	for i := 0; i < adjectives; i++ {
		if capitalize {
			result += strings.Title(Adjectives[rand.Intn(len(Adjectives))]) + delimiter

		} else {
			result += strings.ToLower(Adjectives[rand.Intn(len(Adjectives))]) + delimiter
		}
	}

	if capitalize {
		result += strings.Title(Animals[rand.Intn(len(Animals))])

	} else {
		result += strings.ToLower(Animals[rand.Intn(len(Animals))])
	}

	return result
}

func MarshalJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	return string(jsonData), err
}

func JSONResponse(writer http.ResponseWriter, data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		http.Error(writer, "Error creating JSON response", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	fmt.Fprintf(writer, "%s", jsonData)
}

func StringToBytes(str string) []byte {
	return []byte(str)
}

func BytesToHexString(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func HexStringToBytes(str string) []byte {
	result, err := hex.DecodeString(str)

	if err != nil {
		return nil
	}

	return result
}

func Argon2(password string, salt string) string {
	return BytesToHexString(argon2.IDKey(StringToBytes(password), StringToBytes(salt), 1, 64*1024, 4, 32))
}

func SHA512(str string) string {
	hash := sha512.New()
	hash.Write(StringToBytes(str))
	return BytesToHexString(hash.Sum([]byte{}))
}
