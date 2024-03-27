package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func FormatHexToInt(s string) int64 {
	parsed, _ := strconv.ParseInt(strings.Replace(s, "0x", "", -1), 16, 32)
	return parsed
}

func FormatIntToHex(n int64) string {
	return fmt.Sprintf("0x%x", n)
}
