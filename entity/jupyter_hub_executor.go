package entity

type JupyterHubExecutor struct {
	SchedulerId    string `json:"scheduler-id" validate:"nonzero,nonnil"`
	CronExpression string `json:"cron-expression"`
	User           string `json:"user" validate:"nonzero,nonnil"`
	Cron           bool   `query:"cron"`
}
