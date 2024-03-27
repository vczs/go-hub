
package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	r := gin.Default()
	r.GET("/stream", func(ctx *gin.Context) {
		stream := getResp()
		defer stream.Close()
		ctx.Stream(func(w io.Writer) bool {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("\nStream finished")
				return false
			}
			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				return false
			}
			ctx.Render(-1, sse.Event{Data: response.Choices[0].Delta.Content})
			return true
		})
	})
	r.Run(":8080")
}

func getResp() *openai.ChatCompletionStream {
	c := openai.NewClient("your token")
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT4,
		MaxTokens: 100,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "简单介绍一下k8s",
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return nil
	}
	return stream
}
