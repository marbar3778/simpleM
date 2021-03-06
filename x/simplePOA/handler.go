package staking

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/marbar3778/simpleM/x/simplePOA/keeper"
	"github.com/marbar3778/simpleM/x/simplePOA/tags"
	"github.com/marbar3778/simpleM/x/simplePOA/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case types.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, msg, k)

		case types.MsgEditValidator:
			return handleMsgEditValidator(ctx, msg, k)

		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) ([]abci.ValidatorUpdate, sdk.Tags) {
	resTags := sdk.NewTags()

	// Calculate validator set changes.
	validatorUpdates := k.ApplyAndReturnValidatorSetUpdates(ctx)

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	return validatorUpdates, resTags
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgCreateValidator(ctx sdk.Context, msg types.MsgCreateValidator, k keeper.Keeper) sdk.Result {
	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetValidator(ctx, msg.ValidatorAddress); found {
		return ErrValidatorOwnerExists(k.Codespace()).Result()
	}

	if _, found := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); found {
		return ErrValidatorPubKeyExists(k.Codespace()).Result()
	}

	// if msg.Value.Denom != k.GetParams(ctx).BondDenom {
	// 	return ErrBadDenom(k.Codespace()).Result()
	// }

	if _, err := msg.Description.EnsureLength(); err != nil {
		return err.Result()
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PubKey)
		if !common.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return ErrValidatorPubKeyTypeNotSupported(k.Codespace(),
				tmPubKey.Type,
				ctx.ConsensusParams().Validator.PubKeyTypes).Result()
		}
	}

	validator := NewValidator(msg.ValidatorAddress, msg.PubKey, msg.Description)

	// validator.MinSelfDelegation = msg.MinSelfDelegation // need to remove msgSelfdelegation

	k.SetValidator(ctx, validator)
	k.SetValidatorByConsAddr(ctx, validator)
	k.SetNewValidatorByPowerIndex(ctx, validator)

	// call the after-creation hook
	k.AfterValidatorCreated(ctx, validator.OperatorAddress)

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	_, err = k.Delegate(ctx, msg.DelegatorAddress, msg.Value.Amount, validator)
	if err != nil {
		return err.Result()
	}

	resTags := sdk.NewTags(
		tags.Category, tags.TxCategory,
		tags.Sender, msg.DelegatorAddress.String(),
		tags.DstValidator, msg.ValidatorAddress.String(),
	)

	return sdk.Result{
		Tags: resTags,
	}
}

func handleMsgEditValidator(ctx sdk.Context, msg types.MsgEditValidator, k keeper.Keeper) sdk.Result {
	// validator must already be registered
	validator, found := k.GetValidator(ctx, msg.ValidatorAddress)
	if !found {
		return ErrNoValidatorFound(k.Codespace()).Result()
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := validator.Description.UpdateDescription(msg.Description)
	if err != nil {
		return err.Result()
	}

	validator.Description = description

	if msg.CommissionRate != nil {
		commission, err := k.UpdateValidatorCommission(ctx, validator, *msg.CommissionRate)
		if err != nil {
			return err.Result()
		}

		// call the before-modification hook since we're about to update the commission
		k.BeforeValidatorModified(ctx, msg.ValidatorAddress)

		validator.Commission = commission
	}

	// if msg.MinSelfDelegation != nil {
	// 	if !(*msg.MinSelfDelegation).GT(validator.MinSelfDelegation) {
	// 		return ErrMinSelfDelegationDecreased(k.Codespace()).Result()
	// 	}
	// 	if (*msg.MinSelfDelegation).GT(validator.Tokens) {
	// 		return ErrSelfDelegationBelowMinimum(k.Codespace()).Result()
	// 	}
	// 	validator.MinSelfDelegation = (*msg.MinSelfDelegation)
	// }

	k.SetValidator(ctx, validator)

	resTags := sdk.NewTags(
		tags.Category, tags.TxCategory,
		tags.Sender, msg.ValidatorAddress.String(),
	)

	return sdk.Result{
		Tags: resTags,
	}
}
