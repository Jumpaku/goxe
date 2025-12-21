package internal

type Config struct {
	TraceStmt bool
	TraceVar  bool
	TraceCall bool

	ShowTimestamp bool
	ShowGoroutine bool

	UniqueString string
	LineWidth    int

	ResolveType ResolveType
	ModuleName  string
}

func (cfg *Config) LibraryPackageName() string {
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return "main"
	}
	return "xtracego_" + cfg.UniqueString
}

func (cfg *Config) LibraryImportPath() string {
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return ""
	}
	return cfg.ModuleName + "/" + cfg.LibraryPackageName()
}

func (cfg *Config) LibraryFileName() string {
	return "xtracego_" + cfg.UniqueString + ".go"
}

func (cfg *Config) ExecutableFileName() string {
	return "main_" + cfg.UniqueString
}

func (cfg *Config) IdentifierPrintlnStatement() string {
	funcName := "PrintlnStatement_" + cfg.UniqueString
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return funcName
	}
	return cfg.LibraryPackageName() + "." + funcName
}

func (cfg *Config) IdentifierPrintlnVariable() string {
	funcName := "PrintlnVariable_" + cfg.UniqueString
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return funcName
	}
	return cfg.LibraryPackageName() + "." + funcName
}

func (cfg *Config) IdentifierPrintlnReturnVariable() string {
	funcName := "PrintlnReturnVariable_" + cfg.UniqueString
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return funcName
	}
	return cfg.LibraryPackageName() + "." + funcName
}

func (cfg *Config) IdentifierPrintlnCall() string {
	funcName := "PrintlnCall_" + cfg.UniqueString
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return funcName
	}
	return cfg.LibraryPackageName() + "." + funcName
}

func (cfg *Config) IdentifierPrintlnReturn() string {
	funcName := "PrintlnReturn_" + cfg.UniqueString
	if cfg.ResolveType == ResolveType_CommandLineArguments {
		return funcName
	}
	return cfg.LibraryPackageName() + "." + funcName
}
