package prompt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aicirt2012/robobackup/src/sys/prompt/rand"
)

var reader io.Reader = os.Stdin

func SelectMultiLineOptions(label string, options [][]string) int {
	serializedOptions := []string{}
	for _, option := range options {
		serializedOptions = append(serializedOptions, strings.Join(option, "\n    "))
	}
	return Select(label, serializedOptions)
}

func Select(label string, options []string) int {
	r := bufio.NewReader(reader)
	log.Println("\n" + label)
	for i, option := range options {
		fmt.Printf("[%v] %v\n", i, option)
	}
	log.Print("Type option: ")
	in, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.Atoi(strings.TrimSpace(in))
	if err != nil || i < 0 || i > len(options)-1 {
		log.Println("Invalid input, try again!")
		return Select(label, options)
	}
	return i
}

func ConfirmWithYes(s string) bool {
	r := bufio.NewReader(reader)
	log.Print(s + " [y/n]: ")
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	in := strings.ToLower(strings.TrimSpace(res))
	return len(in) > 0 && in[0] == 'y'
}

func ConfirmWithNumber() {
	rand := rand.ThreeDigits()
	r := bufio.NewReader(reader)
	fmt.Printf("Please confirm execution with '%v': ", rand)
	in, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if rand != strings.TrimSpace(in) {
		log.Println("Confirmation failed!")
		os.Exit(1)
	}
}

func PressAnyKeyToClose() {
	r := bufio.NewReader(reader)
	log.Print("\nPress any key to close ")
	_, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
}
