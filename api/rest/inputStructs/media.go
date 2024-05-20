package inputStructs

type Media struct {
	RecordTTL int64 `mapstructure:"recordTTL"`
	EventTTL  int64 `mapstructure:"eventTTL"`
}

type MediaInput struct {
	Body Media
}

type MediaResp struct {
	Body Media
}
