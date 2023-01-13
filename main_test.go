package csvmanager

import (
	"fmt"
	"testing"
)

func TestColumn(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 200)
	col := rds.Col("high").Float()

	if col[0] != 20492.10 {
		t.Error("not the same column data")
	}

}

func TestRow(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	row := rds.Row(0).Float()

	if row[1] != 20445.50 {
		fmt.Println(row[1])
		t.Error("wrong row data")
	}
}

func TestRows(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	rows := rds.Rows(1, 2, rds.SizeofRows())
	nwRow := rows[1].Float()
	expected := []float64{1667264400000, 20548.50, 20445.40, 20502.00, 11047.159, 1667267999999, 226531027.34130, 74358, 5874.933, 120471227.79020, 0}
	if nwRow[4] != expected[4] {
		t.Error("wrong row data")
	}
}

func TestFloat(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	Column := rds.Col("test").Float()

	if len(Column) == 0 {
		t.Error("no columns returned")
	}
	if Column[0] != 1231 {
		t.Error("unexpected column")

	}
}
