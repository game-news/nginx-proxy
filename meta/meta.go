package meta

type User struct {
	Uid             int64
	Username        string
	IsAnonymous     bool
	IsAuthenticated bool
	IsAdmin         bool
}

type DigData struct {
	RemoteAddr        string
	RemoteUser        string
	TimeLocal         string
	HttpMethod        string
	HttpUrl           string
	HttpVersion       string
	Status            string
	BodyBytesSent     string
	HttpReferer       string
	HttpUserAgent     string
	HttpXForwardedFor string
	HttpToken         string
}
