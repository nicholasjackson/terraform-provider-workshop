package main

import (
	"context"
	"strings"
)

type DaggerCI struct {
}

func New() *DaggerCI {
	return &DaggerCI{}
}

func (d *DaggerCI) FunctionalTest(
	ctx context.Context,
	// source directory containing the tests
	src *Directory,
	// working directory for the tests, relative to the source directory
	WorkingDirectory,
	// runtime to use for the tests, either docker or podman
	Runtime string,
) error {
	wd := strings.TrimPrefix(WorkingDirectory, "/")
	pl := dag.Pipeline("functional-test-" + wd + "-" + Runtime)

	// get the architecture of the current machine
	platform, err := pl.DefaultPlatform(ctx)
	if err != nil {
		panic(err)
	}

	arch := strings.Split(string(platform), "/")[1]

	_, err = pl.Jumppad().
		TestBlueprintWithVersion(
			ctx,
			src,
			"0.12.1",
			JumppadTestBlueprintWithVersionOpts{WorkingDirectory: WorkingDirectory, Architecture: arch},
		)

	return err
}
