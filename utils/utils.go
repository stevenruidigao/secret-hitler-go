package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	// "math/rand"
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

func RandInt(start uint32, stop uint32) uint64 {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	// fmt.Println(binary.BigEndian.Uint32(bytes))
	// fmt.Println(start, stop-start, binary.BigEndian.Uint64(bytes), binary.BigEndian.Uint64(bytes)*(stop-start), start+(binary.BigEndian.Uint64(bytes)*(stop-start))/18446744073709551615)
	return uint64(start) + uint64(binary.BigEndian.Uint32(bytes))*uint64(stop-start)/4294967295
	// 11537050802889002836
	// 18446744073709551615
	// 11750154010745732844
	// 499711006744772735
	// 4848356644344510913
	// 13435754436415632515
}

func GenerateCombination(adjectives int, delimiter string, capitalize bool) string {
	result := ""

	for i := 0; i < adjectives; i++ {

		if capitalize {
			result += strings.Title(Adjectives[RandInt(0, uint32(len(Adjectives)))]) + delimiter

		} else {
			result += strings.ToLower(Adjectives[RandInt(0, uint32(len(Adjectives)))]) + delimiter
		}
	}

	if capitalize {
		result += strings.Title(Animals[RandInt(0, uint32(len(Animals)))])

	} else {
		result += strings.ToLower(Animals[RandInt(0, uint32(len(Animals)))])
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
