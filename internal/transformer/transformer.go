package transformer

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/gdamore/tcell/v2"
	"github.com/sharadregoti/devops/model"
	"github.com/tidwall/gjson"
)

func GetResourceInTableFormat(t *model.ResourceTransfomer, resources []interface{}) ([]*model.TableRow, error) {
	gjson.AddModifier("age", func(json, arg string) string {
		return getAge(json[1 : len(json)-1])
	})

	tableResult := make([]*model.TableRow, 0)

	headerRow := new(model.TableRow)
	for _, o := range t.Operations {
		headerRow.Data = append(headerRow.Data, strings.ToUpper(o.Name))
	}
	headerRow.Color = tcell.ColorYellow

	tableResult = append(tableResult, headerRow)
	// copy(tableResult[len(tableResult)-1], headerRow)

	for _, v := range resources {
		res, err := TransformResource(t, v)
		if err != nil {
			return nil, err
		}

		tableResult = append(tableResult, res)
		// copy(tableResult[len(tableResult)-1], res)
	}

	return tableResult, nil
}

func TransformResource(t *model.ResourceTransfomer, data interface{}) (*model.TableRow, error) {
	resultRow := new(model.TableRow)
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
				// Evaluate gjson expression
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
	resultRow.Data = dataRow
	resultRow.Color = tcell.ColorWhite

	for _, s := range t.Styles {
		for _, c := range s.Conditions {
			// Evaluate the condition
			// fmt.Println("here")
			output, err := expr.Eval(c, data)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate style condition: %v", err)
			}

			result, ok := output.(bool)
			if !ok {
				return nil, fmt.Errorf("condition evaluation resulted into an unknown type %v, expected boolean", result)
			}

			if !result {
				break
			}
			switch s.RowBackgroundColor {
			case "red":
				resultRow.Color = tcell.ColorRed
			case "yellow":
				resultRow.Color = tcell.ColorYellow
			case "blue":
				resultRow.Color = tcell.ColorBlue
			case "orange":
				resultRow.Color = tcell.ColorOrange
			case "green":
				resultRow.Color = tcell.ColorGreen
			case "aqua":
				resultRow.Color = tcell.ColorAqua
			}
		}
	}

	return resultRow, nil
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
