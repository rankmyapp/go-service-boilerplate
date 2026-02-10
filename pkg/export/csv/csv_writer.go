package csv

import (
	"bytes"
	"encoding/csv"
)

func writeCSV(rows [][]string, delimiter rune) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	writer.Comma = delimiter
	writer.UseCRLF = false

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
