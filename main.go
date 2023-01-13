package csvmanager

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type ColList struct {
	header  string
	colData []string
}
type RowList struct {
	rowData []string
}

type Frame struct {
	Cols []ColList
	Rws  []RowList
}

type Types interface {
	Float() []float64
	Int() []int
	Bool() []bool
	Interface(v interface{})
}

// Bool returns an array of boolean values
func (c ColList) Bool() []bool {
	nwDataBool := make([]bool, 0, len(c.colData))
	for _, v := range c.colData {
		val, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("%v cannot convert to float64: %v", v, err))
		}
		nwDataBool = append(nwDataBool, val)
	}
	return nwDataBool
}

// Bool returns an array of boolean values
func (r RowList) Bool() []bool {
	nwDataBool := make([]bool, 0, len(r.rowData))
	for _, v := range r.rowData {
		val, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("%v cannot convert to float64: %v", v, err))
		}
		nwDataBool = append(nwDataBool, val)
	}
	return nwDataBool
}

func (f Frame) Col(colName string) ColList {

	for _, record := range f.Cols {
		if colName == record.header {
			return record
		}
	}
	return ColList{}
}

func (f Frame) ColLength() int {
	return len(f.Cols)
}

// FLoat always returns a float64 type
func (c ColList) Float() []float64 {
	nwDataFLoat := make([]float64, 0, len(c.colData))
	for _, v := range c.colData {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("%v cannot convert to float64: %v", v, err))
		}
		nwDataFLoat = append(nwDataFLoat, val)
	}
	return nwDataFLoat
}

// FLoat always returns a float64 type
func (r RowList) Float() []float64 {
	nwDataFLoat := make([]float64, 0, len(r.rowData))
	for _, v := range r.rowData {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("%v cannot convert to float64: %v", v, err))
		}
		nwDataFLoat = append(nwDataFLoat, val)
	}
	return nwDataFLoat
}

// Int always returns an int type values
//
// Convert to any other type of integer
func (c ColList) Int() []int {
	nwDataInt := make([]int, 0, len(c.colData))
	for _, v := range c.colData {
		val, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("%v cannot convert to float64: %v", v, err))
		}
		nwDataInt = append(nwDataInt, val)
	}
	return nwDataInt
}

// Int always returns an int type values
//
// Convert to any other type of integer
func (r RowList) Int() []int {
	nwDataInt := make([]int, 0, len(r.rowData))
	for _, v := range r.rowData {
		val, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("%v cannot convert to float64: %v", v, err))
		}
		nwDataInt = append(nwDataInt, val)
	}
	return nwDataInt
}

func ReadCsv(filePath string, bufSize ...int) (Frame, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Frame{}, err
	}
	defer file.Close()

	var rd *csv.Reader
	if len(bufSize) > 0 {
		rdder := bufio.NewReaderSize(file, bufSize[0])
		rd = csv.NewReader(rdder)
	} else {
		rd = csv.NewReader(file)
	}

	records, err := rd.ReadAll()
	if err != nil {
		return Frame{}, err
	}
	f := Frame{
		Cols: genCols(records),
		Rws:  genRows(records),
	}
	return f, nil
}

/**/
// Row returns only one row
func (f Frame) Row(rowLine int) RowList {
	if rowLine < 0 || rowLine > len(f.Rws) {
		panic("provided index out of bounds")
	}
	return f.Rws[rowLine]
}

// Row returns the rows specified by Range.
//
// Range should be endtered in order: start,counter,end(min,range,max),
// else it returns all rows exluding the Header.
//
// Only 0,2, or 3 values can be specified.
func (f Frame) Rows(Range ...int) []RowList {
	if len(Range) == 0 {
		return f.Rws
	}
	if Range[0] < 0 || Range[0] > len(f.Rws) {
		panic("provided index out of bounds")
	}
	if Range[1] < 0 || Range[1] > len(f.Rws) {
		panic("provided index out of bounds")
	}
	if Range[1] < 0 || Range[2] > len(f.Rws) {
		panic("provided index out of bounds")
	}
	if len(Range) == 2 {
		return f.Rws[Range[0]:Range[1]]
	}
	if len(Range) == 3 {
		Length := make([]RowList, 0)
		for i := Range[0]; i <= Range[2]; i += Range[1] {
			Length = append(Length, f.Rws[i])
		}
		return Length
	}
	if len(Range) > 3 {
		panic("you can't specify more than 3 values")
	}
	return nil
}

func (r RowList) RowsLength() int {
	return len(r.rowData)
}

func (f Frame) SizeofRows() int {
	return len(f.Rws)
}
