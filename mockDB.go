package main

import (
	"errors"
)

var (
	users       = make(map[uint64]User)
	cookies     = make(map[string]uint64)
	swipedUsers = make(map[uint64][]uint64)
)

type MockDB struct {
	//DB int
}

func (MockDB) getUser(email string) (User, error) {
	if len(users) == 0 {
		return User{}, errors.New("users is empty map")
	}

	currentUser := User{}
	okUser := false
	for _, value := range users {
		if value.Email == email {
			currentUser = value
			okUser = true
		}
	}
	if !okUser {
		return User{}, errors.New("User not found")
	}

	return currentUser, nil
}

func (MockDB) createUser(logUserData LoginUser) (User, error) {
	newID := uint64(len(users) + 1)

	users[newID] = makeUser(newID, logUserData.Email, logUserData.Password)

	return users[newID], nil
}

func (MockDB) addSwipedUsers(currentUserId uint64, swipedUserId uint64) error {
	if len(users) == 0 {
		return errors.New("users is empty map")
	}

	if currentUserId != swipedUserId {
		swipedUsers[currentUserId] = append(swipedUsers[currentUserId], swipedUserId)
	}
	return nil
}

func (MockDB) getNextUserForSwipe(currentUserId uint64) (User, error) {
	/*if len(users) == 0 {
		return User{}, errors.New("users is empty map")
	}
	if len(swipedUsers) == 0 {
		for key, value := range users {
			if key != currentUserId {
				return value, nil
			}
		}
		return User{}, errors.New("haven't any other users for swipe")
	}

	// find all users swiped by the current user
	var allSwipedUsersForCurrentUser []uint64
	for key, value := range swipedUsers {
		if key == currentUserId {
			allSwipedUsersForCurrentUser = value
		}
	}

	// find a user who has not yet been swiped by the current user
	for key, value := range users {
		if key == currentUserId {
			continue
		}
		if !existsIn(key, allSwipedUsersForCurrentUser) {
			return value, nil
		}
	}

	return User{}, errors.New("haven't any other users for swipe")*/
	return makeUser(9999, "testNextUserForSwipe@mail.ru", "123456qQ"), nil
}

func existsIn(value uint64, target []uint64) bool {
	exists := false
	for i := range target {
		if value == target[i] {
			exists = true
		}
	}

	return exists
}

type MockSessionDB struct {
}

func (MockSessionDB) getUserByCookie(sessionCookie string) (User, error) {
	if len(cookies) == 0 {
		return User{}, errors.New("cookies is empty map")
	}

	currentUserId, okCookie := cookies[sessionCookie]
	if !okCookie {
		return User{}, errors.New("cookie not found")
	}

	currentUser, okUser := users[currentUserId]
	if !okUser {
		return User{}, errors.New("user not found")
	}

	return currentUser, nil
}

func (MockSessionDB) newSessionCookie(sessionCookie string, userId uint64) error {
	cookies[sessionCookie] = userId
	return nil
}

func (MockSessionDB) deleteSessionCookie(sessionCookie string) error {
	delete(cookies, sessionCookie)
	return nil
}
