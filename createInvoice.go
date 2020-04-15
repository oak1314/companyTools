package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const sourceData string = "./作成用データ/請求書元データ.xlsx"
const templateData string = "./作成用データ/請求書テンプレート.xlsx"

// Invoice 請求書
type Invoice struct {
	cellCompany string
	cellMonth   string
	cellNames   []string
	cellMoney   []string
}

func main() {

	f, err := excelize.OpenFile(sourceData)
	if err != nil {
		println(err.Error())
		return
	}

	names := make([]string,0)
	moneys := make([]string,0)
	inv := new(Invoice)
	// var inv Invoice
	// inv := Invoice{}
	currentLine := 3
	currentCompany := ""
	currentMonth := ""
	for {
		cellCompany, err := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", currentLine))
		if err != nil {
			println(err.Error())
			return
		}

		cellMonth, err := f.GetCellValue("Sheet1", fmt.Sprintf("E%d", currentLine))
		if err != nil {
			println(err.Error())
			return
		}

		if currentLine == 3 {
			// 一回目
			currentCompany = cellCompany
			currentMonth = cellMonth
		}

		// 会社が変わる
		if cellCompany != currentCompany && currentCompany != "" {
			inv.cellCompany = currentCompany
			inv.cellMonth = currentMonth
			inv.cellNames = names
			inv.cellMoney = moneys
			createNewExcel(*inv)

			currentCompany = cellCompany
			currentMonth = cellMonth
			names = names[:0]
			moneys = moneys[:0]
		}

		cellName, err := f.GetCellValue("Sheet1", fmt.Sprintf("C%d", currentLine))
		if err != nil {
			println(err.Error())
			return
		}
		names = append(names, cellName)

		cellMoney, err := f.GetCellValue("Sheet1", fmt.Sprintf("D%d", currentLine))
		if err != nil {
			println(err.Error())
			return
		}
		moneys = append(moneys, cellMoney)
		currentLine++

		// お終い
		if cellCompany == "" || cellMonth == "" || cellName == "" || cellMoney == "" {
			break
		}
	}

}

func createNewExcel(inv Invoice) {
	f, err := excelize.OpenFile(templateData)
	if err != nil {
		println(err.Error())
		return
	}

	f.SetCellDefault("Sheet1", "B9", inv.cellCompany+"　御中")
	f.SetCellDefault("Sheet1", "B3", inv.cellMonth)

	for index, cellName := range inv.cellNames {
		f.SetCellDefault("Sheet1", fmt.Sprintf("A%d", 3+index), cellName)
	}
	for index, money := range inv.cellMoney {
		f.SetCellDefault("Sheet1", fmt.Sprintf("E%d", 20+index), money)
		if index == 0 {
			continue
		}

		cell, err := f.GetCellValue("Sheet1", fmt.Sprintf("B%d", 20+index-1))
		if err != nil {
			println(err.Error())
			return
		}
		f.SetCellDefault("Sheet1", fmt.Sprintf("B%d", 20+index), cell)

		// cell, err = f.GetCellFormula("Sheet1", fmt.Sprintf("C%d", 20+index-1))
		// if err != nil {
		// 	println(err.Error())
		// 	return
		// }
		// f.SetCellFormula("Sheet1", fmt.Sprintf("C%d", 20+index), cell)

		f.SetCellDefault("Sheet1", fmt.Sprintf("D%d", 20+index), "1")
	}

	f.UpdateLinkedValue()

	// Save xlsx file by the given path.
	if err := f.SaveAs("請求書_" + inv.cellCompany + "_" + inv.cellMonth + "月.xlsx"); err != nil {
		fmt.Println(err)
	}
}
