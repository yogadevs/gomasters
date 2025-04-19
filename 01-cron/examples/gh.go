package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	cron "gomasters/01-cron"
)

type GitHubTask struct {
	Language string
}

func (gt *GitHubTask) Exec() {
	repos, err := getTopRepos(gt.Language)
	if err != nil {
		log.Printf("Error getting top %s reps: %v", gt.Language, err)
		return
	}

	log.Printf("Top %s reps:", gt.Language)
	for i, repo := range repos {
		if i >= 3 {
			break
		}
		log.Printf("  %d. %s - %s (%d)",
			i+1,
			repo.FullName,
			repo.Description,
			repo.Stars)
	}
}

type GitHubRepo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
}

func getTopRepos(language string) ([]GitHubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=language:%s&sort=stars&order=desc", language)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Go-Cron")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []GitHubRepo `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func main() {
	planner := cron.NewPlanner()

	planner.Add(&GitHubTask{Language: "go"}, time.Now().Add(5*time.Second))
	planner.Add(&GitHubTask{Language: "ruby"}, time.Now().Add(10*time.Second))
	planner.Add(&GitHubTask{Language: "javascript"}, time.Now().Add(10*time.Second))

	select {}
}
