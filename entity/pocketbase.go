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

type User struct {
	Avatar          string   `json:"avatar"`
	CollectionId    string   `json:"collectionId"`
	CollectionName  string   `json:"collectionName"`
	Company         string   `json:"company"`
	Created         string   `json:"created"`
	CreatedBy       string   `json:"createdBy"`
	Email           string   `json:"email"`
	EmailVisibility bool     `json:"emailVisibility"`
	FirstName       string   `json:"firstName"`
	Groups          []string `json:"groups"`
	Id              string   `json:"id"`
	JPort           int      `json:"jPort"`
	JToken          string   `json:"jToken"`
	LastName        string   `json:"lastName"`
	Role            string   `json:"role"`
	Updated         string   `json:"updated"`
	Username        string   `json:"username"`
	Verified        bool     `json:"verified"`
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
