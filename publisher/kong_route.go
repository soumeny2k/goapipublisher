package publisher

type KongRoute struct {
	Name              string       `json:"name"`
	Path              string       `json:"path"`
	PathPattern       string       `json:"pathPattern"`
	Method            string       `json:"method"`
	Retries           uint         `json:"retries"`
	RateLimit         uint         `json:"rateLimit"`
	ConnectionTimeout uint         `json:"connectionTimeout"`
	CacheEnabled      bool         `json:"cacheEnabled"`
	Headers           []KongHeader `json:"headers"`
}

type KongHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
