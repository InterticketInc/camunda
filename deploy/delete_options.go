package deploy

// DeleteOptions options query map for delete deployment request
type DeleteOptions struct {
	// Cascade true, if all process instances, historic process instances and jobs for this deployment should be deleted.
	Cascade bool `url:"cascade,omitempty"`
	// SkipCustomListeners true, if only the built-in ExecutionListeners should be notified with the end event.
	SkipCustomListeners bool `url:"skipCustomListeners,omitempty"`
	//SkipIoMappings true, if all input/output mappings should not be invoked.
	SkipIoMappings bool `url:"skipIoMappings,omitempty"`
}

type DeleteOption interface{}

func DeleteCascade(cascade bool) DeleteOption {
	return &struct{}{}
}
