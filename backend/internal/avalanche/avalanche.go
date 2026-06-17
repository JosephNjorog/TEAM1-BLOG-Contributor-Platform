// Package avalanche sends USDC payments on the Avalanche C-Chain from the
// platform's treasury wallet, and reports back the Core-wallet-compatible
// view of those payments: a transaction hash plus its confirmation status.
package avalanche

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// usdcDecimals is fixed for the USDC ERC-20 contract on Avalanche C-Chain.
const usdcDecimals = 6

// transferSelector is the first 4 bytes of keccak256("transfer(address,uint256)").
var transferSelector = []byte{0xa9, 0x05, 0x9c, 0xbb}

type Sender interface {
	// Send submits a USDC transfer for amountUSD to toAddress and returns
	// the transaction hash.
	Send(ctx context.Context, toAddress string, amountUSD float64) (txHash string, err error)
	// Confirm blocks until the transaction is mined (or ctx is done) and
	// reports whether it succeeded.
	Confirm(ctx context.Context, txHash string) (confirmed bool, err error)
}

func NewSender(rpcURL, treasuryKeyHex, usdcContract string, chainID int64, mock bool) (Sender, error) {
	if mock {
		return &mockSender{}, nil
	}
	privateKey, err := crypto.HexToECDSA(treasuryKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid AVALANCHE_TREASURY_PRIVATE_KEY: %w", err)
	}
	return &liveSender{
		rpcURL:       rpcURL,
		privateKey:   privateKey,
		fromAddress:  crypto.PubkeyToAddress(privateKey.PublicKey),
		usdcContract: common.HexToAddress(usdcContract),
		chainID:      big.NewInt(chainID),
	}, nil
}

// mockSender simulates an onchain transfer for local development: a
// plausible-looking tx hash with no real chain interaction, confirmed after
// a short delay so the platform's two-step "initiated" -> "confirmed" UX
// still has something to show.
type mockSender struct{}

func (m *mockSender) Send(_ context.Context, _ string, _ float64) (string, error) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return fmt.Sprintf("0x%x", b), nil
}

func (m *mockSender) Confirm(ctx context.Context, _ string) (bool, error) {
	select {
	case <-time.After(3 * time.Second):
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

type liveSender struct {
	rpcURL       string
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	usdcContract common.Address
	chainID      *big.Int
}

func (s *liveSender) Send(ctx context.Context, toAddress string, amountUSD float64) (string, error) {
	client, err := ethclient.DialContext(ctx, s.rpcURL)
	if err != nil {
		return "", fmt.Errorf("dial avalanche rpc: %w", err)
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, s.fromAddress)
	if err != nil {
		return "", fmt.Errorf("get nonce: %w", err)
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("suggest gas price: %w", err)
	}

	data := encodeUSDCTransfer(common.HexToAddress(toAddress), amountUSD)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: s.fromAddress, To: &s.usdcContract, Data: data})
	if err != nil {
		gasLimit = 100_000 // a generous fallback for a single ERC-20 transfer
	}

	tx := types.NewTransaction(nonce, s.usdcContract, big.NewInt(0), gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainID), s.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign transaction: %w", err)
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("send transaction: %w", err)
	}
	return signedTx.Hash().Hex(), nil
}

func (s *liveSender) Confirm(ctx context.Context, txHash string) (bool, error) {
	client, err := ethclient.DialContext(ctx, s.rpcURL)
	if err != nil {
		return false, fmt.Errorf("dial avalanche rpc: %w", err)
	}
	defer client.Close()

	hash := common.HexToHash(txHash)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if err == nil {
			return receipt.Status == types.ReceiptStatusSuccessful, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-ticker.C:
		}
	}
}

// encodeUSDCTransfer ABI-encodes a call to transfer(address,uint256). Scales
// from dollars to USDC's 6 decimal places via cents first to avoid float
// rounding artifacts on the final integer amount.
func encodeUSDCTransfer(to common.Address, amountUSD float64) []byte {
	cents := big.NewInt(int64(amountUSD*100 + 0.5))
	scale := big.NewInt(10)
	scale.Exp(scale, big.NewInt(usdcDecimals-2), nil)
	amount := new(big.Int).Mul(cents, scale)

	data := make([]byte, 0, 4+32+32)
	data = append(data, transferSelector...)
	data = append(data, common.LeftPadBytes(to.Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(amount.Bytes(), 32)...)
	return data
}
