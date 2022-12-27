package csvmanager

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type Frame struct {
	Headers []string
	Data    [][]string
}

func ReadCsv(filePath string, bufSize ...int) (Frame, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Frame{}, err
	}
	defer file.Close()

	rd := new(csv.Reader)

	if len(bufSize) > 0 {
		rdder := bufio.NewReaderSize(file, bufSize[0])
		rd = csv.NewReader(rdder)
	} else {
		rd = csv.NewReader(file)
	}

	records, err := rd.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(records[0])

	f := Frame{
		Headers: records[0],
		Data:    records[1:],
	}

	return f, nil
}

func (f Frame) Column(colName string) []string {
	if !slices.Contains(f.Headers, colName) {
		panic("column not found in records")
	}
	colData := make([]string, 0, len(f.Data))
	colIndex := slices.Index(f.Headers, colName)
	if colIndex == -1 {
		panic("column not found in records")
	}

	for _, record := range f.Data {
		colData = append(colData, record[colIndex])
	}
	return colData
}

// Row returns the rows specified by Range.
//
// Range should be endtered in order: start,counter,end,
// else it returns All rows exluding the Header.
//
// Only 0,2, or 3 values can be specified.
func (f Frame) Rows(Range ...int) [][]string {

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

func (f Frame) ColumnLenght() int {
	return len(f.Data)
}