package twitter

type SearchResult struct {
	Text      string `json:"text"`
	CreatedAt int64  `json:"createdAt"`
	Language  string `json:"language"`
	AuthorID  string `json:"authorID"`
	Geo       Geo    `json:"geo"`
	ID        string `json:"id"`
	Source    string `json:"Source"`
}

type Geo struct {
	PlaceID     string    `json:"placeID"`
	Coordinates []float64 `json:"coordinates"`
}
