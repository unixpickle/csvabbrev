package main

import (
	"encoding/csv"
	"io"
)

func DeflateStream(r io.Reader, w io.Writer, buffer bool) error {
	var lastRecord []string
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()
	return ReadCSV(r, func(record []string) error {
		deflated := deflateRecord(record, lastRecord)
		lastRecord = record
		if err := csvWriter.Write(deflated); err != nil {
			return err
		}
		if !buffer {
			csvWriter.Flush()
			if err := csvWriter.Error(); err != nil {
				return err
			}
		}
		return nil
	})
}

func deflateRecord(record, lastRecord []string) []string {
	res := make([]string, len(record))
	for i, x := range record {
		if lastRecord != nil && x == lastRecord[i] && x != "" {
			res[i] = "^"
		} else if IsCarrotOnly(x) {
			res[i] = "^" + x
		} else {
			res[i] = x
		}
	}
	return res
}
