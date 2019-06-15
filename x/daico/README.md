# DAICO

The idea of a DAICO is a Decentralized Autonomous Initial Coin Offereing.

1. A company that wants to do an ICO comes to the platform, proposes a new project to be funded. (Roadmap)
2. An investor can come and put money into the DAO of the company.
3. The investor has acces to the tokens of the company who proposed the ICO.
4. The company has to propose uses of the funds.
5. The participants in the DAO vote on the proposal.
6. If the vote passes the company gets the funds, if they vote no on the proposal the funds are not released.
7. If the company is not performing the participants can vote to get their funds back.

### Types

```go
type Participants struct {
  userName string
  userAddress sdk.AccAddress
  amountPaticipated sdk.Coins // coins represent vote power
  ICOparticipated Proposal
}
```

Amount paid out on accepted proposals is taken based on (total amount / amount of personal tokens)

```go
type Proposer struct {
  description string
  companyName string
  fundsUsed string
  ...
}
```
