package server

import (
	"net/http"
	"strconv"

	"github.com/davidonium/serverplate/internal/domain"
	"github.com/davidonium/serverplate/internal/templates"
)

func homeHandler(pairStore domain.PairStore) appHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		stats, err := pairStore.Stats(ctx, domain.RandomPairFilters{})
		if err != nil {
			return err
		}

		c := templates.HomePage(templates.HomeViewModel{
			PossiblePairCount: stats.PairCount,
		})
		return component(w, r, http.StatusOK, c)
	}
}

func configStatsHandler(pairStore domain.PairStore) appHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		lengthEnabled := r.FormValue("length_enabled") == "on"
		lengthMode := r.FormValue("length_mode")
		lengthValue, _ := strconv.Atoi(r.FormValue("length_value"))

		filters := domain.RandomPairFilters{}
		if lengthEnabled {
			filters.Length = lengthValue
			filters.LengthMode = domain.LengthMode(lengthMode)
		}

		stats, err := pairStore.Stats(ctx, filters)
		if err != nil {
			return err
		}

		c := templates.ConfigurationStatsPartial(templates.ConfigurationStatsPartialViewModel{
			PossiblePairCount: stats.PairCount,
		})
		return component(w, r, http.StatusOK, c)
	}
}
