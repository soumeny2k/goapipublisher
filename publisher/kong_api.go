package publisher

type KongApi struct {
	Name              string        `json:"name"`
	RateLimit         uint          `json:"rateLimit"`
	ConnectionTimeout uint          `json:"connectionTimeout"`
	Protocol          string        `json:"protocol" default:"http"`
	Type              string        `json:"type" default:"REST"`
	Backends          []KongBackend `json:"backends"`
}

type KongBackend struct {
	Url    string `json:"url"`
	Weight uint   `json:"weight"`
}
