package responses

type Pagination struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
}
