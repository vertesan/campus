package jwto_test

import (
	"testing"
	"vertesan/campus/network/jwto"
)

func TestIsJwtExpired(t *testing.T) {
  expd := jwto.IsJwtExpired("eyJhbGciOiJSUzI1NiIsImtpZCI6IjNjOTNjMWEyNGNhZjgyN2I4ZGRlOWY4MmQyMzE1MzY1MDg4YWU2MTIiLCJ0eXAiOiJKV1QifQ.eyJwdWlkIjoiOVlLNU5FSEUiLCJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vYm5lLWNhbXB1cyIsImF1ZCI6ImJuZS1jYW1wdXMiLCJhdXRoX3RpbWUiOjE3MTYxMTYzNTYsInVzZXJfaWQiOiJhODRkMzI3YS0zM2Q0LTRjYzQtODE5Ny0xMWRhY2YyYTAzZGEiLCJzdWIiOiJhODRkMzI3YS0zM2Q0LTRjYzQtODE5Ny0xMWRhY2YyYTAzZGEiLCJpYXQiOjE3MTY0NjU4MjIsImV4cCI6MTcxNjQ2OTQyMiwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6e30sInNpZ25faW5fcHJvdmlkZXIiOiJjdXN0b20ifX0.RmhuwvI0h9zyxe1XLr81RF3CFDqoYQ7TFMolxtpniJF6oGNPmzO-m7A3vz1NOATo0hjTTjJ7hiPaPKhdYGHzLz0bbkrdga392ierBbYCJR0uRNae7r0jzXExHn6HUvK6QU45lorSBRvLKAbhRgJr4LwpNzyv-DzoTNOqttYgXJGoWchTK0edx0-sMZmOlzbdvmEaN2NynkPVYzdftJa69fmLY3DILrHxz2q3nDc4L6kiI4NAQwKi0rFcbYhroh79fGGfkdOkT_vvl_PM9jzBo-aqJUQemEDWmiPYxxtD3zfUs5qOJfj7vGonioMPN-a7ZQ09uBb2L-AAkn7IP_6QFA")
	if !expd {
		t.Fatal("failed")
	}
}
