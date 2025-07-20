package helpers

import (
	"net/http"
	"semita/config"
	"sync"

	"github.com/gorilla/sessions"
)

var storeOnce sync.Once
var store *sessions.CookieStore

func getStore() *sessions.CookieStore {
	var appConfig = config.AppConfig()

	storeOnce.Do(func() {
		store = sessions.NewCookieStore([]byte(appConfig.Key))
	})
	return store
}

func GetFlashNotifications(response http.ResponseWriter, request *http.Request) (string, string) {
	var session, _ = getStore().Get(request, "flash-core_session")

	var alertId = ""
	var alertMensaje = ""

	if session.Values["alert_id"] != nil {
		alertId, _ = session.Values["alert_id"].(string)
		delete(session.Values, "alert_id")
	}

	if session.Values["alert_mensaje"] != nil {
		alertMensaje, _ = session.Values["alert_mensaje"].(string)
		delete(session.Values, "alert_mensaje")
	}

	_ = session.Save(request, response)

	return alertId, alertMensaje
}

func CreateFlashNotification(response http.ResponseWriter, request *http.Request, alertID string, alertMessage string) {
	var session, err = getStore().Get(request, "flash-core_session")

	if err != nil {
		Logs("ERROR", err.Error())
		http.Error(response, "Error al crear la sesión", http.StatusInternalServerError)
		return
	}

	session.Values["alert_id"] = alertID
	session.Values["alert_mensaje"] = alertMessage

	err = session.Save(request, response)
	if err != nil {
		return
	}
}
