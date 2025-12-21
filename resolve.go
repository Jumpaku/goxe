package xtracego

import (
	"strings"

	"github.com/Jumpaku/xtracego/internal"
)

type ResolvedResult struct {
	IsModule    bool
	ModuleName  string
	GoModFile   string
	SourceFiles []string
	PackageDir  string
}

func ResolvePackages(packageArgs []string) (resolved ResolvedResult, err error) {
	pkg, err := internal.ResolvePackage(strings.Join(packageArgs, ","))
	if err != nil {
		return ResolvedResult{}, err
	}
	return ResolvedResult{
		IsModule:    pkg.ResolveType != internal.ResolveType_CommandLineArguments,
		ModuleName:  pkg.Module,
		GoModFile:   pkg.GoModFile,
		SourceFiles: pkg.SourceFiles,
		PackageDir:  pkg.PackageDir,
	}, nil
}
