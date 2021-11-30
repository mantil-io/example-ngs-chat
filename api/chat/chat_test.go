package chat

import (
	"fmt"
	"log"
	"testing"
)

func TestErrors(t *testing.T) {
	log.SetFlags(log.Llongfile)
	err := one()
	fmt.Printf("err: %s", err)
}

func Errorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	log.Output(2, fmt.Sprintf("%s", err))
	return err
}

func one() error {
	return Err(two())
}

func two() error {
	err := three()
	return Errorf("two failed because three returned %w", err)
}

func three() error {
	return Errorf("error in three")
}
