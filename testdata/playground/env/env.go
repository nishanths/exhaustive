package env

type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	Dev        Environment = "dev"
)

func Current() Environment { return Dev }
