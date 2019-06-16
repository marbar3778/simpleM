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

// Create a Proposal to create funds
type MsgCreateProposal struct {
	Description     string         `json:"description"`
	CompanyName     string         `json:"company_name"`
	ProposerAddress sdk.AccAddress `json:"proposerAddress"`
	TokenDenom      string         `json:"token_denom"`
	TokenAmount     int            `json:"token_amount"`
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
	if len(msg.TokenDenom) == 0 {
		return sdk.ErrUnknownRequest("Token Denom is not defined")
	}
	if msg.TokenAmount == 0 {
		return sdk.ErrUnknownRequest("Token amount must be greater than 0")
	}
	return nil
}

// Become a participant
type MsgBecomeParticipant struct {
	UserName          string            `json:"username"`
	UserAddress       sdk.AccAddress    `json:"user_address"`
	ProposalReference ProposalReference `json:"proposalReference"`
	ParticipantAmount sdk.Coins         `json:"participant_amount"`
}

func (msg MsgBecomeParticipant) Route() string { return RouterKey }
func (msg MsgBecomeParticipant) Type() string  { return "become_participant" }

func (msg MsgBecomeParticipant) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.UserAddress}
}

func (msg MsgBecomeParticipant) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBecomeParticipant) ValidateBasic() sdk.Error {
	if len(msg.ProposalReference.Reference) == 0 {
		return sdk.ErrUnknownRequest("Proposal ID is not present")
	}
	if len(msg.UserName) == 0 {
		return sdk.ErrUnknownRequest("User name is not present")
	}
	if msg.UserAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.UserAddress.String())
	}
	if msg.ParticipantAmount.IsZero() {
		return sdk.ErrInsufficientCoins("Coins needs to be greater than zero")
	}
	return nil
}

// Vote on
