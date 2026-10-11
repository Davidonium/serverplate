package server

import (
	"net/http"

	"github.com/davidonium/serverplate/internal/domain"
	"github.com/davidonium/serverplate/internal/templates"
)

func statsHandler(pairStore domain.PairStore) appHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		stats, err := pairStore.Stats(ctx, domain.RandomPairFilters{})
		if err != nil {
			return err
		}

		vm := templates.StatsViewModel{
			DatabaseSizeBytes: stats.DatabaseSizeBytes,
			PairCount:         stats.PairCount,
			AdjectiveCount:    stats.AdjectiveCount,
			NounCount:         stats.NounCount,
		}
		c := templates.StatsPage(vm)
		return component(w, r, http.StatusOK, c)
	}
}
