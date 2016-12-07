package pet

const (
	ERR_NOERR    = 0    //没有错误
	ERR_UNKNOWN  = 1001 //未知错误
	ERR_INTERNAL = 1002 //内部错误
	ERR_MYSQL    = 1003 //mysql错误
	ERR_REDIS    = 1004 //redis错误

	ERR_PATH           = 2000 //redis错误
	ERR_INVALID_PARAM  = 2001 //请求参数错误
	ERR_INVALID_FORMAT = 2002 //格式错误
	ERR_REQUIRE_PARAM = 2003 //格式错误
	ERR_JSON_STYLE     = 2011 //json格式错误
	ERR_DATA           = 2012 //错误数据
)
