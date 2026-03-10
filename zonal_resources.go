package cloudscale

type ZonalResource struct {
	Zone ZoneStub `json:"zone"`
}

type ZonalResourceRequest struct {
	Zone string `json:"zone,omitempty"`
}
