package api

var (
	XAPIHeader        string = "X-API-KEY"
	ContentTypeHeader        = "application/json"
)

type ListParams struct {
	PageIndex int    `validate:"gte=0"`
	Limit     int    `validate:"gt=0,lte=100"`
	Order     string `validate:"oneof=ASC DESC asc desc"`
	OrderBy   string `validate:"oneof=name createdAt"`
	Filter    map[string]string
}
