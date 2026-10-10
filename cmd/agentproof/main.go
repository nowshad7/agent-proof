// Command agentproof maps explicit acceptance requirements to executable
// evidence and produces portable verification receipts for agent-written
// changes.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/nowshad7/agent-proof/internal/proof"
)

const usage = `agentproof — reproducible verification for agent-written changes

Usage:
  agentproof version
  agentproof verify [flags]
  agentproof report --input <receipt.json> [--format json|md|html]

verify flags:
  --trust-local        Required. Acknowledge that checks run locally with your
                       privileges. This is not a sandbox.
  --repo <dir>         Repository to verify (default ".").
  --commit <rev>       Revision to resolve and snapshot (default "HEAD").
  --config <path>      Use a separately reviewed policy instead of the
                       committed agent-proof.json.
  --out <dir>          Directory for receipt files (default ".agent-proof").
  --runner <name>      Optional runner identity recorded in the receipt.

Exit codes:
  0  all required acceptance passed
  1  acceptance failed or required evidence was unverified/skipped
  2  configuration, execution, parsing or report error
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, usage)
		return 2
	}
	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println("agentproof " + proof.Version)
		return 0
	case "help", "--help", "-h":
		fmt.Print(usage)
		return 0
	case "verify":
		return runVerify(args[1:])
	case "report":
		return runReport(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}

func runVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	trust := fs.Bool("trust-local", false, "acknowledge local execution")
	repo := fs.String("repo", ".", "repository directory")
	commit := fs.String("commit", "HEAD", "revision to snapshot")
	config := fs.String("config", "", "separately reviewed policy path")
	out := fs.String("out", ".agent-proof", "output directory")
	runner := fs.String("runner", "", "runner identity")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if !*trust {
		fmt.Fprintln(os.Stderr, "error: verify requires --trust-local. Checks run locally with your privileges; this is not a sandbox.")
		return 2
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	receipt, err := proof.Verify(ctx, proof.Options{
		RepoDir:    *repo,
		Revision:   *commit,
		ConfigPath: *config,
		Runner:     *runner,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	if err := writeReceipts(*out, receipt); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	printSummary(receipt, *out)

	switch receipt.Status {
	case proof.StatusPass:
		return 0
	case proof.StatusError:
		return 2
	default:
		return 1
	}
}

func writeReceipts(dir string, r *proof.Receipt) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	jsonBytes, err := proof.RenderJSON(r)
	if err != nil {
		return err
	}
	files := map[string][]byte{
		"receipt.json": jsonBytes,
		"receipt.md":   proof.RenderMarkdown(r),
		"receipt.html": proof.RenderHTML(r),
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", name, err)
		}
	}
	return nil
}

func printSummary(r *proof.Receipt, out string) {
	fmt.Printf("AgentProof %s\n", r.Status)
	fmt.Printf("  commit %s\n", r.Commit)
	for _, rr := range r.Requirements {
		fmt.Printf("  [%s] %s\n", rr.State, rr.ID)
	}
	fmt.Printf("Receipts written to %s/\n", out)
}

func runReport(args []string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	input := fs.String("input", "", "receipt.json path")
	format := fs.String("format", "md", "json|md|html")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: --input is required")
		return 2
	}
	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	receipt, err := proof.ParseReceipt(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}
	switch *format {
	case "json":
		b, err := proof.RenderJSON(receipt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 2
		}
		os.Stdout.Write(b)
	case "md", "markdown":
		os.Stdout.Write(proof.RenderMarkdown(receipt))
	case "html":
		os.Stdout.Write(proof.RenderHTML(receipt))
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format %q\n", *format)
		return 2
	}
	return 0
}
