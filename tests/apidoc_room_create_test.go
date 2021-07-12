package tests

import (
	"testing"

	"github.com/tidwall/gjson"

	"github.com/matrix-org/complement/internal/b"
	"github.com/matrix-org/complement/internal/client"
	"github.com/matrix-org/complement/internal/match"
	"github.com/matrix-org/complement/internal/must"
)

func TestRoomCreate(t *testing.T) {
	deployment := Deploy(t, b.BlueprintAlice)
	defer deployment.Destroy(t)
	authedClient := deployment.Client(t, "hs1", "@alice:hs1")
	t.Run("Parallel", func(t *testing.T) {
		// sytest: POST /createRoom makes a public room
		t.Run("POST /createRoom makes a public room", func(t *testing.T) {
			t.Parallel()
			reqBody := client.WithJSONBody(t, map[string]interface{}{
				"visibility":      "public",
				"room_alias_name": "30-room-create-alias-random",
			})
			res := authedClient.MustDoFunc(t, "POST", []string{"_matrix", "client", "r0", "createRoom"}, reqBody)
			must.MatchResponse(t, res, match.HTTPResponse{
				JSON: []match.JSON{
					match.JSONKeyPresent("room_id"),
					match.JSONKeyPresent("room_alias"),
					match.JSONKeyTypeEqual("room_id", gjson.String),
				},
			})
		})
		// sytest: POST /createRoom makes a private room
		t.Run("POST /createRoom makes a private room", func(t *testing.T) {
			t.Parallel()
			reqBody := client.WithJSONBody(t, map[string]interface{}{
				"visibility": "private",
			})
			res := authedClient.MustDoFunc(t, "POST", []string{"_matrix", "client", "r0", "createRoom"}, reqBody)
			must.MatchResponse(t, res, match.HTTPResponse{
				JSON: []match.JSON{
					match.JSONKeyPresent("room_id"),
					match.JSONKeyTypeEqual("room_id", gjson.String),
				},
			})
		})
		// sytest: POST /createRoom makes a private room with invites
		t.Run("POST /createRoom makes a private room with invites", func(t *testing.T) {
			t.Parallel()
			userInvite := deployment.RegisterUser(t, "hs1", "create_room", "superuser")
			reqBody := client.WithJSONBody(t, map[string]interface{}{
				"visibility": "private",
				"invite":     []string{userInvite.UserID},
			})
			res := authedClient.MustDoFunc(t, "POST", []string{"_matrix", "client", "r0", "createRoom"}, reqBody)
			must.MatchResponse(t, res, match.HTTPResponse{
				JSON: []match.JSON{
					match.JSONKeyPresent("room_id"),
					match.JSONKeyTypeEqual("room_id", gjson.String),
				},
			})
		})
		// sytest: POST /createRoom rejects attempts to create rooms with numeric versions
		t.Run("POST /createRoom rejects attempts to create rooms with numeric versions", func(t *testing.T) {
			t.Parallel()
			reqBody := client.WithJSONBody(t, map[string]interface{}{
				"visibility":   "private",
				"room_version": 1,
				"preset":       "public_chat",
			})
			res := authedClient.DoFunc(t, "POST", []string{"_matrix", "client", "r0", "createRoom"}, reqBody)
			must.MatchResponse(t, res, match.HTTPResponse{
				StatusCode: 400,
				JSON: []match.JSON{
					match.JSONKeyPresent("errcode"),
					match.JSONKeyEqual("errcode", "M_BAD_JSON"),
				},
			})
		})
		// sytest: POST /createRoom rejects attempts to create rooms with unknown versions
		t.Run("POST /createRoom rejects attempts to create rooms with unknown versions", func(t *testing.T) {
			t.Parallel()
			reqBody := client.WithJSONBody(t, map[string]interface{}{
				"visibility":   "private",
				"room_version": "ahfgwjyerhgiuveisbruvybseyrugvi",
				"preset":       "public_chat",
			})
			res := authedClient.DoFunc(t, "POST", []string{"_matrix", "client", "r0", "createRoom"}, reqBody)
			must.MatchResponse(t, res, match.HTTPResponse{
				StatusCode: 400,
				JSON: []match.JSON{
					match.JSONKeyPresent("errcode"),
					match.JSONKeyEqual("errcode", "M_UNSUPPORTED_ROOM_VERSION"),
				},
			})
		})
	})
}
