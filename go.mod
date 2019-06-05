module github.com/marbar3778/simpleM

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.35.0
	github.com/gorilla/mux v1.7.0
	github.com/marbar3778/taskStream v0.0.0-20190519080545-0a3a06e973f2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.5
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
