# Sistema de Roles y Permisos - Documentación

Este sistema de roles y permisos está inspirado en el paquete `spatie/laravel-permission` de Laravel pero implementado completamente en Go con Gin y SQL puro.

## Características

- ✅ Asignación de múltiples roles a usuarios
- ✅ Asignación de múltiples permisos a roles
- ✅ Asignación directa de permisos a usuarios
- ✅ Verificación de roles y permisos
- ✅ Middleware para proteger rutas
- ✅ Helpers para usar en controladores y vistas
- ✅ API completa para gestión
- ✅ Comando de seeding para datos iniciales

## Estructura de Base de Datos

### Tablas creadas

- `roles` - Almacena los roles del sistema
- `permissions` - Almacena los permisos del sistema  
- `user_roles` - Relación muchos a muchos entre usuarios y roles
- `role_permissions` - Relación muchos a muchos entre roles y permisos
- `user_permissions` - Permisos directos asignados a usuarios

## Comandos

### Ejecutar migraciones

```bash
go run main.go migrate
```

### Poblar la base de datos con roles y permisos iniciales

```bash
go run main.go seed:roles-permissions
```

Esto creará:

- **Roles**: super-admin, admin, editor, moderator, user
- **Permisos**: create-users, edit-users, delete-users, view-users, create-roles, edit-roles, etc.

## Uso en Middleware

### Verificar un rol específico

```go
router.GET("/admin", middleware.RequireRole("admin"), handler)
```

### Verificar cualquiera de varios roles

```go
router.GET("/editors", middleware.RequireAnyRole([]string{"admin", "editor"}), handler)
```

### Verificar que tenga todos los roles especificados

```go
router.GET("/special", middleware.RequireAllRoles([]string{"admin", "editor"}), handler)
```

### Verificar un permiso específico

```go
router.GET("/users", middleware.RequirePermission("view-users"), handler)
```

### Verificar cualquiera de varios permisos

```go
router.GET("/content", middleware.RequireAnyPermission([]string{"edit-posts", "publish-posts"}), handler)
```

### Verificar que tenga todos los permisos especificados

```go
router.GET("/manage", middleware.RequireAllPermissions([]string{"create-users", "edit-users"}), handler)
```

### Verificar rol O permiso

```go
router.GET("/settings", middleware.CheckRoleOrPermission("super-admin", "manage-settings"), handler)
```

## Uso en Controladores

### Importar helpers

```go
import "semita/app/helpers"
```

### Verificar roles

```go
func (c *Controller) MyHandler(ctx *gin.Context) {
    if !helpers.HasRole(ctx.Request, "admin") {
        // Usuario no es admin
        return
    }
    
    if helpers.HasAnyRole(ctx.Request, []string{"admin", "editor"}) {
        // Usuario es admin O editor
    }
    
    if helpers.IsUserAdmin(ctx.Request) {
        // Usuario es admin o super-admin
    }
}
```

### Verificar permisos

```go
func (c *Controller) MyHandler(ctx *gin.Context) {
    if !helpers.HasPermission(ctx.Request, "edit-users") {
        // Usuario no puede editar usuarios
        return
    }
    
    if helpers.CanManageUsers(ctx.Request) {
        // Usuario puede gestionar usuarios
    }
    
    if helpers.HasAllPermissions(ctx.Request, []string{"create-posts", "edit-posts"}) {
        // Usuario tiene ambos permisos
    }
}
```

### Obtener roles y permisos del usuario

```go
func (c *Controller) MyHandler(ctx *gin.Context) {
    userRoles, success := helpers.GetUserRoles(ctx.Request)
    if success {
        // userRoles es un []string con los nombres de los roles
    }
    
    userPermissions, success := helpers.GetUserPermissions(ctx.Request)
    if success {
        // userPermissions es un []string con los nombres de los permisos
    }
}
```

## API Endpoints

### Roles

```
GET    /api/roles                    - Listar todos los roles
GET    /api/roles/:id                - Obtener un rol específico
POST   /api/roles                    - Crear nuevo rol
PUT    /api/roles/:id                - Actualizar rol
DELETE /api/roles/:id                - Eliminar rol
POST   /api/roles/assign-user        - Asignar rol a usuario
POST   /api/roles/revoke-user        - Revocar rol de usuario
GET    /api/roles/user/:user_id      - Obtener roles de un usuario
```

### Permisos

```
GET    /api/permissions                    - Listar todos los permisos
GET    /api/permissions/:id                - Obtener un permiso específico
POST   /api/permissions                    - Crear nuevo permiso
PUT    /api/permissions/:id                - Actualizar permiso
DELETE /api/permissions/:id                - Eliminar permiso
POST   /api/permissions/assign-user        - Asignar permiso a usuario
POST   /api/permissions/assign-role        - Asignar permiso a rol
POST   /api/permissions/revoke-user        - Revocar permiso de usuario
POST   /api/permissions/revoke-role        - Revocar permiso de rol
GET    /api/permissions/user/:user_id      - Obtener permisos de un usuario
GET    /api/permissions/role/:role_id      - Obtener permisos de un rol
```

### Verificaciones de Usuario

```
GET    /api/user-permissions/user/:user_id                - Información completa del usuario
GET    /api/user-permissions/current-user                 - Información del usuario logueado
GET    /api/user-permissions/user/:user_id/check-role     - Verificar rol (query: role, guard)
GET    /api/user-permissions/user/:user_id/check-permission - Verificar permiso (query: permission, guard)
GET    /api/user-permissions/current-user/check-role      - Verificar rol del usuario logueado
GET    /api/user-permissions/current-user/check-permission - Verificar permiso del usuario logueado
```

## Ejemplos de Payloads API

### Crear rol

```json
POST /api/roles
{
    "name": "manager",
    "guard_name": "web",
    "description": "Manager role"
}
```

### Crear permiso

```json
POST /api/permissions
{
    "name": "view-reports",
    "guard_name": "web", 
    "description": "View system reports"
}
```

### Asignar rol a usuario

```json
POST /api/roles/assign-user
{
    "user_id": 1,
    "role_id": 2
}
```

### Asignar permiso a rol

```json
POST /api/permissions/assign-role
{
    "role_id": 1,
    "permission_id": 3
}
```

### Asignar permiso directo a usuario

```json
POST /api/permissions/assign-user
{
    "user_id": 1,
    "permission_id": 3
}
```

## Modelos Disponibles

### Roles

```go
// Obtener todos los roles
roles, err := models.GetAllRoles()

// Obtener rol por ID
role, err := models.GetRoleByID(1)

// Obtener rol por nombre
role, err := models.GetRoleByName("admin", "web")

// Crear rol
roleData := structs.CreateRoleStruct{
    Name: "manager",
    GuardName: "web",
    Description: "Manager role",
}
role, err := models.CreateRole(roleData)

// Asignar rol a usuario
err := models.AssignRoleToUser(userID, roleID)

// Verificar si usuario tiene rol
hasRole, err := models.UserHasRoleByName(userID, "admin", "web")
```

### Permisos

```go
// Obtener todos los permisos
permissions, err := models.GetAllPermissions()

// Crear permiso
permissionData := structs.CreatePermissionStruct{
    Name: "view-reports",
    GuardName: "web",
    Description: "View reports",
}
permission, err := models.CreatePermission(permissionData)

// Asignar permiso a rol
err := models.AssignPermissionToRole(roleID, permissionID)

// Verificar si usuario tiene permiso
hasPermission, err := models.UserHasPermission(userID, "view-reports", "web")

// Obtener todos los permisos de un usuario (directos + heredados)
permissions, err := models.GetUserAllPermissions(userID)
```

## Guards

El sistema soporta múltiples "guards" (contextos de seguridad). Por defecto usa "web", pero puedes crear guards como "api", "admin", etc.

```go
// Verificar rol en guard específico
hasRole := helpers.HasRole(request, "admin", "api")

// Crear rol con guard específico
roleData := structs.CreateRoleStruct{
    Name: "api-admin",
    GuardName: "api",
}
```

## Rutas Web de Administración

```
GET /admin                      - Dashboard (requiere view-dashboard)
GET /admin/users               - Lista de usuarios (requiere view-users)
GET /admin/users/:id           - Detalle de usuario (requiere view-users)
GET /admin/roles               - Lista de roles (requiere view-roles)
GET /admin/permissions         - Lista de permisos (requiere view-permissions)
```

## Helpers Disponibles

### Verificaciones básicas

- `HasRole(request, roleName, guardName)` - Verificar rol específico
- `HasAnyRole(request, roleNames, guardName)` - Verificar cualquier rol
- `HasAllRoles(request, roleNames, guardName)` - Verificar todos los roles
- `HasPermission(request, permissionName, guardName)` - Verificar permiso específico
- `HasAnyPermission(request, permissionNames, guardName)` - Verificar cualquier permiso
- `HasAllPermissions(request, permissionNames, guardName)` - Verificar todos los permisos

### Shortcuts útiles

- `IsUserAdmin(request)` - Es admin o super-admin
- `IsUserSuperAdmin(request)` - Es super-admin
- `CanManageUsers(request)` - Puede gestionar usuarios
- `CanManageRoles(request)` - Puede gestionar roles
- `CanManagePermissions(request)` - Puede gestionar permisos
- `CanAccessDashboard(request)` - Puede acceder al dashboard

### Obtener información

- `GetUserRoles(request)` - Obtener roles del usuario
- `GetUserPermissions(request)` - Obtener permisos del usuario

## Jerarquía de Roles por Defecto

1. **super-admin**: Todos los permisos
2. **admin**: Permisos administrativos (usuarios, roles, configuración)
3. **editor**: Permisos de contenido (posts, dashboard)
4. **moderator**: Permisos básicos de moderación
5. **user**: Sin permisos especiales

## Extensibilidad

El sistema es completamente extensible:

1. **Agregar nuevos permisos**: Crear en la base de datos o via API
2. **Crear nuevos roles**: Via comando seed o API
3. **Middleware personalizado**: Combinar verificaciones existentes
4. **Guards personalizados**: Para diferentes contextos (web, api, admin)
5. **Helpers personalizados**: Crear verificaciones específicas para tu aplicación

## Seguridad

- Todas las consultas usan SQL preparado (prevención de SQL injection)
- Verificaciones en múltiples capas (middleware, controladores, vistas)
- Sistema de guards para diferentes contextos
- Roles y permisos con timestamps para auditoría
- Cascada de eliminación para mantener integridad referencial
