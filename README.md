# ☎️ Sistema de Base de Datos - Centro de Atención Telefónica (CAT)

Este proyecto implementa un sistema de gestión de datos para un **Centro de Atención Telefónica (CAT)**. Permite almacenar y manipular información sobre clientes, operadores, trámites, colas de atención y rendimiento del personal mediante una base de datos relacional y el lenguaje de programación Go.

## 🛠️ Tecnologías utilizadas

- **Go (Golang)** – Lógica del sistema y manejo de archivos.
- **PostgreSQL** – Base de datos relacional para la persistencia estructurada.
- **BoltDB** – Alternativa de base de datos embebida usada en Go.
- **JSON** – Fuente de datos para cargar entidades (clientes, operadores, etc.).
- **SQL Scripts** – Creación de tablas, claves foráneas, triggers y procedimientos almacenados.

## 🧠 Tecnologías y conceptos aplicados
- Lenguaje: Go (main.go, database.go, boltdb.go)
- Manejo de archivos JSON – Carga de información desde archivos clientes.json, operadores.json, etc.
- Base de datos relacional – Scripts SQL para:
- Crear tablas (create_tables.sql)
- Agregar PKs y FKs (add_pksandfks.sql)
- Procedimientos almacenados y triggers (create_sps_trgs.sql)
- Eliminación estructurada (drop_pksandfks.sql)
- Persistencia – Implementación de acceso a datos tanto con PostgreSQL como con BoltDB (base de datos embebida en Go)
- Modularización en Go – Separación clara por capas (main, database, etc.)
- Uso de go.mod y go.sum – Gestión de dependencias Go Modules

## 📂 Estructura del proyecto

```text
bbdd-CAT-main/
│
├── main.go                      # Entrada principal de la aplicación
├── database.go                  # Lógica para conectar con la base de datos
├── boltdb.go                    # Lógica alternativa con BoltDB
├── go.mod / go.sum              # Módulos y dependencias
│
├── files/                       # Archivos JSON de entrada
│   ├── clientes.json
│   ├── operadores.json
│   ├── datos_de_prueba.json
│   └── archivos_json.xlsx       # Documento de apoyo
│
└── sql/                         # Scripts SQL para PostgreSQL
    ├── create_tables.sql
    ├── add_pksandfks.sql
    ├── create_sps_trgs.sql
    └── drop_pksandfks.sql
```

## 📌 Funcionalidades principales
Carga de datos desde archivos JSON.
Creación de estructura de base de datos con SQL.
Inserción y gestión de registros mediante Go.
Opcionalidad de usar base embebida (BoltDB) o externa (PostgreSQL).
Modularización del sistema para mantenerlo escalable y mantenible.

## 👩‍💻 Autoría
Aldana Micaela Filiberto y Juliana nuñez
Estudiantes de Licenciatura en Sistemas.

---

Este proyecto fue desarrollado con fines académicos para integrar conocimientos de bases de datos, lenguaje Go y estructuras de datos reales en un entorno de simulación de atención telefónica.
