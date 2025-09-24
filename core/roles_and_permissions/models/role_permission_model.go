package models

import dataStructs "boilerplate_backend_api_gin/app/data/models"

// Role struct representa un rol en el sistema
type RoleStruct struct {
	ID        string `json:"uuid"`
	Name      string `json:"name"`
	GuardName string `json:"guard_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Permission struct representa un permiso en el sistema
type PermissionStruct struct {
	ID        string `json:"uuid"`
	Name      string `json:"name"`
	GuardName string `json:"guard_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// RoleWithPermissions representa un rol con sus permisos asociados
type RoleWithPermissions struct {
	RoleStruct
	Permissions []PermissionStruct `json:"permissions"`
}

// UserWithRolesAndPermissions representa un usuario con sus roles y permisos
type UserWithRolesAndPermissions struct {
	dataStructs.UserStruct
	Roles             []RoleStruct       `json:"roles"`
	DirectPermissions []PermissionStruct `json:"direct_permissions"`
	AllPermissions    []PermissionStruct `json:"all_permissions"`
}

// CreateRoleStruct para crear nuevos roles
type CreateRoleStruct struct {
	Name        string   `json:"name" binding:"required"`
	GuardName   string   `json:"guard_name"`
	Permissions []string `json:"permissions,omitempty"`
}

// JsonApiRoleRequest estructura para manejar requests JSON API de roles
type JsonApiRoleRequest struct {
	Data JsonApiRoleData `json:"data" binding:"required"`
}

type JsonApiRoleData struct {
	Type       string                `json:"type" binding:"required"`
	ID         string                `json:"id,omitempty"`
	Attributes JsonApiRoleAttributes `json:"attributes" binding:"required"`
}

type JsonApiRoleAttributes struct {
	Name        string   `json:"name" binding:"required"`
	Permissions []string `json:"permissions,omitempty"`
}

// CreatePermissionStruct para crear nuevos permisos
type CreatePermissionStruct struct {
	Name      string `json:"name" binding:"required"`
	GuardName string `json:"guard_name"`
}

// AssignRoleRequest para asignar roles a usuarios
type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

// AssignPermissionRequest para asignar permisos a usuarios o roles
type AssignPermissionRequest struct {
	PermissionID string `json:"permission_id" binding:"required"`
	UserID       string `json:"user_id,omitempty"`
	RoleID       string `json:"role_id,omitempty"`
}

// RolePermissionCheck para verificaciones de permisos
type RolePermissionCheck struct {
	HasRole       bool     `json:"has_role"`
	HasPermission bool     `json:"has_permission"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
}
