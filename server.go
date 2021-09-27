package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
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
	StatusEmailAlreadyExists  = 1001
)

func sendResp(resp JSON, w *http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}
	(*w).WriteHeader(http.StatusOK)
	(*w).Write(byteResp)
}

func createSessionCookie(user LoginUser) http.Cookie {
	expiration := time.Now().Add(10 * time.Hour)

	data := user.Password + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   false,
		HttpOnly: true,
	}

	return cookie
}

func (env *Env) currentUser(w http.ResponseWriter, r *http.Request) {
	var resp JSON
	session, err := r.Cookie("sessionId")
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	currentUser, err := env.getUserByCookie(session.Value)
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

	identifiableUser, err := env.db.getUser(logUserData.Email)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	status := StatusOK
	if identifiableUser.isCorrectPassword(logUserData.Password) {
		cookie := createSessionCookie(logUserData)
		err = env.sessionDB.newSessionCookie(cookie.Value, identifiableUser.ID)
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

	identifiableUser, _ := env.db.getUser(logUserData.Email)
	if !identifiableUser.isEmpty() {
		resp.Status = StatusEmailAlreadyExists
		sendResp(resp, &w)
		return
	}

	user, err := env.db.createUser(logUserData)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, &w)
		return
	}

	cookie := createSessionCookie(logUserData)
	err = env.sessionDB.newSessionCookie(cookie.Value, user.ID)
	if err != nil {
		resp.Status = StatusInternalServerError
		sendResp(resp, &w)
		return
	}

	http.SetCookie(w, &cookie)

	resp.Status = StatusOK
	sendResp(resp, &w)
}

func (env *Env) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sessionId")

	if err != nil {
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
	currentUser, err := env.getUserByCookie(session.Value)
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
		getUser(email string) (User, error)
		getUserByID(userID uint64) (User, error)
		createUser(logUserData LoginUser) (User, error)
		addSwipedUsers(currentUserId, swipedUserId uint64) error
		getNextUserForSwipe(currentUserId uint64) (User, error)
	}
	sessionDB interface {
		getUserIDByCookie(sessionCookie string) (uint64, error)
		newSessionCookie(sessionCookie string, userId uint64) error
		deleteSessionCookie(sessionCookie string) error
	}
}

func (env Env) getUserByCookie(sessionCookie string) (User, error) {
	userID, err := env.sessionDB.getUserIDByCookie(sessionCookie)
	if err != nil {
		return User{}, errors.New("error sessionDB: getUserIDByCookie")
	}

	user, err := env.db.getUserByID(userID)
	if err != nil {
		return User{}, errors.New("error db: getUserByID")
	}

	return user, nil
}

var (
	db = NewMockDB()
)
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
	db.users[1] = marvin
	db.users[2] = marvin2
	db.users[3] = marvin3
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
		db:        db, // NewMockDB()
		sessionDB: NewSessionDB(),
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET")
	router.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")
	//router.HandleFunc("/api/v1/createprofile", env.loginHandler).Methods("POST")
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
