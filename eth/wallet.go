package eth

import (
    "context"
    "crypto/ecdsa"
    "fmt"
    "math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

const INFURA_URL = "https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID"

func CreateWallet() (privateKeyHex, address string, err error) {
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        return "", "", err
    }
    privateKeyBytes := crypto.FromECDSA(privateKey)
    privateKeyHex = fmt.Sprintf("%x", privateKeyBytes)
    address = crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
    return privateKeyHex, address, nil
}

func GetBalance(address string) (string, error) {
    client, err := ethclient.Dial(INFURA_URL)
    if err != nil {
        return "", err
    }
    defer client.Close()

    account := common.HexToAddress(address)
    balance, err := client.BalanceAt(context.Background(), account, nil)
    if err != nil {
        return "", err
    }
    return balance.String(), nil
}

func SendTransaction(privateKeyHex, toAddress, amount string) (string, error) {
    client, err := ethclient.Dial(INFURA_URL)
    if err != nil {
        return "", err
    }
    defer client.Close()

    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        return "", err
    }

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return "", err
    }

    value := new(big.Int)
    value.SetString(amount, 10)
    gasLimit := uint64(21000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        return "", err
    }

    to := common.HexToAddress(toAddress)
    tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, nil)

    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        return "", err
    }

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        return "", err
    }

    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        return "", err
    }

    return signedTx.Hash().Hex(), nil
}
