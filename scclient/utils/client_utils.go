package utils

func GetAuthToken(message interface{}) string {
	itemsMap := message.(map[string]interface{})
	data := itemsMap["data"]
	return data.(map[string]interface{})["token"].(string)
}

func GetIsAuthenticated(message interface{}) bool {
	itemsMap := message.(map[string]interface{})
	data := itemsMap["data"]
	return data.(map[string]interface{})["isAuthenticated"].(bool)
}
