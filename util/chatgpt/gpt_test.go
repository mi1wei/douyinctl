package chatgpt

import (
	"fmt"
	"testing"
)

func TestGptClient_Send(t *testing.T) {
	gptClient := NewClient()
	for i := 0; i < len(gptClient.openaiClient); i++ {
		resp, err := gptClient.Send("hi")
		fmt.Println(resp, err)
	}

}
