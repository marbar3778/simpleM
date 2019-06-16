package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transaction tags for NFT messages
var (
	TxCategory  = "daico"
	Category    = sdk.TagCategory
	Proposer    = sdk.TagSender
	ProposerID  = "proposer-id"
	Participant = sdk.TagSender
)
