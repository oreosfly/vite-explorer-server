package account

import (
	"github.com/gin-gonic/gin"
	typeRequest "github.com/vitelabs/vite-explorer-server/type/request"
	"github.com/vitelabs/vite-explorer-server/util"
	serviceAccount "github.com/vitelabs/vite-explorer-server/service/account"
	"github.com/vitelabs/go-vite/common/types"
	"errors"
	"github.com/vitelabs/go-vite/vite"
	"github.com/vitelabs/go-vite/ledger"
	"math/big"
	"github.com/vitelabs/vite-explorer-server/type/response"
	"github.com/vitelabs/vite-explorer-server/vitelog"
	errors2 "github.com/pkg/errors"

	"math/rand"
)

func Detail(c *gin.Context)  {
	var accountDetailQuery typeRequest.AccountDetail

	if err := c.Bind(&accountDetailQuery); err != nil {
		util.RespondError(c, 400, err)
		return
	}

	if !types.IsValidHexAddress(accountDetailQuery.AccountAddress) {
		util.RespondFailed(c, 1, errors.New("AccountAddress is invalid"), "")
		return
	}

	accountAddress, err:= types.HexToAddress(accountDetailQuery.AccountAddress)

	if err != nil {
		util.RespondFailed(c, 6, err, "")
		return
	}

	account, err := serviceAccount.GetAccount(&accountAddress)
	if err != nil {
		vitelog.Logger.Error(errors2.Wrap(err, "accountController.detail"))
		util.RespondSuccess(c, nil,"")
		return
	}


	util.RespondSuccess(c, account,"")
}

func NewTestToken (c *gin.Context) {
	var accountNewTestToken typeRequest.AccountNewTestToken
	if err := c.Bind(&accountNewTestToken); err != nil {
		util.RespondError(c, 400, err)
		return
	}

	if !types.IsValidHexAddress(accountNewTestToken.AccountAddress) {
		util.RespondFailed(c, 1, errors.New("AccountAddress is invalid"), "")
		return
	}

	toAddr, _ := types.HexToAddress(accountNewTestToken.AccountAddress)

	vite := c.MustGet("vite").(*vite.Vite)
	genesisAddr := c.MustGet("genesisAddr").(types.Address)
	randomAmount := rand.Intn(150)

	amount := big.NewInt(50)
	amount.Add(amount, big.NewInt(int64(randomAmount)))

	for i := 0; i < 18; i++ {
		amount.Mul(amount, big.NewInt(10))
	}

	createTxErr := vite.Ledger().Ac().CreateTx(&ledger.AccountBlock{
		AccountAddress: &genesisAddr,
		To: &toAddr,
		TokenId: &ledger.MockViteTokenId,
		Amount: amount,
	})

	if createTxErr != nil {
		util.RespondFailed(c, 2, errors.New("Create transaction failed. Error is " + createTxErr.Error()), "")
		return
	}
	util.RespondSuccess(c, response.NewNewTestToken(amount, ledger.MockViteTokenId),"")
}