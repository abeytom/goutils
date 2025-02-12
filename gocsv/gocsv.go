package gocsv

import (
	"encoding/csv"
	"fmt"
	"github.com/abeytom/goutils/gofile"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
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

func ReadByRecord(fp string, each func([]string) error) error {
	tf, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer tf.Close()
	reader := csv.NewReader(tf)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// use as a poison to stop reading
		err = each(record)
		if err != nil {
			return err
		}
	}
}

func NewCsvWriter(fp string) (*os.File, *csv.Writer, error) {
	parent := filepath.Dir(fp)
	if !gofile.IsDir(parent) {
		err := os.MkdirAll(parent, 0755)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "cannot create dir [%v]", parent)
		}
	}
	f, err := os.Create(fp)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot open the file [%v]", fp)
	}
	return f, csv.NewWriter(f), nil
}

func NewCsvWriterOrPanic(fp string) (*os.File, *csv.Writer) {
	file, writer, err := NewCsvWriter(fp)
	if err != nil {
		panic(err)
	}
	return file, writer
}

func WriteRecord(writer *csv.Writer, keys ...string) error {
	return writer.Write(keys)
}
