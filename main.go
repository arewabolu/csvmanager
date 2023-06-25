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
// The position of the item on the row could also be set to a particular field
// using struct tags
func (r rowList) Interface(v interface{}) error {
	result := reflect.ValueOf(v).Elem() //
	if result.Kind() != reflect.Struct {
		return errors.New("need struct, v is not a struct")
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

		// get sructs field tpe
		// get query parameter name
		key := field.Tag.Get("position")
		pos, err := strconv.Atoi(key)
		if err != nil || key == "" {
			pos = i
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

		var out interface{}

		switch fieldKind {
		case reflect.Bool:
			out, err = strconv.ParseBool(r.rowData[pos])
			if err != nil {
				return errors.New("failed to parse bool" + field.Name)
			}
		case reflect.Int:
			i, err := strconv.ParseInt(r.rowData[pos], 10, fieldType.Bits())
			if err != nil {
				return errors.New("failed to parse int" + field.Name)
			}
			out = int(i)
		case reflect.Int8:
			i, err := strconv.ParseInt(r.rowData[pos], 10, 8)
			if err != nil {
				return errors.New("failed to parse int8" + field.Name)
			}
			out = int8(i)
		case reflect.Int16:
			i, err := strconv.ParseInt(r.rowData[pos], 10, 16)
			if err != nil {
				return errors.New("failed to parse int16" + field.Name)
			}
			out = int16(i)
		case reflect.Int32:
			i, err := strconv.ParseInt(r.rowData[pos], 10, 32)
			if err != nil {
				return errors.New("failed to parse int32" + field.Name)
			}
			out = int32(i)
		case reflect.Int64:
			if out, err = strconv.ParseInt(r.rowData[pos], 10, 64); err != nil {
				return errors.New("failed to parse int64" + field.Name)
			}
		case reflect.Uint:
			u, err := strconv.ParseUint(r.rowData[pos], 10, fieldType.Bits())
			if err != nil {
				return errors.New("failed to parse uint" + field.Name)
			}
			out = uint(u)
		case reflect.Uint8:
			u, err := strconv.ParseUint(r.rowData[pos], 10, 8)
			if err != nil {
				return errors.New("failed to parse uint8" + field.Name)
			}
			out = uint8(u)
		case reflect.Uint16:
			u, err := strconv.ParseUint(r.rowData[pos], 10, 16)
			if err != nil {
				return errors.New("failed to parse uint16" + field.Name)
			}
			out = uint16(u)
		case reflect.Uint32:
			u, err := strconv.ParseUint(r.rowData[pos], 10, 32)
			if err != nil {
				return errors.New("failed to parse uint32" + field.Name)
			}
			out = uint32(u)
		case reflect.Uint64:
			if out, err = strconv.ParseUint(r.rowData[pos], 10, 64); err != nil {
				return errors.New("failed to parse uint64" + field.Name)
			}
			//result float points 2/3/4/5/6/7?
		case reflect.Float32:
			f, err := strconv.ParseFloat(r.rowData[pos], fieldType.Bits())
			if err != nil {
				return errors.New("failed to parse float32" + field.Name)
			}
			out = float32(f)
		case reflect.Float64:
			out, err = strconv.ParseFloat(r.rowData[pos], fieldType.Bits())
			if err != nil {
				return errors.New("failed to parse float64" + field.Name)
			}
		case reflect.String:
			out = r.rowData[pos]
		default:
			return errors.New("this function doesn't support convertion" + "for the field '%s'" + field.Name)
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
func (c colList) Bool() []bool {
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
func (r rowList) Bool() []bool {
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

func (f Frame) Col(colName string) colList {

	for _, record := range f.cols {
		if colName == record.header {
			return record
		}
	}
	return colList{}
}

// returns column specified by the parameter
// columns index alsways start from 0.
// ColWithPosition is recommended to be used when header is false
func (f Frame) ColWithPosition(colPos int) colList {
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
func (c colList) Float() []float64 {
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
func (r rowList) Float() []float64 {
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
func (c colList) Int() []int {
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
func (r rowList) Int() []int {
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

func ReadCsv(filePath string, perm fs.FileMode, header bool, bufSize ...int) (Frame, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, perm)
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
func (f Frame) Row(rowLine int) rowList {
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
func (f Frame) Rows(Range ...int) []rowList {
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
		Length := make([]rowList, 0)
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
func (r rowList) RowsLength() int {
	return len(r.rowData)
}

// returns the number of rows in the read file
func (f Frame) SizeofRows() int {
	return len(f.rws)
}

// String returns an array of string values for given column
func (c colList) String() []string {
	nwDataString := make([]string, 0, len(c.colData))
	nwDataString = append(nwDataString, c.colData...)
	return nwDataString
}

// String returns an array of string values for given row(s)
func (r rowList) String() []string {
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

func PrependRow(filePath string, perm fs.FileMode, header bool, nwData []string) error {
	data, err := ReadCsv(filePath, perm, header)
	if err != nil {
		return err
	}
	dataRws := data.Rows
	strRows := extractRowsString(dataRws())

	nwFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, perm)
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
