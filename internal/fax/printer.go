package fax

import (
	"os"

	"github.com/securityguy/escpos"
)

func NewPrinter(devicePath string) (*escpos.Escpos, error) {
	file, err := os.OpenFile(devicePath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	return escpos.New(file), nil
}
