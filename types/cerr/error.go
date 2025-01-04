package cerr

type error string

var NoDoc = error("no doc find")

func (r error) Error() string {
	return string(r)
}
