package csvloader

import (
	"encoding/csv"
	"os"
)

type CSVLoader struct {
	filePath string
}

func NewCSVLoader(filePath string) *CSVLoader {
	return &CSVLoader{filePath: filePath}
}

func (c *CSVLoader) Load() ([][]string, error) {
	file, err := os.Open(c.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
