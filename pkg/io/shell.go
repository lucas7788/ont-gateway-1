package io

// ShellInput for input
type ShellInput struct {
	Shell string
}

// ShellOutput for output
type ShellOutput struct {
	BaseResp
	Out string
}
