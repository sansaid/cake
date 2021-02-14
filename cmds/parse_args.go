package cmds

// func matchPrefix(str string) string {

// }

func ParseArgs(args []string) (*Config, error) {
	return &Config{
		registry: "",
		imageTagMap: map[string]string{
			"test": "test",
		},
	}, nil
}
