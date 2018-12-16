package config

type MessagingConfig struct {
	RunCleanerEveryHours       uint64
	RunRemoverEveryMinutes     uint64
	RunNotificatorEverySeconds uint64
}

func (this MessagingConfig) Validate() error {
	//TODO may be validate something?
	return nil
}
