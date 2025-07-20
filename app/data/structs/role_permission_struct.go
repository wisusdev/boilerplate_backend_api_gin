package structs

// Role struct representa un rol en el sistema
type RoleStruct struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	GuardName   string `json:"guard_name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Permission struct representa un permiso en el sistema
type PermissionStruct struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	GuardName   string `json:"guard_name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// RoleWithPermissions representa un rol con sus permisos asociados
type RoleWithPermissions struct {
	RoleStruct
	Permissions []PermissionStruct `json:"permissions"`
}

// UserWithRolesAndPermissions representa un usuario con sus roles y permisos
type UserWithRolesAndPermissions struct {
	UserStruct
	Roles             []RoleStruct       `json:"roles"`
	DirectPermissions []PermissionStruct `json:"direct_permissions"`
	AllPermissions    []PermissionStruct `json:"all_permissions"`
}

// CreateRoleStruct para crear nuevos roles
type CreateRoleStruct struct {
	Name        string `json:"name" binding:"required"`
	GuardName   string `json:"guard_name"`
	Description string `json:"description"`
}

// CreatePermissionStruct para crear nuevos permisos
type CreatePermissionStruct struct {
	Name        string `json:"name" binding:"required"`
	GuardName   string `json:"guard_name"`
	Description string `json:"description"`
}

// AssignRoleRequest para asignar roles a usuarios
type AssignRoleRequest struct {
	UserID int `json:"user_id" binding:"required"`
	RoleID int `json:"role_id" binding:"required"`
}

// AssignPermissionRequest para asignar permisos a usuarios o roles
type AssignPermissionRequest struct {
	PermissionID int `json:"permission_id" binding:"required"`
	UserID       int `json:"user_id,omitempty"`
	RoleID       int `json:"role_id,omitempty"`
}

// RolePermissionCheck para verificaciones de permisos
type RolePermissionCheck struct {
	HasRole       bool     `json:"has_role"`
	HasPermission bool     `json:"has_permission"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
}
