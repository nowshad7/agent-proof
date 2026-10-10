// Package export writes records to CSV while preserving Unicode text.
package export

import (
	"bytes"
	"encoding/csv"
)

// ToCSV encodes rows as RFC 4180 CSV. Input text is written unchanged, so
// Unicode content is preserved.
func ToCSV(rows [][]string) (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
