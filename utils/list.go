package utils

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
)

func ListSessions(cCtx *cli.Context) {
	var sessions []string
	hasFlags := cCtx.Bool("tmux") || cCtx.Bool("zoxide")

	if !hasFlags || cCtx.Bool("tmux") {
		tmuxSessions, err := getTmuxSessions()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		sessions = append(sessions, tmuxSessions...)
	}

	if !hasFlags || cCtx.Bool("zoxide") {
		zoxideResults, err := getZoxideResults()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		sessions = append(sessions, zoxideResults...)
	}

	fmt.Println(strings.Join(sessions, "\n"))
}

func getTmuxSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_last_attached} #{session_name}")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sessionList := strings.TrimSpace(string(output))
	sessionItems := strings.Split(sessionList, "\n")

	// Custom sorting by session_last_attached in reverse order
	sort.SliceStable(sessionItems, func(i, j int) bool {
		return sessionItems[i] > sessionItems[j]
	})

	sessions := make([]string, len(sessionItems))
	for i, item := range sessionItems {
		fields := strings.Fields(item)
		sessions[i] = fields[1]
	}

	return sessions, nil
}

func getZoxideResults() ([]string, error) {
	cmd := exec.Command("zoxide", "query", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	resultList := strings.TrimSpace(string(output))
	results := strings.Split(resultList, "\n")
	return results, nil
}