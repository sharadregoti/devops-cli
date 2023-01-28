package utils

import (
	"fmt"
	"strings"
)

func GetTableTitle(title string, total int) string {
	return fmt.Sprintf("%s (%d)", title, total)
}

func ParseTableTitle(title string) (string, string) {
	arr := strings.Split(title, " ")
	return strings.ToLower(arr[0]), strings.TrimPrefix(strings.TrimSuffix(arr[1], ")"), "(")
}
