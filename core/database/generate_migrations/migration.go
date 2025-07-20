package generate_migrations

import "semita/core/database/database_connections"

// Migration interface que define los métodos que debe implementar cada migración
type Migration interface {
	Up(db database_connections.SQLAdapter) error
	Down(db database_connections.SQLAdapter) error
	GetName() string
	GetTimestamp() string
}

// BaseMigration estructura base que pueden embeber las migraciones
type BaseMigration struct {
	Name      string
	Timestamp string
}

func (m *BaseMigration) GetName() string {
	return m.Name
}

func (m *BaseMigration) GetTimestamp() string {
	return m.Timestamp
}
