package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
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
