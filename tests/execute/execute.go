package execute

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func Binary(dir string, binary string, inputs ...string) {
	var wg sync.WaitGroup
	cmd := exec.Command(binary)
	cmd.Dir = dir
	r, w, _ := os.Pipe()
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, input := range inputs {
			fmt.Fprintln(w, input)
			time.Sleep(3 * time.Duration(time.Second))
		}
		w.Close()
	}()

	wg.Wait()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}
	r.Close()
}
