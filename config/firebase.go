package config

type FirebaseConfig struct {
	ServerKey string
}

func (this FirebaseConfig) Validate() error {
	//TODO may be validate something?
	return nil
}
