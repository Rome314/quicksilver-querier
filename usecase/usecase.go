package usecase

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcCore "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type Client interface {
	GetAllAccounts(ctx context.Context) ([]*types.Any, error)
	GetAllIBCChannels(ctx context.Context) ([]*ibcCore.IdentifiedChannel, error)
	GetAllICSReceipts(ctx context.Context) ([]icstypes.Receipt, error)
	GetValidatorDelegations(ctx context.Context, validatorAddr string) (stakingtypes.DelegationResponses, error)
	GetAllValidators(ctx context.Context) ([]stakingtypes.Validator, error)
}

type Logger interface {
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

type UseCase struct {
	Cli    Client
	Logger Logger
}

func NewUseCase(client Client, logger Logger) *UseCase {
	return &UseCase{
		client, logger,
	}
}

// GetPendingStakingReceipts gets all pending staking receipts
func (uc *UseCase) GetPendingStakingReceipts(ctx context.Context) ([]icstypes.Receipt, error) {
	uc.Logger.Infof("Getting all ICS receipts")
	allReceipts, err := uc.Cli.GetAllICSReceipts(ctx)
	if err != nil {
		uc.Logger.Errorf("Failed to get all receipts: %e", err.Error())
		return nil, fmt.Errorf("failed to get all receipts: %w", err)
	}

	var pendingReceipts []icstypes.Receipt
	for _, receipt := range allReceipts {
		if receipt.Completed == nil || receipt.Completed.IsZero() {
			pendingReceipts = append(pendingReceipts, receipt)
		}
	}

	uc.Logger.Infof(fmt.Sprintf("Found %d pending receipts", len(pendingReceipts)))
	return pendingReceipts, nil
}

// GetChannelsStatuses gets the statuses of all IBC channels
func (uc *UseCase) GetChannelsStatuses(ctx context.Context) ([]*ChannelStatus, error) {
	uc.Logger.Infof("Getting all IBC channels")
	channels, err := uc.Cli.GetAllIBCChannels(ctx)
	if err != nil {
		uc.Logger.Errorf("Failed to get all IBC channels: %e", err.Error())
		return nil, err
	}

	uc.Logger.Infof("Parsing channel statuses")
	channelsStatuses := parseChannels(channels)

	uc.Logger.Infof(fmt.Sprintf("Found %d channel statuses", len(channelsStatuses)))
	return channelsStatuses, nil
}

// parseChannels parses the statuses of the given channels
func parseChannels(channels []*ibcCore.IdentifiedChannel) []*ChannelStatus {
	var statuses []*ChannelStatus
	for _, ch := range channels {
		statuses = append(statuses, ChannelStatusFromIdentifiedChannel(ch))
	}
	return statuses
}

// GetAllVestingAccounts gets all vesting accounts
func (uc *UseCase) GetAllVestingAccounts(ctx context.Context) ([]*AnyVestingAccount, error) {
	uc.Logger.Infof("Getting all accounts")
	accounts, err := uc.Cli.GetAllAccounts(ctx)
	if err != nil {
		uc.Logger.Errorf("Failed to get all accounts: %e", err.Error())
		return nil, err
	}

	uc.Logger.Infof("Extracting vesting accounts")
	vestingAccounts, err := extractVestingAccounts(accounts)
	if err != nil {
		uc.Logger.Errorf("Failed to extract vesting accounts: %e", err.Error())
		return nil, err
	}

	uc.Logger.Infof(fmt.Sprintf("Found %d vesting accounts", len(vestingAccounts)))
	return vestingAccounts, nil
}

// extractVestingAccounts extracts vesting accounts from the given accounts
func extractVestingAccounts(accounts []*types.Any) ([]*AnyVestingAccount, error) {
	vestingAccounts := make([]*AnyVestingAccount, 0, len(accounts))
	for _, acc := range accounts {
		vestingAccount, err := AnyVestingAccountFromProtoAny(acc)
		if err != nil {
			if err == NotVestingAccount {
				continue
			}
			return nil, err
		}
		vestingAccounts = append(vestingAccounts, vestingAccount)
	}
	return vestingAccounts, nil
}
