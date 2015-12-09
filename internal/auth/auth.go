package auth

import "net/http"

var password = ""

func SetPassword(passwd string) {
	password = passwd
}

func SaveSession(w http.ResponseWriter, passwd string) {
	cookie := &http.Cookie{
		Name:     "password",
		Value:    passwd,
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

func Logged(r *http.Request) bool {
	cookie, err := r.Cookie("password")
	if err == nil {
		return password == cookie.Value
	}
	return false
}
