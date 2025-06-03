package main

import (
	"fmt"

	_ "github.com/lib/pq"
)

func main() {

	for {
		printMenu()

		var option int
		fmt.Scanln(&option)

		fmt.Println("")

		switch option {
		case 1:
			CreateDatabase()
		case 2:
			CreateTables()
		case 3:
			AddPKsAndFks()
		case 4:
			DropPKsAndFks()
		case 5:
			LoadData()
		case 6:
			CreateSPsAndTRGs()
		case 7:
			IniciarPruebas()
		case 8:
			LoadDataBoltDB()
		case 0:
			return
		}
		fmt.Println("")
	}
}

func printMenu() {
	fmt.Println(` 
++++++++++++++++++++++++++++++++++++++++++++++++
+   1 → Crear base de datos                    +
+   2 → Crear tablas                           +
+   3 → Agregar PK's y Fk's                    +
+   4 → Eliminar PKs y FKs                     +
+   5 → Cargar datos                           +
+   6 → Crear stored procedures y triggers     +
+   7 → Iniciar pruebas                        +
+   8 → Cargar datos en BoltDB                 +
+   0 → salir                                  +
++++++++++++++++++++++++++++++++++++++++++++++++`)

	fmt.Println("Seleccione una opción...")
	fmt.Println("")
}
