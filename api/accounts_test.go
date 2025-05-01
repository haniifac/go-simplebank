package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/haniifac/simplebank/db/mock"
	db "github.com/haniifac/simplebank/db/sqlc"
	"github.com/haniifac/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	account1 := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// stub the GetAccount method to return the account1 when called with the same ID
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatcherAccount(t, recorder, account1)
			},
		},
		{
			name:      "NotFound",
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatcherAccount(t, recorder, db.Account{})
			},
		},
		{
			name:      "InternalError",
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatcherAccount(t, recorder, db.Account{})
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				// requireBodyMatcherAccount(t, recorder, db.Account{})
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

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			// Send request
			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(util.RandomInt(1, 1000)),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomMoney(),
	}
}

func requireBodyMatcherAccount(t *testing.T, body *httptest.ResponseRecorder, account db.Account) {
	// t.Helper()
	bodyResponse, err := io.ReadAll(body.Body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(bodyResponse, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)

}
