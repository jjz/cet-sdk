package client

import (
	"github.com/coinexchain/dex/modules/market"
	"github.com/coinexchain/dex/modules/market/client/cli"

	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
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
	// Group asset queries under a subcommand
	assQueryCmd := &cobra.Command{
		Use:   market.MarketKey,
		Short: "Querying commands for the market module",
	}

	return assQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	assTxCmd := &cobra.Command{
		Use:   market.MarketKey,
		Short: "market transactions subcommands",
	}

	assTxCmd.AddCommand(client.PostCommands(
		cli.CreateMarketCmd(market.MarketKey, mc.cdc),
		cli.CreateGTEOrderTxCmd(market.MarketKey, mc.cdc),
	)...)

	return assTxCmd
}