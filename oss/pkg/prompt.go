package pkg

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

type Prompt struct {
}

func NewPrompt() *Prompt {
	return &Prompt{}
}

func (p *Prompt) TriggerPrompt(label string) string {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Success: "{{ . | blue }} ",
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("%s :", label),
		Templates: templates,
	}

	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}
