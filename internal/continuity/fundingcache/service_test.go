package fundingcache_test

import (
	"errors"
	"testing"

	"carecontinuity/internal/continuity/fundingcache"
)

func TestFundingRollbackCachePublicBehavior(t *testing.T) {
	coordinator := fundingcache.NewCoordinator()
	commitFailure := errors.New("database commit rejected")
	if err := coordinator.Apply("record", "uncommitted", func() error { return commitFailure }); !errors.Is(err, commitFailure) {
		t.Fatalf("expected commit error, got %v", err)
	}
	if value, ok := coordinator.Lookup("record"); ok {
		t.Fatalf("failed transaction leaked value %q", value)
	}
	if err := coordinator.Apply("record", "committed", func() error { return nil }); err != nil {
		t.Fatalf("valid commit failed: %v", err)
	}
	if value, ok := coordinator.Lookup("record"); !ok || value != "committed" {
		t.Fatalf("committed value missing: %q %v", value, ok)
	}
}
