package csvmanager

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	gohaskell "github.com/arewabolu/GoHaskell"
	"golang.org/x/exp/slices"
)

type ColList struct {
	header  string
	colData []string
}
type RowList struct {
	rowData []string
}

type Frame struct {
	Headers []string
	Cols    []ColList
	Rws     []RowList
	Err     error
}

type Types interface {
	Float() []float64
	Int() []int
	Bool() []bool
	Interface(v interface{})
}

type Error interface {
	Err() error
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

// Check if a header exists in the frame
// returns and error if false otherwise
// it returns the headers position in the slice
func (f Frame) CheckHeader(header string) (int, error) {
	pos, ok := slices.BinarySearch(f.Headers, header)
	if !ok {
		return -1, errors.New("header not found")
	}
	return pos, nil
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

// List all headers in the given frame
func (f Frame) ListHeaders() []string {
	return f.Headers
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
		Headers: records[0],
		Cols:    genCols(records),
		Rws:     genRows(records),
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

// returns the number of rows in the read file
func (f Frame) SizeofRows() int {
	return len(f.Rws)
}

// Create a new csv file.
func WriteNewCSV(filePath string, rowData []string) Frame {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0700)
	if err != nil {
		return Frame{Err: err}
	}

	defer file.Close()

	wr := csv.NewWriter(file)
	defer wr.Flush()

	wr.Write(rowData)
	//	for _, i := range entry {
	//		fV := strconv.FormatFloat(i, 'f', 2, 64)
	//		err := wr.Write([]string{fV})
	//		if err != nil {
	//			panic(err)
	//}
	//	}
	return Frame{Err: nil}
}

// ReplaceRow is used to edit an existing row in a csv file.
//
// It does not create a new file, it only updates the existing file with the edited row.
func ReplaceRow(filePath string, pos int, nwData []string) Frame {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0700)
	if err != nil {
		return Frame{Err: err}
	}
	defer file.Close()
	rdder := csv.NewReader(file)
	rdRecords, err := rdder.ReadAll()
	if err != nil {
		return Frame{Err: err}
	}
	if pos >= len(rdRecords) || pos < 0 {
		return Frame{Err: errors.New("replacing nonexistent row is not supported")}
	}

	nwRecords := gohaskell.Pop(rdRecords, pos)
	nwRecords = gohaskell.Put(nwRecords, nwData, pos)
	nwFile, err := os.Create(filePath)
	if err != nil {
		return Frame{Err: errors.New("could not overrite existing file")}
	}

	wr := csv.NewWriter(nwFile)
	defer wr.Flush()

	for _, rec := range nwRecords {
		wr.Write(rec)
	}

	return Frame{Err: nil}
}
