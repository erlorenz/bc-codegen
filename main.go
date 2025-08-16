package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/erlorenz/bc-codegen/generate"
	"github.com/erlorenz/bc-codegen/metadata"
)

type CLI struct {
	InputFile string
	OutputDir string
	Name      string
	Language  string
}

func main() {
	cli := &CLI{
		OutputDir: "generated",
		Name:      "schemas",
		Language:  "typescript",
	}

	flag.StringVar(&cli.OutputDir, "outdir", "generated", "Output directory for generated files")
	flag.StringVar(&cli.Name, "name", "schema", "Name of the output file (without extension)")
	flag.StringVar(&cli.Language, "lang", "typescript", "Language to generate (typescript)")
	flag.Parse()

	// Make it v2.schema.ts insted of v2.ts
	if cli.Name != "schema" {
		cli.Name = cli.Name + ".schema"
	}
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

	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		return err
	}

	var generator generate.Generator
	switch c.Language {
	case "typescript":
		generator, err = generate.NewTypeScript(c.OutputDir, c.Name)
		if err != nil {
			return fmt.Errorf("failed to create TypeScript generator: %w", err)
		}
	default:
		return fmt.Errorf("unsupported language: %s", c.Language)
	}

	return generator.Generate(*model)
}
