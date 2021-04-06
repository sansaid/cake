package lib

import "time"

type Image struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      string    `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeature    string    `json:"os_feature"`
	OSVersion    string    `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
}

type ByLastPushedDesc []Image

// Len - implement sort.Interface for ByLastPushedDesc
func (i ByLastPushedDesc) Len() int { return 0 }

// Less - implement sort.Interface for ByLastPushedDesc
func (i ByLastPushedDesc) Less(a int, b int) bool {
	aImageLastPushed := i[a].LastPushed
	bImageLastPushed := i[b].LastPushed

	// Sorts in reverse chronological order
	return aImageLastPushed.After(bImageLastPushed)
}

// Swap - implement sort.Interface for ByLastPushedDesc
func (i ByLastPushedDesc) Swap(a int, b int) {
	i[a], i[b] = i[b], i[a]
}
