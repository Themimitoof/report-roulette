package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
)

var (
	// Conf via env variables or flags
	host    string = "gitlab.com"
	token   string
	noProof bool
	silent  bool

	// Data to process/include in the roulette
	toLookup       []string
	gitLabUsers    []GitLabUser
	nonGitlabUsers []string
)

func main() {
	if h := os.Getenv("GITLAB_HOST"); h != "" {
		host = h
	}

	token = os.Getenv("GITLAB_TOKEN")
	if token == "" {
		fmt.Println("You need to specify the $GITLAB_TOKEN variable before. Please check: https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html")
		os.Exit(255)
		return
	}

	flag.BoolVar(&noProof, "n", false, "Don't display the list of people in the roulette.")
	flag.BoolVar(&silent, "s", false, "Output the name directly instead of a phrase.")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Please give a list of members of groups.")
		os.Exit(1)
	}

	for _, val := range slices.CompactFunc(flag.Args(), strings.EqualFold) {
		switch val[0] {
		case '+':
			nonGitlabUsers = append(nonGitlabUsers, val[1:])

		case '@':
			if groupUsers := lookupGroupUsers(val); len(groupUsers) != 0 {
				gitLabUsers = append(gitLabUsers, groupUsers...)
			}

		default:
			if user := lookupUser(val); user != (GitLabUser{}) {
				gitLabUsers = append(gitLabUsers, user)
			} else {
				if !slices.Contains[[]string, string](nonGitlabUsers, strings.ToLower(val)) {
					nonGitlabUsers = append(nonGitlabUsers, val)
				}
			}
		}
	}

	roulette := append([]string{}, nonGitlabUsers...)
	gitLabUsers = slices.CompactFunc(gitLabUsers, func(g GitLabUser, v GitLabUser) bool {
		return strings.ToLower(g.Username) == strings.ToLower(v.Username)
	})

	for _, user := range gitLabUsers {
		displayUser := fmt.Sprintf("%s (%s)", user.Name, user.Username)

		if !slices.Contains[[]string, string](roulette, strings.ToLower(user.Username)) && !slices.Contains[[]string, string](roulette, displayUser) {
			roulette = append(roulette, displayUser)
		}
	}

	if len(roulette) == 0 {
		fmt.Println("The roulette is empty. Pulling the trigger is useless.")
		return
	}

	slices.Sort(roulette)
	choosenUser := roulette[rand.Intn(len(roulette))]

	if silent {
		fmt.Printf("%s\n", choosenUser)
		return
	}

	// Display the list of members in the roulette
	if !noProof {
		fmt.Println("Users in the roulette:")
		for _, user := range roulette {
			fmt.Printf("  - %s\n", user)
		}

		fmt.Println("")
	}

	fmt.Printf("The roulette stopped on %s!\n", choosenUser)
}
