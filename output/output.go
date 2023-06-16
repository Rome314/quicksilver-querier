package output

import (
	csvoutput "QuicksilverDumper/output/csv"
)

type Outputer interface {
	WriteToFile(path string) error
}

func GetCSVOutputer(v csvoutput.CsvConvertable) (Outputer, error) {
	return csvoutput.NewCsvOutputer(v)
}
