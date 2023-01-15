package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/MeganViga/SimpleBank2/db/mock"
	db "github.com/MeganViga/SimpleBank2/db/sqlc"
	"github.com/MeganViga/SimpleBank2/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T){
	account := randomAccount()
	t.Log("Initial:",account)
	arg := db.CreateAccountParams{
		Owner: account.Owner,
		Balance: account.Balance,
		Currency: account.Currency,
	}
	controller := gomock.NewController(t)
	defer controller.Finish()
	store := mockdb.NewMockStore(controller)

	//stub 
	store.EXPECT().CreateAccount(gomock.Any(),gomock.Eq(arg)).Times(1).Return(account,nil)

	//creating server using mock store
	server := NewServer(store)

	// instead of making real api call , we can use recorder
	recorder := httptest.NewRecorder()

	url := "/users"

	var buf bytes.Buffer
    err := json.NewEncoder(&buf).Encode(arg)
    if err != nil {
        log.Fatal(err)
    }

	//creating request
	request, err := http.NewRequest(http.MethodPost,url,&buf)
	fmt.Println("Request:",request)
	require.NoError(t,err)
	//making request call
	server.router.ServeHTTP(recorder,request)

	//checking response
	t.Log(recorder)
	require.Equal(t,http.StatusOK,recorder.Code)
	//t.Log(recorder)
}

func TestGetAccountApi(t *testing.T){
	account := randomAccount()
	controller := gomock.NewController(t)
	defer controller.Finish()
	store := mockdb.NewMockStore(controller)

	//stub
	store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(account,nil)

	//creating server using mockstore
	server := NewServer(store)

	//instead of making actual api call, we can use recorder
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/getuser/%d",account.ID)
	//creating request
	request, err := http.NewRequest(http.MethodGet,url,nil)
	require.NoError(t,err)

	//making request call
	server.router.ServeHTTP(recorder,request)

	//checking response
	require.Equal(t,http.StatusBadGateway,recorder.Code)

}

func randomAccount()db.Account{
	return db.Account{
		ID: int64(util.RandInt(1,1000)),
		Owner: util.RandomOwner(),
		Balance: int64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}
}