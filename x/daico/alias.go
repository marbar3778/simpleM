package daico

import (
	"github.com/marbar3778/simpleM/x/daico/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	ModuleCdc       = types.ModuleCdc
	RegosterCodec   = types.RegisterCodec
	NewParticipant  = types.NewParticipant
	NewParticipants = types.NewParticipants
)

type (
	Participant       = types.Participant
	Participants      = types.Participants
	Proposal          = types.Proposal
	Pool              = types.Pool
	ProposalReference = types.ProposalReference
)
