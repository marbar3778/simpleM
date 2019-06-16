package daico

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper struct {
	proposalStoreKey    sdk.StoreKey
	participantStoreKey sdk.StoreKey
	poolStorKey         sdk.StoreKey
	transientStoreKey   sdk.StoreKey
	paramsKeeper        params.Keeper
	paramsCodeSpace     params.Subspace
	cdc                 *codec.Codec
}

func NewKeeper(sKey, tKey, bKey, pkey sdk.StoreKey, pKeeper params.Keeper, pCodeSpace params.Subspace, cdc *codec.Codec) Keeper {
	return Keeper{
		proposalStoreKey:    sKey,
		participantStoreKey: bKey,
		poolStorKey:         pkey,
		transientStoreKey:   tKey,
		paramsKeeper:        pKeeper,
		paramsCodeSpace:     pCodeSpace,
	}
}

var (
	pKeyPrefix = []byte{0x00} // key for NFT collections
)

// -------------------------------------
// Proposals

// Get Proposal
func (k Keeper) GetProposal(ctx sdk.Context, id string) (p Proposal, ok bool) {
	store := ctx.KVStore(k.proposalStoreKey)
	i := store.Get([]byte(id))
	if i == nil {
		return Proposal{}, false
	}
	var pa Proposal
	k.cdc.MustUnmarshalBinaryBare(i, &pa)
	return pa, true
}

// Add a proposal to the store
func (k Keeper) SetProposal(ctx sdk.Context, pr Proposal) {
	store := ctx.KVStore(k.proposalStoreKey)
	store.Set([]byte(pr.ID), k.cdc.MustMarshalBinaryBare(pr))
}

func (k Keeper) CreateProposal(ctx sdk.Context, pr Proposal) {

}

// Remove the proposal form the store
func (k Keeper) RemoveProposal(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.proposalStoreKey)
	store.Delete([]byte(id))
}

// -------------------------------------
// Participant

// Get participant

func (k Keeper) GetParticipant(ctx sdk.Context, id string) (pa Participant, ok bool) {
	store := ctx.KVStore(k.participantStoreKey)
	i := store.Get([]byte(id))
	if i == nil {
		return Participant{}, false
	}
	k.cdc.MustUnmarshalBinaryBare(i, &pa)
	return pa, true
}

// // Iterate all the participants
// func (k Keeper) IterateParticipants(ctx sdk.Context, handler func(pa Participants) (stop bool)) {
// 	store := ctx.KVStore(k.participantStoreKey)
// 	iterator := sdk.KVStorePrefixIterator(store, pKeyPrefix)
// 	defer iterator.Close();

// 	for ; iterator.Valid(); iterator.Next() {
// 		var participant Participant
// 		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &participant)
// 		if handler(participant) {
// 			break
// 		}
// 	}
// }

// // Get Participant
// func (k Keeper) GetParticipants(ctx sdk.Context) (participants Participants ) {
// 	k.IterateParticipants(ctx,
// 	func (pa Participants)(stop bool) {
// 		participants = append(participants, pa)
// 		return false
// 	})
// 	return
// }

// Set Participants
func (k Keeper) SetParticipant(ctx sdk.Context, pa Participant) {
	store := ctx.KVStore(k.participantStoreKey)
	store.Set([]byte(pa.GetID()), k.cdc.MustMarshalBinaryBare(pa))
}

// New backer
func (k Keeper) NewParticipant(ctx sdk.Context, pa Participant) Participant {
	participant := NewParticipant(pa.UserName, []ProposalReference, pa.UserAddress)

	k.SetParticipant(ctx, participant)
	return participant
}

// Become a backer to a proposal
func (k Keeper) BecomeParticipant(ctx sdk.Context, pa Participant, proposalID string) {
	p, ok := k.GetParticipant(ctx, pa.GetID())
	if !ok {
		p = k.NewParticipant(ctx, pa)
		k.SetParticipant(ctx, p)
	}

	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok {
		panic("Proposal does not exist")
	}

	proposal.Participants = append(proposal.Participants, p.GetAddress())
	// increment value in pool
	k.SetProposal(ctx, proposal)

}

// Get Pool
func (k Keeper) GetPool(ctx sdk.Context, ID string) (pool Pool) {
	store := ctx.KVStore(k.poolStorKey)
	i := store.Get([]byte(ID))
	if i == nil {
		panic("Pool does not exist")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(i, &pool)
	return
}

// set the pool with the funds
func (k Keeper) SetPool(ctx sdk.Context, pool Pool) {
	store := ctx.KVStore(k.poolStorKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(pool)
	store.Set([]byte(pool.GetID()), b)
}

// Destroy Pool
func (k Keeper) DestroyPool(ctx sdk.Context, pool Pool) {
	// iterate through all the participants in the proposal
	// allocated / raised,
	// percentage of pa funds given is returned in respect to what was not allocated.
}
