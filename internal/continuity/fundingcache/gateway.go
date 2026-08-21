package fundingcache

type Coordinator struct{ store *StateStore }

func NewCoordinator() *Coordinator {
	policy := PublishPolicy{Mode: "eager", CacheReads: true}
	return &Coordinator{store: NewStateStore(policy)}
}

func (c *Coordinator) Apply(key, value string, commit func() error) error {
	return c.store.Apply(key, value, commit)
}
func (c *Coordinator) Lookup(key string) (string, bool) { return c.store.Lookup(key) }
