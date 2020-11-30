package filesystem

type DriverConfig struct {
	Driver    string `json:"driver"`
	Domain    string `json:"domain"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Transport string `json:"transport"`
}

type FileInfo struct {
	Title    string `json:"title"`
	Hash     string `json:"hash"`
	Size     int `json:"size"`
	Type     string `json:"type"`
	RawName  string `json:"raw_name"`
	FileName string `json:"file_name"`
	Path     string `json:"path"`
	SaveUrl  string `json:"save_url"`
	Url      string `json:"url"`
}
