package lib

import (
	"fmt"
	"os"
)

type ErrorCode struct {
	description string
	code        int
}

func ExitErr(errCode ErrorCode, err error) {
	fmt.Println(errCode.description)

	if err != nil {
		fmt.Println(err)
	}

	os.Exit(errCode.code)
}

// ErrNoRegistry - an error used to specify if no registry is defined when running cake
var ErrNoRegistry ErrorCode = ErrorCode{
	description: "Registry must be specified, either through -registry flag, CAKE_REGISTRY" +
		" environment variable, or through a .cake file.",
	code: 1,
}

// ErrUnrecognisedSubcommands - an error to specify subcommands are uncregonised
var ErrUnrecognisedSubcommands ErrorCode = ErrorCode{
	description: "Expected at least one of the following subcommands: start, stop",
	code:        2,
}

// ErrCakeNotFound - thrown when a call to a cake instance is required but is not found
var ErrCakeNotFound ErrorCode = ErrorCode{
	description: "No cake instance started",
	code:        3,
}

// ErrGettingRepoTags - thrown when connection to repo is disrupted
var ErrGettingRepoTags ErrorCode = ErrorCode{
	description: "Error while connecting to repo",
	code:        4,
}

// ErrReadingRepoTags - thrown when repo response could not be read
var ErrReadingRepoTags ErrorCode = ErrorCode{
	description: "Error while reading repo tags",
	code:        5,
}

// ErrJsonDecode - thrown when error decoding JSON
var ErrJsonDecode ErrorCode = ErrorCode{
	description: "Error while decoding JSON",
	code:        6,
}
