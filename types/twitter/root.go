package twitter

type SearchResult struct {
	Text      string `json:"text"`
	CreatedAt int64  `json:"createdAt"`
	Language  string `json:"language"`
	AuthorID  string `json:"authorID"`
	Geo       Geo    `json:"geo"`
	PostID    string `json:"postID"`
	Source    string `json:"Source"`

	UserInfo User     `json:"userInfo"`
	Location Location `json:"location"`
}

type Geo struct {
	PlaceID     string    `json:"placeID"`
	Coordinates []float64 `json:"coordinates"`
}

type User struct {
	UserName string `json:"userName"`
	UserID   string `json:"userID"`
}

type Location struct {
	CountryCode string `json:"countryCode"`
	FullName    string `json:"fullName"`
}
