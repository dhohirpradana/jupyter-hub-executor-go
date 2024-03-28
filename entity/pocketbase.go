package entity

type AuthBody struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Admin struct {
		ID      string `json:"id"`
		Created string `json:"created"`
		Updated string `json:"updated"`
		Avatar  int    `json:"avatar"`
		Email   string `json:"email"`
	} `json:"admin"`
	Token string `json:"token"`
}

type SchedulerResponse struct {
	CollectionId   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Created        string `json:"created"`
	Cron           string `json:"cron"`
	Id             string `json:"id"`
	JobId          string `json:"jobId"`
	LastRun        string `json:"lastRun"`
	Log            string `json:"log"`
	Name           string `json:"name"`
	NextRun        string `json:"nextRun"`
	PathNotebook   string `json:"pathNotebook"`
	Schedule       bool   `json:"schedule"`
	Status         string `json:"status"`
	Time           string `json:"time"`
	Updated        string `json:"updated"`
	User           string `json:"user"`
}
