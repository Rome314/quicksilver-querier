package usecase

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	ibcCore "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
)

type ValidatorWithDelegators struct {
	ValidatorAddress string
	DelegatorAddress string
	Delegations      sdk.Coins
	TotalShares      sdk.Dec
}

type ChannelStatus struct {
	SourceChannelId       string
	SourcePortId          string
	CounterpartyChannelId string
	CounterpartyPortId    string
	State                 string
}

func ChannelStatusFromIdentifiedChannel(channel *ibcCore.IdentifiedChannel) *ChannelStatus {
	return &ChannelStatus{
		SourceChannelId:       channel.ChannelId,
		SourcePortId:          channel.PortId,
		CounterpartyChannelId: channel.Counterparty.ChannelId,
		CounterpartyPortId:    channel.Counterparty.PortId,
		State:                 channel.State.String(),
	}
}

type VestingAccountType string

const (
	Delayed         VestingAccountType = "Delayed"
	Continuous      VestingAccountType = "Continuous"
	Periodic        VestingAccountType = "Periodic"
	PermanentLocked VestingAccountType = "PermanentLocked"
)

type AnyVestingAccount struct {
	AccountType     VestingAccountType
	Delayed         *vestingtypes.DelayedVestingAccount
	Continuous      *vestingtypes.ContinuousVestingAccount
	Periodic        *vestingtypes.PeriodicVestingAccount
	PermanentLocked *vestingtypes.PermanentLockedAccount
}

func (v *AnyVestingAccount) String() string {
	switch v.AccountType {
	case Delayed:
		return v.Delayed.String()
	case Continuous:
		return v.Continuous.String()
	case Periodic:
		return v.Periodic.String()
	case PermanentLocked:
		return v.PermanentLocked.String()
	default:
		return "Unknown"
	}
}

var NotVestingAccount = errors.New("not a vesting account")

func AnyVestingAccountFromProtoAny(any *types.Any) (*AnyVestingAccount, error) {
	account := &AnyVestingAccount{}

	var err error
	switch any.TypeUrl {
	case "/cosmos.vesting.v1beta1.DelayedVestingAccount":
		account.AccountType = Delayed
		acc := &vestingtypes.DelayedVestingAccount{}
		err = acc.Unmarshal(any.Value)
		account.Delayed = acc
	case "/cosmos.vesting.v1beta1.ContinuousVestingAccount":
		account.AccountType = Continuous
		acc := &vestingtypes.ContinuousVestingAccount{}
		err = acc.Unmarshal(any.Value)
		account.Continuous = acc
	case "/cosmos.vesting.v1beta1.PeriodicVestingAccount":
		account.AccountType = Periodic
		acc := &vestingtypes.PeriodicVestingAccount{}
		err = acc.Unmarshal(any.Value)
		account.Periodic = acc
	case "/cosmos.vesting.v1beta1.PermanentLockedAccount":
		account.AccountType = PermanentLocked
		acc := &vestingtypes.PermanentLockedAccount{}
		err = acc.Unmarshal(any.Value)
		account.PermanentLocked = acc
	default:
		err = NotVestingAccount
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}
