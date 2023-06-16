package main

import (
	"fmt"

	grpcclient "QuicksilverDumper/client/grpc"
	"QuicksilverDumper/output"
	csvoutput "QuicksilverDumper/output/csv"
	"QuicksilverDumper/usecase"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "App to query data from a quicksilver-node",
	Long: `A Go application that can query the following data from a quicksilver-node:
- Validator to delegator mapping
- Vesting accounts details categorized by type
- IBC channels statuses between two specified chains
- Client state for the channels
- All "pending" receipts in the x/interchainstaking module`,
}

// Command Names
const (
	GetPendingStakingReceiptsCmdName     = "pending-staking-receipts"
	GetChannelsStatusesCmdName           = "channels-statuses"
	GetAllVestingAccountsCmdName         = "vesting-accounts"
	GetAllValidatorsAndDelegatorsCmdName = "validators-delegators"
)

var node string
var format string
var outputFile string

var getPendingStakingReceiptsCmd = &cobra.Command{
	Use:   "pending-staking-receipts",
	Short: "Query all 'pending' receipts in the x/interchainstaking module",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infoln("GetPendingStakingReceipts called")
		err := executeCommand(cmd, args)
		if err != nil {
			cmd.ErrOrStderr().Write([]byte(err.Error()))
		}
		logger.Infoln("GetPendingStakingReceipts finished")
	},
}

var getChannelsStatusesCmd = &cobra.Command{
	Use:   "channels-statuses",
	Short: "Query all IBC channels between two specified chains for their STATUS",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infoln("GetChannelsStatuses called")
		err := executeCommand(cmd, args)
		if err != nil {
			cmd.ErrOrStderr().Write([]byte(err.Error()))
		}
		logger.Infoln("GetChannelsStatuses finished")
	},
}

var getAllVestingAccountsCmd = &cobra.Command{
	Use:   "vesting-accounts",
	Short: "Query vesting accounts details categorized by type",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infoln("GetAllVestingAccounts called")
		err := executeCommand(cmd, args)
		if err != nil {
			cmd.ErrOrStderr().Write([]byte(err.Error()))
		}
		logger.Infoln("GetAllVestingAccounts finished")
	},
}

var getAllValidatorsAndDelegatorsCmd = &cobra.Command{
	Use:   "validators-delegators",
	Short: "Query validator to delegator mapping and delegator to validator mapping",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infoln("GetAllValidatorsAndDelegators called")
		err := executeCommand(cmd, args)
		if err != nil {
			cmd.ErrOrStderr().Write([]byte(err.Error()))
		}
		logger.Infoln("GetAllValidatorsAndDelegators finished")

	},
}

func executeCommand(cmd *cobra.Command, args []string) error {

	client, err := grpcclient.NewGRPCClient(node)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}

	uc := usecase.NewUseCase(client, logger)

	var res csvoutput.CsvConvertable
	switch cmd.Use {
	case GetPendingStakingReceiptsCmdName:
		result, err := uc.GetPendingStakingReceipts(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get pending staking receipts: %w", err)
		}
		res = csvoutput.GetPendingStakingReceiptsResponse(result)
	case GetChannelsStatusesCmdName:
		result, err := uc.GetChannelsStatuses(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get channels statuses: %w", err)
		}
		res = csvoutput.GetChannelsStatusesResponse(result)

	case GetAllVestingAccountsCmdName:
		result, err := uc.GetAllVestingAccounts(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get all vesting accounts: %w", err)
		}
		res = csvoutput.GetAllVestingAccountsResponse(result)

	case GetAllValidatorsAndDelegatorsCmdName:
		result, err := uc.GetAllValidatorsAndDelegators(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get all validators and delegators: %w", err)
		}
		res = csvoutput.GetAllValidatorsAndDelegatorsResponse(result)

	default:
		return fmt.Errorf("unknown command: %s", cmd.Short)

	}

	logger.Infof("writing %s result into %s in %s format", cmd.Use, outputFile, format)
	outputer, err := output.GetCSVOutputer(res)
	if err != nil {
		return fmt.Errorf("failed to get outputer: %w", err)
	}
	if err := outputer.WriteToFile(outputFile); err != nil {
		return fmt.Errorf("failed to output: %w", err)
	}
	return nil
}
