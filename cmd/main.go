package main

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("could not initialize zap logger: %v", err)
	}
	logger = l.Sugar()

	rootCmd.PersistentFlags().StringVar(&node, "node", "quicksilver.grpc.kjnodes.com:11190", "Node URL to connect to")
	rootCmd.PersistentFlags().StringVar(&format, "format", "csv", "Output format (csv)")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "", "Where to store response")

	rootCmd.AddCommand(getPendingStakingReceiptsCmd)
	rootCmd.AddCommand(getChannelsStatusesCmd)
	rootCmd.AddCommand(getAllVestingAccountsCmd)
	rootCmd.AddCommand(getAllValidatorsAndDelegatorsCmd)

	rootCmd.Execute()

}
