package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength  = 70
	MaxIdentityLength = 3000
	MaxWebsiteLength  = 140
	MaxDetailsLength  = 280
)

// Validator defines the total amount of bond shares and their exchange rate to
// coins. Slashing results in a decrease in the exchange rate, allowing correct
// calculation of future undelegations without iterating over delegators.
// When coins are delegated to this validator, the validator is credited with a
// delegation whose number of bond shares is based on the amount of coins delegated
// divided by the current exchange rate. Voting power can be calculated as total
// bonded shares multiplied by exchange rate.
type Authority struct {
	OperatorAddress         sdk.ValAddress `json:"operator_address"` // address of the validator's operator; bech encoded in JSON
	ConsPubKey              crypto.PubKey  `json:"consensus_pubkey"` // the consensus public key of the validator; bech encoded in JSON
	Jailed                  bool           `json:"jailed"`           // has the validator been jailed from bonded status?
	Status                  sdk.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Description             Description    `json:"description"`      // description terms for the validator
	UnbondingHeight         int64          `json:"unbonding_height"` // if unbonding, height at which this validator has begun unbonding
	UnbondingCompletionTime time.Time      `json:"unbonding_time"`   // if unbonding, min time for the validator to complete unbonding
	Commission              Commission     `json:"commission"`       // commission parameters
	Power                   sdk.Int        `json:"power"`            // Power associated with POA, will be default unless set
}

// Validators is a collection of Validator
type Authorities []Authority

func (v Authorities) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// ToSDKValidators -  convenience function convert []Validators to []sdk.Validators
func (a Authorities) ToSDKValidators() (authorities []Authority) {
	for _, val := range a {
		authorities = append(authorities, val)
	}
	return
}

// NewValidator - initialize a new validator
func NewAuthority(operator sdk.ValAddress, pubKey crypto.PubKey, description Description) Authority {
	return Authority{
		OperatorAddress:         operator,
		ConsPubKey:              pubKey,
		Jailed:                  false,
		Status:                  sdk.Unbonded,
		Description:             description,
		UnbondingHeight:         int64(0),
		UnbondingCompletionTime: time.Unix(0, 0).UTC(),
		Power:                   sdk.ZeroInt(),
	}
}

// return the redelegation
func MustMarshalValidator(cdc *codec.Codec, authority Authority) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(authority)
}

// unmarshal a redelegation from a store value
func MustUnmarshalValidator(cdc *codec.Codec, value []byte) Authority {
	authority, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}
	return authority
}

// unmarshal a redelegation from a store value
func UnmarshalValidator(cdc *codec.Codec, value []byte) (authority Authority, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &authority)
	return
}

// String returns a human readable string representation of a validator.
func (a Authority) String() string {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(a.ConsPubKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`Validator
  Operator Address:           %s
  Validator Consensus Pubkey: %s
  Jailed:                     %v
  Status:                     %s
  Power:                     	%s
  Description:                %s
  Unbonding Height:           %d
  Unbonding Completion Time:  %v
  Commission:                 %s`, a.OperatorAddress, bechConsPubKey,
		a.Jailed, a.Status, a.Power,
		a.Description, a.UnbondingHeight, a.UnbondingCompletionTime, a.Commission)
}

// this is a helper struct used for JSON de- and encoding only
type bechValidator struct {
	OperatorAddress         sdk.ValAddress `json:"operator_address"` // the bech32 address of the validator's operator
	ConsPubKey              string         `json:"consensus_pubkey"` // the bech32 consensus public key of the validator
	Jailed                  bool           `json:"jailed"`           // has the validator been jailed from bonded status?
	Status                  sdk.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	DelegatorShares         sdk.Dec        `json:"delegator_shares"` // total shares issued to a validator's delegators
	Description             Description    `json:"description"`      // description terms for the validator
	UnbondingHeight         int64          `json:"unbonding_height"` // if unbonding, height at which this validator has begun unbonding
	UnbondingCompletionTime time.Time      `json:"unbonding_time"`   // if unbonding, min time for the validator to complete unbonding
	Commission              Commission     `json:"commission"`       // commission parameters
	Power                   sdk.Int        `json:"power"`            // power of authority`
	// Tokens                  sdk.Int        `json:"tokens"`              // delegated tokens (incl. self-delegation)
	// MinSelfDelegation       sdk.Int        `json:"min_self_delegation"` // minimum self delegation
}

// MarshalJSON marshals the validator to JSON using Bech32
func (a Authority) MarshalJSON() ([]byte, error) {
	bechConsPubKey, err := sdk.Bech32ifyConsPub(a.ConsPubKey)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSON(bechValidator{
		OperatorAddress:         a.OperatorAddress,
		ConsPubKey:              bechConsPubKey,
		Jailed:                  a.Jailed,
		Status:                  a.Status,
		Description:             a.Description,
		UnbondingHeight:         a.UnbondingHeight,
		UnbondingCompletionTime: a.UnbondingCompletionTime,
		Commission:              a.Commission,
		Power:                   a.Power,
		// MinSelfDelegation:       v.MinSelfDelegation,
		// Tokens:                  v.Tokens,
	})
}

// UnmarshalJSON unmarshals the validator from JSON using Bech32
func (a *Authority) UnmarshalJSON(data []byte) error {
	bv := &bechValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, bv); err != nil {
		return err
	}
	consPubKey, err := sdk.GetConsPubKeyBech32(bv.ConsPubKey)
	if err != nil {
		return err
	}
	*a = Authority{
		OperatorAddress:         bv.OperatorAddress,
		ConsPubKey:              consPubKey,
		Jailed:                  bv.Jailed,
		Status:                  bv.Status,
		Description:             bv.Description,
		UnbondingHeight:         bv.UnbondingHeight,
		UnbondingCompletionTime: bv.UnbondingCompletionTime,
		Commission:              bv.Commission,
	}
	return nil
}

// only the vitals
func (a Authority) TestEquivalent(v2 Authority) bool {
	return a.ConsPubKey.Equals(v2.ConsPubKey) &&
		bytes.Equal(a.OperatorAddress, v2.OperatorAddress) &&
		a.Status.Equal(v2.Status) &&
		a.Power.Equal(v2.Power) &&
		a.Description == v2.Description &&
		a.Commission.Equal(v2.Commission)
}

// return the TM validator address
func (v Authority) ConsAddress() sdk.ConsAddress {
	return sdk.ConsAddress(v.ConsPubKey.Address())
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

// NewDescription returns a new Description with the provided values.
func NewDescription(moniker, identity, website, details string) Description {
	return Description{
		Moniker:  moniker,
		Identity: identity,
		Website:  website,
		Details:  details,
	}
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, sdk.Error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}
	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}
	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}
	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return Description{
		Moniker:  d2.Moniker,
		Identity: d2.Identity,
		Website:  d2.Website,
		Details:  d2.Details,
	}.EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, sdk.Error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, ErrDescriptionLength(DefaultCodespace, "moniker", len(d.Moniker), MaxMonikerLength)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, ErrDescriptionLength(DefaultCodespace, "identity", len(d.Identity), MaxIdentityLength)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, ErrDescriptionLength(DefaultCodespace, "website", len(d.Website), MaxWebsiteLength)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, ErrDescriptionLength(DefaultCodespace, "details", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Authority) ABCIValidatorUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
		Power:  v.TendermintPower(),
	}
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Authority) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.ConsPubKey),
		Power:  0,
	}
}

// UpdateStatus updates the location of the shares within a validator
// to reflect the new status
func (v Authority) UpdateStatus(pool Pool, NewStatus sdk.BondStatus) (Authority, Pool) {

	switch v.Status {
	case sdk.Unbonded:

		switch NewStatus {
		case sdk.Unbonded:
			return v, pool
		case sdk.Bonded:
			pool = pool.notBondedTokensToBonded(v.Power)
		}
	case sdk.Unbonding:

		switch NewStatus {
		case sdk.Unbonding:
			return v, pool
		case sdk.Bonded:
			pool = pool.notBondedTokensToBonded(v.Power)
		}
	case sdk.Bonded:

		switch NewStatus {
		case sdk.Bonded:
			return v, pool
		default:
			pool = pool.bondedTokensToNotBonded(v.Power)
		}
	}

	v.Status = NewStatus
	return v, pool
}

// // removes tokens from a validator
// func (v Validator) RemoveTokens(pool Pool, tokens sdk.Int) (Validator, Pool) {
// 	if tokens.IsNegative() {
// 		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
// 	}
// 	if v.Tokens.LT(tokens) {
// 		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Power, tokens))
// 	}
// 	v.Tokens = v.Tokens.Sub(tokens)
// 	// TODO: It is not obvious from the name of the function that this will happen. Either justify or move outside.
// 	if v.Status == sdk.Bonded {
// 		pool = pool.bondedTokensToNotBonded(tokens)
// 	}
// 	return v, pool
// }

// SetInitialCommission attempts to set a validator's initial commission. An
// error is returned if the commission is invalid.
func (v Authority) SetInitialCommission(commission Commission) (Authority, sdk.Error) {
	if err := commission.Validate(); err != nil {
		return v, err
	}

	v.Commission = commission
	return v, nil
}

// AddTokensFromDel adds tokens to a validator
// CONTRACT: Tokens are assumed to have come from not-bonded pool.
func (v Authority) AddTokensFromDel(pool Pool, amount sdk.Int) (Authority, Pool, sdk.Dec) {

	// calculate the shares to issue
	var issuedShares sdk.Dec

	if v.Status == sdk.Bonded {
		pool = pool.notBondedTokensToBonded(amount)
	}

	v.Power = v.Power.Add(amount)

	return v, pool, issuedShares
}

// RemoveDelShares removes delegator shares from a validator.
// NOTE: because token fractions are left in the valiadator,
//       the exchange rate of future shares of this validator can increase.
// CONTRACT: Tokens are assumed to move to the not-bonded pool.
// func (v Validator) RemoveDelShares(pool Pool, delShares sdk.Dec) (Validator, Pool, sdk.Int) {

// 	remainingShares := v.DelegatorShares.Sub(delShares)
// 	var issuedTokens sdk.Int
// 	if remainingShares.IsZero() {

// 		// last delegation share gets any trimmings
// 		issuedTokens = v.Tokens
// 		v.Tokens = sdk.ZeroInt()
// 	} else {

// 		// leave excess tokens in the validator
// 		// however fully use all the delegator shares
// 		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
// 		v.Tokens = v.Tokens.Sub(issuedTokens)
// 		if v.Tokens.IsNegative() {
// 			panic("attempting to remove more tokens than available in validator")
// 		}
// 	}

// 	v.DelegatorShares = remainingShares
// 	if v.Status == sdk.Bonded {
// 		pool = pool.bondedTokensToNotBonded(issuedTokens)
// 	}

// 	return v, pool, issuedTokens
// }

// SharesFromTokensTruncated returns the truncated shares of a delegation given
// a bond amount. It returns an error if the validator has no tokens.
// func (v Validator) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, sdk.Error) {
// 	if v.Power.IsZero() {
// 		return sdk.ZeroDec(), ErrInsufficientShares(DefaultCodespace)
// 	}

// 	return v.GetDelegatorShares().MulInt(amt).QuoTruncate(v.GetTokens().ToDec()), nil
// }

// // get the bonded tokens which the validator holds
// func (v Validator) BondedTokens() sdk.Int {
// 	if v.Status == sdk.Bonded {
// 		return v.Tokens
// 	}
// 	return sdk.ZeroInt()
// }

// get the Tendermint Power
// a reduction of 10^6 from validator tokens is applied
func (v Authority) TendermintPower() int64 {
	if v.Status == sdk.Bonded {
		return v.PotentialTendermintPower()
	}
	return 0
}

// potential Tendermint power
func (v Authority) PotentialTendermintPower() int64 {
	return sdk.TokensToTendermintPower(v.Power)
}

// ensure fulfills the sdk validator types
// var _ sdk.Validator = Validator{}

// nolint - for sdk.Validator
func (v Authority) IsJailed() bool               { return v.Jailed }
func (v Authority) GetMoniker() string           { return v.Description.Moniker }
func (v Authority) GetStatus() sdk.BondStatus    { return v.Status }
func (v Authority) GetOperator() sdk.ValAddress  { return v.OperatorAddress }
func (v Authority) GetConsPubKey() crypto.PubKey { return v.ConsPubKey }
func (v Authority) GetConsAddr() sdk.ConsAddress { return sdk.ConsAddress(v.ConsPubKey.Address()) }
func (v Authority) GetPower() sdk.Int            { return v.Power }
func (v Authority) GetTendermintPower() int64    { return v.TendermintPower() }
func (v Authority) GetCommission() sdk.Dec       { return v.Commission.Rate }

// func (v Validator) GetBondedTokens() sdk.Int { return v.BondedTokens() }

// func (v Validator) GetMinSelfDelegation() sdk.Int { return v.MinSelfDelegation }
