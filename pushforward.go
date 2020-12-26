package mob_push_sdk

type PushForward struct {
	NextType       int       `json:"nextType"`
	Url            string    `json:"url"`
	Scheme         string    `json:"scheme"`
	SchemeDataList []PushMap `json:"schemeDataList"`
}

func (push *Push) setForward(forward PushForward) *Push {
	push.PushForward = &forward
	return push
}
