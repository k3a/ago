package abo

import (
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// reader enhancing comon reader with methods used in ABO format
type reader struct {
	io.Reader
}

const formatDDMMYY = "020106"

func cleanStr(b []byte) string {
	return strings.TrimSpace(string(b))
}

// ReadStrWindows1250 reads windows-1250-formated string and returns UTF-8 string
func (rdr *reader) ReadStrWindows1250(buff []byte) (string, error) {
	if _, err := rdr.Read(buff); err != nil {
		return "", err
	}

	dec := charmap.Windows1250.NewDecoder()
	out, err := dec.Bytes(buff)
	if err != nil {
		return "", err
	}

	return cleanStr(out), nil
}

// ReadStr reads ASCII/UTF-8 string
func (rdr *reader) ReadStr(buff []byte) (string, error) {
	if _, err := rdr.Read(buff); err != nil {
		return "", err
	}

	return cleanStr(buff), nil
}

// ReadInt reads integer value
func (rdr *reader) ReadInt(buff []byte) (int, error) {
	if _, err := rdr.Read(buff); err != nil {
		return 0, err
	}

	return strconv.Atoi(cleanStr(buff))
}

// ReadMonetaryAmount reads monetary value
func (rdr *reader) ReadMonetaryAmount(buff []byte) (float64, error) {
	halere, err := rdr.ReadInt(buff)
	if err != nil {
		return 0, err
	}

	return float64(halere) / 100, nil
}

// ReadTime reads DDMMYY time
func (rdr *reader) ReadTime(buff []byte) (time.Time, error) {
	if _, err := rdr.Read(buff); err != nil {
		return time.Time{}, err
	}

	return time.Parse(formatDDMMYY, cleanStr(buff))
}

func newAboReader(rdr io.Reader) *reader {
	if rdr, isAlready := rdr.(*reader); isAlready {
		return rdr
	}
	return &reader{rdr}
}

type writer struct {
	io.Writer
}

func (wr *writer) WritePad(inBytes []byte, byteLen int, paddingByte byte, padLeft bool) error {
	if len(inBytes) > byteLen {
		// shortened data
		_, err := wr.Write(inBytes[:byteLen])
		return err
	}

	// prep padding
	pad := make([]byte, byteLen-len(inBytes))
	for i := range pad {
		pad[i] = paddingByte
	}

	// left padding
	if padLeft {
		_, err := wr.Write(pad)
		if err != nil {
			return err
		}
	}

	// write data
	_, err := wr.Write(inBytes)
	if err != nil {
		return err
	}

	// right padding
	if !padLeft {
		_, err := wr.Write(pad)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wr *writer) WriteStr(str string, byteLen int) error {
	return wr.WritePad([]byte(str), byteLen, ' ', false)
}

func (wr *writer) WriteStrWindows1250(str string, byteLen int) error {
	enc := charmap.Windows1250.NewEncoder()
	out, err := enc.String(str)
	if err != nil {
		return err
	}

	return wr.WriteStr(out, byteLen)
}

func (wr *writer) WriteInt(i int, byteLen int) error {
	return wr.WritePad([]byte(strconv.Itoa(i)), byteLen, '0', true)
}

func (wr *writer) WriteMonetaryAmount(amount float64, byteLen int) error {
	return wr.WriteInt(int(amount*100), byteLen)
}

func (wr *writer) WriteTime(tm time.Time) error {
	return wr.WriteStr(tm.Format(formatDDMMYY), 6)
}

func (wr *writer) WriteLineEnd() error {
	_, err := wr.Write([]byte("\n"))
	return err
}

func newWriter(wr io.Writer) *writer {
	if wr, isAlready := wr.(*writer); isAlready {
		return wr
	}
	return &writer{wr}
}
