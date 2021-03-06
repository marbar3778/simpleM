package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/marbar3778/simpleM/x/simplePOA/client/cli"
	"github.com/marbar3778/simpleM/x/simplePOA/types"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the staking module",
	}
	stakingQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQueryValidator(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidators(mc.storeKey, mc.cdc),
		cli.GetCmdQueryParams(mc.storeKey, mc.cdc),
		cli.GetCmdQueryPool(mc.storeKey, mc.cdc))...)

	return stakingQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Staking transaction subcommands",
	}

	stakingTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdCreateValidator(mc.cdc),
		cli.GetCmdEditValidator(mc.cdc),
		// cli.GetCmdUnbond(mc.storeKey, mc.cdc),
	)...)

	return stakingTxCmd
}
