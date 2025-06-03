package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

func CrearConexionNoSQL() (*bolt.DB, error) {
	db, err := bolt.Open("baez_filiberto_nuñez_schillaci_db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
	//abre transaccion de escritura
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

	err = b.Put(key, val)
	if err != nil {
		return err
	}

	//cierra transaccion
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func ReadUnique(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
	var buf []byte

	//abre una transaccion de lectura
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf = b.Get(key)
		return nil
	})

	return buf, err
}

func InsertarClientesNoSQL(db *bolt.DB) {
	clientes := []Cliente{
		{1, "Ken", "Thompson", 5153057, "1995-05-05", "15-2889-7948", "ken@thompson.org"},
		{2, "Dennis", "Ritchie", 25610126, "1955-04-11", "15-7811-5045", "dennis@ritchie.org"},
		{3, "Donald", "Knuth", 9168297, "1984-04-05", "15-2780-6005", "don@knuth.org"},
		{4, "Rob", "Pike", 4915593, "1946-08-16", "15-1114-9719", "rob@pike.org"},
		{5, "Douglas", "McIlroy", 33187055, "1939-06-09", "15-9625-0245", "douglas@mcilroy.org"},
		{6, "Brian", "Kernighan", 13897948, "1992-11-22", "15-6410-6066", "brian@kernighan.org"},
		{7, "Bill", "Joy", 34115045, "1954-02-04", "15-4215-8655", "bill@joy.org"},
		{8, "Marshall Kirk", "McKusick", 9806005, "1995-12-27", "15-5197-4379", "marshall_kirk@mckusick.org"},
		{9, "Theo", "de Raadt", 5149719, "1950-02-07", "15-6470-9444", "theo@deraadt.org"},
		{10, "Cristina", "Kirchner", 6250245, "1990-08-17", "15-5291-0113", "cfk@fpv.gov.ar"},
		{11, "Diego", "Maradona", 19158655, "1985-02-27", "15-3361-4854", "diego@dios.com.ar"},
		{12, "Martín", "Palermo", 5974379, "1918-06-09", "15-9877-3169", "martin@palermo.com.ar"},
		{13, "Guillermo", "Barros Schelotto", 3910113, "1982-05-03", "15-5020-5695", "guille@melli.com.ar"},
		{14, "Susú", "Pecoraro", 7547862, "1935-04-03", "15-6695-9505", "susu@pecoraro.com.ar"},
		{15, "Norma", "Aleandro", 26614854, "1992-03-18", "15-9155-4115", "norma@aleandro.com.ar"},
		{16, "Soledad", "Silveyra", 7773169, "1957-07-28", "15-9184-4522", "sole@silveyra.com.ar"},
		{17, "Libertad", "Lamarque", 32205695, "1971-03-07", "15-6363-9690", "libertad@lamarque.com.ar"},
		{18, "Ana María", "Picchio", 19020903, "1946-08-06", "15-4819-2117", "ana.maria@picchio.com.ar"},
		{19, "Niní", "Marshall", 10535508, "1951-09-07", "15-9799-6045", "nini@marshall.com"},
		{20, "Claudia", "Lapacó", 30934609, "1961-08-03", "15-2005-4879", "claudia@lapaco.com.ar"},
	}

	for _, cliente := range clientes {
		data, err := json.Marshal(cliente)
		if err != nil {
			log.Fatal(err)
		}

		err = CreateUpdate(db, "cliente", []byte(strconv.Itoa(cliente.IDCliente)), data)
		if err != nil {
			log.Fatal(err)
		}

		resultado, err := ReadUnique(db, "cliente", []byte(strconv.Itoa(cliente.IDCliente)))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", resultado)
	}
}

func InsertarOperadoresNoSQL(db *bolt.DB) {
	operadores := []Operador{
		{1, "Wilhelm", "Strinitz", 5054058, "2018-05-14", true},
		{2, "Emanuel", "Lasker", 24610127, "2018-12-24", true},
		{3, "Jose Raul", "Capablanca", 9068298, "2019-11-19", true},
	}

	for _, operador := range operadores {
		data, err := json.Marshal(operador)
		if err != nil {
			log.Fatal(err)
		}

		err = CreateUpdate(db, "operador", []byte(strconv.Itoa(operador.IDOperadore)), data)
		if err != nil {
			log.Fatal(err)
		}

		resultado, err := ReadUnique(db, "operador", []byte(strconv.Itoa(operador.IDOperadore)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resultado)
	}
}

func InsertarTramites(db *bolt.DB) {
	tramites := []Tramite{
		{2, 4, 1, "reclamo", "2024-11-25 19:01:24", "El monto de la última factura fue debitado dos veces en la tarjeta de crédito", "2024-11-25 19:01:24", "Los dos cobros corresponden a facturas de meses diferentes", "rechazado"},
		{1, 4, 1, "consulta", "2024-11-25 19:01:24", "¿Es posible suspender temporalmente el servicio por 2 meses? (vacaciones)", "2024-11-25 19:01:24", "Es posible suspender el servicio, avisando con 20 dias de anticipación", "solucionado"},
	}

	for _, tramite := range tramites {
		data, err := json.Marshal(tramite)
		err = CreateUpdate(db, "tramite", []byte(strconv.Itoa(tramite.IDTramite)), data)
		if err != nil {
			log.Fatal(err)
		}

		resultado, err := ReadUnique(db, "tramite", []byte(strconv.Itoa(tramite.IDTramite)))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", resultado)
	}
}

func InsertarLlamados(db *bolt.DB) {
	llamados := []Cola_atencion{
		{2, 8, "2024-11-25 19:01:24", 0, "", "", "desistido"},
		{1, 4, "2024-11-25 19:01:24", 1, "2024-11-25 19:01:24", "2024-11-25 19:01:24", "finalizado"},
		{3, 12, "2024-11-25 19:01:24", 2, "2024-11-25 19:01:24", "2024-11-25 19:01:24", "desistido"},
		{4, 16, "2024-11-25 19:01:24", 3, "2024-11-25 19:01:24", "2024-11-25 19:01:24", "finalizado"},
		{5, 20, "2024-11-25 19:01:24", 1, "2024-11-25 19:01:24", "2024-11-25 19:01:24", "finalizado"},
	}

	for _, llamado := range llamados {
		data, err := json.Marshal(llamado)
		err = CreateUpdate(db, "llamado", []byte(strconv.Itoa(llamado.ID_cola_atencion)), data)
		if err != nil {
			log.Fatal(err)
		}

		resultado, err := ReadUnique(db, "llamado", []byte(strconv.Itoa(llamado.ID_cola_atencion)))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", resultado)
	}
}

func LoadDataBoltDB() {
	fmt.Printf("Cargando datos en BoltDB...\n")

	db, _ := CrearConexionNoSQL()
	defer db.Close()

	InsertarClientesNoSQL(db)
	InsertarOperadoresNoSQL(db)
	InsertarTramites(db)
	InsertarLlamados(db)

	fmt.Printf("Datos cargados con éxito en BoltDB!")
}
