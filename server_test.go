package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestCurrentUser(t *testing.T) {
	//t.Parallel() TODO Tests Parallel

	h := Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	var body io.Reader

	id := uint64(len(users) + 1)
	users[id] = makeUser(id, "testCurrentUser1@mail.ru", "123456qQ")
	cookies["123"] = id

	idStr := strconv.FormatUint(id, 10)
	expected := `{"status":200,"body":{"id":`+ idStr +`,"name":"","email":"testCurrentUser1@mail.ru","age":0,"description":"","imgSrc":"","tags":null}}`

	r := httptest.NewRequest("GET", "/api/v1/currentuser", body)
	c := http.Cookie{
		Name:     "sessionId",
		Value:    "123",
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()

	h.currentUser(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if w.Body.String() != expected {
		t.Error("session is invalid")
	}
}

func TestLogin(t *testing.T) {
	//t.Parallel()

	h := Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	body := bytes.NewReader([]byte(`{"email":"testLogin1@mail.ru","password":"123456qQ"}`))

	id := uint64(123123)
	users[id] = makeUser(id, "testLogin1@mail.ru", "123456qQ")

	r := httptest.NewRequest("POST", "/api/v1/login/", body)
	w := httptest.NewRecorder()

	h.loginHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	cookie := false
	for _, userID := range cookies {
		if userID == id {
			cookie = true
			break
		}
	}
	if !cookie {
		t.Error("session was not created")
	}
}

func TestSignup(t *testing.T) {
	//t.Parallel()

	h := Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	body := bytes.NewReader([]byte(`{"email":"testSignup1@mail.ru","password":"123456qQ"}`))

	expectedUsers := make(map[uint64]User)
	for k, v := range users {
		expectedUsers[k] = v
	}
	id := uint64(len(expectedUsers) + 1)
	expectedUsers[id] = makeUser(id, "testSignup1@mail.ru", "123456qQ")

	r := httptest.NewRequest("POST", "/api/v1/signup/", body)
	w := httptest.NewRecorder()

	h.signupHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if !reflect.DeepEqual(users, expectedUsers) {
		t.Error("user was not created")
	}

	cookie := false
	for _, userID := range cookies {
		if userID == id {
			cookie = true
			break
		}
	}
	if !cookie {
		t.Error("session was not created")
	}
}

func TestLogout(t *testing.T) {
	//t.Parallel() TODO Tests Parallel

	h := Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	var body io.Reader

	id := uint64(len(users) + 1)
	users[id] = makeUser(id, "testLogout1@mail.ru", "123456qQ")
	cookies["123"] = id

	r := httptest.NewRequest("GET", "/api/v1/logout", body)
	c := http.Cookie{
		Name:     "sessionId",
		Value:    "123",
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()

	h.logoutHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if _, ok := cookies["123"]; ok {
		t.Error("user session not ended")
	}
}

func TestNextUser(t *testing.T) {
	//t.Parallel() TODO Tests Parallel

	h := Env{
		db:        MockDB{},
		sessionDB: MockSessionDB{},
	}

	body := bytes.NewReader([]byte(`{"id":1234}`))

	id := uint64(len(users) + 1)
	users[id] = makeUser(id, "testNextUser1@mail.ru", "123456qQ")
	cookies["123"] = id

	expected := `{"status":200,"body":{"id":9999,"name":"","email":"testNextUserForSwipe@mail.ru","age":0,"description":"","imgSrc":"","tags":null}}`

	r := httptest.NewRequest("POST", "/api/v1/nextswipeuser", body)
	c := http.Cookie{
		Name:     "sessionId",
		Value:    "123",
	}
	r.AddCookie(&c)
	w := httptest.NewRecorder()

	h.nextUserHandler(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	if _, ok := swipedUsers[id]; !ok {
		t.Error("swipe not saved")
	}

	if w.Body.String() != expected {
		t.Error("invalid data nextSwipe")
	}
}
