package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/erlorenz/bc-codegen/generate"
	"github.com/erlorenz/bc-codegen/metadata"
)

type CLI struct {
	InputFile  string
	OutputFile string
	Language   string
}

func main() {
	cli := &CLI{
		OutputFile: "schema.ts",
		Language:   "typescript",
	}

	flag.StringVar(&cli.OutputFile, "out", "schema.ts", "Output file path")
	flag.StringVar(&cli.OutputFile, "o", "schema.ts", "Output file path (short)")
	flag.StringVar(&cli.Language, "lang", "typescript", "Language to generate (typescript)")
	flag.Parse()
	// Get metadata file from positional argument
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: exactly one metadata XML file argument is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <metadata.xml>\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}
	cli.InputFile = args[0]

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *CLI) Run() error {
	model, err := metadata.Parse(c.InputFile)
	if err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	var generator generate.Generator
	switch c.Language {
	case "typescript":
		generator, err = generate.NewTypeScript(c.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create TypeScript generator: %w", err)
		}
	default:
		return fmt.Errorf("unsupported language: %s", c.Language)
	}

	return generator.Generate(*model)
}
