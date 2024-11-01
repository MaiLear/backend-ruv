package files

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	db "file/database"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var incompleteArrays [][]string

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// fillIncompleteArray llena un arreglo incompleto
func fillIncompleteArray(rows []string, tabRow int) []string {
	quantity := tabRow - len(incompleteArrays[0])

	for i := 0; i < quantity; i++ {
		incompleteArrays[0] = append(incompleteArrays[0], rows[i])
		if i+1 == quantity {
			break
		}
	}

	return incompleteArrays[0]
}

// convertTextToArray convierte texto en un arreglo según el carácter
func convertTextToArray(text string, character string, defaultValue string) []string {
	array := make([]string, 0)
	counterCharacter := 0
	newWord := ""
	targetChar := []rune(character)[0] // Toma el primer carácter de la cadena.

	for _, char := range text {
		if char == targetChar {
			counterCharacter++
			continue
		}

		switch {
		case counterCharacter == 1:
			array = append(array, newWord)
			newWord = ""
			counterCharacter = 0
		case counterCharacter > 1:
			array = append(array, newWord)
			for i := 0; i < counterCharacter-1; i++ {
				array = append(array, defaultValue)
			}
			newWord = ""
			counterCharacter = 0
		}

		newWord += string(char)
	}

	// Asegurarse de agregar la última palabra si existe
	if newWord != "" {
		array = append(array, newWord)
	}

	return array
}

// filterFields filtra los campos de un texto dado según las opciones y el número de filas de tabla.
func filterFields(text []string, options []string, tabRow int) []string {
	deleteValue := options[0]
	emptyValue := options[1]
	result := []string{}

	regex := regexp.MustCompile(fmt.Sprintf(`.+%s`, regexp.QuoteMeta(deleteValue)))
	indexInsert := -1
	offset := tabRow
	oldValue := ""

	for index, row := range text {
		row = strings.ReplaceAll(row, "\n", "")
		row = strings.ReplaceAll(row, "\t", "")
		oldValue = row

		if regex.MatchString(row) {
			indexInsert = index + 1 // Establece la posición para insertar el deleteValue en la siguiente fila
			result = append(result, strings.Replace(row, deleteValue, "", -1))
		} else if index == indexInsert {
			result = append(result, deleteValue) // Inserta el deleteValue en la siguiente fila
			result = append(result, oldValue)    // Luego inserta el valor original que fue reemplazado
			indexInsert = -1                     // Reinicia el índice de inserción
		} else if strings.TrimSpace(row) == "" {
			result = append(result, emptyValue)
		} else {
			result = append(result, row) // Si no es un caso de reemplazo, inserta el valor original
		}

		if len(result) == offset{
			lastIndex := offset - 1
			if result[lastIndex] == deleteValue{
				result = append(result[:lastIndex], append([]string{"NULL"}, result[lastIndex:]...)...)
			}
		}
	}

	return result
}

func processTextToArray(text string, tabRow int, addColumn string) [][]string{
	// Llama a convertTextToArray y maneja el error
	rows := convertTextToArray(text, "»", "NULL")
	incompleteArrays = [][]string{}

	// Llama a filterFields y maneja el error
	filterRows := filterFields(rows, []string{
		addColumn,
		"NULL",
	}, 53)

	// Imprime filas originales y filtradas
	// fmt.Println("rows:", rows)
	// fmt.Println("filter rows:", filterRows)

	// Inicializa el array doble
	doubleArray := [][]string{}
	offset := 0

	for i := 0; i < len(filterRows); i += tabRow {
		if len(incompleteArrays) > 0 {
			arrayComplete := fillIncompleteArray(filterRows, tabRow)
			// Agrega el array completo a doubleArray
			doubleArray = append(doubleArray, arrayComplete)
			incompleteArrays = [][]string{} // Vacía el array de incompletos
		}

		col := filterRows[i:min(i+tabRow, len(filterRows))] // Obtener un sub-slice

		if len(col) > tabRow {
			col = col[:tabRow] // Limita la fila a 53 elementos
		}else if len(col) == tabRow {
			doubleArray = append(doubleArray, col)
		} else {
			incompleteArrays = append(incompleteArrays, col)
		}

		offset += tabRow
	}
	return doubleArray
}

func GetCsvData(fileName string) {
	file, err := os.Open("./uploads/" + fileName)

	if err != nil {
		fmt.Println("Error al abrir el archivo", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	//Tamaño buffer 125mb
	buffer := make([]byte, 125*1024*1024)

	decoder := charmap.ISO8859_1.NewDecoder()
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			decoderChunk, _, decodeErr := transform.Bytes(decoder, buffer[:n])
			if decodeErr != nil {
				fmt.Println("Error modifing fragment:", decodeErr)
				return
			}
			data := processTextToArray(string(decoderChunk),53,"UNIDAD VICTIMAS")

			db.Insert(data)

			// for i, row := range data {
			// 	fmt.Printf("Fila %d: %v\n", i, row)
			// }
			

			fmt.Println(data)
		}

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error al leer el archivo:", err)
			return
		}
	}
}
