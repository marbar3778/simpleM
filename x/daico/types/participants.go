package types

import (
	"fmt"
	"github.com/rs/xid"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Participant struct {
	Id                 string         `json:"id"`
	UserName           string         `json:"username"`
	UserAddress        sdk.AccAddress `json:"user_address"`
	AmountParticipated sdk.Coins      `json:"amount_participated"` // coins represent vote power
	ProposalReference  []string       `json:"proposalReference"`
}

func NewParticipant(userName string, pRef []string, userAddress sdk.AccAddress, amountParticipated sdk.Coins) Participant {
	guid := xid.New().String()
	return Participant{
		Id:                 guid,
		UserName:           userName,
		UserAddress:        userAddress,
		AmountParticipated: amountParticipated,
		ProposalReference:  pRef,
	}
}

// fmt.Stringer
func (p Participant) String() string {
	strings := fmt.Sprintf(`
	id: %s,
	userName: %s,
	userAddress: %s,
	amountParticipated: %s`, p.Id, p.UserName,
		p.UserAddress.String(),
		p.AmountParticipated.String())

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
