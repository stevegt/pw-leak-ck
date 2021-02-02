package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var version string

func main() {
	showMasked := gotFlag("-m")

	// handle ^C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	term := int(os.Stdin.Fd())
	termState, err := terminal.GetState(term)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-sigs
		terminal.Restore(term, termState)
		os.Exit(0)
	}()

	fmt.Println("enter passwords, one per line (^C to quit):")
	for {
		// read password from terminal
		fmt.Print("> ")
		pw, err := terminal.ReadPassword(term)
		if err != nil {
			log.Fatal(err)
		}
		if len(pw) == 0 {
			fmt.Println()
			continue
		}

		// get leak count
		leaks := ck(pw)

		// show results
		var masked string
		if showMasked && len(pw) > 1 {
			masked = string(pw[0:1]) + strings.Repeat("*", len(pw)-2) + string(pw[len(pw)-1:])
		} else {
			masked = ""
		}
		if leaks > 0 {
			fmt.Printf("%s leaked %d times\n", masked, leaks)
		} else {
			fmt.Printf("%s no known leaks\n", masked)
		}
	}

	return
}

func ck(pw []byte) (leaks int) {
	// generate hash
	bin := sha1.Sum(pw)
	hex := strings.ToUpper(fmt.Sprintf("%X", bin))
	first5 := hex[:5]

	// send first 5 bytes to server
	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", first5)
	var client = &http.Client{}
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// search response for a line that matches full hash
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		gothex := first5 + parts[0]
		if hex == gothex {
			leaks, err = strconv.Atoi(parts[1])
			if err != nil {
				log.Fatalf("%v: %v", err, line)
			}
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func gotFlag(flag string) bool {
	if len(os.Args) < 2 {
		return false
	}
	for _, arg := range os.Args[1:] {
		if flag == arg {
			return true
		}
	}
	return false
}
