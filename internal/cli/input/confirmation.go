package input

type confirmation string

const yes confirmation = "y"
const no confirmation = "n"

type ConfirmationConfig struct {
	Prompt  string
	Default confirmation
}

type confirmationInput struct {
	cfg ConfirmationConfig
}

func NewConfirmationInput(cfg ConfirmationConfig) *confirmationInput {
	if cfg.Default == "" {
		cfg.Default = no
	}
	return &confirmationInput{
		cfg: cfg,
	}
}

func (i *confirmationInput) Read() (bool, error) {
	input := NewOptionInput(OptionInputConfig{
		Prompt:  i.cfg.Prompt,
		Default: string(i.cfg.Default),
		Options: []string{string(yes), string(no)},
	})

	value, err := input.Read()
	if err != nil {
		return false, err
	}

	return value == string(yes), nil
}
