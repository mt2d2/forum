package model

import (
	"reflect"
	"testing"
)

func mockUserTest() *User {
	return &User{
		Id:       1,
		Username: "test",
		Email:    "test",
		Password: []byte{},
		PasswordHash: []byte{0x24, 0x32, 0x61, 0x24, 0x31, 0x30,
			0x24, 0x72, 0x45, 0x73, 0x37, 0x75, 0x64, 0x69, 0x4b,
			0x77, 0x4b, 0x65, 0x46, 0x69, 0x4c, 0x63, 0x33, 0x49,
			0x68, 0x4d, 0x65, 0x63, 0x65, 0x51, 0x6c, 0x49, 0x68,
			0x4d, 0x46, 0x51, 0x4e, 0x70, 0x30, 0x6a, 0x79, 0x69,
			0x51, 0x78, 0x4d, 0x75, 0x77, 0x31, 0x51, 0x43, 0x36,
			0x37, 0x4f, 0x6e, 0x4f, 0x47, 0x6a, 0x63, 0x51, 0x75},
	}
}

func mockUserTester() *User {
	return &User{
		Id:       2,
		Username: "tester",
		Email:    "test@test.com",
		Password: []byte{},
		PasswordHash: []byte{0x24, 0x32, 0x61, 0x24, 0x31, 0x30,
			0x24, 0x55, 0x2f, 0x31, 0x58, 0x4e, 0x51, 0x67, 0x54,
			0x50, 0x54, 0x52, 0x6d, 0x37, 0x34, 0x6e, 0x45, 0x6c,
			0x49, 0x51, 0x47, 0x39, 0x75, 0x6b, 0x66, 0x6f, 0x79,
			0x6a, 0x4b, 0x75, 0x47, 0x2e, 0x6a, 0x55, 0x47, 0x37,
			0x65, 0x34, 0x58, 0x64, 0x48, 0x57, 0x33, 0x43, 0x70,
			0x64, 0x6b, 0x67, 0x65, 0x47, 0x51, 0x6c, 0x4a, 0x6d},
	}
}

func TestEmptyuser(t *testing.T) {
	if !reflect.DeepEqual(NewUser(), &User{-1, "", "", []byte{}, []byte{}}) {
		t.Error("user not empty")
	}
}

func TestFindOneUserByUsername(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	userTest, err := FindOneUserByUsername(db, "test")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(userTest, *mockUserTest()) {
		t.Error("wrong user")
	}

	userTester, err := FindOneUserByUsername(db, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(userTester, *mockUserTester()) {
		t.Error("wrong user")
	}
}

func TestFindOneUserById(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	userTest, err := FindOneUserById(db, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(userTest, *mockUserTest()) {
		t.Error("wrong user")
	}

	userTester, err := FindOneUserById(db, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(userTester, *mockUserTester()) {
		t.Error("wrong user")
	}
}

func TestSaveUser(t *testing.T) {
	db, err := GetMockupDB()
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}

	user := mockUserTest()
	user.PasswordHash = []byte{}

	err = SaveUser(db, user)
	if err == nil {
		t.Error("a user must have a hashed password")
	}

	user = mockUserTest()
	err = SaveUser(db, user)
	if err != nil {
		t.Fatal(err)
	}
	user.Id = 3
	userTest, err := FindOneUserById(db, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(userTest, *user) {
		t.Error("wrong user")
	}
}
