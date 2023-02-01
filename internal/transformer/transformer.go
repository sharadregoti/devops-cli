package transformer

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/shared"
	"github.com/tidwall/gjson"
)

func GetSelfContainedResource(dataPaths []string, resource interface{}) ([]interface{}, error) {
	// Extract real data from parent object

	strData, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource data: %w", err)
	}

	containedResources := []interface{}{}
	for _, p := range dataPaths {
		value := gjson.Get(string(strData), p)
		switch v := value.Value().(type) {
		case []interface{}:
			containedResources = append(containedResources, v...)
		case interface{}:
			containedResources = append(containedResources, v)
		case nil:
			continue
		default:
			return nil, fmt.Errorf("failed to extract nested resource: json_path %v data is not of type of array got %v", p, reflect.TypeOf(v))
		}
	}
	return containedResources, nil
}

func GetArgs(resource interface{}, args map[string]interface{}) map[string]interface{} {
	strData, _ := json.Marshal(resource)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to marshal resource data: %w", err)
	// }

	nestArgs := map[string]interface{}{}
	for k, v := range args {
		strV, ok := v.(string)
		if ok {
			gjsonValue := gjson.Get(string(strData), strV)
			if !gjsonValue.Exists() {
				nestArgs[k] = strV
				continue
			}
			nestArgs[k] = gjsonValue.Value()
		} else {
			nestArgs[k] = v
		}
	}

	return nestArgs
}

func GetResourceInTableFormat(logger hclog.Logger, t *shared.ResourceTransfomer, resources []interface{}) ([]*model.TableRow, []map[string]interface{}, error) {
	gjson.AddModifier("age", func(json, arg string) string {
		return getAge(json[1 : len(json)-1])
	})

	gjson.AddModifier("pick", func(json, arg string) string {
		return getPick(json, arg)
	})

	tableResult := make([]*model.TableRow, 0)

	headerRow := new(model.TableRow)
	for _, o := range t.Operations {
		headerRow.Data = append(headerRow.Data, strings.ToUpper(o.Name))
	}
	headerRow.Color = tcell.ColorYellow

	tableResult = append(tableResult, headerRow)
	// copy(tableResult[len(tableResult)-1], headerRow)

	nestArgs := []map[string]interface{}{}
	for _, resource := range resources {
		res, nestArg, err := TransformResource(logger, t, resource)
		if err != nil {
			return nil, nil, err
		}

		tableResult = append(tableResult, res)
		// copy(tableResult[len(tableResult)-1], res)

		if t.Nesting.IsNested {
			nestArgs = append(nestArgs, nestArg)
		}
	}

	return tableResult, nestArgs, nil
}

func TransformResource(logger hclog.Logger, t *shared.ResourceTransfomer, data interface{}) (*model.TableRow, map[string]interface{}, error) {
	resultRow := new(model.TableRow)
	dataRow := make([]string, 0)

	strData, err := json.Marshal(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal resource data: %w", err)
	}

	nestArgs := map[string]interface{}{}
	if t.Nesting.IsNested {
		for k, v := range t.Nesting.Args {
			strV, ok := v.(string)
			if ok {
				gjsonValue := gjson.Get(string(strData), strV)
				nestArgs[k] = gjsonValue.Value()
			} else {
				nestArgs[k] = v
			}
		}
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
			program, err := expr.Compile(c, expr.Env(data))
			if err != nil {
				logger.Debug("skipping style condition as failed to compile style condition", err)
				continue
			}

			output, err := expr.Run(program, data)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to evaluate style condition: %v", err)
			}

			result, ok := output.(bool)
			if !ok {
				return nil, nil, fmt.Errorf("condition evaluation resulted into an unknown type %v, expected boolean", result)
			}

			if !result {
				break
			}
			switch s.RowBackgroundColor {
			case "white":
				resultRow.Color = tcell.ColorWhite
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

	return resultRow, nestArgs, nil
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
		result = "NA"
	}

	return fmt.Sprintf(`"%s"`, result)
}

func getPick(ts string, args string) string {
	result := []string{}

	var temp []map[string]interface{}
	_ = json.Unmarshal([]byte(ts), &temp)

	for _, myMap := range temp {
		port := []string{}
		for _, field := range strings.Split(args, ",") {
			rs, ok := myMap[field]
			if !ok || rs == "" {
				continue
			}
			port = append(port, fmt.Sprintf("%v", rs))
		}
		result = append(result, strings.Join(port, "->"))
	}

	return fmt.Sprintf(`"%s"`, strings.Join(result, " "))
}
