package main

import (
	"github.com/xuri/excelize/v2"
	"log"
)

func setLogValues(file *excelize.File, logValues map[string]string, key string, col int, row int) {
	if value, ok := logValues[key]; ok {
		cellName, err := excelize.CoordinatesToCellName(col+1, row+2)
		if err != nil {
			log.Println(err.Error())
			return
		}
		err = file.SetCellDefault("Daten Kontrolle", cellName, value)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func setCPRValues(
	file *excelize.File,
	cprValues map[string]map[string]SimPadCPREventParameter,
	eventName string,
	paramName string,
	col int,
	row int,
) {

	if param, ok := cprValues[eventName][paramName]; ok {
		cellName, err := excelize.CoordinatesToCellName(col+1, row+2)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = file.SetCellDefault("Daten Kontrolle", cellName, param.Value)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func setValuesForEntry(file *excelize.File, data *ProcessedSimPadData) error {

	rows, err := file.GetRows("Daten Kontrolle")
	if err != nil {
		return err
	}

	cols, err := file.GetCols("Daten Kontrolle")
	if err != nil {
		return err
	}

	// New Func
	var firstFreeCol int
	for i, col := range cols[3:] {
		if col[1] == "" {
			firstFreeCol = i + 3
			break
		}
	}

	for r, row := range rows[1:] {
		type1 := row[0]
		type2 := row[1]

		setLogValues(file, data.Log, type2, firstFreeCol, r)
		setCPRValues(file, data.CPR, type1, type2, firstFreeCol, r)
	}
	return nil
}
