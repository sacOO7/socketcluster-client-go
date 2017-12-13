package parser


func Parse(rid int, cid int, event interface{}) MessageType {
	if event != nil {
		if event == "#publish" {
			return PUBLISH

		} else if event == "#removeAuthToken" {
			return REMOVETOKEN

		} else if event == "#setAuthToken" {
			return SETTOKEN

		} else {
			return EVENT
		}
	} else if rid == 1 {
		return ISAUTHENTICATED

	} else {
		return ACKRECEIVE
	}
}
