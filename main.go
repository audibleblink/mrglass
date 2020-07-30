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

	USERHASH is a newline-seperated file in the the format of 
	'USER:HASH'

	CRACKED is a newline-seperated file or pipe whose entries 
	are in the the format of 
	'HASH:PLAINTEXT_PASSWORD'

Examples:
	mrglass hashes_with_usernames.txt cracked.txt
	hashstack lists cracked 1 1 | mrglass hashes_with_usernames.txt
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

	correlated, errs := correlate(crackScanner, userHashFile)
	if errs != nil {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	display(correlated)
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
func correlate(crackScanner *bufio.Scanner, users *os.File) (map[string][]string, []error) {
	userHash := loadHashMap(users)
	loot := make(map[string][]string)
	var errs []error
	for crackScanner.Scan() {
		line := crackScanner.Text()
		hashAndPass := strings.SplitN(line, ":", 2)
		if len(hashAndPass) != 2 {
			lineErr := fmt.Errorf("<CRACKED>: bad format on line %s\n", hashAndPass[0])
			errs = append(errs, lineErr)
			continue
		}
		hash, pass := hashAndPass[0], hashAndPass[1]
		users := userHash[hash]
		loot[pass] = append(loot[pass], users...)
	}
	return loot, errs
}

// loadHashMap puts in memory a map of password hashes with usernames
// as the values for easy retreival
func loadHashMap(hashes *os.File) map[string][]string {
	scanner := bufio.NewScanner(hashes)
	userHash := make(map[string][]string)
	for scanner.Scan() {
		userAndHash := strings.SplitN(scanner.Text(), ":", 2)
		if len(userAndHash) != 2 {
			fmt.Fprintln(os.Stderr, "<USERHASH>: bad format")
			continue
		}
		user, hash := userAndHash[0], userAndHash[1]
		userHash[hash] = append(userHash[hash], user)
	}
	return userHash
}

// hasPipe tells us whether or not mrglass is part of a pipeline
func hasPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		bail(err)
	}
	return info.Mode()&os.ModeNamedPipe != 0
}

func display(loot map[string][]string) {
	for pass, users := range loot {
		for _, user := range users {
			fmt.Printf("%s:%s\n", user, pass)
		}
	}
}

func bail(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	help()
}
