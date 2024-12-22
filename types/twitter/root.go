package twitter

type SearchResult struct {
	UserName  string `json:"userName"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"createdAt"`
}
