package endpoint

type RestHTTP interface {
	RestHTTP() string
}

type RestHTTPURL string

func (url RestHTTPURL) RestHTTP() string {
	return string(url)
}
