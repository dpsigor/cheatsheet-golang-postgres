package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/dpsigor/cheatsheet-golang-postgres/db/mock"
	db "github.com/dpsigor/cheatsheet-golang-postgres/db/sqlc"
	"github.com/dpsigor/cheatsheet-golang-postgres/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (db.User, string) {
	return db.User{
		Username: util.RandomString(10),
		FullName: util.RandomString(10),
		Email:    util.RandomEmail(),
	}, util.RandomString(10)
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func TestCreateUser(t *testing.T) {
	user, pw := randomUser(t)
	hashedPw, err := util.HashPassword(pw)
	require.NoError(t, err)

	tests := []struct {
		name       string
		body       gin.H
		buildStubs func(store *mockdb.MockStore)
		checkRes   func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  pw,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: hashedPw,
					FullName:       user.FullName,
					Email:          user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), eqCreateUserParamsMatcher{arg: arg, password: pw}).
					Times(1).
					Return(user, nil)
			},
			checkRes: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var gotUser db.User
				err := json.NewDecoder(recorder.Body).Decode(&gotUser)
				require.NoError(t, err)
				require.Equal(t, user, gotUser)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			bs, err := json.Marshal(tc.body)
			require.NoError(t, err)
			payload := bytes.NewReader(bs)
			request, err := http.NewRequest(http.MethodPost, "/users", payload)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkRes(t, recorder)
		})
	}
}
