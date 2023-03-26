package server

type RequestShorten struct {
	URL string `json:"url"`
}

type ResponseShorten struct {
	Result string `json:"result"`
}

type RequestShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalID    string `json:"original_url"`
}

type ResponseShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
