package config

type Source struct {
	Books      string `json:"books"`
	Highlights string `json:"highlights"`
	Kindle     string `json:"kindle"`
}

type AWS struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Region    string `json:"region"`
}

type GoogleBooks struct {
	APIKey string `json:"apiKey"`
}
