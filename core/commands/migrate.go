package commands

import (
	"boilerplate_backend_api_gin/core/cli"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"boilerplate_backend_api_gin/database/migrations"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var MigrateCmd = cli.Command{
	Name:        "migrate",
	Description: "Ejecuta las migraciones de base de datos",
	Execute: func(args []string) {
		migrations.WithMigrator(func(migrator *generate_migrations.Migrator) {
			if err := migrator.Migrate(); err != nil {
				log.Fatal("Error running database:", err)
			}
			fmt.Println("Migrations completed successfully!")
		})
	},
}

var MigrateFreshCmd = cli.Command{
	Name:        "migrate:fresh",
	Description: "Elimina y vuelve a crear todas las tablas",
	Execute: func(args []string) {
		migrations.WithMigrator(func(migrator *generate_migrations.Migrator) {
			if err := migrator.Fresh(); err != nil {
				log.Fatal("Error refreshing database:", err)
			}
		})
	},
}

var MigrateRollbackCmd = cli.Command{
	Name:        "migrate:rollback",
	Description: "Revierte la última migración",
	Execute: func(args []string) {
		migrations.WithMigrator(func(migrator *generate_migrations.Migrator) {
			if err := migrator.Rollback(); err != nil {
				log.Fatal("Error rolling back database:", err)
			}
			fmt.Println("Rollback completed successfully!")
		})
	},
}

var MakeMigrationCmd = cli.Command{
	Name:        "make:migration",
	Description: "Crea un archivo de migración nuevo",
	Execute: func(args []string) {
		name := args[0]
		createMigrationFile(name)
		fmt.Println("Migration file created successfully!")
	},
}

// Copia la función createMigrationFile, toPascalCase y getTableName aquí desde tu código actual

func createMigrationFile(name string) {
	const migrationTpl = `package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
)

type {{.StructName}} struct {
	generate_migrations.BaseMigration
}

func New{{.StructName}}() *{{.StructName}} {
	return &{{.StructName}}{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "{{.MigrationName}}",
			Timestamp: "{{.Timestamp}}",
		},
	}
}

func (m *{{.StructName}}) Up(db database_connections.SQLAdapter) error {
	query := {{backtick}}
		CREATE TABLE {{.TableName}} (
			id INT PRIMARY KEY AUTO_INCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	{{backtick}}
	_, err := db.Exec(query)
	return err
}

func (m *{{.StructName}}) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS {{.TableName}}")
	return err
}
`
	timestamp := time.Now().Format("2006_01_02_150405")
	safeName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	filename := timestamp + "_" + safeName + ".go"
	dir := filepath.Join("database", "migrations")
	fullpath := filepath.Join(dir, filename)
	structName := toPascalCase(safeName)
	tableName := getTableName(safeName)
	data := map[string]string{
		"StructName":    structName,
		"MigrationName": safeName,
		"Timestamp":     timestamp,
		"TableName":     tableName,
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	fPath, err := os.Create(fullpath)
	if err != nil {
		log.Fatalf("No se pudo crear el archivo de migración: %v", err)
	}
	defer fPath.Close()
	funcMap := template.FuncMap{
		"backtick": func() string { return "`" },
	}
	tmpl, err := template.New("migration").Funcs(funcMap).Parse(migrationTpl)
	if err != nil {
		log.Fatalf("No se pudo parsear la plantilla: %v", err)
	}
	if err := tmpl.Execute(fPath, data); err != nil {
		log.Fatalf("No se pudo escribir la plantilla: %v", err)
	}
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

func getTableName(safeName string) string {
	parts := strings.Split(safeName, "_")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "table" && i > 0 {
			return parts[i-1]
		}
	}
	return safeName
}
