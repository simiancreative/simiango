package cryptkeeper

import "os"

var keyName = "TOKEN_SECRET"
var keyValue = "cGFzc3BocmFzZXdoaWNobmVlZHN0b2JlMzJieXRlcw=="

func SetKeyName(newName string) {
	keyName = newName
}

func SetKey(newKey string) {
	keyValue = newKey
}

func init() {
	keyValue = os.Getenv(keyName)
}
