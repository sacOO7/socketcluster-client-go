package models

type AuthData struct {
	AuthToken * string `json:"authToken"`
}

type HandShake struct {
	Event string `json:"event"`
	Data AuthData `json:"data"`
	Cid int `json:"cid"`
}

func GetHandshakeObject() HandShake{
	return HandShake {
		Event : "#handshake", 
		Data : AuthData {
		 AuthToken : nil, 
		 } ,
		Cid : 1,
	 }
}	