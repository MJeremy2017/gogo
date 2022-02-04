package user


import (
	"testing"
	"net/http/httptest"
	"net/http"
	"strings"
	"reflect"
	"encoding/json"
	"fmt"
	"io"
)

type MockUserService struct {
	RegisterFunc func(user User) (string, error)

	UsersRegistered []User
}

func (m *MockUserService) Register(user User) (insertedID string, err error) {
	m.UsersRegistered = append(m.UsersRegistered, user)
	return m.RegisterFunc(user)
}

func TestRegisterUser(t *testing.T) {
	t.Run("can register valid user", func(t *testing.T) {
		user := User{Name: "JZ"}
		expectedInsertedID := "something"

		mockUserService := &MockUserService{
			RegisterFunc: func(user User) (string, error) {
				return expectedInsertedID, nil
			},
		}
		userServer := NewUserServer(mockUserService)

		req := httptest.NewRequest(http.MethodGet, "/", userToJSON(user))
		res := httptest.NewRecorder()

		userServer.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusCreated)

		if res.Body.String() != expectedInsertedID {
			t.Errorf("expected body of %q but got %q", res.Body.String(), expectedInsertedID)
		}

		if len(mockUserService.UsersRegistered) != 1 {
			t.Fatalf("expected 1 user added but got %d", len(mockUserService.UsersRegistered))
		}

		if !reflect.DeepEqual(mockUserService.UsersRegistered[0], user) {
			t.Errorf("expected user registered to be %+v, but got %+v", user, 
				mockUserService.UsersRegistered[0])
		}
	})

	t.Run("returns 400 bad request if body is not valid JSON", func(t *testing.T) {
		server := NewUserServer(nil)

		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("trouble"))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)
		assertStatus(t, res.Code, http.StatusBadRequest)
	})
}


func userToJSON(user User) io.Reader {
	// only encode public fields
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("fail to marshal %+v, %v", user, err)
	}

	return strings.NewReader(string(b))
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("want status code %v, but got %v", got, want)
	}
}








