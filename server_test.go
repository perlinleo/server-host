package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	t.Parallel()

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	id := uint64(len(testDB.users) + 1)
	testDB.users[id] = makeUser(id, "testCurrentUser1@mail.ru", "123456qQ")
	testSessionDB.cookies["123"] = id

	var body io.Reader

	idStr := strconv.FormatUint(id, 10)
	expected := `{"status":200,"body":{"id":`+ idStr +`,"name":"","email":"testCurrentUser1@mail.ru","age":0,"description":"","imgSrc":"","tags":null}}`

	r := httptest.NewRequest("GET", "/api/v1/currentuser", body)
	c := http.Cookie{
		Name:     "sessionId",
		Value:    "123",
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()

	env.currentUser(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if w.Body.String() != expected {
		t.Error("session is invalid")
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	body := bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`))

	id := uint64(len(testDB.users) + 1)
	testDB.users[id] = makeUser(id, "testLogin1@mail.ru", "123456qQ")

	r := httptest.NewRequest("POST", "/api/v1/login/", body)
	w := httptest.NewRecorder()

	env.loginHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if !testSessionDB.isSessionByUserID(id) {
		t.Error("session was not created")
	}
}

func TestSignup(t *testing.T) {
	t.Parallel()

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	email := "testSignup1@mail.ru"
	password := "123456qQ"
	body := bytes.NewReader([]byte(`{"email":"` + email + `","password":"` + password + `"}`))

	expectedID := uint64(1)
	expectedUsers := makeUser(expectedID, email, password)

	r := httptest.NewRequest("POST", "/api/v1/signup/", body)
	w := httptest.NewRecorder()

	env.signupHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	newUser, _ := testDB.getUser(email)
	if !reflect.DeepEqual(newUser, expectedUsers) {
		t.Error("user was not created")
	}

	if !testSessionDB.isSessionByUserID(expectedID) {
		t.Error("session was not created")
	}
}

func TestLogout(t *testing.T) {
	t.Parallel()

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	var body io.Reader

	id := uint64(len(testDB.users) + 1)
	testDB.users[id] = makeUser(id, "testLogout1@mail.ru", "123456qQ")
	testSessionDB.cookies["123"] = id

	r := httptest.NewRequest("GET", "/api/v1/logout", body)
	c := http.Cookie{
		Name:     "sessionId",
		Value:    "123",
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()

	env.logoutHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if _, ok := testSessionDB.cookies["123"]; ok {
		t.Error("user session not ended")
	}
}

func TestNextUser(t *testing.T) {
	t.Parallel()

	testDB := NewMockDB()
	testSessionDB := NewSessionDB()

	env := &Env{
		db:        testDB,
		sessionDB: testSessionDB,
	}

	swipedUserID := uint64(1234)
	swipedUserIDStr := strconv.FormatUint(swipedUserID, 10)
	body := bytes.NewReader([]byte(`{"id":` + swipedUserIDStr + `}`))

	currenUser, _ := testDB.createUser(LoginUser{
		Email: "testCurrUser1@mail.ru",
		Password: "123456qQ\"",
	})
	testSessionDB.cookies["123"] = currenUser.ID

	nextUser, _ := testDB.createUser(LoginUser{
		Email: "testNextUser1@mail.ru",
		Password: "123456qQ\"",
	})

	q := JSON{
		Status: StatusOK,
		Body: nextUser,
	}
	expected, _ := json.Marshal(q)

	r := httptest.NewRequest("POST", "/api/v1/nextswipeuser", body)
	c := http.Cookie{
		Name:     "sessionId",
		Value:    "123",
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()

	env.nextUserHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if !testDB.isSwiped(currenUser.ID, swipedUserID) {
		t.Error("swipe not saved")
		fmt.Println(testDB.swipedUsers)
		fmt.Println(currenUser)
		fmt.Println(nextUser)
	}

	if !reflect.DeepEqual(w.Body.Bytes(), expected) {
		t.Error("invalid data nextSwipe")
		fmt.Println(w.Body.String())
	}
}
