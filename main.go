package main

import (
	"crypto/tls"
	"encoding/json"
	
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"
)

// type User struct {
// 	ID          int
// 	Name        string
// 	Age         int
// 	Description string
// 	Img         string
// 	Tags        []string
// }

var (
	users   = make(map[string]string)
	cookies = make(map[string]string)
)

type StatusLogedInJSON struct {
	Status string `json:"status"` // status 400 200
	// body
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (statusJSON *StatusLogedInJSON) ChangeStatus() {
	// 200 404 not found
	if statusJSON.Status == "error" {
		statusJSON.Status = "ok"
	} else {
		statusJSON.Status = "error"
	}
}

func homePageHandler(rw http.ResponseWriter, r *http.Request) {
	// _, err := r.Cookie("session_id")
	// loggedIn := (err != http.ErrNoCookie)
	rw.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	t, _ := template.ParseFiles("static/main.html")
	_ = t.Execute(rw, nil)
}

// func cookieHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// }

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// перенести в роутер
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		fmt.Println(11111)
		return
	}
	if cookies[session.Name] == session.Value {
		m := StatusLogedInJSON{"ok"}
		b, err := json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return

	}
	fmt.Println(session)

	b, _ := ioutil.ReadAll(r.Body)
	jsn := string(b)
	var logUserData LoginUser
	err = json.Unmarshal([]byte(jsn), &logUserData)
	if err != nil {
		// error 400 in json
		fmt.Println(err)
	}
	m := StatusLogedInJSON{"error"}
	for key, value := range users {
		if key == logUserData.Email && value == logUserData.Password {
			m.ChangeStatus()
		}
	}
	b, err = json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if m.Status == "ok" {
		// cookie
		expiration := time.Now().Add(10 * time.Hour)
		cookie := http.Cookie{
			Name:     "session_id",
			Value:    logUserData.Email,
			Expires:  expiration,
			Secure:   true,
			HttpOnly: true,
		}
		// r.Cookie()
		// fmt.Println(r.Cookie("session_id"))
		cookies["session_id"] = logUserData.Email

		http.SetCookie(w, &cookie)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	users["mumeu222@mail.ru"] = "VBif222!"

	mux := http.NewServeMux()
	mux.HandleFunc("/", homePageHandler)
	mux.HandleFunc("/api/v1/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	cfg := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
    }

	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	)
	mux.Handle("/data/", staticHandler)

	server := http.Server{
		Addr:         ":80",
		Handler:      mux,
		ReadTimeout:  http.DefaultClient.Timeout,
		WriteTimeout: http.DefaultClient.Timeout,
		TLSConfig:    cfg,
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	fmt.Println("starting server at :80")
	server.ListenAndServe()
}
