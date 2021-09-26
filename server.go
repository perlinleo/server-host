package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func (env *Env) currentUser(w http.ResponseWriter, r *http.Request) {
	var resp JSON
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

	resp.Status = StatusOK
	resp.Body = currentUser

	sendResp(resp, &w)
}

func (env *Env) loginHandler(w http.ResponseWriter, r *http.Request) {
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
			Secure:   false,
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
		Secure:   false,
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

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
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

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET")
	router.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")
	router.HandleFunc("/api/v1/signup", env.signupHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", env.logoutHandler).Methods("GET")
	router.HandleFunc("/api/v1/nextswipeuser", env.nextUserHandler).Methods("POST")

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServe())
}
