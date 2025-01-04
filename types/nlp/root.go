package nlp

type NlpDoc struct {
	Summary                  string `json:"summary"`
	Preference               int64  `json:"preference"`
	TotalAggregatedDocuments int64  `json:"totalAggregatedDocuments"`
	CreatedAt                int64  `json:"createdAt"` //unix time
	// TODO -> 무언가 추가로 넣을만한 데이터가 있을까?
}
