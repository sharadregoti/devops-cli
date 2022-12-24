package transformer

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/sharadregoti/devops/model"
	"github.com/tidwall/gjson"
)

func GetResourceInTableFormat(t *model.ResourceTransfomer, resources []interface{}) ([][]string, error) {
	gjson.AddModifier("age", func(json, arg string) string {
		return getAge(json[1 : len(json)-1])
	})

	tableResult := make([][]string, 0)

	headerRows := make([]string, 0)
	for _, o := range t.Operations {
		headerRows = append(headerRows, strings.ToUpper(o.Name))
	}

	tableResult = append(tableResult, make([]string, len(headerRows)))
	copy(tableResult[len(tableResult)-1], headerRows)

	for _, v := range resources {
		res, err := TransformResource(t, v)
		if err != nil {
			return nil, err
		}

		tableResult = append(tableResult, make([]string, len(res)))
		copy(tableResult[len(tableResult)-1], res)
	}

	return tableResult, nil
}

func TransformResource(t *model.ResourceTransfomer, data interface{}) ([]string, error) {
	dataRow := make([]string, 0)

	strData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource data: %w", err)
	}

	// Get column names in title case
	for _, o := range t.Operations {
		var pathExecResults []interface{} = make([]interface{}, 0)
		for _, j := range o.JSONPaths {
			if j.Path != "" {
				value := gjson.Get(string(strData), j.Path).String()
				if value == "" {
					value = "NA"
				}
				pathExecResults = append(pathExecResults, value)
			}
		}

		if o.OutputFormat == "" {
			o.OutputFormat = "%v"
		}

		dataRow = append(dataRow, fmt.Sprintf(o.OutputFormat, pathExecResults...))
	}

	return dataRow, nil
}

func getAge(ts string) string {
	// result = carbon.ParseByLayout(ts, time.RFC3339).Age()
	// _, err := time.Parse(ts, time.RFC3339)
	// Parse the time string using the time.Parse function
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		fmt.Println(err)
		return "invalid"
	}

	// Calculate the number of seconds, minutes, hours, and days since the given time
	seconds := time.Since(t).Seconds()
	minutes := time.Since(t).Minutes()
	hours := time.Since(t).Hours()
	days := time.Since(t).Hours() / 24

	// Round the values to print only whole numbers
	seconds = math.Round(seconds)
	minutes = math.Round(minutes)
	hours = math.Round(hours)
	days = math.Round(days)

	// Print the values if they are greater than 0
	result := ""
	if days > 0 {
		result = fmt.Sprintf("%.0fd", days)
	} else if hours > 0 {
		result = fmt.Sprintf("%.0fh", hours)
	} else if minutes > 0 {
		result = fmt.Sprintf("%.0fm", minutes)
	} else if seconds > 0 {
		result = fmt.Sprintf("%.0fs", seconds)
	} else {
		result = "Yo"
	}

	return fmt.Sprintf(`"%s"`, result)
}
