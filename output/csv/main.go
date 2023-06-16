package csvoutput

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type CsvOutputer struct {
	value CsvConvertable
}

func NewCsvOutputer(val interface{}) (CsvOutputer, error) {
	if value, ok := val.(CsvConvertable); !ok {
		return CsvOutputer{}, fmt.Errorf("value is not CsvConvertable")
	} else {
		return CsvOutputer{value: value}, nil
	}
}

func (c CsvOutputer) WriteToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	// Open a CSV file
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("could not create CSV file: %e", err)

	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err = writer.Write(c.value.GetHeaders()); err != nil {
		return fmt.Errorf("could not write to CSV file: %e", err)
	}
	if err = writer.WriteAll(c.value.GetValues()); err != nil {
		return fmt.Errorf("could not write to CSV file: %e", err)
	}
	return nil
}
