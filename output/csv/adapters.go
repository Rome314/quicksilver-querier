package csvoutput

import (
	"fmt"
	"strconv"
	"strings"

	"QuicksilverDumper/usecase"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var _ CsvConvertable = GetPendingStakingReceiptsResponse{}
var _ CsvConvertable = GetChannelsStatusesResponse{}
var _ CsvConvertable = GetAllVestingAccountsResponse{}
var _ CsvConvertable = GetAllValidatorsAndDelegatorsResponse{}

type CsvConvertable interface {
	GetHeaders() []string
	GetValues() [][]string
}

type GetPendingStakingReceiptsResponse []icstypes.Receipt

func (g GetPendingStakingReceiptsResponse) GetHeaders() []string {
	return []string{"ChainId", "Sender", "Txhash", "Coins", "FirstSeen", "Completed"}
}

func (g GetPendingStakingReceiptsResponse) GetValues() [][]string {
	values := make([][]string, 0, len(g))
	for _, receipt := range g {

		amountsStr := make([]string, 0, len(receipt.Amount))
		for _, coin := range receipt.Amount {
			amountsStr = append(amountsStr, fmt.Sprintf("%s,%d", coin.Denom, coin.Amount))
		}

		firstSeen, completed := "null", "null"
		if receipt.FirstSeen != nil {
			firstSeen = receipt.FirstSeen.String()
		}
		if receipt.Completed != nil {
			completed = receipt.Completed.String()
		}

		value := []string{
			receipt.ChainId,
			receipt.Sender,
			receipt.Txhash,
			fmt.Sprintf("[%s]", strings.Join(amountsStr, ",")),
			firstSeen,
			completed,
		}
		values = append(values, value)
	}
	return values
}

type GetChannelsStatusesResponse []*usecase.ChannelStatus

func (g GetChannelsStatusesResponse) GetHeaders() []string {
	return []string{"SourceChannelId", "SourcePortId", "CounterpartyChannelId", "CounterpartyPortId", "State"}
}

func (g GetChannelsStatusesResponse) GetValues() [][]string {
	values := make([][]string, 0, len(g))
	for _, channelStatus := range g {
		value := []string{
			channelStatus.SourceChannelId,
			channelStatus.SourcePortId,
			channelStatus.CounterpartyChannelId,
			channelStatus.CounterpartyPortId,
			channelStatus.State,
		}

		values = append(values, value)

	}
	return values
}

type GetAllVestingAccountsResponse []*usecase.AnyVestingAccount

func (g GetAllVestingAccountsResponse) GetHeaders() []string {
	return []string{"Account Type", "Account Address", "Original Vesting", "Delegated Free", "Delegated Vesting", "End Time", "Start Time", "Periods"}
}

func (g GetAllVestingAccountsResponse) GetValues() [][]string {
	values := make([][]string, 0, len(g))
	for _, account := range g {
		var value []string
		switch account.AccountType {
		case usecase.Delayed:
			value = []string{
				string(account.AccountType),
				account.Delayed.Address,
				account.Delayed.OriginalVesting.String(),
				account.Delayed.DelegatedFree.String(),
				account.Delayed.DelegatedVesting.String(),
				strconv.FormatInt(account.Delayed.EndTime, 10),
				"N/A",
				"N/A",
			}
		case usecase.Continuous:
			value = []string{
				string(account.AccountType),
				account.Continuous.Address,
				account.Continuous.OriginalVesting.String(),
				account.Continuous.DelegatedFree.String(),
				account.Continuous.DelegatedVesting.String(),
				strconv.FormatInt(account.Continuous.EndTime, 10),
				strconv.FormatInt(account.Continuous.StartTime, 10),
				"N/A",
			}
		case usecase.Periodic:
			periods := make([]string, len(account.Periodic.VestingPeriods))
			for i, period := range account.Periodic.VestingPeriods {
				periods[i] = fmt.Sprintf("{length: %d, amount: %s}", period.Length, period.Amount.String())
			}
			value = []string{
				string(account.AccountType),
				account.Periodic.Address,
				account.Periodic.OriginalVesting.String(),
				account.Periodic.DelegatedFree.String(),
				account.Periodic.DelegatedVesting.String(),
				strconv.FormatInt(account.Periodic.EndTime, 10),
				strconv.FormatInt(account.Periodic.StartTime, 10),
				strings.Join(periods, "; "),
			}
		case usecase.PermanentLocked:
			value = []string{
				string(account.AccountType),
				account.PermanentLocked.Address,
				account.PermanentLocked.OriginalVesting.String(),
				account.PermanentLocked.DelegatedFree.String(),
				account.PermanentLocked.DelegatedVesting.String(),
				"N/A",
				"N/A",
				"N/A",
			}
		}
		values = append(values, value)
	}
	return values
}

type GetAllValidatorsAndDelegatorsResponse []usecase.ValidatorWithDelegators

func (g GetAllValidatorsAndDelegatorsResponse) GetHeaders() []string {
	return []string{"ValidatorAddress", "DelegatorAddress", "Delegations", "TotalShares"}
}

func (g GetAllValidatorsAndDelegatorsResponse) GetValues() [][]string {
	values := make([][]string, 0, len(g))
	for _, vwd := range g {
		// Write the data
		delegatoinsStr := make([]string, 0, len(vwd.Delegations))
		for _, coin := range vwd.Delegations {
			delegatoinsStr = append(delegatoinsStr, fmt.Sprintf("{%s,%d}", coin.Denom, coin.Amount.Int64()))
		}
		delegations := fmt.Sprintf("[%s]", strings.Join(delegatoinsStr, ","))
		values = append(values, []string{vwd.ValidatorAddress, vwd.DelegatorAddress, delegations, vwd.TotalShares.String()})

	}
	return values
}
