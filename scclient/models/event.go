package models


// For first handshake event we are sending response. 
// {"event":"#handshake","data":{"authToken":"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InNhY2hpbiIsImlhdCI6MTUxMzE3Mzk3OCwiZXhwIjoxNTEzMjYwMzc4fQ.W7O1MQoZfqgOwECfz2Cwbp2S3MOVC1DP-8GUTQ2MmdQ"},"cid":1}

// {"rid":1,"data":{"id":"5bkufJ8IuAdShw5zAAAA","isAuthenticated":false,"pingTimeout":20000,"authError":{"name":"AuthTokenInvalidError","message":"invalid signature"}}}

//There are server events related to auth token

// {"event":"#removeAuthToken","data":null,"cid":1}

// {"event":"#setAuthToken","data":{"token":"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VybmFtZSI6InNhY2hpbiIsImlhdCI6MTUxMzE3MzUyOCwiZXhwIjoxNTEzMjU5OTI4fQ.7WcvRpSRFxMj6rw8FFqxe9gjd-HaMK1pxIvxqDoOpuY"},"cid":2}


//// This is a request reply mech 

// {"event":"chat","data":"Hi","cid":3}

// {"rid":3,"error":"This is error","data":"This is success"}


type EmitEvent struct {
	Event string `json:"event"`
	Data interface {} `json:"data"`
	Cid int `json:"cid"`
}

type ReceiveEvent struct {
	Data interface {} `json:"data"`
	Error interface {} `json:"error"`
	Rid int `json:"rid"`	
}


func GetEmitEventObject(eventname string, data interface {}, messageId int) EmitEvent{
	return EmitEvent {
		Event : eventname, 
		Data : data,
		Cid : messageId,
	 }
}


func GetReceiveEventObject(data interface {},error interface {}, messageId int) ReceiveEvent{
	return ReceiveEvent { 
		Data : data,
		Error : error,
		Rid : messageId,
	 }
}
