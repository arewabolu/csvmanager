package csvmanager

import "os"

type colList struct {
	header  string
	colData []string
}

type rowList struct {
	rowData []string
}

type WriteFrame struct {
	//list all csv headers here
	Headers []string
	// Set Column to true if []string is a list of columns.
	//
	// If not set Row automatically defaults to true.
	Column bool
	// Set Row to true if []string is a list of rows.
	//
	//Note: Default is true
	Row bool
	//columns should contain a list of all columns
	//which must be properly formatted
	Arrays [][]string
	//File should be a file with right permissions to be written to
	File *os.File
}

type Frame struct {
	headers []string
	cols    []colList
	rws     []rowList
}

type Types interface {
	Float() []float64
	Int() []int
	Bool() []bool
	String() []string
	Interface(v interface{}) error
}
