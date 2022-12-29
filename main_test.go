package csvmanager

import (
	"fmt"
	"testing"
)

func TestColumn(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 200)
	col := rds.Column("high")

	if col[0] != "204192.10" {
		fmt.Println(col)
		t.Error("not the same column data")
	}

}

func TestRow(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	row := rds.Row(1)

	if row[1] != "20445.50" {
		t.Error("wrong row data")
	}
}

func TestRows(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	rows := rds.Rows(1, 2, rds.RowLength())
	nwRow := rows[1]
	expected := []string{"1667264400000", "20548.50", "20445.40", "20502.00", "11047.159", "1667267999999", "226531027.34130", "74358,5874.933", "120471227.79020", "0"}
	if nwRow[4] != expected[4] {
		fmt.Println(nwRow[0])
		t.Error("wrong row data")
	}
}