package poker

import (
	"net/http"
	"testing"
	"net/http/httptest"
)


func TestRecordingWinsAndRetrivingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, "[]")
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("test get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(t, response.Code, http.StatusOK)
		assertBodyEqual(t, response.Body.String(), "3")
	})

	t.Run("test get league", func(t *testing.T) {
		want := []Player{
			{player, 3},
		}

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetLeagueRequest())

		assertStatus(t, response.Code, http.StatusOK)
		got := getLeagueFromResponse(t, response.Body)
		assertLeague(t, got, want)
	})

}
