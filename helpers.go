package csvmanager

import (
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
)

type rTesting struct {
	One int    `json:"one"`
	Two string `json:"two"`
}

func newCol(header string, records [][]string) []string {
	colIndex := slices.Index(records[0], header)

	colData := make([]string, 0, len(records)-1)
	for _, record := range records[1:] {
		colData = append(colData, record[colIndex])
	}
	return colData
}

func genCols(records [][]string) []ColList {
	var colist []ColList
	for _, header := range records[0] {
		colData := newCol(header, records)
		colist = append(colist, ColList{header: header, colData: colData})
	}
	return colist
}

func genRows(records [][]string) []RowList {
	var rowlist []RowList
	for _, record := range records[1:] {
		rowlist = append(rowlist, RowList{rowData: record})
	}
	return rowlist
}

func jsonNew(v interface{}) {
	vVal := reflect.ValueOf(v).Elem()
	vValType := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		field := vValType.Field(i)
		fieldType := field.Type
		key, _ := field.Tag.Lookup(field.Name)
		if key == "" {
			key = strings.ToLower(field.Name)
		}
		switch fieldType.Kind() {
		case reflect.Int:
		}
	}
}
