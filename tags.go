package cloudscale

type TagMap map[string]string

type TaggedResource struct {
	Tags TagMap `json:"tags"`
}

type TaggedResourceRequest struct {
	Tags TagMap `json:"tags,omitempty"`
}
