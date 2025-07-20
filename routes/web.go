package routes

import (
	"semita/app/http/controllers/web"
	"semita/app/http/middleware"

	"github.com/gin-gonic/gin"
)

func Web() *gin.Engine {
	router := gin.Default()

	// IMPORTANTE: El middleware debe estar ANTES de todas las rutas
	router.Use(middleware.MethodOverride(), middleware.LanguageMiddleware())

	// Ahora define todas las rutas
	router.GET("/", web.HomeIndex)

	// Auth routes
	router.GET("/auth/login", middleware.RedirectGuest(web.AuthLogin))
	router.POST("/auth/login", middleware.RedirectGuest(web.AuthLoginPost))
	router.GET("/auth/logout", middleware.RequireAuth(web.AuthLogout))
	router.GET("/auth/register", middleware.RedirectGuest(web.AuthRegister))
	router.POST("/auth/register", middleware.RedirectGuest(web.AuthRegisterPost))
	router.GET("/auth/forgot-password", web.AuthForgotPassword)
	router.POST("/auth/forgot-password", web.AuthForgotPasswordPost)
	router.GET("/auth/reset-password", web.AuthResetPassword)
	router.POST("/auth/reset-password", web.AuthResetPasswordPost)

	router.GET("/users", middleware.RequireAuth(web.UserIndex))
	router.GET("/users/create", middleware.RequireAuth(web.UserCreate))
	router.POST("/users/store", middleware.RequireAuth(web.UserStore))
	router.GET("/users/show/:id", middleware.RequireAuth(web.UserShow))
	router.GET("/users/edit/:id", middleware.RequireAuth(web.UserEdit))
	router.POST("/users/update/:id", middleware.RequireAuth(web.UserUpdate))
	router.POST("/users/delete/:id", middleware.RequireAuth(web.UserDelete))

	// Inicializar controlador administrativo
	adminController := &web.AdminController{}

	// Rutas administrativas protegidas con roles y permisos
	admin := router.Group("/admin")
	admin.Use(middleware.RequireAuth(func(c *gin.Context) { c.Next() }))
	{
		// Dashboard principal - requiere permiso para ver dashboard
		admin.GET("/", middleware.RequirePermission("view-dashboard"), adminController.Dashboard)

		// Gestión de usuarios - requiere permiso para ver usuarios
		users := admin.Group("/users")
		users.Use(middleware.RequirePermission("view-users"))
		{
			users.GET("/", adminController.UsersIndex)
			users.GET("/:id", adminController.UserShow)
		}

		// Gestión de roles - requiere permiso para ver roles
		roles := admin.Group("/roles")
		roles.Use(middleware.RequirePermission("view-roles"))
		{
			roles.GET("/", adminController.RolesIndex)
		}

		// Gestión de permisos - requiere permiso para ver permisos
		permissions := admin.Group("/permissions")
		permissions.Use(middleware.RequirePermission("view-permissions"))
		{
			permissions.GET("/", adminController.PermissionsIndex)
		}

		// Ejemplo avanzado - requiere ser admin
		admin.GET("/advanced", middleware.RequireRole("admin"), adminController.AdvancedPermissionExample)

		// Ejemplo con múltiples roles permitidos
		admin.GET("/editors-only", middleware.RequireAnyRole([]string{"admin", "editor", "super-admin"}), func(c *gin.Context) {
			c.HTML(200, "admin/editors.html", gin.H{
				"message": "Solo editores, admins y super-admins pueden ver esto",
			})
		})

		// Ejemplo con múltiples permisos requeridos
		admin.GET("/content-management", middleware.RequireAllPermissions([]string{"manage-posts", "edit-posts"}), func(c *gin.Context) {
			c.HTML(200, "admin/content.html", gin.H{
				"message": "Necesitas ambos permisos: manage-posts y edit-posts",
			})
		})

		// Ejemplo con rol O permiso
		admin.GET("/settings", middleware.CheckRoleOrPermission("super-admin", "manage-settings"), func(c *gin.Context) {
			c.HTML(200, "admin/settings.html", gin.H{
				"message": "Eres super-admin O tienes el permiso manage-settings",
			})
		})
	}

	// Ruta para cambiar idioma
	router.POST("/set-lang", middleware.SetLangHandler)

	return router
}
