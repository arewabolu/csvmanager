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

func genCols(records [][]string) []ColList {
	if len(records) == 0 {
		return []ColList{}
	}
	var colist []ColList
	for _, header := range records[0] {
		colData := newCol(header, records)
		colist = append(colist, ColList{header: header, colData: colData})
	}
	return colist
}

func genRows(records [][]string, start int) []RowList {
	if len(records) == 0 {
		return []RowList{}
	}
	var rowList []RowList
	for _, record := range records[start:] {
		rowList = append(rowList, RowList{rowData: record})
	}
	return rowList
}

func extractItems(slice [][]string, n int) []string {
	var result []string

	for _, innerSlice := range slice {
		if n < len(innerSlice) {
			result = append(result, innerSlice[n])
		} else {
			result = append(result, "")
		}
	}

	return result
}

func extractRowsString(rows []RowList) [][]string {
	data := make([][]string, 0, len(rows))
	for i := 0; i < len(rows); i++ {
		rows := rows[i].String()
		data = append(data, rows)
	}
	return data
}
