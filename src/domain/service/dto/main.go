package dto

type PageBasic struct {
	PageOffset int    `json:"offset"`
	PageSize   int    `json:"size"`
	Keyword    string `json:"keyword"`
}
