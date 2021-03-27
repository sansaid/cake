package lib

import "fmt"

type ErrorCode struct {
	description string
	code        int
}

func ExitErr(err ErrorCode) {
	fmt.Println(err.description)
	os.Exit(err.code)
}

// ErrNoRegistry - an error used to specify if no registry is defined when running cake
ErrNoRegistry := ErrorCode{
	description: "Registry must be specified, either through -registry flag, CAKE_REGISTRY" +
	" environment variable, or through a .cake file.",
	code: 1,
}

// ErrUnrecognisedSubcommands - an error to specify subcommands are uncregonised
ErrUnrecognisedSubcommands := ErrorCode{
	description: "Expected at least one of the following subcommands: start, stop",
	code: 2,
}

// ErrCakeNotFound - thrown when a call to a cake instance is required but is not found
ErrCakeNotFound := ErrorCode{
	description: "No cake instance started",
	code: 3,
}

// ErrGettingRepoTags - thrown when connection to repo is disrupted
ErrGettingRepoTags := ErrorCode{
	description: "Error while connecting to repo",
	code: 3,
}

// ErrReadingRepoTags - thrown when repo response could not be read
ErrReadingRepoTags := ErrorCode{
	description: "Error while reading repo tags",
	code: 4,
}