package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"boilerplate_backend_api_gin/core/database/schema"
)

type CreateFailedJobsTable struct {
	generate_migrations.BaseMigration
}

func NewCreateFailedJobsTable() *CreateFailedJobsTable {
	return &CreateFailedJobsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_failed_jobs_table",
			Timestamp: "2019_08_19_000000",
		},
	}
}

func (m *CreateFailedJobsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("failed_jobs", func(table *schema.Blueprint) {
		table.Id()
		table.String("uuid", 255).Unique()
		table.Text("connection")
		table.Text("queue")
		table.Text("payload")
		table.Text("exception")
		table.Timestamp("failed_at").UseCurrent()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateFailedJobsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS failed_jobs")
	return err
}
