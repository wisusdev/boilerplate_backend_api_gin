package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/database/migrations"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

// Función auxiliar para conectar y registrar migraciones
func withMigrator(action func(migrator *generate_migrations.Migrator)) {
	fmt.Println("🔌 Conectando a la base de datos...")
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	fmt.Println("🛠️  Inicializando migrator...")
	migrator := generate_migrations.NewMigrator(db)

	fmt.Println("📝 Registrando migraciones...")
	migrator.Register(migrations.NewCreateUsersTable())
	migrator.Register(migrations.NewCreatePasswordResetsTable())
	migrator.Register(migrations.NewCreateOAuthClientsTable())
	migrator.Register(migrations.NewCreateOAuthTokensTable())
	migrator.Register(migrations.NewCreateOAuthScopesTable())
	migrator.Register(migrations.NewCreateRolesTable())
	migrator.Register(migrations.NewCreatePermissionsTable())
	migrator.Register(migrations.NewCreateUserRolesTable())
	migrator.Register(migrations.NewCreateRolePermissionsTable())
	migrator.Register(migrations.NewCreateUserPermissionsTable())

	fmt.Println("🚀 Ejecutando acción del migrator...")
	action(migrator)
}

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Ejecuta las migraciones de base de datos",
	Run: func(cmd *cobra.Command, args []string) {
		withMigrator(func(migrator *generate_migrations.Migrator) {
			if err := migrator.Migrate(); err != nil {
				log.Fatal("Error running database:", err)
			}
			fmt.Println("Migrations completed successfully!")
		})
	},
}

var MigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Elimina y vuelve a crear todas las tablas",
	Run: func(cmd *cobra.Command, args []string) {
		withMigrator(func(migrator *generate_migrations.Migrator) {
			if err := migrator.Fresh(); err != nil {
				log.Fatal("Error refreshing database:", err)
			}
		})
	},
}

var MigrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Revierte la última migración",
	Run: func(cmd *cobra.Command, args []string) {
		withMigrator(func(migrator *generate_migrations.Migrator) {
			if err := migrator.Rollback(); err != nil {
				log.Fatal("Error rolling back database:", err)
			}
			fmt.Println("Rollback completed successfully!")
		})
	},
}

var MakeMigrationCmd = &cobra.Command{
	Use:   "make:migration",
	Short: "Crea un archivo de migración nuevo",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		createMigrationFile(name)
		fmt.Println("Migration file created successfully!")
	},
}

func init() {
	MigrateCmd.AddCommand(MigrateFreshCmd)
	MigrateCmd.AddCommand(MigrateRollbackCmd)
	MigrateCmd.AddCommand(MakeMigrationCmd)
}

// Copia la función createMigrationFile, toPascalCase y getTableName aquí desde tu código actual

func createMigrationFile(name string) {
	const migrationTpl = `package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
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
