package chatgpt

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

const (
	video_description = "假设你在抖音上发表了一个短视频, 短视频的剧情是 %s, 你准备好了吗"
	REQUEST_PREFIX    = "请用幽默的方式回复 粉丝的评论: %s"
)

func (c *GptClient) FirstMessage(videoDesc string) (string, error) {
	return c.Send(fmt.Sprintf(video_description, videoDesc))
}

func (c *GptClient) Send(message string) (string, error) {
	defer func() {
		c.index++
		c.index = c.index % len(c.openaiClient)
		fmt.Println(c.index)
	}()

	c.contents = append(c.contents, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})
	resp, err := c.openaiClient[c.index].CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: c.contents,
		},
	)
	if err != nil {
		return "", err
	}
	c.contents = append(c.contents, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp.Choices[0].Message.Content,
	})
	//fmt.Println(fmt.Sprintf("prompt tokens:%d,completion tokens:%d,total tokens %d", resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens))
	return resp.Choices[0].Message.Content, nil
}
