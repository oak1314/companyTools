package main

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os/exec"
	"strings"
)

const sourceData string = "./作成用データ/請求書元データ.xlsx"
const templateData string = "./作成用データ/請求書テンプレート.xlsx"

// Invoice 請求書
type Invoice struct {
	cellCompany  string
	cellMonth    string
	cellPayMonth string
	cellNames    []string
	cellMoney    []string
}

func main() {

	f, err := excelize.OpenFile(sourceData)
	if err != nil {
		println(err.Error())
		return
	}

	names := make([]string, 0)
	moneys := make([]string, 0)
	inv := new(Invoice)
	// var inv Invoice
	// inv := Invoice{}
	currentLine := 3
	currentCompany := ""
	currentMonth := ""
	currentPayMonth := ""
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

		cellPayMonth, err := f.GetCellValue("Sheet1", fmt.Sprintf("J%d", currentLine))
		if err != nil {
			println(err.Error())
			return
		}

		if currentLine == 3 {
			// 一回目
			currentCompany = cellCompany
			currentMonth = cellMonth
			currentPayMonth = cellPayMonth
		}

		// 会社が変わる
		if cellCompany != currentCompany && currentCompany != "" {
			// 一個前の会社フィアルを生成
			inv.cellCompany = currentCompany
			inv.cellMonth = currentMonth
			inv.cellPayMonth = currentPayMonth
			inv.cellNames = names
			inv.cellMoney = moneys
			createNewExcel(*inv)

			// 現在の会社を代入
			currentCompany = cellCompany
			currentMonth = cellMonth
			currentPayMonth = cellPayMonth
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
	CmdPython()
}

// call python script
func CmdPython() (err error) {
	args := []string{"convertToPDF.py"}
	out, err := exec.Command("python", args...).Output()
	if err != nil {
		return
	}
	result := string(out)
	if strings.Index(result, "success") != 0 {
		err = errors.New(fmt.Sprintf("main.py error：%s", result))
	}
	return
}

func createNewExcel(inv Invoice) {
	f, err := excelize.OpenFile(templateData)
	if err != nil {
		println(err.Error())
		return
	}
	// 請求会社名
	f.SetCellDefault("Sheet1", "B9", inv.cellCompany+"　御中")
	// 印刷エリア外(月)
	f.SetCellDefault("Sheet1", "B3", inv.cellMonth)
	// お振込期限
	if inv.cellPayMonth == "1" {
		f.SetCellFormula("Sheet1", "B38", "EOMONTH(DATE(YEAR(NOW()),$B$3,1),1)")
	}else if inv.cellPayMonth == "1.5" {
		f.SetCellFormula("Sheet1", "B38", "EOMONTH(DATE(YEAR(NOW()),$B$3,1),1)+DAY(15)")
	}else if inv.cellPayMonth == "2" {
		f.SetCellFormula("Sheet1", "B38", "EOMONTH(DATE(YEAR(NOW()),$B$3,1),2)")
	}
	//f.CalcCellValue("Sheet1", "B38")

	// 印刷エリア外(氏名)
	for index, cellName := range inv.cellNames {
		f.SetCellDefault("Sheet1", fmt.Sprintf("A%d", 3+index), cellName)
	}
	for index, money := range inv.cellMoney {
		// 単価
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