package readablebytes

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Exported units abbreviations
const (
	// Decimal

	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB

	// Binary

	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
	PiB = 1024 * TiB
)

var units = []byte("kmgtpezy")
var decUnits = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
var binUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

// HumanSize returns a human-readable size in decimal units (eg. "32KB", "32MB").
func HumanSize(size float64) string {
	return HumanDecimalSize(size)
}

// HumanBinarySize returns a human-readable size in binary units (eg. "32kiB", "32MiB").
func HumanBinarySize(size float64) string {
	return humanSizeWithPrecision(size, 4, 1024.0, binUnits)
}

// HumanDecimalSize returns a human-readable size in decimal units (eg. "32KB", "32MB").
func HumanDecimalSize(size float64) string {
	return humanSizeWithPrecision(size, 4, 1000.0, decUnits)
}

// FromHumanString returns an int64 bytes size from a human-readable string
func FromHumanString(size string) (int64, error) {
	return parseString(size)
}

func getSizeAndUnit(size float64, base float64, _map []string) (float64, string) {
	i := 0
	unitsLimit := len(_map) - 1
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	return size, _map[i]
}

func humanSizeWithPrecision(size float64, precision int, base float64, _map []string) string {
	size, unit := getSizeAndUnit(size, base, _map)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

func parseString(sizeStr string) (int64, error) {
	sizeStr = strings.ToLower(sizeStr)
	strLen := len(sizeStr)
	if strLen == 0 {
		return -1, fmt.Errorf("invalid size: '%s'", sizeStr)
	}
	var unitPrefixPos, lastNumberPos int
	var binary bool
	for i := strLen - 1; i >= 0; i-- {
		if unicode.IsDigit(rune(sizeStr[i])) {
			lastNumberPos = i
			break
		}

		if sizeStr[i] == 'i' {
			binary = true
			continue
		}

		if sizeStr[i] != ' ' {
			unitPrefixPos = i
		}
	}

	size, err := strconv.ParseFloat(sizeStr[:lastNumberPos+1], 64)
	if err != nil {
		return -1, err
	}

	if size < 0 {
		return -1, fmt.Errorf("size is less than zero")
	}

	if unitPrefixPos > 0 {
		index := bytes.IndexByte(units, sizeStr[unitPrefixPos])
		if index != -1 {
			base := 1000
			if binary {
				base = 1024
			}
			size *= math.Pow(float64(base), float64(index+1))
		}
	}

	return int64(size), nil
}
