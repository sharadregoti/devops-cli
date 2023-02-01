package model

import "github.com/gdamore/tcell/v2"

type TableRow struct {
	Data  []string    `json:"data" yaml:"data"`
	Color tcell.Color `json:"color" yaml:"color"`
}