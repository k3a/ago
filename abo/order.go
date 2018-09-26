package abo

import (
	"io"
	"os"
	"time"
)

// Item is a payment order item
type Item struct {
	Recipient struct {
		AccountNumPrefix int
		AccountNum       int
		BankCode         int
	}
	Amount              float64
	VS                  int
	KS                  int
	SS                  int
	MessageForRecipient string
}

// Group groups payment order items to be made from a single fund source
type Group struct {
	Payer struct {
		AccountNumPrefix int
		AccountNum       int
	}
	DueDate time.Time

	Items []*Item
}

// Order is an order in ABO/KPC format
type Order struct {
	CreationDate time.Time
	Client       struct {
		Name          string
		AccountNumber int
		BankCode      int
	}
	Groups []*Group
}

// Write writes an item to a writer
func (it *Item) Write(inWr io.Writer) error { //nolint:gocyclo,doesn't make sense
	wr := newWriter(inWr)

	// recipient account number (format: 000000-0000000000)
	if err := wr.WriteInt(it.Recipient.AccountNumPrefix, 6); err != nil {
		return err
	}
	if err := wr.WriteStr("-", 1); err != nil {
		return err
	}
	if err := wr.WriteInt(it.Recipient.AccountNum, 10); err != nil {
		return err
	}

	// field separator
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// amount
	if err := wr.WriteMonetaryAmount(it.Amount, 15); err != nil {
		return err
	}

	// field separator
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// VS
	if err := wr.WriteInt(it.VS, 10); err != nil {
		return err
	}

	// field separator
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// bank
	if err := wr.WriteInt(it.Recipient.BankCode, 4); err != nil {
		return err
	}

	// KS
	if err := wr.WriteInt(it.KS, 4); err != nil {
		return err
	}

	// field separator
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// SS (optional)
	if it.SS != 0 {
		if err := wr.WriteInt(it.SS, 10); err != nil {
			return err
		}

		// field separator
		if err := wr.WriteStr(" ", 1); err != nil {
			return err
		}
	}

	// msg for recipient (optional)
	if len(it.MessageForRecipient) > 0 {
		maxMsgLen := 4 * 35
		trimmedMsg := it.MessageForRecipient
		if len(it.MessageForRecipient) > maxMsgLen {
			trimmedMsg = trimmedMsg[:maxMsgLen]
		}

		if err := wr.WriteStr("AV:"+trimmedMsg, 3+ /*maxMsgLen*/ len(trimmedMsg)); err != nil {
			return err
		}
	}

	// new line
	if err := wr.WriteLineEnd(); err != nil {
		return err
	}

	return nil
}

// Write writes the order group to a writer
func (gr *Group) Write(inWr io.Writer) error { //nolint:gocyclo,doesn't make sense
	wr := newWriter(inWr)

	// start of group (msg type + sep)
	if err := wr.WriteStr("2 ", 1+1); err != nil {
		return err
	}

	// payer account (000000-0000000000)
	if err := wr.WriteInt(gr.Payer.AccountNumPrefix, 6); err != nil {
		return err
	}
	if err := wr.WriteStr("-", 1); err != nil {
		return err
	}
	if err := wr.WriteInt(gr.Payer.AccountNum, 10); err != nil {
		return err
	}

	// field separator
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// total amount in the group
	totalAmount := 0.0
	for _, it := range gr.Items {
		totalAmount += it.Amount
	}
	if err := wr.WriteMonetaryAmount(totalAmount, 14); err != nil {
		return err
	}

	// field separator
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// due date
	if err := wr.WriteTime(gr.DueDate); err != nil {
		return err
	}

	// new line
	if err := wr.WriteLineEnd(); err != nil {
		return err
	}

	// write items
	for _, it := range gr.Items {
		if err := it.Write(wr); err != nil {
			return err
		}
	}

	// end of group
	if err := wr.WriteStr("3 +", 3); err != nil {
		return err
	}

	// new line
	if err := wr.WriteLineEnd(); err != nil {
		return err
	}

	return nil
}

// AddItem adds a payment order to the group
func (gr *Group) AddItem(recpAccPrefix, recpAccNum, recpBankCode int, amount float64, vs, ks, ss int, msgForRecp string) *Item {
	it := new(Item)

	it.Recipient.AccountNumPrefix = recpAccPrefix
	it.Recipient.AccountNum = recpAccNum
	it.Recipient.BankCode = recpBankCode
	it.Amount = amount
	it.VS = vs
	it.KS = ks
	it.SS = ss
	it.MessageForRecipient = msgForRecp

	gr.Items = append(gr.Items, it)

	return it
}

// AddItemSimple is shorter version of AddItem.
// Adds a payment order to the group
func (gr *Group) AddItemSimple(recpAccPrefix, recpAccNum, recpBankCode int, amount float64, vs int, msgForRecp string) *Item {
	return gr.AddItem(recpAccPrefix, recpAccNum, recpBankCode, amount, vs, 0, 0, msgForRecp)
}

func (or *Order) writeAccounting(inWr io.Writer) error { //nolint:gocyclo,doesnt make sense
	wr := newWriter(inWr)

	// start of accounting (msg type + sep)
	if err := wr.WriteStr("1 ", 1+1); err != nil {
		return err
	}

	// type of data + sep
	if err := wr.WriteInt(1501, 4); err != nil {
		return err
	}
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	// accounting file number + sep
	if err := wr.WriteInt(0, 6); err != nil {
		return err
	}
	if err := wr.WriteStr(" ", 1); err != nil {
		return err
	}

	if err := wr.WriteInt(or.Client.BankCode, 4); err != nil {
		return err
	}

	// new line
	if err := wr.WriteLineEnd(); err != nil {
		return err
	}

	// write groups
	for _, gr := range or.Groups {
		if err := gr.Write(wr); err != nil {
			return err
		}
	}

	// end of accounting
	if err := wr.WriteStr("5 +", 3); err != nil {
		return err
	}

	// new line
	if err := wr.WriteLineEnd(); err != nil {
		return err
	}

	return nil
}

// Write writes the order to a writer
func (or *Order) Write(inWr io.Writer) error {
	wr := newWriter(inWr)

	// Message type
	if err := wr.WriteStr("UHL1", 4); err != nil {
		return newErr("unable to write msg header: %v", err)
	}

	// creation date
	if err := wr.WriteTime(or.CreationDate); err != nil {
		return newErr("unable to write creation date: %v", err)
	}

	// name of client
	if err := wr.WriteStrWindows1250(or.Client.Name, 20); err != nil {
		return newErr("unable to write client name: %v", err)
	}

	// client acc num
	if err := wr.WriteInt(or.Client.AccountNumber, 10); err != nil {
		return newErr("unable to write client account number: %v", err)
	}

	// Interval of accounting files - start
	if err := wr.WriteInt(0, 3); err != nil {
		return newErr("unable to write acc interval start: %v", err)
	}

	// Interval of accounting files - end
	if err := wr.WriteInt(999, 3); err != nil {
		return newErr("unable to write acc interval end: %v", err)
	}

	// Code fixed part
	if err := wr.WriteInt(0, 6); err != nil {
		return newErr("unable to write fixed part: %v", err)
	}

	// Code secret part
	if err := wr.WriteInt(0, 6); err != nil {
		return newErr("unable to write secret part: %v", err)
	}

	// new line
	if err := wr.WriteLineEnd(); err != nil {
		return newErr("unable to write line end: %v", err)
	}

	return or.writeAccounting(wr)
}

// WriteToFile writes the order to a .kpc file
func (or *Order) WriteToFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return newErr("unable to open target file %s: %s", path, err)
	}
	defer f.Close()

	return or.Write(f)
}

// AddGroup adds a payment group. It is a group of payment orders sent
// from a single source of funds.
func (or *Order) AddGroup(payerAccountPrefix, payerAccountNum int, dueDate time.Time) *Group {
	gr := new(Group)

	gr.Payer.AccountNumPrefix = payerAccountPrefix
	gr.Payer.AccountNum = payerAccountNum
	gr.DueDate = dueDate

	or.Groups = append(or.Groups, gr)

	return gr
}
