package cli

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/marbar3778/simpleM/x/simplePOA/types"
)

func buildCommissionMsg(rateStr, maxRateStr, maxChangeRateStr string) (commission types.CommissionMsg, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, errors.New("must specify all validator commission parameters")
	}

	rate, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := sdk.NewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := sdk.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = types.NewCommissionMsg(rate, maxRate, maxChangeRate)
	return commission, nil
}
