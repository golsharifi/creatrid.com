package blockchain

import (
	"context"
	"log"
	"time"

	"github.com/creatrid/creatrid/internal/store"
)

// ConfirmationWorker polls for pending anchors and confirms them on-chain.
type ConfirmationWorker struct {
	store   *store.Store
	service *AnchorService
}

// NewConfirmationWorker creates a new worker.
func NewConfirmationWorker(st *store.Store, svc *AnchorService) *ConfirmationWorker {
	return &ConfirmationWorker{store: st, service: svc}
}

// Start begins the polling loop.
func (w *ConfirmationWorker) Start(ctx context.Context) {
	log.Println("Blockchain confirmation worker started")
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processPending(ctx)
		case <-ctx.Done():
			log.Println("Blockchain confirmation worker stopped")
			return
		}
	}
}

func (w *ConfirmationWorker) processPending(ctx context.Context) {
	anchors, err := w.store.ListPendingAnchors(ctx)
	if err != nil {
		log.Printf("Confirmation worker: failed to list pending anchors: %v", err)
		return
	}

	for _, anchor := range anchors {
		if anchor.TxHash == nil {
			continue
		}

		blockNumber, ts, confirmed, err := w.service.CheckTransaction(ctx, *anchor.TxHash)
		if err != nil {
			log.Printf("Confirmation worker: tx %s failed: %v", *anchor.TxHash, err)
			if err := w.store.UpdateAnchorStatus(ctx, anchor.ID, "failed", err.Error(), 0, nil); err != nil {
				log.Printf("Confirmation worker: failed to update anchor %s: %v", anchor.ID, err)
			}
			continue
		}

		if confirmed {
			if err := w.store.UpdateAnchorStatus(ctx, anchor.ID, "confirmed", "", blockNumber, &ts); err != nil {
				log.Printf("Confirmation worker: failed to confirm anchor %s: %v", anchor.ID, err)
			} else {
				log.Printf("Confirmation worker: confirmed anchor %s at block %d", anchor.ID, blockNumber)
			}
		}
	}
}
