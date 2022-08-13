package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Sentiment Analysis",
		})
	})

	router.POST("/send", func(c *gin.Context) {
		prompt := c.PostForm("prompt")

		fmt.Println(prompt)
		response := sendPrompt(prompt)

		c.HTML(http.StatusOK, "response.tmpl", gin.H{
			"title":    "URL Shortener",
			"response": response,
		})
	})

	router.Run(":8080")
}

func sendPrompt(phrase string) string {
	godotenv.Load()
	c := gogpt.NewClient(os.Getenv("API_KEY"))
	ctx := context.Background()

	examples := []string{
		"Phrase: I love being an IT support engineer!\nSentiment: Positive.\n\n",
		"Phrase: I am a big fan of the new technology!\nSentiment: Positive.\n\n",
		"Phrase: I can't believe what happened to the old technology!\nSentiment: Negative.\n\n",
		"Phrase: I am not a fan of the new technology!\nSentiment: Negative.\n\n",
		// "Phrase: I hate everything about the new technology!\nSentiment: Negative.\n\n",
		// "Phrase: Why can't my computer just work?\nSentiment: Negative.\n\n",
		// "Phrase: I love the new Macbook!\nSentiment: Positive.\n\n",
	}

	userPrompt := fmt.Sprintf("Phrase: '%s' \nSentiment:", phrase)

	req := gogpt.CompletionRequest{
		Model:       "curie",
		MaxTokens:   5,
		Prompt:      examples[0] + examples[1] + examples[2] + examples[3] + userPrompt,
		Stop:        []string{"."},
		Temperature: 0.4,
	}

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	aiResponse := fmt.Sprintf(resp.Choices[0].Text)
	return aiResponse
}
