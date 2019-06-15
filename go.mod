module github.com/marbar3778/simpleM

go 1.12

require (
	github.com/btcsuite/btcd v0.0.0-20190213025234-306aecffea32 // indirect
	github.com/cosmos/cosmos-sdk v0.35.0
	github.com/ethereum/go-ethereum v1.8.23 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/gorilla/mux v1.7.0
	github.com/rs/xid v1.2.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
	golang.org/x/crypto v0.0.0-20190313024323-a1f597ede03a // indirect
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.28.2-0.20190615082617-941effc14d01
