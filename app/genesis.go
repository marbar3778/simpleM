package app

import (
	"encoding/json"
)

type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for gaia.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}
