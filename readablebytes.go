package readablebytes

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

type unitMap map[string]int64

var (
	decMap    = unitMap{"k": KB, "m": MB, "g": GB, "t": TB, "p": PB}
	binMap    = unitMap{"k": KiB, "m": MiB, "g": GiB, "t": TiB, "p": PiB}
	sizeRegex = regexp.MustCompile(`(?i)^(\d+(\.\d+)*) ?([kmgtp])?(i?)b?$`)
)

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
	matches := sizeRegex.FindStringSubmatch(sizeStr)
	if len(matches) != 5 {
		return -1, fmt.Errorf("invalid size: '%s'", sizeStr)
	}

	size, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return -1, err
	}

	unitPrefix := strings.ToLower(matches[3])
	unitsMap := decMap
	if matches[4] != "" {
		unitsMap = binMap
	}
	if mul, ok := unitsMap[unitPrefix]; ok {
		size *= float64(mul)
	}

	return int64(size), nil
}
