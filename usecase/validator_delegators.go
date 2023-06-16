package usecase

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type sharesAndCoins struct {
	TotalShares sdk.Dec
	Coins       map[string]sdk.Int
}

// GetAllValidatorsAndDelegators gets all validators and their delegators
func (uc *UseCase) GetAllValidatorsAndDelegators(ctx context.Context) (validators []ValidatorWithDelegators, err error) {
	uc.Logger.Infof("Getting all validators")
	allValidators, err := uc.Cli.GetAllValidators(ctx)
	if err != nil {
		uc.Logger.Errorf("Failed to get all validators: %e", err.Error())
		return nil, err
	}

	for _, validator := range allValidators {
		uc.Logger.Infof(fmt.Sprintf("Getting delegations for validator: %s", validator.OperatorAddress))
		delegations, err := uc.Cli.GetValidatorDelegations(ctx, validator.OperatorAddress)
		if err != nil {
			uc.Logger.Errorf("Failed to get validator delegations: %e", err.Error())
			return nil, err
		}

		uc.Logger.Infof("Getting validator delegators")
		validators = append(validators, uc.getValidatorDelegators(delegations)...)
	}

	uc.Logger.Infof(fmt.Sprintf("Found %d validators with delegators", len(validators)))
	return validators, nil
}

// getValidatorDelegators gets the delegators for a validator
func (uc *UseCase) getValidatorDelegators(delegations stakingtypes.DelegationResponses) (validators []ValidatorWithDelegators) {
	uc.Logger.Infof("Grouping delegations")
	validator, delegators := groupDelegations(delegations)
	for delegator, dc := range delegators {

		vw := ValidatorWithDelegators{
			ValidatorAddress: validator,
			DelegatorAddress: delegator,
			TotalShares:      dc.TotalShares,
			Delegations:      sdk.Coins{},
		}

		for denom, amount := range dc.Coins {
			vw.Delegations = append(vw.Delegations, sdk.NewCoin(denom, amount))
		}

		validators = append(validators, vw)
	}
	return
}

// groupDelegations groups delegations by delegator
func groupDelegations(delegations stakingtypes.DelegationResponses) (string, map[string]*sharesAndCoins) {
	validatorAddress := ""
	delegators := map[string]*sharesAndCoins{}

	for _, delegation := range delegations {

		dc, ok := delegators[delegation.Delegation.DelegatorAddress]
		if !ok {
			dc = &sharesAndCoins{
				Coins:       map[string]sdk.Int{},
				TotalShares: sdk.NewDec(0),
			}
			delegators[delegation.Delegation.DelegatorAddress] = dc
		}

		if amount, ok := dc.Coins[delegation.Balance.Denom]; !ok {
			dc.Coins[delegation.Balance.Denom] = delegation.Balance.Amount
		} else {
			dc.Coins[delegation.Balance.Denom] = amount.Add(delegation.Balance.Amount)
		}

		dc.TotalShares = dc.TotalShares.Add(delegation.Delegation.Shares)
		if validatorAddress == "" {
			validatorAddress = delegation.Delegation.ValidatorAddress
		}

	}
	return validatorAddress, delegators
}
