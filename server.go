package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func getAgeFromDate(date string) (uint, error) {
	// 2012-12-12
	// 0123456789
	// userDay, err := strconv.Atoi(date[8:])
	// if err != nil {
	// 	return 0, errors.New("failed on userDay")
	// }
	// userMonth, err := strconv.Atoi(date[5:7])
	// if err != nil {
	// 	return 0, errors.New("failed on userMonth")
	// }

	userYear, err := strconv.Atoi(date[:4])
	if err != nil {
		return 0, errors.New("failed on userYear")
	}

	age := (uint)(time.Now().Year() - userYear)

	return age, nil
}

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
	setupCORSResponse(&w, r)

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

	identifiableUser, err := env.db.getUser(logUserData.Email)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	status := StatusOK
	if identifiableUser.isCorrectPassword(logUserData.Password) {
<<<<<<< HEAD
		cookie := createSessionCookie(logUserData)
		err = env.sessionDB.newSessionCookie(cookie.Value, identifiableUser.ID)
=======
		expiration := time.Now().Add(24 * time.Hour)

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
>>>>>>> fd832483177b4e1b1b126afbf9747668d30103e7
		if err != nil {
			resp.Status = StatusInternalServerError
			sendResp(resp, &w)
			return
		}

		http.SetCookie(w, &cookie)
	} else {
		status = StatusNotFound
	}

	resp.Body = identifiableUser
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

func (env *Env) editHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	var user User
	err = json.Unmarshal(byteReq, &user)
	if err != nil {
		fmt.Println("unmarhal error")
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	session, err := r.Cookie("sessionId")
	if err != nil {
		fmt.Println("session error")
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

	currentUser.Name = user.Name
	currentUser.Date = user.Date
	currentUser.Description = user.Description
	currentUser.Tags = user.Tags
	currentUser.Age, err = getAgeFromDate(user.Date)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	err = env.db.updateUser(currentUser)
	if err != nil {
		resp.Status = StatusNotFound
		sendResp(resp, &w)
		return
	}

	resp.Status = StatusOK
	resp.Body = currentUser

<<<<<<< HEAD
=======
	resp.Body = newUser
	resp.Status = status
	sendResp(resp, &w)
}

func (env *Env) editHandler(w http.ResponseWriter, r *http.Request) {
	var resp JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	var UserData User
	err = json.Unmarshal(byteReq, &UserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	status := StatusOK

	_, updateError := env.db.updateUserModel(UserData)
	if updateError != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, &w)
		return
	}

	// куки
	expiration := time.Now().Add(10 * time.Hour)

	data := UserData.Password + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   false,
		HttpOnly: true,
	}

	updateUser, updateError := env.db.updateUserModel(UserData)

	http.SetCookie(w, &cookie)
	// куки

	resp.Body = updateUser
	resp.Status = status
>>>>>>> fd832483177b4e1b1b126afbf9747668d30103e7
	sendResp(resp, &w)
}

func (env *Env) logoutHandler(w http.ResponseWriter, r *http.Request) {
	setupCORSResponse(&w, r)

	session, err := r.Cookie("sessionId")

<<<<<<< HEAD
	if err != nil {
=======
	if err == http.ErrNoCookie {

>>>>>>> fd832483177b4e1b1b126afbf9747668d30103e7
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

type Env struct {
	db interface {
<<<<<<< HEAD
		getUser(email string) (User, error)
		getUserByID(userID uint64) (User, error)
		createUser(logUserData LoginUser) (User, error)
		addSwipedUsers(currentUserId, swipedUserId uint64) error
		getNextUserForSwipe(currentUserId uint64) (User, error)
		updateUser(user User) error
=======
		getUserModel(string) (User, error)
		insertUserModel(User) (int, error)
		updateUserModel(User) (User, error)
>>>>>>> fd832483177b4e1b1b126afbf9747668d30103e7
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
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
	marvin2 := User{
		ID:          2,
		Name:        "Mikhail2",
		Email:       "mumeu222@mail.ru2",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
	marvin3 := User{
		ID:          3,
		Name:        "Mikhail3",
		Email:       "mumeu222@mail.ru3",
		Password:    "af57966e1958f52e41550e822dd8e8a4", //VBif222!
		Date:        "2012-12-12",
		Age:         20,
		Description: "Hahahahaha",
		ImgSrc:      "/img/Yachty-tout.jpg",
		Tags:        []string{"soccer", "anime"},
	}
	db.users[1] = marvin
	db.users[2] = marvin2
	db.users[3] = marvin3
}

func main() {
	env := &Env{
		db:        db, // NewMockDB()
		sessionDB: NewSessionDB(),
	}

	router := mux.NewRouter()

<<<<<<< HEAD
	router.PathPrefix("/api/v1/").HandlerFunc(env.corsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET")
	router.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")
	//router.HandleFunc("/api/v1/createprofile", env.loginHandler).Methods("POST")
	router.HandleFunc("/api/v1/edit", env.editHandler).Methods("POST")
	router.HandleFunc("/api/v1/signup", env.signupHandler).Methods("POST")
	router.HandleFunc("/api/v1/logout", env.logoutHandler).Methods("GET")
	router.HandleFunc("/api/v1/nextswipeuser", env.nextUserHandler).Methods("POST")

=======
	mux.HandleFunc("/api/v1/currentuser", env.currentUser).Methods("GET")
	mux.HandleFunc("/api/v1/login", env.loginHandler).Methods("POST")
	mux.HandleFunc("/api/v1/signup", env.signupHandler).Methods("POST")
	mux.HandleFunc("/api/v1/edit", env.editHandler).Methods("POST")
	mux.HandleFunc("/api/v1/logout", env.logoutHandler).Methods("GET")

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	mux.PathPrefix("/").Handler(spa)

>>>>>>> fd832483177b4e1b1b126afbf9747668d30103e7
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

<<<<<<< HEAD
	log.Fatal(srv.ListenAndServeTLS("./monkeys-drip.com+3.pem", "./monkeys-drip.com+3-key.pem"))
=======
	log.Fatal(srv.ListenAndServe())
>>>>>>> fd832483177b4e1b1b126afbf9747668d30103e7
}
