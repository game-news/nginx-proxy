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

type UrlData struct {
	Data    DigData
	User    User
	UrlType string
	UrlId   string
}

type StorageBlock struct {
	CounterType  string
	StorageModel string
	UData        UrlData
}

type CmdParams struct {
	LogFilePath string
	RoutineNum  int
	LineNumName string
}
