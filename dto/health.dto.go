package dto

type HealthData struct {
	Status       string            `json:"status"`
	Uptime       string            `json:"uptime"`
	Started      string            `json:"started"`
	Env          string            `json:"env"`
	Hostname     string            `json:"hostname"`
	Dependencies map[string]string `json:"dependencies"`
}
