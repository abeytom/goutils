package gocsv

import (
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

func GetColumnIndexOrPanic(records [][]string, k string) int {
	index, err := GetColumnIndex(records, k)
	if err != nil {
		panic(err)
	}
	return index
}

func GetColumnIndex(records [][]string, k string) (int, error) {
	index := FindColumnIndex(records, k)
	if index >= 0 {
		return index, nil
	}
	return -1, errors.New(fmt.Sprintf("cannot find the key [%v] from %v", k, records[0]))
}

func FindColumnIndex(records [][]string, k string) int {
	if len(records) == 0 {
		return -1
	}
	row := records[0]
	for i, key := range row {
		if k == key {
			return i
		}
	}
	return -1
}

func ReadAllCsv(fp string) ([][]string, error) {
	tf, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer tf.Close()
	reader := csv.NewReader(tf)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
