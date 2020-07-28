package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func help() {
	message := `
Usage: 
	mrglass <USERHASH> <CRACKED>

	<CRACKED> | mrglass <USERHASH> 

	mrglass takes USERHASH as a newline-seperated file in the
	the format of 'USER:HASH'

	mrglass takes CRACKED as a newline-seperated file or 
	pipe whose entries are in the the format of 'HASH:PLAINTEXT'

Examples:
	mrglass hashes_with_usernames.txt cracked.txt

	hashstack lists cracked company | mrglass hashes_with_usernames.txt

	mrglass hashes_with_usernames.txt <(hashstack lists cracked $pID $lID)

	`

	fmt.Println(message)
	os.Exit(1)
}

func main() {
	crackScanner, err := newCrackScanner(os.Args)
	if err != nil {
		bail(err)
	}

	userHashFile, err := os.Open(os.Args[1])
	if err != nil {
		bail(err)
	}

	correlate(crackScanner, userHashFile)
}

// newCrackScanner conditionally evaluates whether our source
// of cracked passwords to correlate will come from a file or pipe
func newCrackScanner(args []string) (crackScanner *bufio.Scanner, err error) {

	if len(args) == 2 && hasPipe() {
		crackScanner = bufio.NewScanner(os.Stdin)

	} else if len(args) == 3 {
		var crackedFile *os.File
		crackedFile, err = os.Open(args[2])
		if err != nil {
			return
		}
		crackScanner = bufio.NewScanner(crackedFile)

	} else {
		err = fmt.Errorf("Invalid arguments")
	}
	return
}

// correlate looks up individual cracked password hashes
// against our memoized hash:user map
func correlate(crackScanner *bufio.Scanner, users *os.File) {
	userHash := loadHashMap(users)
	for crackScanner.Scan() {
		line := crackScanner.Text()
		hashAndPass := strings.SplitN(line, ":", 2)
		hash, pass := hashAndPass[0], hashAndPass[1]
		users := userHash[hash]
		for _, user := range users {
			fmt.Printf("%s:%s\n", user, pass)
		}
	}
}

// hasPipe tells us whether or not mrglass is part of a pipeline
func hasPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		bail(err)
	}
	return info.Mode()&os.ModeNamedPipe != 0
}

// loadHashMap puts in memory a map of password hashes with usernames
// as the values for easy retreival
func loadHashMap(hashes *os.File) (userHash map[string][]string) {
	scanner := bufio.NewScanner(hashes)
	userHash = make(map[string][]string)
	for scanner.Scan() {
		userAndHash := strings.SplitN(scanner.Text(), ":", 2)
		user, hash := userAndHash[0], userAndHash[1]
		userHash[hash] = append(userHash[hash], user)
	}
	return userHash
}

func bail(err error) {
	fmt.Println(err.Error())
	help()
}
