package main

import (
	"bytes"
	"io/ioutil"
	"os"

	humcommon "github.com/xeedio/linux-https-user-management"
)

func fileContains(line []byte, filePath string) (bool, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		humcommon.Log().Warnf("Error reading file %s: %v", filePath, err)
		return false, err
	}
	return bytes.Contains(data, line), nil
}

func appendLineToFile(line []byte, filePath string) error {
	present, err := fileContains(line, filePath)
	if err != nil {
		humcommon.Log().Warnf("Error from fileContains: %v", err)
		return err
	}

	if present {
		return nil
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		humcommon.Log().Warnf("Error opening %s: %v", filePath, err)
		return err
	}
	if _, err := f.Write(line); err != nil {
		f.Close() // ignore error; Write error takes precedence
		humcommon.Log().Warnf("Error writing line %s to file: %v", line, err)
		return err
	}
	if err := f.Close(); err != nil {
		humcommon.Log().Warnf("Error closing file: %v", err)
		return err
	}

	return nil
}
