package indicator

const (
	ExecStatObject = "exec_stat"
)

type ExecStat struct {
	EnvName     string `json:"env_name"`
	EnvCode     string `json:"env_code"`
	ProjectName string `json:"project_name"`
	ProjectCode string `json:"project_code"`
	Num         int    `json:"num"`
	FailNum     int    `json:"fail_num"`
	Count       int    `json:"count"`
	FailCount   int    `json:"fail_count"`
	Begin       string `json:"begin"`
	End         string `json:"end"`
}
