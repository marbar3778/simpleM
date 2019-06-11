package types

import (
	"bytes"
	"encoding/json"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgCreateAuthority{}
	_ sdk.Msg = &MsgEditAuthority{}
)

//______________________________________________________________________

// MsgCreateValidator - struct for bonding transactions
type MsgCreateAuthority struct {
	Description      Description    `json:"description"`
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	PubKey           crypto.PubKey  `json:"pubkey"`
	Power            sdk.Int        `json:"power"`
}

type msgCreateAuthorityJSON struct {
	Description      Description    `json:"description"`
	Commission       CommissionMsg  `json:"commission"`
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	PubKey           string         `json:"pubkey"`
	Power            sdk.Int        `json:"power"`
}

// Default way to create validator. Delegator address and validator address are the same
func NewMsgCreateAuthority(
	valAddr sdk.ValAddress, pubKey crypto.PubKey, power sdk.Int,
	description Description, commission CommissionMsg,
) MsgCreateAuthority {

	// // if no power is set then is is defaulted to 1
	// if power.IsZero() {
	// 	power = sdk.Int{1} // TODO: make 1 a *big.Int
	// }

	return MsgCreateAuthority{
		Description:      description,
		DelegatorAddress: sdk.AccAddress(valAddr),
		ValidatorAddress: valAddr,
		PubKey:           pubKey,
		Power:            power,
	}
}

//nolint
func (msg MsgCreateAuthority) Route() string { return RouterKey }
func (msg MsgCreateAuthority) Type() string  { return "create_validator" }

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgCreateAuthority) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	addrs := []sdk.AccAddress{msg.DelegatorAddress}

	if !bytes.Equal(msg.DelegatorAddress.Bytes(), msg.ValidatorAddress.Bytes()) {
		// if validator addr is not same as delegator addr, validator must sign
		// msg as well
		addrs = append(addrs, sdk.AccAddress(msg.ValidatorAddress))
	}
	return addrs
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON
// serialization of the MsgCreateAuthority type.
func (msg MsgCreateAuthority) MarshalJSON() ([]byte, error) {
	return json.Marshal(msgCreateAuthorityJSON{
		Description:      msg.Description,
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
		PubKey:           sdk.MustBech32ifyConsPub(msg.PubKey),
		Power:            msg.Power,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface to provide custom
// JSON deserialization of the MsgCreateAuthority type.
func (msg *MsgCreateAuthority) UnmarshalJSON(bz []byte) error {
	var msgCreateValJSON msgCreateAuthorityJSON
	if err := json.Unmarshal(bz, &msgCreateValJSON); err != nil {
		return err
	}

	msg.Description = msgCreateValJSON.Description
	msg.DelegatorAddress = msgCreateValJSON.DelegatorAddress
	msg.ValidatorAddress = msgCreateValJSON.ValidatorAddress
	var err error
	msg.PubKey, err = sdk.GetConsPubKeyBech32(msgCreateValJSON.PubKey)
	if err != nil {
		return err
	}
	msg.Power = msgCreateValJSON.Power

	return nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateAuthority) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgCreateAuthority) ValidateBasic() sdk.Error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAddress.Empty() {
		return ErrNilDelegatorAddr(DefaultCodespace)
	}
	if msg.ValidatorAddress.Empty() {
		return ErrNilValidatorAddr(DefaultCodespace)
	}
	if !sdk.AccAddress(msg.ValidatorAddress).Equals(msg.DelegatorAddress) {
		return ErrBadValidatorAddr(DefaultCodespace)
	}
	if msg.Power.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount(DefaultCodespace)
	}
	if msg.Description == (Description{}) {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "description must be included")
	}

	return nil
}

// MsgEditValidator - struct for editing a validator
type MsgEditAuthority struct {
	Description
	AuthorityAddress sdk.ValAddress `json:"address"`
}

func NewMsgEditAuthority(authorityAddr sdk.ValAddress, description Description, newRate *sdk.Dec, newMinSelfDelegation *sdk.Int) MsgEditAuthority {
	return MsgEditAuthority{
		Description:      description,
		AuthorityAddress: authorityAddr,
	}
}

//nolint
func (msg MsgEditAuthority) Route() string { return RouterKey }
func (msg MsgEditAuthority) Type() string  { return "edit_validator" }
func (msg MsgEditAuthority) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.AuthorityAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgEditAuthority) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgEditAuthority) ValidateBasic() sdk.Error {
	if msg.AuthorityAddress.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil validator address")
	}

	if msg.Description == (Description{}) {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "transaction must include some information to modify")
	}

	return nil
}

// TODO: check how a validator should be removed.
// func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) MsgUndelegate {
// 	return MsgUndelegate{
// 		DelegatorAddress: delAddr,
// 		ValidatorAddress: valAddr,
// 		Amount:           amount,
// 	}
// }

// //nolint
// func (msg MsgUndelegate) Route() string                { return RouterKey }
// func (msg MsgUndelegate) Type() string                 { return "begin_unbonding" }
// func (msg MsgUndelegate) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.DelegatorAddress} }

// // get the bytes for the message signer to sign on
// func (msg MsgUndelegate) GetSignBytes() []byte {
// 	bz := ModuleCdc.MustMarshalJSON(msg)
// 	return sdk.MustSortJSON(bz)
// }

// // quick validity check
// func (msg MsgUndelegate) ValidateBasic() sdk.Error {
// 	if msg.DelegatorAddress.Empty() {
// 		return ErrNilDelegatorAddr(DefaultCodespace)
// 	}
// 	if msg.ValidatorAddress.Empty() {
// 		return ErrNilValidatorAddr(DefaultCodespace)
// 	}
// 	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
// 		return ErrBadSharesAmount(DefaultCodespace)
// 	}
// 	return nil
// }
