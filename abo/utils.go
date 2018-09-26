package abo

import "fmt"

func newErr(format string, args ...interface{}) error {
	return fmt.Errorf("abo: "+format, args...)
}
