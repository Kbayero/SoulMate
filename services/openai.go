package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/KbaYero/SoulMate/database"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIResponse struct {
	Results    []bool `json:"results"`
	Percentage int    `json:"percentage"`
}

func GetResponses(questions []string, player1Answers []string, player2Answers []string, name1, name2 string) ([]database.Result, int) {
	var results []database.Result

	for i, question := range questions {
		player1Answer := player1Answers[i]
		player2Answer := player2Answers[i]

		results = append(results, database.Result{
			Question:      question,
			Player1Answer: player1Answer,
			Player2Answer: player2Answer,
		})
	}

	prompt := buildPrompt(questions, player1Answers, player2Answers)

	openAiResponse, err := sendPromptToOpenAI(prompt)
	if err != nil {
		fmt.Println("Error sending prompt to OpenAI:", err)
		return results, 0
	}

	for i, result := range openAiResponse.Results {
		results[i].IsCorrect = result
		results[i].Player1Name = name1
		results[i].Player2Name = name2
		err := database.GetDB().Create(&results[i])
		if err != nil {
			fmt.Println("Error saving result to database:", err)
		}
	}

	return results, openAiResponse.Percentage
}

func buildPrompt(questions []string, player1Answers []string, player2Answers []string) string {
	prompt := `You are provided with a list of questions and answers from two different players. For each question, compare the answers from Player 1 and Player 2 and determine whether they convey the same meaning. The answers may have different wording, spelling errors, or paraphrasing. Consider them correct if they express the same idea.

Please return a JSON object in the following format:
{
  "results": [bool], // List of booleans indicating if the answers are the same for each question
  "percentage": int  // Percentage of questions where the answers are considered the same
}

Only provide the JSON object without any additional text.

Here is the data:

`

	for i, question := range questions {
		prompt += fmt.Sprintf("Question %d: %s\n", i+1, question)
		prompt += fmt.Sprintf("Player 1 Answer: %s\n", player1Answers[i])
		prompt += fmt.Sprintf("Player 2 Answer: %s\n\n", player2Answers[i])
	}

	return prompt
}

func sendPromptToOpenAI(prompt string) (OpenAIResponse, error) {

	client := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		panic(err.Error())
	}

	content := chatCompletion.Choices[0].Message.Content

	var openAiResponse OpenAIResponse
	err = json.Unmarshal([]byte(content), &openAiResponse)
	if err != nil {
		start := strings.Index(content, "{")
		end := strings.LastIndex(content, "}")
		if start != -1 && end != -1 && start < end {
			jsonContent := content[start : end+1]
			err = json.Unmarshal([]byte(jsonContent), &openAiResponse)
			if err != nil {
				return OpenAIResponse{}, fmt.Errorf("error parsing JSON response: %v", err)
			}
		} else {
			return OpenAIResponse{}, fmt.Errorf("invalid JSON response: %v", err)
		}
	}

	return openAiResponse, nil
}
