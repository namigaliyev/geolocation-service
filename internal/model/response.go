package model

type PageResponse struct {
	Items          interface{} `json:"items"`
	HasNextPage    bool        `json:"hasNextPage"`
	TotalPageCount int         `json:"totalPageCount"`
}
