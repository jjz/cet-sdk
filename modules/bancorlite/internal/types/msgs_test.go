package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var addrOwner = sdk.AccAddress("owner")
var addrNull = sdk.AccAddress("")
var addrUser = sdk.AccAddress("user")

func TestMsgBancorInit_ValidateBasic(t *testing.T) {
	type fields struct {
		Owner            sdk.AccAddress
		Stock            string
		Money            string
		InitPrice        sdk.Dec
		MaxSupply        sdk.Int
		MaxPrice         sdk.Dec
		EnableCancelTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{
			"positive",
			fields{
				addrOwner,
				"abc",
				"cet",
				sdk.NewDec(0),
				sdk.NewInt(100),
				sdk.NewDec(10),
				100},
			nil,
		},
		{
			"negative owner",
			fields{
				addrNull,
				"abc",
				"cet",
				sdk.NewDec(0),
				sdk.NewInt(100),
				sdk.NewDec(10),
				1000,
			},
			sdk.ErrInvalidAddress("missing owner address"),
		},
		{
			"negative token",
			fields{
				addrOwner,
				"cet",
				"abc",
				sdk.NewDec(0),
				sdk.NewInt(100),
				sdk.NewDec(10),
				1000,
			},
			ErrInvalidSymbol(),
		},
		{
			"negative supply",
			fields{
				addrOwner,
				"abc",
				"cet",
				sdk.NewDec(0),
				sdk.NewInt(0),
				sdk.NewDec(10),
				1000,
			},
			ErrNonPositiveSupply(),
		},
		{
			"negative price",
			fields{
				addrOwner,
				"abc",
				"cet",
				sdk.NewDec(0),
				sdk.NewInt(100),
				sdk.NewDec(0),
				1000,
			},
			ErrNonPositivePrice(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBancorInit{
				Owner:            tt.fields.Owner,
				Stock:            tt.fields.Stock,
				Money:            tt.fields.Money,
				InitPrice:        tt.fields.InitPrice,
				MaxSupply:        tt.fields.MaxSupply,
				MaxPrice:         tt.fields.MaxPrice,
				EnableCancelTime: tt.fields.EnableCancelTime,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgBancorInit.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBancorTrade_ValidateBasic(t *testing.T) {
	type fields struct {
		Sender     sdk.AccAddress
		Stock      string
		Money      string
		Amount     int64
		IsBuy      bool
		MoneyLimit int64
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Error
	}{
		{
			name: "positive",
			fields: fields{
				Sender:     addrUser,
				Stock:      "abc",
				Money:      "cet",
				Amount:     10,
				IsBuy:      true,
				MoneyLimit: 10,
			},
			want: nil,
		},
		{
			name: "negative sender",
			fields: fields{
				Sender:     addrNull,
				Stock:      "abc",
				Money:      "cet",
				Amount:     10,
				IsBuy:      true,
				MoneyLimit: 10,
			},
			want: sdk.ErrInvalidAddress("missing sender address"),
		},
		{
			name: "negative token",
			fields: fields{
				Sender:     addrUser,
				Stock:      "cet",
				Money:      "abc",
				Amount:     10,
				IsBuy:      true,
				MoneyLimit: 10,
			},
			want: ErrInvalidSymbol(),
		},
		{
			name: "negative amount",
			fields: fields{
				Sender:     addrUser,
				Stock:      "abc",
				Money:      "cet",
				Amount:     0,
				IsBuy:      true,
				MoneyLimit: 10,
			},
			want: ErrNonPositiveAmount(),
		},
		{
			name: "negative amount exceed max",
			fields: fields{
				Sender:     addrUser,
				Stock:      "abc",
				Money:      "cet",
				Amount:     MaxTradeAmount + 1,
				IsBuy:      true,
				MoneyLimit: 10,
			},
			want: ErrTradeAmountIsTooLarge(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgBancorTrade{
				Sender:     tt.fields.Sender,
				Stock:      tt.fields.Stock,
				Money:      tt.fields.Money,
				Amount:     tt.fields.Amount,
				IsBuy:      tt.fields.IsBuy,
				MoneyLimit: tt.fields.MoneyLimit,
			}
			if got := msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgBancorTrade.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}