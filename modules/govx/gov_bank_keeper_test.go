package govx

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/coinexchain/dex/testutil"
	dex "github.com/coinexchain/dex/types"
)

func setUpKeeper() (sdk.Context, GovBankKeeper) {
	db := dbm.NewMemDB()
	cdc := codec.New()
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	skey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	distKey := sdk.NewKVStoreKey(distribution.StoreKey)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	fckKey := sdk.NewKVStoreKey(auth.FeeStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(skey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(distKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckKey, sdk.StoreTypeIAVL, db)

	_ = ms.LoadLatestVersion()
	paramsKeeper := params.NewKeeper(cdc, skey, tkey)

	ak := auth.NewAccountKeeper(cdc, authKey, paramsKeeper.Subspace(auth.StoreKey), auth.ProtoBaseAccount)
	fck := auth.NewFeeCollectionKeeper(cdc, fckKey)
	bk := bank.NewBaseKeeper(ak, paramsKeeper.Subspace(bank.DefaultParamspace), sdk.CodespaceRoot)
	sk := staking.NewKeeper(cdc, keyStaking, tkey, bk, paramsKeeper.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	dk := distribution.NewKeeper(cdc, distKey, paramsKeeper.Subspace(distribution.StoreKey), bk, sk, fck, types.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	govBankKeeper := NewKeeper(bk, ak, dk)

	//initialize FeePool
	feePool := types.FeePool{
		CommunityPool: sdk.NewDecCoins(dex.NewCetCoins(0)),
	}
	govBankKeeper.dk.SetFeePool(ctx, feePool)

	return ctx, govBankKeeper
}

func TestBurntDepositsToCommunityPool(t *testing.T) {
	ctx, keeper := setUpKeeper()

	depositAddr := gov.DepositedCoinsAccAddr
	acc := keeper.ak.GetAccount(ctx, depositAddr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, depositAddr)
	}

	totalAmt := sdk.Coins{
		sdk.Coin{
			Denom:  "cet",
			Amount: sdk.NewInt(900),
		},
	}
	err := acc.SetCoins(totalAmt)
	if err != nil {
		panic(err)
	}
	keeper.ak.SetAccount(ctx, acc)

	burntAmt := sdk.Coins{
		sdk.Coin{
			Denom:  "cet",
			Amount: sdk.NewInt(300),
		},
	}

	_, _ = keeper.SendCoins(ctx, gov.DepositedCoinsAccAddr, gov.BurnedDepositCoinsAccAddr, burntAmt)

	depositAcc := keeper.ak.GetAccount(ctx, depositAddr)
	require.Equal(t, sdk.NewInt(600), depositAcc.GetCoins().AmountOf("cet"))

	communityPool := keeper.dk.GetFeePool(ctx).CommunityPool
	require.Equal(t, sdk.NewDecCoins(dex.NewCetCoins(300)), communityPool)
}
func TestSendCoinsNormal(t *testing.T) {
	ctx, keeper := setUpKeeper()

	depositAddr := gov.DepositedCoinsAccAddr
	acc := keeper.ak.GetAccount(ctx, depositAddr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, depositAddr)
	}

	totalAmt := sdk.Coins{
		sdk.Coin{
			Denom:  "cet",
			Amount: sdk.NewInt(900),
		},
	}
	err := acc.SetCoins(totalAmt)
	if err != nil {
		panic(err)
	}
	keeper.ak.SetAccount(ctx, acc)

	burntAmt := sdk.Coins{
		sdk.Coin{
			Denom:  "cet",
			Amount: sdk.NewInt(600),
		},
	}

	_, _, toAddr := testutil.KeyPubAddr()
	_, err = keeper.SendCoins(ctx, depositAddr, toAddr, burntAmt)

	require.Nil(t, err)

	depositAcc := keeper.ak.GetAccount(ctx, depositAddr)
	require.Equal(t, sdk.NewInt(300), depositAcc.GetCoins().AmountOf("cet"))

	toAcc := keeper.ak.GetAccount(ctx, toAddr)
	require.Equal(t, sdk.NewInt(600), toAcc.GetCoins().AmountOf("cet"))

	communityPool := keeper.dk.GetFeePool(ctx).CommunityPool
	require.Nil(t, communityPool)

}

func TestSetCoins(t *testing.T) {
	ctx, keeper := setUpKeeper()
	coins := dex.NewCetCoins(10e8)
	_, _, addr := testutil.KeyPubAddr()

	//nil account
	err := setCoins(ctx, keeper.ak, addr, coins)
	acc := keeper.ak.GetAccount(ctx, addr)
	require.Nil(t, err)
	require.Equal(t, coins, acc.GetCoins())

	//invalid amount
	coins[0].Amount = sdk.ZeroInt()
	err = setCoins(ctx, keeper.ak, addr, coins)
	acc = keeper.ak.GetAccount(ctx, addr)
	require.NotNil(t, err)

}
func TestSubtractCoins(t *testing.T) {
	ctx, keeper := setUpKeeper()
	coins := dex.NewCetCoins(10e8)
	_, _, addr := testutil.KeyPubAddr()

	//nil account
	_, err := subtractCoins(ctx, keeper.ak, addr, coins)
	require.NotNil(t, err)

	acc := auth.BaseAccount{
		Address: addr,
		Coins:   dex.NewCetCoins(20e8),
	}
	keeper.ak.SetAccount(ctx, &acc)

	//valid sub
	left, err := subtractCoins(ctx, keeper.ak, addr, coins)
	require.Nil(t, err)
	require.Equal(t, left, coins)

}
