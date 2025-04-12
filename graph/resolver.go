package graph

import (
    "context"
    "github.com/yourusername/crypto-wallet/eth"
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
    return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
    return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateWallet(ctx context.Context) (*Wallet, error) {
    privateKey, address, err := eth.CreateWallet()
    if err != nil {
        return nil, err
    }
    balance, _ := eth.GetBalance(address)
    return &Wallet{Address: address, PrivateKey: privateKey, Balance: balance}, nil
}

func (r *mutationResolver) SendTransaction(ctx context.Context, privateKey string, to string, amount string) (*Transaction, error) {
    hash, err := eth.SendTransaction(privateKey, to, amount)
    if err != nil {
        return nil, err
    }
    return &Transaction{Hash: hash}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetBalance(ctx context.Context, address string) (string, error) {
    return eth.GetBalance(address)
}
