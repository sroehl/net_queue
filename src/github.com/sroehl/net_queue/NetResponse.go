package net_queue

const (
	SUCCESS  = "success"
	ERROR    = "error"
	NO_MSG   = "no_msg"
	FAIL     = "fail"
	HAS_MORE = "has_more"
	UNKNOWN  = "unknown"
)

type NetResponse struct {
	Status string
	Msg    string
	Index  int
}

func new_netresponse(status string, msg string) NetResponse {
	return NetResponse{
		Status: status,
		Msg:    msg,
		Index:  -1,
	}
}
