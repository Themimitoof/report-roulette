package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var httpClient http.Client = http.Client{}

func lookupUser(username string) GitLabUser {
	apiUrl := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/v4/users",
	}

	qs := apiUrl.Query()
	qs.Set("search", username)
	qs.Set("active", "true")
	qs.Set("locked", "false")
	qs.Set("exclude_external", "true")
	qs.Set("without_project_bots", "true")
	apiUrl.RawQuery = qs.Encode()

	req, _ := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, _ := httpClient.Do(req)
	if resp.StatusCode != 200 {
		fmt.Printf("WARN: Unable to lookup for user '%s'.\n", username)
		return GitLabUser{}
	}

	body, _ := io.ReadAll(resp.Body)
	var users []GitLabUser
	json.Unmarshal(body, &users)

	for _, user := range users {
		if user.Username == strings.ToLower(username) {
			return user
		}
	}

	return GitLabUser{}
}

func getGroupId(groupName string) GitLabGroup {
	apiUrl := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/v4/groups",
	}

	qs := apiUrl.Query()
	qs.Set("search", groupName)
	apiUrl.RawQuery = qs.Encode()

	req, _ := http.NewRequest(http.MethodGet, apiUrl.String(), nil)
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, _ := httpClient.Do(req)
	if resp.StatusCode != 200 {
		fmt.Printf("WARN: Unable to lookup for group '%s'.\n", groupName)
		return GitLabGroup{}
	}

	body, _ := io.ReadAll(resp.Body)
	var groups []GitLabGroup
	json.Unmarshal(body, &groups)

	// Compare with the fullpath
	for _, group := range groups {
		if group.FullPath == strings.ToLower(groupName) {
			return group
		}
	}

	// Fallback to the path if no group has been found
	for _, group := range groups {
		if group.Path == strings.ToLower(groupName) {
			return group
		}
	}

	return GitLabGroup{}
}

func lookupGroupUsers(groupName string) []GitLabUser {
	groupName = groupName[1:] // Remove @ at the beginning of the group name
	group := getGroupId(groupName)

	if group == (GitLabGroup{}) {
		return nil
	}

	membersUrl := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   fmt.Sprintf("/api/v4/groups/%d/members", group.Id),
	}

	req, _ := http.NewRequest(http.MethodGet, membersUrl.String(), nil)
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, _ := httpClient.Do(req)
	if resp.StatusCode != 200 {
		fmt.Printf("WARN: Unable to lookup for group '%s'.\n", groupName)
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var users []GitLabUser
	json.Unmarshal(body, &users)

	var filteredUsers []GitLabUser
	for _, user := range users {

		if !strings.Contains(user.Username, "_bot") && user.MembershipState == "active" {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers
}
