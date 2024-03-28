package entity

type ENV struct {
	PocketbaseIdentity        string `json:"identity"`
	PocketbasePassword        string `json:"password"`
	PocketbaseLoginUrl        string
	PocketbaseSchedulerUrl    string
	PocketbaseNotificationUrl string
	PocketbaseUserUrl         string

	JupyterUrl   string
	JupyterWs    string
	JupyterToken string

	ElasticUrl string
	EventUrl   string
}
