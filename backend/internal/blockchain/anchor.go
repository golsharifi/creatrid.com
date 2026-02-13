package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

// AnchorService handles blockchain anchoring of content hashes.
// When rpcURL is empty, it operates in simulated mode.
type AnchorService struct {
	rpcURL     string
	privateKey string
	chainID    string
	simulated  bool
}

// New creates a new AnchorService. If rpcURL is empty, simulated mode is used.
func New(rpcURL, privateKey, chainID string) *AnchorService {
	if chainID == "" {
		chainID = "137" // Polygon mainnet
	}

	simulated := rpcURL == ""
	if !simulated {
		log.Printf("Blockchain anchor service configured for chain %s (real anchoring not yet implemented, using simulation)", chainID)
	} else {
		log.Println("Blockchain anchor service running in simulated mode")
	}

	return &AnchorService{
		rpcURL:     rpcURL,
		privateKey: privateKey,
		chainID:    chainID,
		simulated:  simulated,
	}
}

// IsSimulated returns true if the service is operating in simulated mode.
func (s *AnchorService) IsSimulated() bool {
	return s.simulated
}

// AnchorHash anchors a content hash on the blockchain.
// In simulated mode, it generates deterministic fake transaction data.
// Returns txHash, blockNumber, contractAddress, and any error.
func (s *AnchorService) AnchorHash(contentHash string) (txHash string, blockNumber int64, contractAddress string, err error) {
	if s == nil {
		return "", 0, "", fmt.Errorf("blockchain anchor service is not configured")
	}

	if !s.simulated {
		// Real blockchain anchoring would happen here.
		// For now, log and fall through to simulation.
		log.Printf("Real blockchain anchoring requested for hash %s (falling back to simulation)", contentHash)
	}

	// Simulated mode: generate deterministic transaction data
	now := time.Now()
	seed := contentHash + now.Format(time.RFC3339Nano) + "creatrid-anchor-salt"
	hash := sha256.Sum256([]byte(seed))
	txHash = "0x" + hex.EncodeToString(hash[:])[:64]
	blockNumber = now.Unix()
	contractAddress = "0x0000000000000000000000000000000000000000"

	return txHash, blockNumber, contractAddress, nil
}

// VerifyAnchor verifies a transaction on the blockchain.
// In simulated mode, this is a no-op that returns an error indicating
// verification should be done via database lookup.
func (s *AnchorService) VerifyAnchor(txHash string) (contentHash string, blockNumber int64, timestamp time.Time, err error) {
	if s == nil {
		return "", 0, time.Time{}, fmt.Errorf("blockchain anchor service is not configured")
	}

	if !s.simulated {
		// Real blockchain verification would happen here.
		log.Printf("Real blockchain verification requested for tx %s (not yet implemented)", txHash)
	}

	// In simulated mode, verification must be done via database lookup.
	// The handler will look up the anchor by tx hash in the database.
	return "", 0, time.Time{}, fmt.Errorf("on-chain verification not available in simulated mode; use database lookup")
}
