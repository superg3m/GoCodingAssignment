package Utility

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type StackAny any
type HeapAny any

func GetDSN(dbEngine, host string, port int, user, password, dbname string) string {
	switch dbEngine {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			user, password, host, port, dbname,
		)
	case "sqlite3":
		return dbname // For SQLite, `dbname` is the file path
	case "sqlserver":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			user, password, host, port, dbname,
		)
	default:
		return ""
	}
}

func NewTime(newTime time.Time) *time.Time {
	ret := new(time.Time)
	*ret = newTime

	return ret
}

func GetBit(number int, bit_to_check int) int {
	return ((number & (1 << bit_to_check)) >> bit_to_check)
}

func SetBit(number *int, bit_to_set int) {
	*number |= (1 << bit_to_set)
}

func UnsetBit(number *int, bit_to_unset int) {
	*number &= (^(1 << bit_to_unset))
}

func ToggleBit(number *int, bit_to_toggle int) {
	*number ^= (1 << bit_to_toggle)
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func CastAny[T any](v any) T {
	var result any

	switch any(new(T)).(type) {
	case *int:
		switch v := v.(type) {
		case int:
			result = v
		case string:
			intValue, err := strconv.Atoi(v)
			AssertOnError(err)
			result = intValue
		case float64:
			result = int(v)
		case bool:
			if v {
				result = 1
			} else {
				result = 0
			}
		}
	case *string:
		switch v := v.(type) {
		case int:
			result = strconv.Itoa(v)
		case float64:
			result = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			result = strconv.FormatBool(v)
		case string:
			result = v
		}
	case *bool:
		switch v := v.(type) {
		case bool:
			result = v
		case string:
			boolValue, err := strconv.ParseBool(v)
			AssertOnError(err)
			result = boolValue
		case int:
			result = v != 0
		case float64:
			result = v != 0.0
		}
	case *float64:
		switch v := v.(type) {
		case float64:
			result = v
		case string:
			floatValue, err := strconv.ParseFloat(v, 64)
			AssertOnError(err)
			result = floatValue
		case int:
			result = float64(v)
		case bool:
			if v {
				result = 1.0
			} else {
				result = 0.0
			}
		}
	}

	ret, ok := result.(T)
	Assert(ok)
	return ret
}

func findFileAtDepth(filename string, maxDepth int) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for depth := 0; depth <= maxDepth; depth++ {
		filePath := filepath.Join(currentDir, filename)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil // File found
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("error checking file %s: %w", filePath, err)
		}

		currentDir = filepath.Dir(currentDir)
	}

	return "", fmt.Errorf("file %s not found within %d directories", filename, maxDepth)
}

func GetSiteSettings() map[string]any {
	filename := "siteSettings.json"
	var MaxSearchDepth = 3

	filePath, err := findFileAtDepth(filename, MaxSearchDepth)
	AssertOnError(err)

	byteData, err := os.ReadFile(filePath)
	AssertOnError(err)

	var settings map[string]any
	err = json.Unmarshal(byteData, &settings)
	AssertOnError(err)

	return settings
}

func Sha256HashString(str string) string {
	h := sha256.New()
	h.Write([]byte(str))

	return hex.EncodeToString(h.Sum(nil))
}
