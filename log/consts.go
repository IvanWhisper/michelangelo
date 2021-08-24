package log

type ContextLogKey string

const (
	DEBUG_REQUEST_ID ContextLogKey = "DebugRequestId" // debug

	REQUEST_ID_KEY ContextLogKey = "rid"
	REQUEST_ID     ContextLogKey = "request-id"

	K_SessionId ContextLogKey = "sessionId"
	K_TraceId   ContextLogKey = "traceId"
	K_SpanId    ContextLogKey = "spanId"
	K_ThreadId  ContextLogKey = "threadId"

	K_Cluster    ContextLogKey = "cluster"
	K_AppName    ContextLogKey = "appName"
	K_ServerAddr ContextLogKey = "serverAddr"
	K_ServerPort ContextLogKey = "serverPort"
	K_Version    ContextLogKey = "version"

	K_LogCategory ContextLogKey = "logCategory"

	K_HttpPath     ContextLogKey = "httpPath"
	K_HttpMethod   ContextLogKey = "httpMethod"
	K_Query        ContextLogKey = "query"
	K_ContentType  ContextLogKey = "contentType"
	K_StatusCode   ContextLogKey = "statusCode"
	K_RequestSize  ContextLogKey = "requestSize"
	K_HttpRequest  ContextLogKey = "httpRequest"
	K_HttpResponse ContextLogKey = "httpResponse"
	K_ClientIp     ContextLogKey = "clientIp"
	K_UserAgent    ContextLogKey = "userAgent"
	K_Errors       ContextLogKey = "errors"

	K_Message ContextLogKey = "msg"

	K_BusinessKeyword   ContextLogKey = "businessKeyword"
	K_BusinessTitle     ContextLogKey = "businessTitle"
	K_BusinessOperation ContextLogKey = "businessOperation"

	K_Datetime ContextLogKey = "datetime"
	K_Caller   ContextLogKey = "caller"

	K_Duration ContextLogKey = "duration"
	K_Level    ContextLogKey = "level"
)
