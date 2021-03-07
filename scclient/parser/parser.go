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

func GetMessageDetails(message interface{}) (data interface{}, rid int, cid int, eventname interface{}, error interface{}) {
	//Converting given message into map, with keys and values to that we can parse it

	itemsMap, ok := message.(map[string]interface{})
	if !ok {
		return
	}

	for itemKey, itemValue := range itemsMap {
		switch itemKey {
		case "data":
			data = itemValue
		case "rid":
			rid = int(itemValue.(float64))
		case "cid":
			cid = int(itemValue.(float64))
		case "event":
			eventname = itemValue
		case "error":
			error = itemValue
		}
	}

	return
}
