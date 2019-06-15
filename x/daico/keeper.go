package daico

import (

	"github.com/rs/xid"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper struct {
	proposalStoreKey    sdk.StoreKey
	participantStoreKey sdk.StoreKey
	poolStorKey sdk.StoreKey
	transientStoreKey   sdk.StoreKey
	paramsKeeper        params.Keeper
	paramsCodeSpace     params.Subspace
	cdc                 *codec.Codec
}

func NewKeeper(sKey, tKey, bKey, pkey sdk.StoreKey, pKeeper params.Keeper, pCodeSpace params.Subspace, cdc *codec.Codec) Keeper {
	return Keeper{
		proposalStoreKey:    sKey,
		participantStoreKey: bKey,
		poolStorKey: pkey,
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
func (k Keeper) GetProposal(ctx sdk.Context, id string) (Proposal, error ) {
	store := ctx.KVStore(k.proposalStoreKey)
	i := store.Get([]byte(id))
	if i == nil {
		panic("Error")
	}
	var pa Proposal
	k.cdc.MustUnmarshalBinaryBare(i, &pa)
	return pa, nil
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

func (k Keeper) GetParticipant(ctx sdk.Context, id string )(pa Participant, ok bool) {
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
func (k Keeper) SetParticipant(ctx sdk.Context, pa Participant){
	store := ctx.KVStore(k.participantStoreKey)
	store.Set([]byte(pa.GetID()), k.cdc.MustMarshalBinaryBare(pa))
}


// Become a backer
func (k Keeper) BecomeParticipant(ctx sdk.Context, value sdk.Coins, participantName string, participantAddr sdk.AccAddress, proposalID string) {
	pro, err := k.GetProposal(ctx, proposalID)
	if err != nil {
		panic("Proposal doesnt exist")
	}
	
	guid := xid.New().String()
	participant := NewParticipant(guid, participantName, []string{pro.GetID()}, participantAddr, value)
	k.SetParticipant(ctx, participant)

	// need to add value to the pool that 

	pro.Participants = append(pro.Participants, participant.GetAddress())
	// increment value in pool
	k.SetProposal(ctx, pro)
}

// Send funds

func (k Keeper) GetFeePool(ctx sdk.Context, ID string)  (pool Pool) {
	store := ctx.KVStore(k.poolStorKey)
	i := store.Get([]byte(ID))
	if i == nil {
		panic("Pool does not exist")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(i, &pool)
	return 
}
