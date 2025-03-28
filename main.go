package csvmanager

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strconv"

	gohaskell "github.com/arewabolu/GoHaskell"
	"golang.org/x/exp/slices"
)

func StringConv[T any](arr []T) []string {
	nwArr := make([]string, 0, len(arr))
	for _, v := range arr {
		nwV := fmt.Sprint(v)
		nwArr = append(nwArr, nwV)
	}
	return nwArr
}

// Write to Csv file after WriteFrame implentation
func (w *WriteFrame) WriteCSV() error {
	wr := csv.NewWriter(w.File)
	defer wr.Flush()
	if len(w.Headers) > 0 {
		err := wr.Write(w.Headers)
		if err != nil {
			return err
		}
	}

	switch {
	case w.Column:
		for i := 0; i < len(w.Arrays[0]); i++ {
			err := wr.Write(extractItems(w.Arrays, i))
			if err != nil {
				return err
			}
		}
	default:
		wr.WriteAll(w.Arrays)
	}

	return nil
}

// Interface should be used to convert rows into  a defined struct
func (r RowList) Interface(value interface{}) error {
	result := reflect.ValueOf(value).Elem() //Elem instead of interface{} to unbox the value/dereference the pointer
	if result.Kind() != reflect.Struct {
		return errors.New("need struct, value is not a struct")
	}
	resultType := result.Type()

	for i := 0; i < result.NumField(); i++ {
		field := resultType.Field(i) //fields type
		fieldVal := result.Field(i)  // fields value
		fieldType := field.Type
		fieldKind := fieldVal.Kind()

		// skip hidden fields
		if field.PkgPath != "" {
			continue
		}

		// make sure we're not using a pointer
		var ptr reflect.Value
		if fieldKind == reflect.Ptr {
			// get the underlying type of the pointer
			fieldType = fieldType.Elem()
			fieldKind = fieldType.Kind()
			// create new pointer to hold the value
			ptr = reflect.New(fieldType)
		}

		if !fieldVal.CanSet() {
			return errors.New("cannot set field " + field.Name)
		}

		var out interface{} // use out to box return value and type
		var err error

		switch fieldKind {
		case reflect.Bool:
			out, err = strconv.ParseBool(r.rowData[i])
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as bool", field.Name, r.rowData[i])
			}
		case reflect.Int:
			val, err := strconv.ParseInt(r.rowData[i], 10, fieldType.Bits())
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as int", field.Name, r.rowData[i])
			}
			out = int(val)
		case reflect.Int8:
			val, err := strconv.ParseInt(r.rowData[i], 10, 8)
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as int8", field.Name, r.rowData[i])
			}
			out = int8(val)
		case reflect.Int16:
			val, err := strconv.ParseInt(r.rowData[i], 10, 16)
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as int16", field.Name, r.rowData[i])
			}
			out = int16(val)
		case reflect.Int32:
			val, err := strconv.ParseInt(r.rowData[i], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as int32", field.Name, r.rowData[i])
			}
			out = int32(val)
		case reflect.Int64:
			if out, err = strconv.ParseInt(r.rowData[i], 10, 64); err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as int64", field.Name, r.rowData[i])
			}
		case reflect.Uint:
			u, err := strconv.ParseUint(r.rowData[i], 10, fieldType.Bits())
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as uint", field.Name, r.rowData[i])
			}
			out = uint(u)
		case reflect.Uint8:
			u, err := strconv.ParseUint(r.rowData[i], 10, 8)
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as uint8", field.Name, r.rowData[i])
			}
			out = uint8(u)
		case reflect.Uint16:
			u, err := strconv.ParseUint(r.rowData[i], 10, 16)
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as uint16", field.Name, r.rowData[i])
			}
			out = uint16(u)
		case reflect.Uint32:
			u, err := strconv.ParseUint(r.rowData[i], 10, 32)
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as uint32", field.Name, r.rowData[i])
			}
			out = uint32(u)
		case reflect.Uint64:
			if out, err = strconv.ParseUint(r.rowData[i], 10, 64); err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as uint64", field.Name, r.rowData[i])
			}
			//result float points 2/3/4/5/6/7?
		case reflect.Float32:
			f, err := strconv.ParseFloat(r.rowData[i], fieldType.Bits())
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as float32", field.Name, r.rowData[i])
			}
			out = float32(f)
		case reflect.Float64:
			out, err = strconv.ParseFloat(r.rowData[i], fieldType.Bits())
			if err != nil {
				return fmt.Errorf("failed to parse field %s, value:%v as float64", field.Name, r.rowData[i])
			}
		case reflect.String:
			out = r.rowData[i]
		default:
			return errors.New("this function doesn't support convertion for the field " + field.Name)
		}

		// if original kind is pointer, save as pointer value
		if fieldVal.Kind() == reflect.Ptr {
			// set value pointer is pointing to
			reflect.Indirect(ptr).Set(reflect.ValueOf(out))
			fieldVal.Set(ptr)
		} else {
			fieldVal.Set(reflect.ValueOf(out))
		}
	}

	return nil
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

	for _, record := range f.cols {
		if colName == record.header {
			return record
		}
	}
	return ColList{}

}

// returns column specified by the parameter
// columns index alsways start from 0.
// ColWithPosition is recommended to be used when header is false
func (f Frame) ColWithPosition(colPos int) ColList {
	return f.cols[colPos]
}

func (f Frame) ColLength() int {
	return len(f.cols)
}

// Check if a header exists in the frame
// returns and error if false otherwise
// it returns the headers position in the slice
func (f Frame) CheckHeader(header string) (int, error) {
	ok := slices.Contains(f.headers, header)
	if !ok {
		return -1, errors.New("header not found")
	}
	pos := slices.Index(f.headers, header)
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
	return f.headers
}

// Reads the csv file in the file path given
func ReadCsv(filePath string, header bool, bufSize ...int) (Frame, error) {
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
	if len(records) == 0 {
		return Frame{}, nil
	}

	var f Frame
	if !header {
		f = Frame{
			cols: genCols(records),
			rws:  genRows(records, 0),
		}
	} else {
		f = Frame{
			headers: records[0],
			cols:    genCols(records),
			rws:     genRows(records, 1),
		}
	}

	return f, nil
}

/**/
// Row returns only one row
func (f Frame) Row(rowLine int) RowList {
	if rowLine < 0 || rowLine > len(f.rws) {
		panic("provided index out of bounds")
	}
	return f.rws[rowLine]
}

// Row returns the rows specified by Range.
//
// Range should be endtered in order: start,counter,end(min,range,max),
// else it returns all rows exluding the Header.
//
// Only 0,2, or 3 values can be specified.
func (f Frame) Rows(Range ...int) []RowList {
	if len(Range) == 0 {
		return f.rws
	}
	if Range[0] < 0 || Range[0] > len(f.rws) {
		panic("provided index out of bounds")
	}
	if Range[1] < 0 || Range[1] > len(f.rws) {
		panic("provided index out of bounds")
	}
	if Range[1] < 0 || Range[2] > len(f.rws) {
		panic("provided index out of bounds")
	}
	if len(Range) == 2 {
		return f.rws[Range[0]:Range[1]]
	}
	if len(Range) == 3 {
		Length := make([]RowList, 0)
		for i := Range[0]; i <= Range[2]; i += Range[1] {
			Length = append(Length, f.rws[i])
		}
		return Length
	}
	if len(Range) > 3 {
		panic("you can't specify more than 3 values")
	}
	return nil
}

// returns the number of rows
func (r RowList) RowsLength() int {
	return len(r.rowData)
}

// returns the number of rows in the read file
func (f Frame) SizeofRows() int {
	return len(f.rws)
}

// String returns an array of string values for given column
func (c ColList) String() []string {
	nwDataString := make([]string, 0, len(c.colData))
	nwDataString = append(nwDataString, c.colData...)
	return nwDataString
}

// String returns an array of string values for given row(s)
func (r RowList) String() []string {
	nwDataString := make([]string, 0, len(r.rowData))
	nwDataString = append(nwDataString, r.rowData...)
	return nwDataString
}

// ReplaceRow is used to edit an existing row in a csv file.
//
// It does not create a new file, it only updates the existing file with the edited row.
func ReplaceRow(filePath string, perm fs.FileMode, pos int, nwData []string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	rdder := csv.NewReader(file)
	rdRecords, err := rdder.ReadAll()
	if err != nil {
		return err
	}
	if pos >= len(rdRecords) || pos < 0 {
		return fmt.Errorf("%s has only %v rows", filePath, len(rdRecords))
	}

	nwRecords := gohaskell.Pop(rdRecords, pos)
	nwRecords = gohaskell.Put(nwRecords, nwData, pos)
	nwFile, err := os.Create(filePath)
	if err != nil {
		return errors.New("could not overrite existing file")
	}

	wr := csv.NewWriter(nwFile)
	defer wr.Flush()

	for _, rec := range nwRecords {
		wr.Write(rec)
	}

	return nil
}

func PrependRow(filePath string, header bool, nwData []string) error {
	data, err := ReadCsv(filePath, header)
	if err != nil {
		return err
	}
	dataRws := data.Rows
	strRows := extractRowsString(dataRws())

	nwFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return errors.New("could not overrite existing file")
	}
	nwRows := make([][]string, 0, len(strRows)+1)
	nwRows = append(nwRows, nwData)
	nwRows = append(nwRows, strRows...)
	wr := WriteFrame{
		Headers: data.headers,
		Row:     true,
		Arrays:  nwRows,
		File:    nwFile,
	}
	err = wr.WriteCSV()
	if err != nil {
		return err
	}
	return nil
}
