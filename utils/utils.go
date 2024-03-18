package utils

import (
	"fmt"
	"os"
	"strconv"

	"time"

	"github.com/google/uuid"
)

func GetEnvFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return i
}

func IsDateValue(stringDate string) bool {
	_, err := time.Parse("2006-01-02", stringDate)
	fmt.Println(err)
	return err == nil
}

func GetUUID(name string) string {
	namespaceDNS := uuid.NameSpaceDNS
	uuidV5 := uuid.NewSHA1(namespaceDNS, []byte(name))
	return uuidV5.String()
}
