package csvmanager

import (
	"golang.org/x/exp/slices"
)

func newCol(header string, records [][]string) []string {
	colIndex := slices.Index(records[0], header)

	colData := make([]string, 0, len(records)-1)
	for _, record := range records[1:] {
		if record[colIndex] == "-" || record[colIndex] == "" || record[colIndex] == " " || record[colIndex] == "Nan" {
			continue
		}
		colData = append(colData, record[colIndex])
	}
	return colData
}

func genCols(records [][]string) []colList {
	var colist []colList
	for _, header := range records[0] {
		colData := newCol(header, records)
		colist = append(colist, colList{header: header, colData: colData})
	}
	return colist
}

func genRows(records [][]string) []rowList {
	var rowlist []rowList
	for _, record := range records[1:] {
		rowlist = append(rowlist, rowList{rowData: record})
	}
	return rowlist
}
