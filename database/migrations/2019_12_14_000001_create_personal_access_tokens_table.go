package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreatePersonalAccessTokensTable struct {
	generate_migrations.BaseMigration
}

func NewCreatePersonalAccessTokensTable() *CreatePersonalAccessTokensTable {
	return &CreatePersonalAccessTokensTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_personal_access_tokens_table",
			Timestamp: "2019_12_14_000001",
		},
	}
}

func (m *CreatePersonalAccessTokensTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("personal_access_tokens", func(table *schema.Blueprint) {
		table.Id()
		table.Morphs("tokenable")
		table.String("name", 255)
		table.String("token", 64).Unique()
		table.Text("abilities").Nullable()
		table.Timestamp("last_used_at").Nullable()
		table.Timestamp("expires_at").Nullable()
		table.Timestamps()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreatePersonalAccessTokensTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS personal_access_tokens")
	return err
}
