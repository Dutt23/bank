package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github/dutt23/bank/db/mock"
	db "github/dutt23/bank/db/sqlc"
	"github/dutt23/bank/token"
	"github/dutt23/bank/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	user := util.RandomOwner()
	account := randomAccount(user)

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
	}{{
		name:      "OK",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			// build stubs
			store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			// Check response
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder.Body, account)
		},
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuth(t, request, tokenMaker, authTypeBearer, user, time.Minute)
		},
	},
		{
			name:      "NOT FOUND",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user, time.Minute)
			},
		},
		{
			name:      "INTERNAL ERROR",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user, time.Minute)
			},
		},
		{
			name:      "BAD REQUEST",
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user, time.Minute)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Mock db store
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			// start server
			// For a golang interface it's not explicit that a child has to implement all functions in parent.
			// The parent can be imbedded inside the child itself
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
