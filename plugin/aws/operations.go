package aws

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sharadregoti/devops/common"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sharadregoti/devops/shared"
)

// func (d *AWS) listResources(args shared.GetResourcesArgs) ([]interface{}, error) {
// Set the command to execute
// command := "cloudquery"
// arguments := []string{"sync", "--cq-dir", "", "--no-migrate", "./"}

// // Execute the command
// _, err := exec.Command(command, arguments...).Output()
// if err != nil {
// 	d.common.Error(a.logger,fmt.Sprintf("failed to get describe output, got %v", err))
// 	return nil, err
// }

// 	return make([]interface{}, 0), nil
// }

func (a *AWS) getResourceTypes() ([]string, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "../../plugin/aws/resource_config/db.sql")
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}
	defer db.Close()

	// Get all the table names from the database
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and print the table names
	var result []string
	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			common.Error(a.logger, err.Error())
			return nil, err
		}

		result = append(result, strings.TrimPrefix(tableName, "aws_"))
	}
	err = rows.Err()
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}
	return result, nil
}

func (a *AWS) listResources(args shared.GetResourcesArgs) ([]interface{}, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "../../plugin/aws/resource_config/db.sql")
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}
	defer db.Close()

	// Get the column names of the "users" table
	// aws_ec2_instances
	rows, err := db.Query(fmt.Sprintf("select * from aws_%s", args.ResourceType))
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}
	defer rows.Close()

	// Get the column names from the rows
	columns, err := rows.Columns()
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}

	// Create a slice to hold the values for each row
	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	result := make([]interface{}, 0)
	// Iterate over the rows and scan the data into the values slice
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			common.Error(a.logger, err.Error())
			return nil, err
		}

		// Create a map with the column names as keys and the values as values
		rowData := make(map[string]interface{})
		for i, column := range columns {
			if i < 4 {
				continue
			}
			rowData[column] = *values[i].(*interface{})
		}

		// Print the map for this row
		result = append(result, rowData)
	}
	err = rows.Err()
	if err != nil {
		common.Error(a.logger, err.Error())
		return nil, err
	}

	return result, nil
}
