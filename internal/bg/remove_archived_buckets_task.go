package bg

import (
	"context"
	"log/slog"
	"time"

	"github.com/davidonium/serverplate/internal/domain"
)

func removeArchivedBucketsTask(
	logger *slog.Logger,
	bucketStore domain.BucketStore,
) func(context.Context) error {
	return func(ctx context.Context) error {
		removedCount, err := bucketStore.RemoveBucketsArchivedForMoreThan(ctx, 3*24*time.Hour)
		if err != nil {
			return err
		}

		if removedCount > 0 {
			logger.Info("removed buckets", slog.Int64("amount", removedCount))
		} else {
			logger.Info("no buckets were removed")
		}

		return nil
	}
}
