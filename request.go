package main

type HelloRequest struct {
	Name    string `param:"name"`
	Message string `json:"message"`
}
