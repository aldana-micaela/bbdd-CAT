package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// estructuras de datos
type Cliente struct {
	IDCliente       int    `json:"id_cliente"`
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	DNI             int    `json:"dni"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
}

type Operador struct {
	IDOperadore  int    `json:"id_operadore"`
	Nombre       string `json:"nombre"`
	Apellido     string `json:"apellido"`
	DNI          int    `json:"dni"`
	FechaIngreso string `json:"fecha_ingreso"`
	Disponible   bool   `json:"disponible"`
}

type DatosDePrueba struct {
	IDOrden               int    `json:"id_orden"`
	Operacion             string `json:"operacion"`
	IDCliente             int    `json:"id_cliente"`
	Id_cola_atencion      int    `json:"id_cola_atencion"`
	Tipo_tramite          string `json:"tipo_tramite"`
	Descripcion_tramite   string `json:"descripcion_tramite"`
	Id_tramite            int    `json:"id_tramite"`
	Estado_cierre_tramite string `json:"estado_cierre_tramite"`
	Respuesta_tramite     string `json:"respuesta_tramite"`
}

// estructuras de datos noSQL
type Tramite struct {
	IDTramite        int    `json:"id_tramite"`
	IDCliente        int    `json:"id_cliente"`
	ID_cola_atencion int    `json:"id_cola_atencion"`
	Tipo_tramite     string `json:"tipo_tramite"`
	F_inicio_gestion string `json:"f_inicio_gestion"`
	Descripcion      string `json:"descripcion"`
	F_fin_gestion    string `json:"f_fin_gestion"`
	Respuesta        string `json:"respuesta"`
	Estado           string `json:"estado"`
}

type Cola_atencion struct {
	ID_cola_atencion  int    `json:"id_cola_atencion"`
	IDCliente         int    `json:"id_cliente"`
	F_inicio_llamado  string `json:"f_inicio_llamado"`
	ID_operadore      int    `json:"id_operadore,omitempty"`
	F_inicio_atencion string `json:"f_inicio_atencion,omitempty"`
	F_fin_atencion    string `json:"f_fin_atencion,omitempty"`
	Estado            string `json:"estado"`
}

func CreateDatabase() {
	fmt.Println("Creando base de datos...")
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`drop database if exists baez_filiberto_nuñez_schillaci_db1;`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`create database baez_filiberto_nuñez_schillaci_db1;`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("La base de datos fue creada con éxito!")

}

func GetContentFromFile(filename string) string {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	content := string(bytes)
	return content
}

func CreateTables() {
	db := abrirConexion()

	fmt.Println("Creando tablas...")

	createTablesSql := GetContentFromFile("sql/create_tables.sql")

	_, err := db.Exec(createTablesSql)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Las tablas fueron creadas con éxito!")

	defer db.Close()
}

func AddPKsAndFks() {
	db := abrirConexion()

	fmt.Println("Agregando primary and foreign keys...")
	addConstraintsSql := GetContentFromFile("sql/add_pksandfks.sql")

	_, err := db.Exec(addConstraintsSql)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Primary and foreign keys fueron agregadas con éxito!")
	defer db.Close()

}

func DropPKsAndFks() {
	db := abrirConexion()
	fmt.Println("Eliminando primary and foreign keys...")
	dropConstraintsSql := GetContentFromFile("sql/drop_pksandfks.sql")

	_, err := db.Exec(dropConstraintsSql)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Primary and foreign keys fueron eliminadas con éxito!")
	defer db.Close()
}

func cargaClientes() []Cliente {
	file, _ := os.ReadFile("files/clientes.json")
	var clientes []Cliente
	err := json.Unmarshal(file, &clientes)
	if err != nil {
		log.Fatalf("JSON unmarshalling CLIENTE failed: %s", err)
		return clientes
	} else {
		fmt.Println("CLIENTES se cargo con exito")
		return clientes
	}

}

func cargaOperadores() []Operador {

	file, _ := os.ReadFile("files/operadores.json")
	var operadores []Operador
	err := json.Unmarshal(file, &operadores)
	if err != nil {
		log.Fatalf("JSON unmarshalling OPERADORE failed: %s", err)
		return operadores
	} else {
		fmt.Println("OPERADORES se cargo con exito")
		return operadores
	}

}

func cargaDatosDePrueba() []DatosDePrueba {
	file, _ := os.ReadFile("files/datos_de_prueba.json")
	var datos_prueba []DatosDePrueba
	err := json.Unmarshal(file, &datos_prueba)
	if err != nil {
		log.Fatalf("JSON unmarshalling DATOS DE PRUEBA failed: %s", err)
		return datos_prueba
	} else {
		fmt.Println("DATOS DE PRUEBA se cargo con exito")
		return datos_prueba
	}
}

func insertClientes(lista []Cliente, db *sql.DB) {

	query := `
			insert into cliente (id_cliente, nombre, apellido, dni, fecha_nacimiento, telefono, email)
			values ($1, $2, $3, $4, $5, $6, $7);`
	for _, cliente := range lista {

		_, err := db.Exec(query,
			cliente.IDCliente,
			cliente.Nombre,
			cliente.Apellido,
			cliente.DNI,
			cliente.FechaNacimiento,
			cliente.Telefono,
			cliente.Email)

		if err != nil {
			log.Fatalf("Error al insertar el cliente: %v", err)
		}
	}
}

func insertOperadores(lista []Operador, db *sql.DB) {
	query := `
			insert into operadore (id_operadore, nombre, apellido, dni, fecha_ingreso, disponible)
			values ($1, $2, $3, $4, $5, $6);`

	for _, operador := range lista {
		_, err := db.Exec(query,
			operador.IDOperadore,
			operador.Nombre,
			operador.Apellido,
			operador.DNI,
			operador.FechaIngreso,
			operador.Disponible)

		if err != nil {
			log.Fatalf("Error al insertar al operador: %v", err)
		}
	}

}

func insertDatosDePrueba(lista []DatosDePrueba, db *sql.DB) {
	query := `
			insert into datos_de_prueba (id_orden, operacion, id_cliente, id_cola_atencion, tipo_tramite, descripcion_tramite, id_tramite, estado_cierre_tramite, respuesta_tramite )
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	for _, dato := range lista {
		_, err := db.Exec(query,
			dato.IDOrden,
			dato.Operacion,
			dato.IDCliente,
			dato.Id_cola_atencion,
			dato.Tipo_tramite,
			dato.Descripcion_tramite,
			dato.Id_tramite,
			dato.Estado_cierre_tramite,
			dato.Respuesta_tramite)

		if err != nil {
			log.Fatalf("Error al insertar los datos de prueba: %v", err)
		}
	}
}

func LoadData() {
	db := abrirConexion()
	fmt.Println("Cargando datos...")

	clientes := cargaClientes()
	operadores := cargaOperadores()
	datosDePrueba := cargaDatosDePrueba()

	insertClientes(clientes, db)
	insertOperadores(operadores, db)
	insertDatosDePrueba(datosDePrueba, db)

	fmt.Println("Datos cargados con éxito!")
	defer db.Close()
}

func CreateSPsAndTRGs() {
	db := abrirConexion()
	fmt.Println("Creando stored procedures y triggers...")
	create_sp_trg := GetContentFromFile("sql/create_sps_trgs.sql")

	_, err := db.Exec(create_sp_trg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Stored Procedures y Triggers fueron agregadas con éxito!")
	defer db.Close()
}

func IniciarPruebas() {
	db := abrirConexion()
	fmt.Println("Iniciando pruebas...")
	_, err := db.Exec(`begin;
	set transaction isolation level serializable;
	
	-- Nuevo llamado id_cliente 21 y 4
	select ingresar_llamado(21);
	select ingresar_llamado(4);
	
	-- Atencion llamada 2 veces
	select atender_llamado();
	select atender_llamado();
	
	-- Nuevo llamado id_cliente 8, 12 y 16
	select ingresar_llamado(8);
	select ingresar_llamado(12);
	select ingresar_llamado(16);
	
	-- Baja llamado id_cola_atencion 2
	select desistir_llamado(2);
	select desistir_llamado(2);
	
	-- Nuevo llamado id_cliente 20
	select ingresar_llamado(20);
	
	-- Alta tramite
	select alta_de_tramite(1, 'consulta', '¿Es posible suspender temporalmente el servicio por 2 meses? (vacaciones)');
	select alta_de_tramite(1, 'reclamo', 'El monto de la ultima factura fue debitado dos veces en la tarjeta de credito');
	
	-- Atencion llamada 3 veces
	select atender_llamado();
	select atender_llamado();
	select atender_llamado();
	
	-- Fin llamado id_cola_atencion 1
	select finalizar_llamado(1);
	
	-- Atencion llamada
	select atender_llamado();
	
	-- Baja llamado id_cola_atencion 3
	select desistir_llamado(3);
	
	-- Cierre tramite
	select cierre_tramite(2, 'rechazado', 'Los dos cobros corresponden a facturas de meses diferentes');
	select cierre_tramite(1, 'solucionado', 'Es posible suspender el servicio, avisando con 20 dias de anticipacion');
	
	-- Fin llamados id_cola_atencion 4 y 5
	select finalizar_llamado(4);
	select finalizar_llamado(5);
	
	commit;`)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pruebas iniciadas con éxito!")
	defer db.Close()
}

func abrirConexion() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=baez_filiberto_nuñez_schillaci_db1 sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}
