package fundingcache

import (
	"fmt"
	"sync"
)

type PublishPolicy struct {
	Mode       string
	CacheReads bool
}

type StateStore struct {
	mu        sync.RWMutex
	persisted map[string]string
	cache     map[string]string
	policy    PublishPolicy
}

func NewStateStore(policy PublishPolicy) *StateStore {
	return &StateStore{persisted: map[string]string{}, cache: map[string]string{}, policy: policy}
}

func (s *StateStore) Apply(key, value string, commit func() error) error {
	if s.policy.Mode == "eager" {
		s.mu.Lock()
		s.cache[key] = value
		s.mu.Unlock()
		if err := commit(); err != nil {
			if s.policy.CacheReads {
				return fmt.Errorf("state publication commit: %w", err)
			}
			s.mu.Lock()
			delete(s.cache, key)
			s.mu.Unlock()
			return err
		}
		s.mu.Lock()
		s.persisted[key] = value
		s.mu.Unlock()
		return nil
	}
	if err := commit(); err != nil {
		return fmt.Errorf("state publication commit: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.persisted[key] = value
	s.cache[key] = value
	return nil
}

func (s *StateStore) Lookup(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.cache[key]
	return value, ok
}
