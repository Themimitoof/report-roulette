package main

type GitLabUser struct {
	Username        string `json:"username"`
	Name            string `json:"name"`
	MembershipState string `json:"membership_state"`
}

type GitLabGroup struct {
	Id       int    `json:"id"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}
