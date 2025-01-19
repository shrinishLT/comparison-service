package models

type ComparisonTask struct {
	RequestID      string   `json:"requestId"`
	ComparisonType string   `json:"comparisonType"`
	ImageURLs      []string `json:"imageUrls"`
}

type ComparisonResult struct {
	RequestID string `json:"requestId"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}
