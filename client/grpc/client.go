package grpcclient

import (
	"context"
	"fmt"
	"net"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcCore "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn          *grpc.ClientConn
	AuthClient    authtypes.QueryClient
	IBCClient     ibcCore.QueryClient
	ICSClient     icstypes.QueryClient
	StakingClient stakingtypes.QueryClient
}

func (g *GRPCClient) GetValidatorDelegations(ctx context.Context, validatorAddr string) (stakingtypes.DelegationResponses, error) {
	p := paginator[*stakingtypes.QueryValidatorDelegationsRequest, *stakingtypes.QueryValidatorDelegationsResponse, stakingtypes.DelegationResponse]{
		req: &stakingtypes.QueryValidatorDelegationsRequest{
			ValidatorAddr: validatorAddr,
			Pagination:    &query.PageRequest{Limit: 1000},
		},
		fn: func(ctx context.Context, request *stakingtypes.QueryValidatorDelegationsRequest) (*stakingtypes.QueryValidatorDelegationsResponse, error) {
			return g.StakingClient.ValidatorDelegations(ctx, request)
		},
		getEntities: func(response *stakingtypes.QueryValidatorDelegationsResponse) []stakingtypes.DelegationResponse {
			return response.DelegationResponses
		},
	}

	return p.All(ctx)
}

func (g *GRPCClient) GetAllValidators(ctx context.Context) ([]stakingtypes.Validator, error) {
	p := paginator[*stakingtypes.QueryValidatorsRequest, *stakingtypes.QueryValidatorsResponse, stakingtypes.Validator]{
		req: &stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{Limit: 1000},
		},
		fn: func(ctx context.Context, request *stakingtypes.QueryValidatorsRequest) (*stakingtypes.QueryValidatorsResponse, error) {
			return g.StakingClient.Validators(ctx, request)
		},
		getEntities: func(response *stakingtypes.QueryValidatorsResponse) []stakingtypes.Validator {
			return response.Validators
		},
	}

	return p.All(ctx)
}

func NewGRPCClient(nodeUrl string) (*GRPCClient, error) {
	// Create a new InterfaceRegistry
	interfaceRegistry := codectypes.NewInterfaceRegistry()

	// Register all interfaces needed
	sdk.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)

	// Create a new ProtoCodec (used for message encoding/decoding)
	codec := sdkcodec.NewProtoCodec(interfaceRegistry)

	// Set the node
	conn, err := grpc.Dial(nodeUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, url string) (net.Conn, error) {
			return net.Dial("tcp", url)
		}),
		grpc.WithDefaultCallOptions(
			grpc.ForceCodec(codec.GRPCCodec()),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	authCli := authtypes.NewQueryClient(conn)
	ibcCli := ibcCore.NewQueryClient(conn)
	icsCli := icstypes.NewQueryClient(conn)
	stakingClient := stakingtypes.NewQueryClient(conn)

	resp := &GRPCClient{
		conn:          conn,
		AuthClient:    authCli,
		IBCClient:     ibcCli,
		ICSClient:     icsCli,
		StakingClient: stakingClient,
	}

	return resp, nil
}

func (g *GRPCClient) GetAllICSReceipts(ctx context.Context) ([]icstypes.Receipt, error) {

	zonesResp, err := g.ICSClient.ZoneInfos(ctx, &icstypes.QueryZonesInfoRequest{
		Pagination: &query.PageRequest{CountTotal: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get zones: %w", err)
	}

	chainIds := make([]string, 0, zonesResp.Pagination.Total)
	for _, zone := range zonesResp.Zones {
		chainIds = append(chainIds, zone.ChainId)
	}

	response := make([]icstypes.Receipt, 0)
	for _, chainId := range chainIds {
		req := &icstypes.QueryReceiptsRequest{
			ChainId:    chainId,
			Pagination: &query.PageRequest{CountTotal: true, Limit: 1000},
		}

		p := paginator[*icstypes.QueryReceiptsRequest, *icstypes.QueryReceiptsResponse, icstypes.Receipt]{
			req: req,
			fn: func(ctx context.Context, request *icstypes.QueryReceiptsRequest) (*icstypes.QueryReceiptsResponse, error) {
				return g.ICSClient.Receipts(ctx, request)
			},
			getEntities: func(response *icstypes.QueryReceiptsResponse) []icstypes.Receipt {
				return response.Receipts
			},
		}

		receipts, err := p.All(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get receipts for chain %s: %w", chainId, err)
		}

		response = append(response, receipts...)
	}

	return response, nil
}

func (g *GRPCClient) GetAllIBCChannels(ctx context.Context) ([]*ibcCore.IdentifiedChannel, error) {
	p := paginator[*ibcCore.QueryChannelsRequest, *ibcCore.QueryChannelsResponse, *ibcCore.IdentifiedChannel]{
		req: &ibcCore.QueryChannelsRequest{
			Pagination: &query.PageRequest{CountTotal: true, Limit: 1000},
		},
		fn: func(ctx context.Context, request *ibcCore.QueryChannelsRequest) (*ibcCore.QueryChannelsResponse, error) {
			return g.IBCClient.Channels(ctx, request)
		},
		getEntities: func(response *ibcCore.QueryChannelsResponse) []*ibcCore.IdentifiedChannel {
			return response.Channels
		},
	}

	return p.All(ctx)
}

func (g *GRPCClient) GetAllAccounts(ctx context.Context) ([]*codectypes.Any, error) {
	req := &authtypes.QueryAccountsRequest{
		Pagination: &query.PageRequest{
			CountTotal: true,
			Limit:      500,
		},
	}

	p := paginator[*authtypes.QueryAccountsRequest, *authtypes.QueryAccountsResponse, *codectypes.Any]{
		req: req,
		fn: func(ctx context.Context, request *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error) {
			return g.AuthClient.Accounts(ctx, request)
		},
		getEntities: func(response *authtypes.QueryAccountsResponse) []*codectypes.Any {
			return response.Accounts
		},
	}

	resp, err := p.All(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
