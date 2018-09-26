package abo

import (
	"os"
	"testing"
)

func TestStatement(t *testing.T) {
	rdr, _ := os.Open("./test/fio.gpc")
	defer rdr.Close()

	abo, err := FromReader(rdr)
	if err != nil {
		t.Fatal(err)
	}

	if len(abo.Transactions) == 0 {
		t.Fatal("no transactions parsed")
	}

	tr := abo.Transactions[0]

	if !(tr.Amount > 1234 && tr.Amount < 1235) {
		t.Fatal("bad amount")
	}

	if tr.Recipient.BankCode != 100 {
		t.Fatal("bad bank id")
	}

	if tr.VS != 1446556401 {
		t.Fatal("bad VS")
	}

	if tr.KS != 308 {
		t.Fatal("bad KS")
	}

	if tr.SS != 7815392681 {
		t.Fatal("bad SS")
	}
}
