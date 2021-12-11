package input

import "fmt"

type passwordInput struct {
}

func NewPasswordInput() *passwordInput {
	return &passwordInput{}
}

func (i *passwordInput) Read() (string, error) {
	input := NewSecretInput(Config{
		Prompt:      "password",
		IsMandatory: true,
	})

	inputConfirm := NewSecretInput(Config{
		Prompt:      "confirm password",
		IsMandatory: true,
	})

	var password string
	for {
		pass, err := input.Read()
		if err != nil {
			return "", nil
		}

		confirmPass, err := inputConfirm.Read()
		if err != nil {
			return "", nil
		}

		if pass == confirmPass {
			password = pass
			break
		}

		fmt.Println("passwords don't match")
	}

	return password, nil
}
