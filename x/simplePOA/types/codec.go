package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateAuthority{}, "cosmos-sdk/MsgCreateAuthority", nil)
	cdc.RegisterConcrete(MsgEditAuthority{}, "cosmos-sdk/MsgEditAuthority", nil)
	// cdc.RegisterConcrete(MsgDelegate{}, "cosmos-sdk/MsgDelegate", nil)
	// cdc.RegisterConcrete(MsgUndelegate{}, "cosmos-sdk/MsgUndelegate", nil)
	// cdc.RegisterConcrete(MsgBeginRedelegate{}, "cosmos-sdk/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
