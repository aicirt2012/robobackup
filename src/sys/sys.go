package sys

import (
	"bufio"
	"os/exec"

	"log"
)

func AssertAdminPermissions() {
	if !hasAdminPermissions() {
		log.Fatal("Admin permissions are required")
	}
}

func hasAdminPermissions() bool {
	cmd := exec.Command("net", "session")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return false
	}
	if err := cmd.Start(); err != nil {
		return false
	}
	in := bufio.NewScanner(stderr)
	for in.Scan() {
		return false
	}
	return true
}

func PrintHeadline(lines ...string) {
	log.Println("")
	log.Println(padRight("", '=', 80))
	for _, line := range lines {
		log.Println(padRight("::::: "+line+" ", ':', 80))
	}
	log.Println(padRight("", '=', 80))
}

func padRight(text string, char rune, length int) string {
	for len(text) < length {
		text += string(char)
	}
	return text
}
