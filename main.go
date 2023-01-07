package csvmanager

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Types interface {
	Float() []float64
	Int() []int
	Bool() []bool
}

type RowList struct {
	rowData []string
}

type ColList struct {
	header  string
	colData []string
}

type Frame struct {
	Cols []ColList
	Rws  []RowList
}

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

//func (f *Frame) Int() []int   {}
//func (f *Frame) Bool() []bool {}

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
	}

	return f, nil

}

func (f Frame) Col(colName string) ColList {
	//if colIndex == -1 {
	//	panic("column not found in records")
	//}

	//colData := make([]string, 0, len(c.colData))
	for _, record := range f.Cols {
		if colName == record.header {
			return record
		}
		//colData = append(colData, record[colIndex])
	}
	return ColList{}
}

// Row returns the rows specified by Range.
//
// Range should be endtered in order: start,counter,end,
// else it returns All rows exluding the Header.
//
// Only 0,2, or 3 values can be specified.
/*func (f Frame) Rows(Range ...int) [][]string {
	if len(Range) == 0 {
		return f.Data
	}
	if Range[0] < 0 || Range[0] > len(f.Data) {
		panic("provided index out of bounds")
	}
	if Range[1] < 0 || Range[1] > len(f.Data) {
		panic("provided index out of bounds")
	}
	if Range[1] < 0 || Range[2] > len(f.Data) {
		panic("provided index out of bounds")
	}
	if len(Range) == 2 {
		return f.Data[Range[0]:Range[1]]
	}
	if len(Range) == 3 {
		Length := make([][]string, 0)
		for i := Range[0]; i < len(f.Data); i += Range[1] {
			Length = append(Length, f.Data[i])
		}
		return Length
	}
	if len(Range) > 3 {
		panic("you can't specify more than 3 values")
	}
	return nil
}

// Row returns only one row
func (f Frame) Row(rowLine int) []string {
	if rowLine < 0 || rowLine > len(f.Data) {
		panic("provided index out of bounds")
	}
	return f.Data[rowLine]
}

func (f Frame) RowLength() int {
	return len(f.Headers)
}
*/
func (f Frame) ColumnLenght() int {
	return len(f.Cols)
}
