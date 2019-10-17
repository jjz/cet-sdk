package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeSpaceBankx sdk.CodespaceType = ModuleName

	// 301 ～ 399
	CodeMemoMissing                     sdk.CodeType = 301
	CodeInsufficientCETForActivationFee sdk.CodeType = 302
	CodeInvalidActivationFee            sdk.CodeType = 303
	CodeInvalidUnlockTime               sdk.CodeType = 304
	CodeTokenForbiddenByOwner           sdk.CodeType = 305
	CodeInvalidLockCoinsFee             sdk.CodeType = 306
	CodeNoInputs                        sdk.CodeType = 307
	CodeNoOutputs                       sdk.CodeType = 308
	CodeInputOutputMismatch             sdk.CodeType = 309
	CodeInvalidLockCoinsFreeTime        sdk.CodeType = 310
	CodeInvalidOperation                sdk.CodeType = 311
	CodeRewardExceedsAmount             sdk.CodeType = 312
	CodeLockedCoinNotFound              sdk.CodeType = 313
)

func ErrMemoMissing() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeMemoMissing, "memo is empty")
}

func ErrorInsufficientCETForActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInsufficientCETForActivationFee, "Insufficient CET for Activating fees")
}

func ErrUnlockTime(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidUnlockTime, msg)
}

func ErrTokenForbiddenByOwner() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeTokenForbiddenByOwner, "transfer has been forbidden by token owner")
}
func ErrNoInputs() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeNoInputs, "no inputs in multisend")
}
func ErrNoOutputs() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeNoOutputs, "no outputs in multisend")
}
func ErrInputOutputMismatch(msg string) sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInputOutputMismatch, msg)
}
func ErrorInvalidActivatingFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidActivationFee, "invalid activated fees")
}

func ErrorInvalidLockCoinsFee() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidLockCoinsFee, "invalid lock coins fee")
}

func ErrorInvalidLockCoinsFreeTime() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidLockCoinsFreeTime, "invalid lock coins free time")
}

func ErrorInvalidOperation() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeInvalidOperation, "invalid operation")
}

func ErrorRewardExceedsAmount() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeRewardExceedsAmount, "reward exceeds amount")
}

func ErrorLockedCoinNotFound() sdk.Error {
	return sdk.NewError(CodeSpaceBankx, CodeLockedCoinNotFound, "locked coin not found")
}
