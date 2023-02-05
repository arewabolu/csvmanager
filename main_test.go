package csvmanager

import (
	"testing"
)

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
func TestColumn(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv", 200)
	col := rds.Col("high").Float()

	if col[0] != 20492.10 {
		t.Error("not the same column data")
	}

}

func TestRow(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	row := rds.Row(1).Float()

	if row[1] != 20445.50 {
		t.Error("wrong row data")
	}
}

func TestRows(t *testing.T) {
	rds, _ := ReadCsv("./BTCUSDT-1h-2022-11.csv")
	rows := rds.Rows(1, 2, 6)
	nwRow := rows[1].Float()
	expected := []float64{1667271600000, 20568.10, 20577.00, 20464.60, 20471.30, 11101.117, 1667275199999, 227631088.55590, 74440, 4800.732, 98433200.32610, 0}
	if nwRow[4] != expected[4] {
		t.Error("wrong row data")
	}
}

func TestReplaceRow(t *testing.T) {
	rep := ReplaceRow("./test.csv", 1, []string{"j", "e", "x"})

	if rep.Err != nil {
		t.Error(rep.Err)
	} else {
		return
	}
}
