package terminal

import (
	"fmt"
	"strings"

	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/log"
)

func Confirm(question string) bool {
	fmt.Printf("%s (y|n)\n> ", question)
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		log.Fatalf("Cannot parse input: %s", color.Red("%v", err))
	}
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}
