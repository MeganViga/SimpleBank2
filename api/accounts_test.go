package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/MeganViga/SimpleBank2/db/mock"
	db "github.com/MeganViga/SimpleBank2/db/sqlc"
	"github.com/MeganViga/SimpleBank2/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


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
	require.Equal(t,http.StatusOK,recorder.Code)

}

func randomAccount()db.Account{
	return db.Account{
		ID: int64(util.RandInt(1,1000)),
		Owner: util.RandomOwner(),
		Balance: int64(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}
}