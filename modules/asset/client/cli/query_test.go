package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/coinexchain/dex/client/cliutil"
	"github.com/coinexchain/dex/modules/asset/internal/types"
)

func TestQueryCmds(t *testing.T) {
	testQueryCmd(t, "params", "custom/asset/parameters", nil)
	testQueryCmd(t, "token abc", "custom/asset/token-info", types.NewQueryAssetParams("abc"))
	testQueryCmd(t, "tokens", "custom/asset/token-list", nil)
	testQueryCmd(t, "whitelist abc", "custom/asset/token-whitelist", types.NewQueryWhitelistParams("abc"))
	testQueryCmd(t, "forbidden-addresses abc", "custom/asset/addr-forbidden", types.NewQueryForbiddenAddrParams("abc"))
	testQueryCmd(t, "reserved-symbols", "custom/asset/reserved-symbols", nil)
}

func testQueryCmd(t *testing.T, subCmds string, expectedPath string, expectedParam interface{}) {
	executed := false
	oldCliQuery := cliutil.CliQuery
	cliutil.CliQuery = func(cdc *codec.Codec, path string, param interface{}) error {
		executed = true
		require.Equal(t, path, expectedPath)
		require.Equal(t, param, expectedParam)
		return nil
	}
	defer func() {
		cliutil.CliQuery = oldCliQuery
	}()

	cmd := GetQueryCmd(types.ModuleCdc)
	cmd.SetArgs(strings.Split(subCmds, " "))
	err := cmd.Execute()
	require.NoError(t, err)
	require.True(t, executed)
}
