package cloudscale

type Server struct {
	HREF            string      `json:"href"`
	UUID            string      `json:"uuid"`
	Name            string      `json:"name"`
	Status          string      `json:"status"`
	Flavor          Flavor      `json:"flavor"`
	Image           Image       `json:"image"`
	Volumes         []Volume    `json:"volumes"`
	Interfaces      []Interface `json:"interfaces"`
	SSHFingerprints []string    `json:"ssh_fingerprints"`
	SSHHostKeys     []string    `json:"ssh_host_keys"`
	AntiAfinityWith []string    `json:"anti-affinity-with"`
}

type Flavor struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	VCPUCount int    `json:"vcpu_count"`
	MemoryGB  int    `json:"memory_gb"`
}

type Image struct {
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	OperatingSystem string `json:"operating_system"`
}

type Volume struct {
	Type       string `json:"ssd"`
	DevicePath string `json:"device_path"`
	SizeGB     int    `json:"SizeGB"`
}

type Interface struct {
	Type     string    `json:"type"`
	Adresses []Address `json:"addresses"`
}

type Address struct {
	Version      int    `json:"version"`
	Address      string `json:"address"`
	PrefixLenght string `json:"prefix_lenght"`
	Gateway      string `json:"gateway"`
	ReversePtr   string `json:"reverse_prt"`
}
