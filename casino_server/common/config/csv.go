package config

import (
	"encoding/csv"
	"os"
)

// 从csv读取数据
func readCSV(file string) (records [][]string, e error) {
	f, err := os.Open(file)
	if err != nil {
		return records, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, e = r.ReadAll()
	if e != nil {
		return records, e
	}

	// 去头
	records = records[1:]

	return
}

// 写数据到csv
func writeCSV(file string, records[][]string) (e error) {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	e = w.WriteAll(records)
	if e != nil {
		return e
	}

	w.Flush()

	return
}
