package handler

import (
	"context"
	"log"

	"github.com/creatrid/creatrid/internal/score"
	"github.com/creatrid/creatrid/internal/store"
)

func recalcScore(ctx context.Context, st *store.Store, userID string) {
	user, err := st.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		log.Printf("Score recalc: failed to load user %s: %v", userID, err)
		return
	}

	connections, err := st.FindConnectionsByUserID(ctx, userID)
	if err != nil {
		log.Printf("Score recalc: failed to load connections for %s: %v", userID, err)
		return
	}

	s := score.Calculate(user, connections)

	if err := st.UpdateUserScore(ctx, userID, s); err != nil {
		log.Printf("Score recalc: failed to update score for %s: %v", userID, err)
	}
}
