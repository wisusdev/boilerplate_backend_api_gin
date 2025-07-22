package seeders

import "semita/core/database/generate_seeders"

func CreateSeederManager() *generate_seeders.SeederManager {
	manager := generate_seeders.NewSeederManager()
	
	manager.RegisterSeeder(NewRolesPermissionsSeeder())
	manager.RegisterSeeder(NewUsersSeeder())

	return manager
}
