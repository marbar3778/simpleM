# simplePOA --> POA module

The end goal of this module it to have it be similar to how the POA network works.

## TODO

- [x] Remove ability to delegate
- [x] creation of validator should not set money
- [x] Remove reference to tokens/coins/money
- [x] Change validators to authorities
- [x]Add power to authority on creation
  - [] If no power is assigned it is defaulted to 10
  - [] Make the gentx in genutil be able to create a validator that is for POA
- [] Ability to update authority set
- [] Governance vote to add validators if desired
- [] Governance vote to remove validators if desired
- [] enable/disable validator changes, based on governance
  - If a validator is accepted, then the creation step should only allow the validator to join with the pubkey he specified in the governane proposal.
- [] Create a vlaidator with a gentx,
  - The genutil module cannot be used. This is module has to hav its own gentx function

## Spec

The primary focus of this module is to enable a second type of consensus in the Cosmos-SDK. This module is built to have validators, but not a reference to tokens or distr.
All fees are collected by the authority(validator). For now anyone can become a validaotr if they choose.

#### Possible options for slashing

- When a authority is created it is associated to a account. The account address cannot be changed. If there is a slashing offense then the amount is taken from the validators wallet.
- The slashing amount is kept for a specified time period and in that time period that `pool` can be slashed for some offense.
