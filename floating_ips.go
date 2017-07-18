package cloudscale

type FloatingIP struct {
	HREF    string `json:"href"`
	Network string `json:"network"`
	NextHop string `json:"next_hop"`
	Server  Server `json:"server"`
}

type FloatingIPRequest struct {
	IPVersion      string `json:"ip_version"`
	Server         string `json:"server"`
	ReversePointer string `json:"reverse_ptr"`
}
