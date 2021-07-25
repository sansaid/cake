package main

import (
	"fmt"
	"os"
)

type ErrorCode struct {
	description string
	id          string
	code        int
}

func ExitErr(errCode ErrorCode, err error) {
	fmt.Printf("%s: %s\n", errCode.id, errCode.description)

	if err != nil {
		fmt.Println(err)
	}

	os.Exit(errCode.code)
}

// ErrNoRegistry - an error used to specify if no registry is defined when running cake
var ErrNoRegistry ErrorCode = ErrorCode{
	description: "Registry must be specified, either through -registry flag, CAKE_REGISTRY" +
		" environment variable, or through a .cake file.",
	id:   "NoRegistry",
	code: 1,
}

// ErrUnrecognisedSubcommands - an error to specify subcommands are uncregonised
var ErrUnrecognisedSubcommands ErrorCode = ErrorCode{
	description: "Expected at least one of the following subcommands: start, stop",
	id:          "UnrecognisedSubcommands",
	code:        2,
}

// ErrCakeNotFound - thrown when a call to a cake instance is required but is not found
var ErrCakeNotFound ErrorCode = ErrorCode{
	description: "No cake instance started",
	id:          "CakeNotFound",
	code:        3,
}

// ErrGettingRepoTags - thrown when connection to repo is disrupted
var ErrGettingRepoTags ErrorCode = ErrorCode{
	description: "Error while connecting to repo",
	id:          "FailedGettingRepoTags",
	code:        4,
}

// ErrReadingRepoTags - thrown when repo response could not be read
var ErrReadingRepoTags ErrorCode = ErrorCode{
	description: "Error while reading repo tags",
	id:          "FailedReadingRepoTags",
	code:        5,
}

// ErrJsonDecode - thrown when error decoding JSON
var ErrJsonDecode ErrorCode = ErrorCode{
	description: "Error while decoding JSON",
	id:          "FailedJsonDecode",
	code:        6,
}

// ErrCreateContainer - thrown when creating container
var ErrCreateContainer ErrorCode = ErrorCode{
	description: "Unknown exception while creating container",
	id:          "FailedCreateContainer",
	code:        7,
}
