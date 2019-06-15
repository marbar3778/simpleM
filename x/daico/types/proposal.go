package types

import (
	"fmt"
	"github.com/rs/xid"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Proposal struct {
	ID                     string           `json:"id"`
	Description            string           `json:"description"`
	CompanyName            string           `json:"company_name"`
	FundsUsed              string           `json:"funds_used"`
	ProposerAddress        sdk.AccAddress   `json:"proposerAddress"`
	Participants           []sdk.AccAddress `json:"partiicpants"`
	PoolReference          string           `json:"pool_reference"`
}

func NewProposal(description, companyName, fundsUsed string, proposerAddress sdk.AccAddress, pariticpants []sdk.AccAddress) Proposal {
	guid := xid.New().String()
	return Proposal{
		ID: guid,
		Description:            description,
		CompanyName:            companyName,
		FundsUsed:              fundsUsed,
		ProposerAddress:        proposerAddress,
		Participants:           pariticpants,
		PoolReference:          "",
	}
}

func (p Proposal) GetID() string {
	return p.ID
}
func (p Proposal) String() string {
	strings := fmt.Sprintf(`
	Descriptoion: %s,
	Company Name: %s,
	Funds Used: %s,
	Proposer Address: %s,
	Pariticpants: %v,
	`, p.Description, p.CompanyName, p.FundsUsed, p.ProposerAddress.String(), len(p.Participants))

	return strings
}

func (p Proposal) AddParticipant(newP sdk.AccAddress) Proposal {
	p.Participants = append(p.Participants, newP)
	return p
}
