package csvmanager

import (
	"os"
	"testing"
)

func TestWriteFrame(t *testing.T) {
	file, _ := os.OpenFile("testwritefile.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	x1 := []string{"ars", "liv", "mu"}
	x2 := []string{"mci", "che", "la"}
	w := &WriteFrame{
		Headers: []string{"Home", "Away"},
		Row:     true,
		Arrays:  [][]string{x1, x2},
		File:    file,
	}
	w.WriteCSV()
	file2, _ := os.OpenFile("test5file.csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	w2 := &WriteFrame{
		Column: true,
		Arrays: [][]string{x1, x2},
		File:   file2,
	}
	w2.WriteCSV()
}

func TestFloat(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, true)
	Column := rds.Col("test").Float()

	if len(Column) == 0 {
		t.Error("no columns returned")
	}
	if Column[0] != 1231 {
		t.Error("unexpected column")

	}
}
func TestColumn(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, false)
	col := rds.Col("high").Float()

	if col[0] != 20492.10 {
		t.Error("not the same column data")
	}

}

func TestStrConv(t *testing.T) {
	int1 := []int{1, 2, 3, 4, 5}
	strInt1 := StringConv(int1)
	if strInt1[2] != "4" {
		t.Error("Unable to convert integer to string")
	}
}

func TestRow(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, false)
	row := rds.Row(1).Float()
	t.Error(row[1])
	if row[1] != 20445.50 {
		t.Error("wrong row data")
	}
}

func TestRows(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, true)
	rows := rds.Rows(1, 2, 6)
	nwRow := rows[1].Float()
	expected := []float64{1667271600000, 20568.10, 20577.00, 20464.60, 20471.30, 11101.117, 1667275199999, 227631088.55590, 74440, 4800.732, 98433200.32610, 0}
	if nwRow[4] != expected[4] {
		t.Error("wrong row data")
	}
}

func TestReplaceRow(t *testing.T) {
	err := ReplaceRow("./test.csv", 0770, 1, []string{"q", "e", "x"})

	if err != nil {
		t.Error(err)
	} else {
		return
	}
}

type RwStr struct {
	One   int
	Two   float64 `position:"5"` // testing tags
	Three float64
}

func TestInterface(t *testing.T) {

	decRwStr := &RwStr{}
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, true)
	err := rds.Row(2).Interface(decRwStr)

	if decRwStr.One == 1667268000000 {
		t.Error(decRwStr)
		t.Error(err)
	}

}

func TestInterface2(t *testing.T) {
	var nwtr interface{}
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, true)
	err := rds.Row(2).Interface(&nwtr)
	if err != nil {
		t.Error(err)
	}
}

func TestColWithPositon(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, true)
	col := rds.ColWithPosition(1).Float()

	if col[0] != 20482.10 {
		t.Error("not the same column data")
	}

}

func TestReadHeader(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 0666, true)
	header, err := rds.CheckHeader("open_time")
	if err != nil {
		t.Error(err)
	}
	if header < 0 {
		t.Error(err)
	}
}
