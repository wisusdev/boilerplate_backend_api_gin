package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateUsersTable struct {
	generate_migrations.BaseMigration
}

func NewCreateUsersTable() *CreateUsersTable {
	return &CreateUsersTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_users_table",
			Timestamp: "2024_01_01_000001",
		},
	}
}

func (m *CreateUsersTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("users", func(table *schema.Blueprint) {
		table.Increments("id")
		table.String("first_name", 255)
		table.String("last_name", 255)
		table.String("username", 255).Unique().Index()
		table.String("avatar", 255).Nullable()
		table.String("language", 10).Default("en")
		table.String("email", 255).Unique()
		table.DateTime("email_verified_at").Nullable()
		table.RememberToken()
		table.String("password", 255)
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateUsersTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}
