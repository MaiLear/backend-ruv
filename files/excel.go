package files

import (
	db "file/database"
	"fmt"
	// "fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func GetExcelData(fileName string) {
	file, err := excelize.OpenFile("./uploads/"+fileName)
	if err != nil {
		log.Fatal("error open the file")
	}

	defer file.Close()
	//Obtener una lista de todas las hojas del archivo
	sheets := file.GetSheetList()
	
	for _,sheetName := range sheets{
		//Obtener las filas de la hoja actual del archivo
		rows,err := file.GetRows(sheetName)
		if err != nil {
			log.Println("error to read sheet row",sheetName,err)
			continue
		}
		
		fmt.Println(rows)
		//Insertar en la base de datos
		db.Insert(rows)

	}

	

	

}
