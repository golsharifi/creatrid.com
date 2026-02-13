package blockchain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// AnchorService handles blockchain anchoring of content hashes on Base L2.
type AnchorService struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	chainID    *big.Int
}

// New creates a new AnchorService connected to the given RPC endpoint.
// Returns an error if the RPC URL or private key is missing/invalid.
func New(rpcURL, privateKeyHex, chainIDStr string) (*AnchorService, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("BLOCKCHAIN_RPC_URL is required")
	}
	if privateKeyHex == "" {
		return nil, fmt.Errorf("BLOCKCHAIN_PRIVATE_KEY is required")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	chainID := big.NewInt(8453) // Base mainnet default
	if chainIDStr != "" {
		if _, ok := chainID.SetString(chainIDStr, 10); !ok {
			return nil, fmt.Errorf("invalid chain ID: %s", chainIDStr)
		}
	}

	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	log.Printf("Blockchain anchor service connected to chain %s, wallet %s", chainID.String(), fromAddr.Hex())

	return &AnchorService{
		client:     client,
		privateKey: privateKey,
		fromAddr:   fromAddr,
		chainID:    chainID,
	}, nil
}

// WalletAddress returns the service's wallet address.
func (s *AnchorService) WalletAddress() string {
	return s.fromAddr.Hex()
}

// AnchorHash submits a content hash to the blockchain as transaction data.
// Returns the tx hash immediately (transaction may still be pending).
func (s *AnchorService) AnchorHash(ctx context.Context, contentHash string) (txHash string, err error) {
	if s == nil {
		return "", fmt.Errorf("blockchain anchor service is not configured")
	}

	// Decode content hash to bytes for tx data
	hashBytes, err := hex.DecodeString(contentHash)
	if err != nil {
		return "", fmt.Errorf("invalid content hash: %w", err)
	}

	// Get nonce
	nonce, err := s.client.PendingNonceAt(ctx, s.fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}

	// Get gas price
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create transaction: send 0 ETH to self with content hash as data
	tx := types.NewTransaction(
		nonce,
		s.fromAddr,    // send to self
		big.NewInt(0), // no value
		50000,         // gas limit (21000 base + data cost + margin)
		gasPrice,
		hashBytes,
	)

	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainID), s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	txHashHex := signedTx.Hash().Hex()
	log.Printf("Blockchain anchor tx sent: %s (content hash: %s)", txHashHex, contentHash)

	return txHashHex, nil
}

// CheckTransaction checks if a transaction has been mined and returns block info.
// Returns blockNumber, timestamp, and whether it's confirmed.
func (s *AnchorService) CheckTransaction(ctx context.Context, txHashHex string) (blockNumber int64, timestamp time.Time, confirmed bool, err error) {
	if s == nil {
		return 0, time.Time{}, false, fmt.Errorf("blockchain anchor service is not configured")
	}

	receipt, err := s.client.TransactionReceipt(ctx, common.HexToHash(txHashHex))
	if err != nil {
		// Transaction not yet mined
		return 0, time.Time{}, false, nil
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return 0, time.Time{}, false, fmt.Errorf("transaction reverted")
	}

	block, err := s.client.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return receipt.BlockNumber.Int64(), time.Time{}, true, nil
	}

	ts := time.Unix(int64(block.Time()), 0)
	return receipt.BlockNumber.Int64(), ts, true, nil
}

// VerifyAnchor verifies a transaction on-chain and returns the data stored in it.
func (s *AnchorService) VerifyAnchor(ctx context.Context, txHashHex string) (contentHash string, blockNumber int64, timestamp time.Time, err error) {
	if s == nil {
		return "", 0, time.Time{}, fmt.Errorf("blockchain anchor service is not configured")
	}

	tx, isPending, err := s.client.TransactionByHash(ctx, common.HexToHash(txHashHex))
	if err != nil {
		return "", 0, time.Time{}, fmt.Errorf("transaction not found: %w", err)
	}
	if isPending {
		return "", 0, time.Time{}, fmt.Errorf("transaction is still pending")
	}

	contentHash = hex.EncodeToString(tx.Data())

	receipt, err := s.client.TransactionReceipt(ctx, common.HexToHash(txHashHex))
	if err != nil {
		return contentHash, 0, time.Time{}, nil
	}

	blockNumber = receipt.BlockNumber.Int64()

	block, err := s.client.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return contentHash, blockNumber, time.Time{}, nil
	}

	timestamp = time.Unix(int64(block.Time()), 0)
	return contentHash, blockNumber, timestamp, nil
}
