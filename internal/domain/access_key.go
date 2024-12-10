package domain

import "time"

type AccessKey struct {
	Key         string                `json:"key"`
	CreatedAt   time.Time             `json:"createdAt"`
	Permissions *AccessKeyPermissions `json:"permissions"`
}

type AccessKeyPermissions struct {
	Send    []string `json:"send"`
	Receive []string `json:"receive"`
	Access  []string `json:"access"`
}

type AccessKeyList struct {
	ResultsLength int          `json:"resultsLength"`
	Results       []*AccessKey `json:"results"`
}
