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
			Timestamp: "2014_10_12_000000",
		},
	}
}

func (m *CreateUsersTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("users", func(table *schema.Blueprint) {
		table.UuidPrimary("id")
		table.String("first_name", 255)
		table.String("last_name", 255).Nullable()
		table.String("username", 255).Nullable().Unique()
		table.String("email", 255).Unique()
		table.DateTime("email_verified_at").Nullable()
		table.String("password", 255)
		table.String("avatar", 255).Nullable()
		table.String("language", 2).Default("es")
		table.RememberToken()
		table.Timestamps()
		table.SoftDeletes()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateUsersTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}
