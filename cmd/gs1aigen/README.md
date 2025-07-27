# GS1 AI Generator

The basis for working with GS1 data are the [Application Identifiers](https://ref.gs1.org/ai/). Those are standardized
together with an extensive description at [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary).

This CLI generates Go Code for GS1 AI Descriptions using official Syntax Dictionary.


Usage:

    gs1aigen [flags]

The flags are:

	-disable-struct-gen
		Disables default AI struct declaration generation
	-out string
		Path to the output file (default "airegistry.go")
	-package string
		Package name to use (default "main")
	-release string
		Syntax Dictionary release to use (default "2025-01-30")
	-struct-name string
		Name of the struct to use for generating (default "ApplicationIdentifier")

The release flag uses the [Git Tags](https://github.com/gs1/gs1-syntax-dictionary/tags) used in the GS1 Syntax
Dictionary project. If you use the flag -disable-struct-gen, ensure to reference a struct via flag struct-name that 
exports the following fields:

- AI string
- Flags string
- Specification []string
- Attributes []string
- Title string