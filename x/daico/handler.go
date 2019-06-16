package daico

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateProposal:
			return handleMsgCreateProposal(ctx, k, msg)
		case MsgBecomeParticipant:
			return handleMsgBecomeParticipant(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateProposal(ctx sdk.Context, k Keeper, msg MsgCreateProposal) sdk.Result {

	pro := k.CreateProposal(ctx, msg.Description, msg.CompanyName, msg.ProposerAddress, msg.TokenDenom, msg.TokenAmount)

	return sdk.Result{
		sdk.NewTags(
			Category, TxCategory,
			Proposer, msg.ProposerAddress.String(),
			ProposerID, pro.GetID(),
		),
	}
}

// Become a pariticpant in a Proposal
func handleMsgBecomeParticipant(ctx sdk.Context, k Keeper, msg MsgBecomeParticipant) sdk.Result {
	participant := k.BecomeParticipant(ctx, msg.UserName, msg.UserAddress, msg.ProposalReference, msg.ParticipantAmount)

	return sdk.Result{
		sdk.NewTags(
			Category, TxCategory,
			ParticipantTag, msg.UserAddress,
		)
	}
}
