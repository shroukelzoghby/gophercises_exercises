package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Story map[string]struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func main() {
	var (
		flagStoryJSON = flag.String("storyJSON", "gopher.json", "Path to JSON file containing the story")
		flagCLI       = flag.Bool("cli", false, "Run the story in the terminal")
	)
	flag.Parse()

	file, err := os.Open(*flagStoryJSON)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	defer file.Close()

	var story Story
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&story); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if *flagCLI {
		runCLI(story)
		return
	}

	fmt.Println("Starting the Choose Your Own Adventure server on :8080")
	http.ListenAndServe(":8080", newStoryMux(story))

}

type StoryHandler struct {
	story Story
}

func (h *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arcName := r.URL.Query().Get("arc")
	if arcName == "" {
		arcName = "intro"
	}

	arc, ok := h.story[arcName]
	if !ok {
		http.Error(w, "Arc not found", http.StatusNotFound)
		return
	}

	temp, err := template.ParseFiles("arc.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, arc)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func newStoryMux(story Story) http.Handler {
	return &StoryHandler{story: story}
}

func runCLI(story Story) {
	reader := bufio.NewReader(os.Stdin)
	arcName := "intro"

	for {
		arc, ok := story[arcName]
		if !ok {
			fmt.Printf("Story arc %q was not found.\n", arcName)
			return
		}

		fmt.Printf("\n%s\n\n", arc.Title)
		for _, paragraph := range arc.Story {
			fmt.Println(paragraph)
		}

		if len(arc.Options) == 0 {
			fmt.Println("\nThe End.")
			return
		}

		fmt.Println("\nChoose an option:")
		for index, option := range arc.Options {
			fmt.Printf("Press %d to %s\n", index+1, option.Text)
		}

		choice, err := readChoice(reader, len(arc.Options))
		if err != nil {
			fmt.Println("Goodbye.")
			return
		}
		arcName = arc.Options[choice-1].Arc
	}
}

func readChoice(reader *bufio.Reader, optionCount int) (int, error) {
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err == nil && choice >= 1 && choice <= optionCount {
			return choice, nil
		}
		fmt.Printf("Please enter a number from 1 to %d.\n", optionCount)
	}
}
