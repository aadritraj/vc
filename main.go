package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

var GENAI_API_KEY string = os.Getenv("GENAI_API_KEY")
var DEFAULT_PROMPT = "You are a model for a 'compiler' of a language. Respond ONLY with plaintext code that corresponds to the prompt provided. Do NOT generate any markdown codeblocks (e.g., ```language\ncode\n```) or any additional comments. You will be given a language hint and a prompt. The language hint is a single word that indicates the language. The prompt is a description of the code to be generated."

var langHintVar string
var infileVar string
var outfileVar string

func init() {
	flag.StringVar(&infileVar, "i", "", "Input file to read the prompt from")
	flag.StringVar(&outfileVar, "o", "", "Output file to write the response to")
	flag.StringVar(&langHintVar, "lang", "go", "Language hint for the model")
	flag.StringVar(&GENAI_API_KEY, "key", GENAI_API_KEY, "API key for the model")

	flag.Parse()
}

func main() {
	if len(infileVar) == 0 {
		log.Fatal("Input file must be specified with -i flag")
	}

	ctx := context.Background()
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: DEFAULT_PROMPT},
			},
		},
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  GENAI_API_KEY,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	outfileAbs, err := filepath.Abs(infileVar)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(outfileAbs); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", outfileAbs)
	}

	infileContents, err := os.ReadFile(outfileAbs)
	if err != nil {
		log.Fatalf("Error reading file '%s': %v", outfileAbs, err)
	}

	prompt := fmt.Sprintf("Language hint: %s\nPrompt: %s", langHintVar, infileContents)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(outfileVar) != 0 {
		outfileAbs, err = filepath.Abs(outfileVar)
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(outfileAbs, []byte(result.Text()), 0644); err != nil {
			log.Fatalf("Error writing to file '%s': %v", outfileAbs, err)
		}

		log.Printf("Response written to '%s'", outfileAbs)
	} else {
		fmt.Println(result.Text())
		log.Println("Response printed to stdout")
	}
}
