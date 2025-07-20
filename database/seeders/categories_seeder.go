package seeders

import (
	"log"
	"semita/core/database/database_connections"
	"semita/core/database/seeders"
)

// CategoriesSeeder seeder para categorías
type CategoriesSeeder struct {
	seeders.BaseSeeder
}

// NewCategoriesSeeder crea una nueva instancia del seeder
func NewCategoriesSeeder() *CategoriesSeeder {
	return &CategoriesSeeder{
		BaseSeeder: seeders.BaseSeeder{
			DB:   database_connections.DatabaseConnectSQL(),
			Name: "categories_seeder",
		},
	}
}

// GetName retorna el nombre del seeder
func (cs *CategoriesSeeder) GetName() string {
	return cs.BaseSeeder.Name
}

// GetDependencies retorna las dependencias del seeder
func (cs *CategoriesSeeder) GetDependencies() []string {
	return []string{} // No tiene dependencias
}

// GetTables retorna las tablas que maneja este seeder
func (cs *CategoriesSeeder) GetTables() []string {
	return []string{"categories"}
}

// Seed ejecuta el seeding de categorías
func (cs *CategoriesSeeder) Seed() error {
	log.Println("Seeding categories...")

	categories := []struct {
		Name        string
		Description string
		Slug        string
	}{
		{
			Name:        "Tecnología",
			Description: "Artículos sobre tecnología, programación y desarrollo",
			Slug:        "tecnologia",
		},
		{
			Name:        "Ciencia",
			Description: "Contenido científico y descubrimientos",
			Slug:        "ciencia",
		},
		{
			Name:        "Deportes",
			Description: "Noticias y análisis deportivos",
			Slug:        "deportes",
		},
		{
			Name:        "Cultura",
			Description: "Arte, música, literatura y entretenimiento",
			Slug:        "cultura",
		},
		{
			Name:        "Negocios",
			Description: "Economía, finanzas y mundo empresarial",
			Slug:        "negocios",
		},
		{
			Name:        "Salud",
			Description: "Bienestar, medicina y vida saludable",
			Slug:        "salud",
		},
		{
			Name:        "Educación",
			Description: "Recursos educativos y aprendizaje",
			Slug:        "educacion",
		},
		{
			Name:        "Viajes",
			Description: "Destinos, experiencias y guías de viaje",
			Slug:        "viajes",
		},
	}

	for _, category := range categories {
		// Crear nueva categoría directamente (ya se limpiaron los datos)
		insertQuery := `
			INSERT INTO categories (name, description, slug, created_at, updated_at) 
			VALUES (?, ?, ?, NOW(), NOW())`

		result, err := cs.BaseSeeder.DB.Exec(insertQuery, category.Name, category.Description, category.Slug)
		if err != nil {
			log.Printf("Error creating category '%s': %v", category.Name, err)
			continue
		}

		id, _ := result.LastInsertId()
		log.Printf("Created category: %s (ID: %d)", category.Name, id)
	}

	log.Println("Categories seeding completed successfully!")
	return nil
}
