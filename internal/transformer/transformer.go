package transformer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sharadregoti/devops/model"
	"github.com/tidwall/gjson"
)

func GetResourceInTableFormat(t *model.ResourceTransfomer, resources []interface{}) [][]string {
	tableResult := make([][]string, 0)

	headerRows := make([]string, 0)
	for _, o := range t.Operations {
		headerRows = append(headerRows, strings.ToUpper(o.Name))
	}

	tableResult = append(tableResult, make([]string, len(headerRows)))
	copy(tableResult[len(tableResult)-1], headerRows)

	for _, v := range resources {
		res := TransformResource(t, v)
		tableResult = append(tableResult, make([]string, len(res)))
		copy(tableResult[len(tableResult)-1], res)
	}

	return tableResult
}

func TransformResource(t *model.ResourceTransfomer, data interface{}) []string {
	dataRow := make([]string, 0)

	strData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Get column names in title case
	for _, o := range t.Operations {
		var pathExecResults []interface{} = make([]interface{}, 0)
		for _, j := range o.JSONPaths {
			if j.Path != "" {
				gjson.AddModifier("age", func(json, arg string) string {
					return getAge(strings.ReplaceAll(json, "\"", ""))
				})
				value := gjson.Get(string(strData), j.Path)
				pathExecResults = append(pathExecResults, value)
			}
		}

		if o.OutputFormat == "" {
			o.OutputFormat = "%v"
		}

		dataRow = append(dataRow, fmt.Sprintf(o.OutputFormat, pathExecResults...))
	}

	return dataRow
}

func getAge(ts string) string {
	result := "nil"
	t, err := time.Parse(ts, time.RFC3339)
	if err != nil {
		fmt.Println(ts, err)
		return "nil"
	}

	// calculate the difference between the current time and the timestamp
	difference := time.Since(t)

	// convert the difference to a human-readable string
	if seconds := difference.Seconds(); seconds < 60 {
		result = fmt.Sprintf("%ds", int(seconds))
	} else if minutes := difference.Minutes(); minutes < 60 {
		result = fmt.Sprintf("%dm", int(minutes))
	} else if hours := difference.Hours(); hours < 24 {
		result = fmt.Sprintf("%dh", int(hours))
	} else {
		result = fmt.Sprintf("%dd", int(difference.Hours()/24))
	}

	return result
}
