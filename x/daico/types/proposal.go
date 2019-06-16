package types

import (
	"fmt"

	"github.com/rs/xid"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Proposal struct {
	ID              string           `json:"id"`
	Description     string           `json:"description"`
	CompanyName     string           `json:"company_name"`
	ProposerAddress sdk.AccAddress   `json:"proposerAddress"`
	Participants    []sdk.AccAddress `json:"participants"`
	PoolReference   string           `json:"pool_reference"`
	TokenDenom      string           `json:"token_denom"`
	TokenAmount     int              `json:"token_amount"`
	Coins           sdk.Coins        `json:"coins"`
}

func NewProposal(description, companyName string, proposerAddress sdk.AccAddress, tokenDenom string, tokenAmount int) Proposal {
	guid := xid.New().String()

	coins := sdk.NewCoins(sdk.NewCoin(tokenDenom, sdk.NewInt(int64(tokenAmount))))
	return Proposal{
		ID:              guid,
		Description:     description,
		CompanyName:     companyName,
		ProposerAddress: proposerAddress,
		Participants:    []sdk.AccAddress{},
		PoolReference:   "",
		TokenDenom:      tokenDenom,
		TokenAmount:     tokenAmount,
		Coins:           coins,
	}
}

func (p Proposal) GetID() string {
	return p.ID
}
func (p Proposal) String() string {
	strings := fmt.Sprintf(`
	Descriptoion: %s,
	Company Name: %s,
	Proposer Address: %s,
	Pariticpants: %v,
	`, p.Description, p.CompanyName, p.ProposerAddress.String(), len(p.Participants))

	return strings
}

func (p Proposal) AddParticipant(newP sdk.AccAddress) Proposal {
	p.Participants = append(p.Participants, newP)
	return p
}
