package domain

type AccessKey struct {
	AccessKey string `json:"accessKey"`
}

type AccessKeyPermissions struct {
	Key         string       `json:"key"`
	Permissions *Permissions `json:"permissions"`
}

type Permissions struct {
	Send    []string `json:"send"`
	Receive []string `json:"receive"`
}

type AccessKeyList struct {
	ResultsLength int                     `json:"resultsLength"`
	Results       []*AccessKeyPermissions `json:"results"`
}
