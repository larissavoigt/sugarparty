package auth

import "net/http"

func SaveSession(w http.ResponseWriter, password string) {
	cookie := &http.Cookie{
		Name:     "password",
		Value:    password,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func DestroySession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "password",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func Admin(r *http.Request, password string) bool {
	cookie, err := r.Cookie("password")
	if err == nil {
		return password == cookie.Value
	}
	return false
}
