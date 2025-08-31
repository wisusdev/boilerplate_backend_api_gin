package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreatePermissionTablesTable struct {
	generate_migrations.BaseMigration
}

func NewCreatePermissionTablesTable() *CreatePermissionTablesTable {
	return &CreatePermissionTablesTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_permission_tables",
			Timestamp: "2024_01_12_170021",
		},
	}
}

func (m *CreatePermissionTablesTable) Up(db database_connections.SQLAdapter) error {
	schemaBuilder := schema.NewSchema()

	// Crear tabla permissions
	sqlPermissions := schemaBuilder.Create("permissions", func(table *schema.Blueprint) {
		table.UuidPrimary("uuid")
		table.String("name", 255)
		table.String("guard_name", 255)
		table.Timestamps()
		table.Unique([]string{"name", "guard_name"})
	})

	if _, err := db.Exec(sqlPermissions); err != nil {
		return err
	}

	// Crear tabla roles
	sqlRoles := schemaBuilder.Create("roles", func(table *schema.Blueprint) {
		table.UuidPrimary("uuid")
		table.String("name", 255)
		table.String("guard_name", 255)
		table.Timestamps()
		table.Unique([]string{"name", "guard_name"})
	})

	if _, err := db.Exec(sqlRoles); err != nil {
		return err
	}

	// Crear tabla model_has_permissions
	sqlModelHasPermissions := schemaBuilder.Create("model_has_permissions", func(table *schema.Blueprint) {
		table.Uuid("permission_id")
		table.String("model_type", 255)
		table.Uuid("model_id")
		table.Index("model_id")
		table.Primary([]string{"permission_id", "model_id", "model_type"})
	})

	if _, err := db.Exec(sqlModelHasPermissions); err != nil {
		return err
	}

	// Crear tabla model_has_roles
	sqlModelHasRoles := schemaBuilder.Create("model_has_roles", func(table *schema.Blueprint) {
		table.Uuid("role_id")
		table.String("model_type", 255)
		table.Uuid("model_id")
		table.Index("model_id")
		table.Primary([]string{"role_id", "model_id", "model_type"})
	})

	if _, err := db.Exec(sqlModelHasRoles); err != nil {
		return err
	}

	// Crear tabla role_has_permissions
	sqlRoleHasPermissions := schemaBuilder.Create("role_has_permissions", func(table *schema.Blueprint) {
		table.Uuid("permission_id")
		table.Uuid("role_id")
		table.Primary([]string{"permission_id", "role_id"})
	})

	if _, err := db.Exec(sqlRoleHasPermissions); err != nil {
		return err
	}

	return nil
}

func (m *CreatePermissionTablesTable) Down(db database_connections.SQLAdapter) error {
	tables := []string{
		"role_has_permissions",
		"model_has_roles",
		"model_has_permissions",
		"roles",
		"permissions",
	}

	for _, table := range tables {
		if _, err := db.Exec("DROP TABLE IF EXISTS " + table); err != nil {
			return err
		}
	}

	return nil
}
