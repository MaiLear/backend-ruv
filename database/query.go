package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// func createTransaction(){

// }

func toInterface(sliceString [][]string) [][]interface{} {
	var convertedInterface [][]interface{}
	for _, row := range sliceString {
		var interfaceRow []interface{}
		for _, col := range row {
			interfaceRow = append(interfaceRow, col)
		}
		convertedInterface = append(convertedInterface, interfaceRow)
	}

	return convertedInterface

}

func Insert(rows [][]string) {
	dbpool, err := Connect("postgres://maicol:maicol@localhost:5432/comidas")
	if err != nil {
		log.Fatal(err)
	}

	//Iniciar una transacción
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		log.Fatal("Failed to begin transaction")
	}
	//Asegura la reversion en caso de error
	defer tx.Rollback(context.Background())

	rowsInterface := toInterface(rows)

	_, err = tx.CopyFrom(context.Background(), pgx.Identifier{"ruv_victimas"}, []string{"ORIGEN", "FUENTE", "PROGRAMA", "ID_PERSONA", "ID_HOGAR", "TIPO_DOCUMENTO", "DOCUMENTO", "PRIMERNOMBRE", "SEGUNDONOMBRE", "PRIMERAPELLIDO", "SEGUNDOAPELLIDO", "FECHANACIMIENTO", "EXPEDICIONDOCUMENTO", "FECHAEXPEDICIONDOCUMENTO", "PERTENENCIAETNICA", "GENERO", "TIPOHECHO", "HECHO", "FECHAOCURRENCIA", "CODDANEMUNICIPIOOCURRENCIA", "ZONAOCURRENCIA", "UBICACIONOCURRENCIA", "PRESUNTOACTOR", "PRESUNTOVICTIMIZANTE", "FECHAREPORTE", "TIPOPOBLACION", "TIPOVICTIMA", "PAIS", "CIUDAD", "CODDANEMUNICIPIORESIDENCIA", "ZONARESIDENCIA", "UBICACIONRESIDENCIA", "DIRECCION", "NUMTELEFONOFIJO", "NUMTELEFONOCELULAR", "EMAIL", "FECHAVALORACION", "ESTADOVICTIMA", "NOMBRECOMPLETO", "IDSINIESTRO", "IDMIJEFE", "TIPODESPLAZAMIENTO", "REGISTRADURIA", "VIGENCIADOCUMENTO", "CONSPERSONA", "RELACION", "CODDANEDECLARACION", "CODDANELLEGADA", "CODIGOHECHO", "DISCAPACIDAD", "DESCRIPCIONDISCAPACIDAD", "FUD_FICHA","AFECTACIONES"}, pgx.CopyFromRows(rowsInterface))

	if err != nil {
		log.Fatalf("CopyFrom failed: %v\n", err)
		tx.Rollback(context.Background())
		return
	}

	// Confirmar la transacción si no hubo errores
	if err := tx.Commit(context.Background()); err != nil {
		log.Fatalf("Transaction commit failed: %v\n", err)
	}

	fmt.Println("Rows copied successfully")
	//Cerrar el grupo de conexiones
	defer dbpool.Close()
}
