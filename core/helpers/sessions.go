package helpers

import (
	"boilerplate_backend_api_gin/config"
	"boilerplate_backend_api_gin/core/common/nulltypes"
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/internationalization"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/sessions"
)

// StringToNullString convierte un string a nulltypes.NullString
func StringToNullString(s string) nulltypes.NullString {
	if s == "" {
		return nulltypes.NullString{String: "", Valid: false}
	}
	return nulltypes.NullString{String: s, Valid: true}
}

// NullStringToString convierte un nulltypes.NullString a string
func NullStringToString(ns nulltypes.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

type UserSessionStruct struct {
	ID            string
	FirstName     string
	LastName      string
	Username      string
	Avatar        nulltypes.NullString
	Language      nulltypes.NullString
	Email         string
	Authenticated bool
}

type AuthSessionStruct struct {
	User            UserSessionStruct
	IsAuthenticated bool
	Title           string
	Data            interface{}
	AlertId         string
	AlertMessage    string
	Lang            string
	Translate       func(string) string
}

func AuthSessionService(response http.ResponseWriter, request *http.Request, title string, data interface{}) *AuthSessionStruct {
	user, isAuthenticated := GetAuthenticatedUser(request)
	alertId, alertMessage := GetFlashNotifications(response, request)

	lang := config.AppConfig().Lang
	if cookie, err := request.Cookie("lang"); err == nil {
		lang = cookie.Value
	}

	translate := func(key string) string {
		return internationalization.Translate(key, lang)
	}

	return &AuthSessionStruct{
		User:            user,
		IsAuthenticated: isAuthenticated,
		Title:           title,
		Data:            data,
		AlertId:         alertId,
		AlertMessage:    alertMessage,
		Lang:            lang,
		Translate:       translate,
	}
}

var sessionStoreOnce sync.Once
var sessionStore *sessions.CookieStore

func GetSessionStore() *sessions.CookieStore {
	var appConfig = config.AppConfig()

	sessionStoreOnce.Do(func() {
		sessionStore = sessions.NewCookieStore([]byte(appConfig.Key))
	})
	return sessionStore
}

func LoginUserSession(context *gin.Context, user UserSessionStruct) error {
	var session, sessionError = GetSessionStore().Get(context.Request, "user-core_session")
	if sessionError != nil {
		http.Error(context.Writer, "Error al crear la sesión", http.StatusInternalServerError)
		return sessionError
	}

	session.Values["user_id"] = user.ID
	session.Values["user_first_name"] = user.FirstName
	session.Values["user_last_name"] = user.LastName
	session.Values["user_username"] = user.Username
	session.Values["user_avatar"] = user.Avatar.String
	session.Values["user_language"] = user.Language.String
	session.Values["user_email"] = user.Email
	session.Values["authenticated"] = true

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   84400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
	}

	return session.Save(context.Request, context.Writer)
}

func GetAuthenticatedUser(request *http.Request) (UserSessionStruct, bool) {
	var session, sessionError = GetSessionStore().Get(request, "user-core_session")
	if sessionError != nil {
		return UserSessionStruct{}, false
	}

	var authenticated, ok = session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		return UserSessionStruct{}, false
	}

	// Validar que todos los valores necesarios existan y sean del tipo correcto
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	firstName, ok := session.Values["user_first_name"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	lastName, ok := session.Values["user_last_name"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	username, ok := session.Values["user_username"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	avatar, ok := session.Values["user_avatar"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	language, ok := session.Values["user_language"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	email, ok := session.Values["user_email"].(string)
	if !ok {
		return UserSessionStruct{}, false
	}

	var user = UserSessionStruct{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Avatar:    nulltypes.NullString{String: avatar, Valid: avatar != ""},
		Language:  nulltypes.NullString{String: language, Valid: language != ""},
		Email:     email,
	}

	return user, true
}

// getUserFromDB obtiene un usuario desde la base de datos directamente para evitar ciclos de importación
func getUserFromDB(userID string) (UserSessionStruct, error) {
	// Obtener conexión a la base de datos
	db := database_connections.DatabaseConnectSQL()
	defer db.Close()

	query := "SELECT id, first_name, last_name, username, avatar, language, email FROM users WHERE id = ?"
	row := db.QueryRow(query, userID)

	var userData struct {
		ID        string
		FirstName string
		LastName  string
		Username  string
		Avatar    *string
		Language  string
		Email     string
	}

	err := row.Scan(&userData.ID, &userData.FirstName, &userData.LastName, &userData.Username, &userData.Avatar, &userData.Language, &userData.Email)
	if err != nil {
		return UserSessionStruct{}, err
	}

	// Convertir avatar a NullString
	var avatar nulltypes.NullString
	if userData.Avatar != nil {
		avatar = nulltypes.NullString{String: *userData.Avatar, Valid: true}
	} else {
		avatar = nulltypes.NullString{String: "", Valid: false}
	}

	// Convertir a UserSessionStruct
	user := UserSessionStruct{
		ID:            userData.ID,
		FirstName:     userData.FirstName,
		LastName:      userData.LastName,
		Username:      userData.Username,
		Avatar:        avatar,
		Language:      StringToNullString(userData.Language),
		Email:         userData.Email,
		Authenticated: true,
	}

	return user, nil
}

// GetAppAuthenticatedUser obtiene el usuario autenticado desde el contexto de Gin,
// manejando tanto sesiones web como tokens de API.
func GetAppAuthenticatedUser(c *gin.Context) (UserSessionStruct, bool) {
	// Primero, intentar obtener el user_id del contexto de Gin (para API con JWT)
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok && id != "" {
			// Si existe, obtener los detalles del usuario desde la base de datos
			user, err := getUserFromDB(id)
			if err != nil {
				Logs("ERROR", "Error getting user from database: "+err.Error())
				return UserSessionStruct{}, false
			}

			return user, true
		}
	}

	// Si no se encuentra en el contexto de Gin, recurrir al método de sesión (para web)
	return GetAuthenticatedUser(c.Request)
}

func LogoutUserSession(response http.ResponseWriter, request *http.Request) error {
	var session, sessionError = GetSessionStore().Get(request, "user-core_session")
	if sessionError != nil {
		http.Error(response, "Error al crear la sesión", http.StatusInternalServerError)
		return sessionError
	}

	// Limpiar todos los valores de sesión correctamente
	session.Values["user_id"] = nil
	session.Values["user_first_name"] = nil
	session.Values["user_last_name"] = nil
	session.Values["user_username"] = nil
	session.Values["user_avatar"] = nil
	session.Values["user_language"] = nil
	session.Values["user_email"] = nil
	session.Values["authenticated"] = false

	session.Options.MaxAge = -1

	return session.Save(request, response)
}

func IsUserAuthenticated(request *http.Request) bool {
	_, authenticated := GetAuthenticatedUser(request)
	return authenticated
}
