package types

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(cdc *codec.Codec) {

}

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}