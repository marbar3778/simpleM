package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
Types of handlers

1. Create a proposal
2. Become a participant of a proposal
3. Vote on funds being sent out
4. Vote on people getting there funds back

*/

type MsgCreateProposal struct {
	Description     string         `json:"description"`
	CompanyName     string         `json:"company_name"`
	ProposerAddress sdk.AccAddress `json:"proposerAddress"`
}

func (msg MsgCreateProposal) Route() string { return RouterKey }
func (msg MsgCreateProposal) Type() string  { return "create_proposal" }

func (msg MsgCreateProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ProposerAddress}
}

func (msg MsgCreateProposal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCreateProposal) ValidateBasic() sdk.Error {
	if len(msg.Description) == 0 {
		return sdk.ErrUnknownRequest("Description cannot be empty")
	}
	if len(msg.CompanyName) == 0 {
		return sdk.ErrUnknownRequest("Company name cannot be empty")
	}
	if msg.ProposerAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.ProposerAddress.String())
	}
	return nil
}
