package abo

import (
	"bytes"
	"testing"
	"time"
)

func TestOrder(t *testing.T) {
	buff := new(bytes.Buffer)

	o := new(Order)

	o.CreationDate = time.Now()
	o.Client.BankCode = 2010

	gr := o.AddGroup(0, 2101135843, time.Now().Add(24*time.Hour))
	//gr.AddItem(321, 7654, 9966, 1122.33, 88888888 /*VS*/, 9922 /*KS*/, 9933 /*SS*/, "" /*msg*/)
	gr.AddItemSimple(0, 1900133399, 2010 /*bank*/, 1.23 /*amnt*/, 88888888 /*VS*/, "" /*msg*/)

	if err := o.Write(buff); err != nil {
		t.Fatal(err)
	}

	//t.Fatal(buff.String())
}
