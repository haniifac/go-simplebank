package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	mockdb "github.com/haniifac/simplebank/db/mock"
	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
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
	return fmt.Sprintf("matches arg %v, password %s", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestGetUser(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
				requireBodyMatcherUser(t, recorder, user)
			},
		},
		{
			name:     "NotFound",
			username: user.Username,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.username)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			// Send request
			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"username": user.Username,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    user.Email,
					Username: user.Username,
					Fullname: user.Fullname,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatcherUser(t, recorder, user)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"password": password,
				"username": user.Username,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"username": user.Username,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"username": user.Username,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{
						Code: pq.ErrorCode("23505"),
					})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
		{
			name: "InvalidUsernameFormat",
			body: gin.H{
				"email":    user.Email,
				"password": password,
				"username": "invalid-username!",
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
		{
			name: "InvalidEmailFormat",
			body: gin.H{
				"email":    "invalid-email",
				"password": password,
				"username": user.Username,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
		{
			name: "PasswordTooShort",
			body: gin.H{
				"email":    user.Email,
				"password": "short",
				"username": user.Username,
				"fullname": user.Fullname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatcherUser(t, recorder, db.User{})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			bodyByes, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(bodyByes))

			// Send request
			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = "password"
	hashpass, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashpass,
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	return
}

func requireBodyMatcherUser(t *testing.T, body *httptest.ResponseRecorder, user db.User) {
	// t.Helper()
	bodyResponse, err := io.ReadAll(body.Body)
	require.NoError(t, err)

	var gotUser db.User
	var userInfo gin.H
	err = json.Unmarshal(bodyResponse, &gotUser)
	require.NoError(t, err)

	err = json.Unmarshal(bodyResponse, &userInfo)
	require.NoError(t, err)

	require.Equal(t, user, gotUser)

}
