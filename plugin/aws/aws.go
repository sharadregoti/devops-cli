package plugin

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

type Transformations struct {
	Operations []*Operations `json:"operations"`
}
type Operations struct {
	SrcColumn string `json:"src_column" yaml:"src_column"`
	Name      string `json:"name" yaml:"name"`
	JSONPath  string `json:"json_path" yaml:"json_path"`
	// internal usage
	index int
}

func getEC2() [][]string {
	rows := readCsvFile("plugin/cq_csv_output/aws_ec2_instances.csv")

	f, err := os.ReadFile("plugin/ec2_instance.yaml")
	fmt.Println(err)

	t := new(Transformations)
	_ = yaml.Unmarshal(f, t)

	indexMapper := map[int]int{} // col : t[index]

	filterdData := make([][]string, len(rows))

	for r, cols := range rows {
		for c, col := range cols {

			// Get column namees from cloudQuery
			if r < 1 {
				for o, op := range t.Operations {
					if op.SrcColumn == col {
						indexMapper[c] = o
						filterdData[r] = append(filterdData[r], op.Name)
						break
					}
				}
				continue
			}

			opIndex, ok := indexMapper[c]
			if !ok {
				continue
			}

			op := t.Operations[opIndex]

			tableValue := col

			if op.JSONPath != "" {
				value := gjson.Get(col, op.JSONPath)
				tableValue = value.String()
			}

			filterdData[r] = append(filterdData[r], tableValue)
		}
	}

	return filterdData
}
