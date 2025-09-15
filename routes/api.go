package routes

import (
	"semita/app/http/controllers/api/v1/auth"
	"semita/app/http/controllers/api/v1/base"
	"semita/app/http/middleware"

	"github.com/gin-gonic/gin"
)

func Api(router *gin.RouterGroup) {
	// Inicializar controladores
	roleController := &base.RoleController{}
	permissionController := &base.PermissionController{}
	userPermissionController := &base.UserPermissionController{}

	authRouter := router.Group("/auth")
	{
		// Auth routes
		authRouter.POST("/login", auth.Login)
		authRouter.POST("/register", auth.Register)
		authRouter.POST("/logout", middleware.AuthMiddleware(), auth.Logout)
		authRouter.POST("/forgot-password", auth.ForgotPassword)
		authRouter.POST("/reset-password", auth.ResetPassword)
		authRouter.POST("/email/resend", middleware.AuthMiddleware(), auth.ResendEmailVerify)
		authRouter.GET("/email/verify/:id/:hash", auth.VerifyEmail)
		authRouter.POST("/refresh-token", middleware.AuthMiddleware(), auth.RefreshToken)
	}

	// Rutas protegidas con autenticación
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(), middleware.ValidateJSONAPIDocument())
	{
		// Rutas de roles
		roles := protected.Group("/roles")
		{
			roles.GET("/", roleController.Index)
			roles.GET("/:id", roleController.Show)
			roles.POST("/", middleware.RequirePermission("create-roles"), roleController.Store)
			roles.PUT("/:id", middleware.RequirePermission("edit-roles"), roleController.Update)
			roles.DELETE("/:id", middleware.RequirePermission("delete-roles"), roleController.Delete)
			roles.POST("/assign-user", middleware.RequirePermission("assign-roles"), roleController.AssignToUser)
			roles.POST("/revoke-user", middleware.RequirePermission("assign-roles"), roleController.RevokeFromUser)
			roles.GET("/user/:user_id", roleController.GetUserRoles)
		}

		// Rutas de permisos
		permissions := protected.Group("/permissions")
		{
			permissions.GET("/", permissionController.Index)
			permissions.GET("/:id", permissionController.Show)
			permissions.POST("/", middleware.RequirePermission("create-permissions"), permissionController.Store)
			permissions.PUT("/:id", middleware.RequirePermission("edit-permissions"), permissionController.Update)
			permissions.DELETE("/:id", middleware.RequirePermission("delete-permissions"), permissionController.Delete)
			permissions.POST("/assign-user", middleware.RequirePermission("assign-permissions"), permissionController.AssignToUser)
			permissions.POST("/assign-role", middleware.RequirePermission("assign-permissions"), permissionController.AssignToRole)
			permissions.POST("/revoke-user", middleware.RequirePermission("assign-permissions"), permissionController.RevokeFromUser)
			permissions.POST("/revoke-role", middleware.RequirePermission("assign-permissions"), permissionController.RevokeFromRole)
			permissions.GET("/user/:user_id", permissionController.GetUserPermissions)
			permissions.GET("/role/:role_id", permissionController.GetRolePermissions)
		}

		// Rutas de verificación de permisos
		userPerms := protected.Group("/user-permissions")
		{
			userPerms.GET("/user/:user_id", userPermissionController.CheckUserPermissions)
			userPerms.GET("/current-user", userPermissionController.CheckCurrentUserPermissions)
			userPerms.GET("/user/:user_id/check-role", userPermissionController.CheckRole)
			userPerms.GET("/user/:user_id/check-permission", userPermissionController.CheckPermission)
			userPerms.GET("/current-user/check-role", userPermissionController.CheckCurrentUserRole)
			userPerms.GET("/current-user/check-permission", userPermissionController.CheckUserPermissions)
		}
	}
}
