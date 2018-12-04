package rules

type HTTPOutput struct {
	Method         string
	Endpoint       string
	Headers        map[string]string
	Params         map[string]string
	Authentication struct {
		Header map[string]string
	}
}

type Rules struct {
	Admin struct {
		Namespaces         []string
		DeadLetters        bool `yaml:"dead-letters"`
		ErrorOnDeadLetters bool `yaml:"error-on-dead-letters"`
		Output             struct {
			HTTP HTTPOutput `yaml:"http"`
		}
	}
	Output []struct {
		Namespace string
		HTTP      HTTPOutput `yaml:"http"`
	}
}
