package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}

func EnsureDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|os.ModePerm)
	} else {
		return err
	}
}

func GetEnvVarAsBool(name string, defaultValue bool) bool {
	s := os.Getenv(name)
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s == "1" || s == "y" || s == "yes" || s == "true" {
		return true
	}
	return defaultValue
}

func IntFromString(s string, defaultValue int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("WARNING: Using default int value for '%s'\n", s)
		n = defaultValue
	}
	return n
}

func DateFromString(s string, defaultDate time.Time) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Printf("WARNING: Using default date value for '%s'\n", s)
		t = defaultDate
	}
	return t
}

func DoesFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			log.Println("ERROR: os.Stat() failed. err:", err)
			return false
		}
	}
	return true
}
