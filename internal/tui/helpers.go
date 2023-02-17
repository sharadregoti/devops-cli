package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// A function which return string for use.plugin
const (
	settingPlugin         = "@use.plugin"
	settingAuthentication = "@use.auth"
)

func getPluginSetting(name string) string {
	return fmt.Sprintf("%s.%s", settingPlugin, name)
}

// get plugin name from setting name using split function
func parsePluginSetting(setting string) string {
	return strings.Split(setting, ".")[2]
}

// Returns identifyingName, name string
func parseAuthenticationSetting(setting string) (string, string) {
	arr := strings.Split(setting, ".")
	return arr[2], arr[3]
}

func getAuthenticationSetting(identifyingName, name string) string {
	return fmt.Sprintf("%s.%s.%s", settingAuthentication, identifyingName, name)
}

// ChatGPT wrote this code, don't know how it works
func getNiceFormat(data map[string]string) string {
	// data := map[string]string{
	// 	"key1":      "value1",
	// 	"key3":      "value3",
	// 	"key2":      "value2",
	// 	"key5":      "value5",
	// 	"key4":      "value4",
	// 	"key6":      "value6",
	// 	"key7":      "value7",
	// 	"key8":      "value8",
	// 	"key9":      "value9",
	// 	"key10":     "value10",
	// 	"key11":     "value11",
	// 	"key9434":   "value9",
	// 	"key1012":   "value10",
	// 	"key101266": "value10",
	// }

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var result string
	var columns [][]string
	var maxLengths []int
	for i, k := range keys {
		v := data[k]
		if i%5 == 0 {
			columns = append(columns, []string{})
			maxLengths = append(maxLengths, len(k+": "+v))
		} else {
			maxLengths[len(maxLengths)-1] = max(maxLengths[len(maxLengths)-1], len(k+": "+v))
		}
		columns[len(columns)-1] = append(columns[len(columns)-1], k+": "+v)
	}

	transposed := transpose(columns)

	for i, row := range transposed {
		for j, cell := range row {
			if j > 0 {
				result += "  "
			}
			if i == 0 {
				result += fmt.Sprintf("%-"+strconv.Itoa(maxLengths[j])+"s", cell)
			} else {
				parts := strings.SplitN(cell, ": ", 2)
				if len(parts) == 2 {
					key, value := parts[0], parts[1]
					if j < len(maxLengths)-1 {
						value = fmt.Sprintf("%-"+strconv.Itoa(maxLengths[j]-len(key)-2)+"s", value)
					}
					result += fmt.Sprintf("%s: %s", key, value)
				}
			}
		}
		result += "\n"
	}

	return result
}

func transpose(matrix [][]string) [][]string {
	numRows := len(matrix)
	numCols := len(matrix[0])
	result := make([][]string, numCols)
	for i := 0; i < numCols; i++ {
		result[i] = make([]string, numRows)
		for j := 0; j < numRows; j++ {
			if i < len(matrix[j]) {
				result[i][j] = matrix[j][i]
			} else {
				result[i][j] = ""
			}
		}
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
