package routes

import (
	"boilerplate_backend_api_gin/app/http/controllers/api/v1/auth"
	"boilerplate_backend_api_gin/app/http/controllers/api/v1/base"
	"boilerplate_backend_api_gin/app/http/middleware"

	"github.com/gin-gonic/gin"
)

func Api(router *gin.RouterGroup) {
	// Inicializar controladores
	roleController := &base.RoleController{}
	permissionController := &base.PermissionController{}
	userPermissionController := &base.UserPermissionController{}
	userController := &base.UserController{}

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
		// Rutas de permisos
		protected.GET("/permissions", permissionController.Index)

		// Rutas de roles
		protected.GET("/roles", roleController.Index)
		protected.GET("/roles/:id", roleController.Show)
		protected.POST("/roles", roleController.Store)
		protected.PATCH("/roles/:id", roleController.Update)
		protected.DELETE("/roles/:id", roleController.Delete)

		// Rutas de usuarios
		protected.GET("/users", userController.Index)
		protected.GET("/users/:id", userController.Show)
		protected.POST("/users", userController.Store)
		protected.PUT("/users/:id", userController.Update)
		protected.DELETE("/users/:id", userController.Delete)

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
