package types

import (
	"fmt"
	"github.com/rs/xid"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Pool struct {
	ID string
	Denom string `json:"denom"`
	AllocatedFunds sdk.DecCoins `json:"allocated_funds"`
	RaisedFunds    sdk.DecCoins `json:"raised_funds"`
}

func NewPool(denom string) Pool {
	guid := xid.New().String()
	return Pool{
		ID: guid,
		Denom: denom,
		AllocatedFunds: sdk.DecCoins{},
		RaisedFunds: sdk.DecCoins{},
	}
}

func (p Pool) String() string {
	strings := fmt.Sprintf(`
		ID: %s,
		denomination: %s,
		allocated fund amount: %s,
		raised funds: %s,
	`, p.ID, p.Denom, p.AllocatedFunds.String(), p.RaisedFunds.String())
	return strings
}

func (p Pool) GetID() string {return p.ID}