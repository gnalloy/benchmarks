package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gnalloy.org/benchmarks/parity"
	"gnalloy.org/benchmarks/parity/matrix"
)

type config struct {
	specPath       string
	matrixName     string
	format         parity.Format
	outPath        string
	dryRun         bool
	dumpSpec       bool
	strictExternal bool
	timeout        time.Duration
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer, stderr io.Writer) error {
	cfg, err := parseFlags(args, stderr)
	if err != nil {
		return err
	}
	spec, err := loadSpec(cfg)
	if err != nil {
		return err
	}
	if cfg.dumpSpec {
		return writeSpec(cfg.outPath, stdout, spec)
	}
	if cfg.strictExternal {
		if err := parity.ValidateExternalHarnesses(spec, parity.ExternalHarnessOptions{}); err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()
	report, err := parity.Runner{DryRun: cfg.dryRun}.Run(ctx, spec)
	if err != nil {
		return err
	}
	return writeReport(cfg.outPath, stdout, report, cfg.format)
}

func parseFlags(args []string, stderr io.Writer) (config, error) {
	cfg := config{format: parity.FormatMarkdown, timeout: 6 * time.Hour}
	fs := flag.NewFlagSet("gnalloy-parity", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&cfg.specPath, "spec", "", "parity spec JSON path")
	fs.StringVar(&cfg.matrixName, "matrix", "", "built-in matrix: linux-full or windows-full")
	fs.Var((*formatValue)(&cfg.format), "format", "report format: markdown or json")
	fs.StringVar(&cfg.outPath, "out", "", "output path; empty writes stdout")
	fs.BoolVar(&cfg.dryRun, "dry-run", false, "expand and report scenarios without running commands")
	fs.BoolVar(&cfg.dumpSpec, "dump-spec", false, "write the selected spec as JSON instead of running it")
	fs.BoolVar(&cfg.strictExternal, "strict-external", false, "require every executable external harness to be present before running")
	fs.DurationVar(&cfg.timeout, "timeout", cfg.timeout, "overall run timeout")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if strings.TrimSpace(cfg.specPath) == "" && strings.TrimSpace(cfg.matrixName) == "" {
		return config{}, fmt.Errorf("missing -spec or -matrix")
	}
	if strings.TrimSpace(cfg.specPath) != "" && strings.TrimSpace(cfg.matrixName) != "" {
		return config{}, fmt.Errorf("-spec and -matrix are mutually exclusive")
	}
	if cfg.timeout <= 0 {
		return config{}, fmt.Errorf("timeout must be positive")
	}
	return cfg, nil
}

func loadSpec(cfg config) (parity.Spec, error) {
	if strings.TrimSpace(cfg.specPath) != "" {
		file, err := os.Open(cfg.specPath)
		if err != nil {
			return parity.Spec{}, err
		}
		defer file.Close()
		return parity.LoadSpec(file)
	}
	switch strings.ToLower(strings.TrimSpace(cfg.matrixName)) {
	case "linux-full":
		return matrix.LinuxFullSpec(), nil
	case "windows-full":
		return matrix.WindowsFullSpec(), nil
	default:
		return parity.Spec{}, fmt.Errorf("unknown matrix %q", cfg.matrixName)
	}
}

func writeReport(path string, stdout io.Writer, report parity.Report, format parity.Format) error {
	if strings.TrimSpace(path) == "" {
		return parity.WriteReport(stdout, report, format)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return parity.WriteReport(file, report, format)
}

func writeSpec(path string, stdout io.Writer, spec parity.Spec) error {
	if strings.TrimSpace(path) == "" {
		return encodeSpec(stdout, spec)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return encodeSpec(file, spec)
}

func encodeSpec(w io.Writer, spec parity.Spec) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(spec)
}

type formatValue parity.Format

func (f *formatValue) String() string {
	if f == nil {
		return ""
	}
	return string(*f)
}

func (f *formatValue) Set(value string) error {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "md", "markdown":
		*f = formatValue(parity.FormatMarkdown)
	case "json":
		*f = formatValue(parity.FormatJSON)
	default:
		return fmt.Errorf("unsupported format %q", value)
	}
	return nil
}
