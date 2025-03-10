package client

type FrameStats struct {
	ResX             int                `json:"resx"`
	ResY             int                `json:"resy"`
	Changed          int                `json:"changed"`
	ServerTime       int                `json:"server_time"`
	Clients          []FrameStatsClient `json:"clients"`
	Analysis         int                `json:"analysis"`
	Screenshot       int                `json:"screenshot"`
	EncodingTotal    int                `json:"encoding_total"`
	VideoScaling     int                `json:"videoscaling"`
	TightJpegEncoder EncoderStats       `json:"tightjpegencoder"`
	TightWebpEncoder EncoderStats       `json:"tightwebpencoder"`
}

type FrameStatsResponse struct {
	Frame            FrameStats         `json:"frame,omitempty"`
	Clients          []FrameStatsClient `json:"clients,omitempty"`
	Analysis         int                `json:"analysis,omitempty"`
	Screenshot       int                `json:"screenshot,omitempty"`
	EncodingTotal    int                `json:"encoding_total,omitempty"`
	VideoScaling     int                `json:"videoscaling,omitempty"`
	TightJpegEncoder EncoderStats       `json:"tightjpegencoder,omitempty"`
	TightWebpEncoder EncoderStats       `json:"tightwebpencoder,omitempty"`
	ErrorMessage     string             `json:"error_message,omitempty"`
}

type FrameStatsClient struct {
	Client     string `json:"client"`
	ClientTime int    `json:"client_time"`
	Ping       int    `json:"ping"`
	Processes  []struct {
		ProcessName string `json:"process_name"`
		Time        int    `json:"time"`
	} `json:"processes"`
}

type EncoderStats struct {
	Time  int `json:"time"`
	Count int `json:"count"`
	Area  int `json:"area"`
}

type BottleneckStatsResponse struct {
	KasmUser map[string][]float64 `json:"kasm_user"`
}
