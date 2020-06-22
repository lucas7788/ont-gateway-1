package provenance

import "time"

// Data provenance
func Data(dataID, author, content, agent string, time time.Time) string {
	return `{
		"dataID":` + dataID + `
	}`
}

func Item(itemID, author, content, agent string, time time.Time) {

}

func Token(tokenID, owner, itemID, content, agent string, time time.Time) {

}
