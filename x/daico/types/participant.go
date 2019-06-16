package types

import (
	"fmt"

	"github.com/rs/xid"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProposalReference struct {
	Reference string    `json:"reference"`
	Amount    sdk.Coins `json:"amount"` // This is the amount the Partiicpant added into the pool
}

type Participant struct {
	Id                string              `json:"id"`
	UserName          string              `json:"username"`
	UserAddress       sdk.AccAddress      `json:"user_address"`
	ProposalReference []ProposalReference `json:"proposalReference"`
}

func NewParticipant(userName string, pRef []ProposalReference, userAddress sdk.AccAddress, amountParticipated sdk.Coins) Participant {
	guid := xid.New().String()
	return Participant{
		Id:                guid,
		UserName:          userName,
		UserAddress:       userAddress,
		ProposalReference: pRef,
	}
}

// fmt.Stringer
func (p Participant) String() string {
	strings := fmt.Sprintf(`
	id: %s,
	userName: %s,
	userAddress: %s,
	`, p.Id, p.UserName,
		p.UserAddress.String())

	return strings
}

func (p Participant) GetID() string              { return p.Id }
func (p Participant) GetAddress() sdk.AccAddress { return p.UserAddress }

type Participants []Participant

func NewParticipants(participants ...Participant) Participants {
	if len(participants) == 0 {
		return Participants{}
	}
	return Participants(participants)
}

func (ps *Participants) Add(p Participant) {
	*ps = append(*ps, p)
}
