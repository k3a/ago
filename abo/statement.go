package abo

import (
	"fmt"
	"io"
	"time"

	"github.com/k3a/ago/abo/currency"
)

// Transaction single transaction
type Transaction struct {
	OwnerAccountNumber int
	Recipient          struct {
		Name             string
		AccountNumPrefix int
		AccountNum       int
		BankCode         int
	}
	ID       int
	Amount   float64
	Currency currency.Currency
	Type     int // 1-debet, 2-credit, 4-storno-debet, 5-storno-credit
	VS       int
	KS       int
	SS       int
	DueDate  time.Time
}

var errNoMoreTransactions = newErr("no more transactions in the input")

func (txn *Transaction) String() string {
	return fmt.Sprintf("ID: %d, Type: %d, Recipient: %s, Recipient Acc: %06d-%d/%04d Amount: %.2f %s, VS: %d, KS: %d, SS: %d, Due Date: %s",
		txn.ID, txn.Type, txn.Recipient.Name, txn.Recipient.AccountNumPrefix, txn.Recipient.AccountNum, txn.Recipient.BankCode,
		txn.Amount, txn.Currency, txn.VS, txn.KS, txn.SS, txn.DueDate)
}

// Parse parses a single transaction from the input stream
func (txn *Transaction) Read(inRdr io.Reader) error { //nolint:gocyclo,doesn't make sense here
	rdr := newAboReader(inRdr)
	buf := make([]byte, 64)

	// 1-3 record type
	recType, err := rdr.ReadStr(buf[:3])
	if err != nil {
		if err == io.EOF {
			return errNoMoreTransactions
		}
		return newErr("problem reading txn record type: %v", err)
	}
	if recType != "075" {
		return newErr("wrong txn record type %s", recType)
	}

	// 4-19 owner acc number
	txn.OwnerAccountNumber, err = rdr.ReadInt(buf[:16])
	if err != nil {
		return newErr("problem reading txn owner account num: %v", err)
	}

	// 20-35 recipient acc num
	txn.Recipient.AccountNumPrefix, err = rdr.ReadInt(buf[:6])
	if err != nil {
		return newErr("problem reading txn counterparty acc num prefix: %v", err)
	}
	txn.Recipient.AccountNum, err = rdr.ReadInt(buf[:10])
	if err != nil {
		return newErr("problem reading txn counterparty acc num: %v", err)
	}

	// 36-48 txn id
	txn.ID, err = rdr.ReadInt(buf[:13])
	if err != nil {
		return newErr("problem reading txn id: %v", err)
	}

	// 49-60 amount
	if txn.Amount, err = rdr.ReadMonetaryAmount(buf[:12]); err != nil {
		return newErr("problem reading txn amount: %v", err)
	}

	// 61 type
	txn.Type, err = rdr.ReadInt(buf[:1])
	if err != nil {
		return newErr("problem reading txn type: %v", err)
	}

	// 62-71 VS
	txn.VS, err = rdr.ReadInt(buf[:10])
	if err != nil {
		return newErr("problem reading txn VS: %v", err)
	}

	rdr.Read(buf[:2]) //nolint:gosec,skip 2 bytes

	// counterparty bank id
	txn.Recipient.BankCode, err = rdr.ReadInt(buf[:4])
	if err != nil {
		return newErr("problem reading txn counterparty bank id: %v", err)
	}

	// KS
	txn.KS, err = rdr.ReadInt(buf[:4])
	if err != nil {
		return newErr("problem reading txn KS: %v", err)
	}

	// SS
	txn.SS, err = rdr.ReadInt(buf[:10])
	if err != nil {
		return newErr("problem reading txn SS: %v", err)
	}

	rdr.Read(buf[:6]) //nolint:gosec,skip 6 bytes

	// counterparty acc name
	txn.Recipient.Name, err = rdr.ReadStrWindows1250(buf[:20])
	if err != nil {
		return newErr("problem reading txn counterparty acc name: %v", err)
	}

	rdr.Read(buf[:1]) //nolint:gosec,skip 1 byte

	// currency
	currencyIdent, err := rdr.ReadInt(buf[:4])
	if err != nil {
		return newErr("problem reading txn currency: %v", err)
	}
	txn.Currency = currency.Currency(uint16(currencyIdent))

	// due date
	txn.DueDate, err = rdr.ReadTime(buf[:6])
	if err != nil {
		return newErr("problem reading txn due date: %v", err)
	}

	rdr.Read(buf[:2]) //nolint:gosec,skip 2 crlf bytes

	return nil
}

// Statement with transactions in ABO/GPC format
type Statement struct {
	Info struct {
		AccountNumber   int
		AccountName     string
		StartDate       time.Time
		EndDate         time.Time
		OpeningBalance  float64
		ClosingBalance  float64
		IncomeSum       float64
		ExpenseSum      float64
		StatementNumber int
	}

	Transactions []*Transaction
}

func (s *Statement) readTransactions(rdr io.Reader) error {
	// start empty
	s.Transactions = []*Transaction{}

	for {
		txn := new(Transaction)
		if err := txn.Read(rdr); err != nil {
			if err == errNoMoreTransactions {
				return nil
			}
			return err
		}

		s.Transactions = append(s.Transactions, txn)
	}
}

// Read reads ABO/GPC statement from a reader
func (s *Statement) Read(inRdr io.Reader) error { //nolint:gocyclo,doesn't make sense here
	buf := make([]byte, 32)
	rdr := newAboReader(inRdr)

	// record type
	recType, err := rdr.ReadStr(buf[:3])
	if err != nil {
		return newErr("problem reading record type: %v", err)
	}
	if recType != "074" {
		return newErr("wrong record type %s", recType)
	}

	// acc number
	s.Info.AccountNumber, err = rdr.ReadInt(buf[:16])
	if err != nil {
		return newErr("problem reading account num: %v", err)
	}

	// acc name
	s.Info.AccountName, err = rdr.ReadStrWindows1250(buf[:20])
	if err != nil {
		return newErr("problem reading account name: %v", err)
	}

	// start date
	s.Info.StartDate, err = rdr.ReadTime(buf[:6])
	if err != nil {
		return newErr("problem reading start date: %v", err)
	}

	// opening balance
	openingBalance, err := rdr.ReadInt(buf[:14])
	if err != nil {
		return newErr("problem reading opening balance: %v", err)
	}
	s.Info.OpeningBalance = float64(openingBalance) / 100

	// opening balance sign
	if _, err = rdr.Read(buf[:1]); err != nil {
		return newErr("problem reading opening balance: %v", err)
	}
	if string(buf[:1]) == "-" {
		s.Info.OpeningBalance = -s.Info.OpeningBalance
	}

	// closing balance
	closingBalance, err := rdr.ReadInt(buf[:14])
	if err != nil {
		return newErr("problem reading closing balance: %v", err)
	}
	s.Info.ClosingBalance = float64(closingBalance) / 100

	// closing balance sign
	if _, err = rdr.Read(buf[:1]); err != nil {
		return newErr("problem reading closing balance: %v", err)
	}
	if string(buf[:1]) == "-" {
		s.Info.ClosingBalance = -s.Info.ClosingBalance
	}

	// expense sum
	if s.Info.ExpenseSum, err = rdr.ReadMonetaryAmount(buf[:14]); err != nil {
		return newErr("problem reading expense sum: %v", err)
	}
	rdr.Read(buf[:1]) //nolint:gosec, skip single byte

	// income sum
	if s.Info.IncomeSum, err = rdr.ReadMonetaryAmount(buf[:14]); err != nil {
		return newErr("problem reading income sum: %v", err)
	}
	rdr.Read(buf[:1]) //nolint:gosec, skip single byte

	// statement number
	s.Info.StatementNumber, err = rdr.ReadInt(buf[:3])
	if err != nil {
		return newErr("problem reading statement number: %v", err)
	}

	// end date
	s.Info.EndDate, err = rdr.ReadTime(buf[:6])
	if err != nil {
		return newErr("problem reading end date: %v", err)
	}

	rdr.Read(buf[:14+2]) //nolint:gosec,skip 14 bytes + 2 crlf bytes

	return s.readTransactions(rdr)
}

// String formats the statement as a human-readable summary string
func (s *Statement) String() string {
	str := fmt.Sprintf("Statement For %d (%s)\n"+
		"Opening Balance %.2f, Closing Balance %.2f\n"+
		"Range %s - %s\n\nTransactions:\n",
		s.Info.AccountNumber, s.Info.AccountName,
		s.Info.OpeningBalance, s.Info.ClosingBalance,
		s.Info.StartDate, s.Info.EndDate)

	for _, tx := range s.Transactions {
		str += "- " + tx.String() + "\n"
	}

	return str
}

// FromReader parses ABO statement from io.Reader
func FromReader(rdr io.Reader) (*Statement, error) {
	// close input if possible
	defer func() {
		if rdrc, ok := rdr.(io.ReadCloser); ok {
			rdrc.Close() //nolint:gosec
		}
	}()

	stmt := new(Statement)

	if err := stmt.Read(rdr); err != nil {
		return nil, err
	}

	return stmt, nil
}
