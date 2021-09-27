package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

func sendResp(resp JSON, w *http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(byteResp)
}

func setupCORSResponse(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Allow-Credentials, Set-Cookie, Access-Control-Allow-Credentials, Access-Control-Allow-Origin")
}

func (env *Env) corsHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)
}

func (env *Env) currentUser(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	var resp JSON
	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		fmt.Println("tyt suka")
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	currentUser, err := env.sessionDB.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	fmt.Println(resp)

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, &w)
}

func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	var logUserData LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	identifiableUser, err := env.db.getUserModel(logUserData.Email)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	status := StatusOK
	if identifiableUser.isCorrectPassword(logUserData.Password) {
		expiration := time.Now().Add(10 * time.Hour)

		data := logUserData.Password + time.Now().String()
		md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

		cookie := http.Cookie{
			Name:     "sessionId",
			Value:    md5CookieValue,
			Expires:  expiration,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
		}

		err = env.sessionDB.newSessionCookie(md5CookieValue, identifiableUser.ID)
		if err != nil {
			resp.Status = StatusInternalServerError
			sendResp(resp, &w)
			return
		}

		http.SetCookie(w, &cookie)
	} else {
		status = StatusNotFound
	}

	resp.Status = status
	sendResp(resp, &w)
}

func (env *Env) signupHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	var logUserData LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	identifiableUser, err := env.db.getUserModel(logUserData.Email)
	if err == nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	status := StatusOK

	newID := len(users) + 1
	newEmail := logUserData.Email
	newPassword := logUserData.Password

	md5NewPassword := fmt.Sprintf("%x", md5.Sum([]byte(newPassword)))

	newUser := User{
		ID:       uint64(newID),
		Email:    newEmail,
		Password: md5NewPassword,
	}

	users[uint64(newID)] = newUser

	// куки
	expiration := time.Now().Add(10 * time.Hour)

	data := logUserData.Password + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
	}

	err = env.sessionDB.newSessionCookie(md5CookieValue, identifiableUser.ID)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, &w)
		return
	}

	http.SetCookie(w, &cookie)
	// куки

	resp.Status = status
	sendResp(resp, &w)
}

func (env *Env) logoutHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	session, err := r.Cookie("sessionId")

	if err == http.ErrNoCookie {
		sendResp(JSON{Status: StatusNotFound}, &w)
		return
	}

	err = env.sessionDB.deleteSessionCookie(session.Value)
	if err != nil {
		sendResp(JSON{Status: StatusInternalServerError}, &w)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (env *Env) nextUserHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	var resp JSON

	// get current user by cookie
	session, err := r.Cookie("sessionId")
	if err == http.ErrNoCookie {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}
	currentUser, err := env.sessionDB.getUserByCookie(session.Value)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	// get swiped user id from json
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}
	var swipedUserData SwipedUser
	err = json.Unmarshal(byteReq, &swipedUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	// add in swaped users map for current user
	err = env.db.addSwipedUsers(currentUser.ID, swipedUserData.Id)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}
	// find next user for swipe
	nextUser, err := env.db.getNextUserForSwipe(currentUser.ID)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	resp.Status = StatusOK
	resp.Body = nextUser

	sendResp(resp, &w)
}

type Env struct {
	db interface {
		getUserModel(string) (User, error)
		addSwipedUsers(uint64, uint64) error
		getNextUserForSwipe(uint64) (User, error)
	}
	sessionDB interface {
		getUserByCookie(sessionCookie string) (User, error)
		newSessionCookie(sessionCookie string, userId uint64) error
		deleteSessionCookie(sessionCookie string) error
	}
}

func init() {
	marvin := User{
		ID:          1,
		Name:        "Mikhail",
		Email:       "mumeu222@mail.ru",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"haha", "hihi"},
	}
	marvin2 := User{
		ID:          2,
		Name:        "Mikhail2",
		Email:       "mumeu222@mail.ru2",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"haha", "hihi"},
	}
	marvin3 := User{
		ID:          3,
		Name:        "Mikhail3",
		Email:       "mumeu222@mail.ru3",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"haha", "hihi"},
	}
	users[1] = marvin
	users[2] = marvin2
	users[3] = marvin3
}

func main() {
	/*db, err := sql.Open("postgres", "postgres://user:pass@localhost/bookstore")
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		db: ModelsDB{DB: db},
	}
	*/

	env := &Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	mux := mux.NewRouter()

	// mux.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET")
	// mux.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")
	// mux.HandleFunc("/api/v1/signup", env.signupHandler).Methods("POST")
	// mux.HandleFunc("/api/v1/logout", env.logoutHandler).Methods("GET")
	// mux.HandleFunc("/api/v1/nextswipeuser", env.nextUserHandler).Methods("POST")

	mux.PathPrefix("/api/v1/").HandlerFunc(env.corsHandler).Methods("OPTIONS")
	mux.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET")
	mux.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")
	mux.HandleFunc("/api/v1/signup", env.signupHandler).Methods("POST", "OPTIONS")
	mux.HandleFunc("/api/v1/logout", env.logoutHandler).Methods("GET", "OPTIONS")
	mux.HandleFunc("/api/v1/nextswipeuser", env.nextUserHandler).Methods("POST", "OPTIONS")

	// spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	// mux.PathPrefix("/").Handler(spa)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"https://127.0.0.1:443"},
	// 	AllowCredentials: true,
	// 	AllowedHeaders:   []string{"Access-Control-Allow-Origin"},
	// })
	// handler := c.Handler(mux)

	// headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "Authorization"})
	// originsOk := handlers.AllowedOrigins([]string{"*"})
	// methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	srv := &http.Server{
		Handler: mux,
		// Handler: handler,
		// Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(mux),
		Addr:         ":8080",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServeTLS("./drip-ilyagu.com+4.pem", "./drip-ilyagu.com+4-key.pem"))
}
