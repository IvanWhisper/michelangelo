package log

type ContextLogKey string

func (key ContextLogKey) ToString() string {
	return string(key)
}

const (
	DebugRequestId ContextLogKey = "DebugRequestId" // debug
	RequestIdKey   ContextLogKey = "rid"
	SessionId      ContextLogKey = "sessionId"
	TraceId        ContextLogKey = "traceId"
	SpanId         ContextLogKey = "spanId"
	ThreadId       ContextLogKey = "threadId"
	Cluster        ContextLogKey = "cluster"
	AppName        ContextLogKey = "appName"
	ServerAddr     ContextLogKey = "serverAddr"
	ServerPort     ContextLogKey = "serverPort"
	Version        ContextLogKey = "version"

	Category ContextLogKey = "logCategory"

	HttpPath     ContextLogKey = "httpPath"
	HttpMethod   ContextLogKey = "httpMethod"
	QueryText    ContextLogKey = "query"
	ContentType  ContextLogKey = "contentType"
	StatusCode   ContextLogKey = "statusCode"
	RequestSize  ContextLogKey = "requestSize"
	HttpRequest  ContextLogKey = "httpRequest"
	HttpResponse ContextLogKey = "httpResponse"
	ClientIp     ContextLogKey = "clientIp"
	UserAgent    ContextLogKey = "userAgent"
	Errors       ContextLogKey = "errors"

	Message ContextLogKey = "msg"

	BusinessKeyword   ContextLogKey = "businessKeyword"
	BusinessTitle     ContextLogKey = "businessTitle"
	BusinessOperation ContextLogKey = "businessOperation"

	Datetime ContextLogKey = "datetime"
	Caller   ContextLogKey = "caller"

	Duration ContextLogKey = "duration"
	LevelKey ContextLogKey = "level"
)
