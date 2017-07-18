package cloudscale

type FloatingIP struct {
	HREF    string `json:"href"`
	Network string `json:"network"`
	NextHop string `json:"next_hop"`
	Server  Server `json:"server"`
}
