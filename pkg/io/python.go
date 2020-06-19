package io

// PythonInput for input
type PythonInput struct {
	Py string
}

// PythonOutput for output
type PythonOutput struct {
	BaseResp
	Out string
}
