package twitter

type huggingUrl string

const (
	// TODO BaseURL 확인 후에, suffix 만 활용하자
	//"https://api-inference.huggingface.co/models/facebook/bart-large-cnn"
	t = huggingUrl("https://api.twitter.com/1.1/tweets")
)

func (h huggingUrl) String() string {
	return string(h)
}
