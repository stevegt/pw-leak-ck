package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

var version string

func main() {
	showMasked := gotFlag("-m")
	fmt.Println("enter passwords, one per line:")
	for {
		ck(showMasked)
	}
}

func ck(showMasked bool) {
	fmt.Print("> ")
	buf, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	if len(buf) == 0 {
		fmt.Println()
		return
	}
	var masked string
	if showMasked {
		masked = string(buf[0:1]) + strings.Repeat("*", len(buf)-2) + string(buf[len(buf)-1:])
	} else {
		masked = ""
	}

	bin := sha1.Sum(buf)
	hex := strings.ToUpper(fmt.Sprintf("%X", bin))
	first5 := hex[:5]

	url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", first5)

	var client = &http.Client{}
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	leaks := 0
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
	if leaks > 0 {
		fmt.Printf("%s leaked %d times\n", masked, leaks)
	} else {
		fmt.Printf("%s no known leaks\n", masked)
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
