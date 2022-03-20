package filesystem

type CloudConfig struct {
	Driver    string `json:"driver"`
	Domain    string `json:"domain"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Transport string `json:"transport"`
}
