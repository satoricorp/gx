package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/inference"
)

func newSetCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Configure local gx settings",
	}
	cmd.AddCommand(newSetKeyCommand(ctx))
	return cmd
}

func newSetKeyCommand(ctx context.Context) *cobra.Command {
	var provider string
	var key string
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Set the API key gx uses for local inference",
		Long: strings.Join([]string{
			"Set the API key gx uses for local inference.",
			"",
			"The key is stored in ~/.gx/inference.json or $GX_HOME/inference.json and replaces any previous inference key.",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = ctx
			input := inferencePromptInput(cmd.InOrStdin())
			resolvedProvider := strings.TrimSpace(provider)
			var err error
			if resolvedProvider == "" {
				resolvedProvider, err = promptInferenceProvider(input, cmd.OutOrStdout())
				if err != nil {
					return err
				}
			}
			normalized, err := inference.NormalizeProvider(resolvedProvider)
			if err != nil {
				return err
			}
			resolvedKey := strings.TrimSpace(key)
			if resolvedKey == "" {
				resolvedKey, err = promptInferenceKey(input, cmd.OutOrStdout(), normalized)
				if err != nil {
					return err
				}
			}
			if err := inference.Save(inference.Credentials{Provider: normalized, APIKey: resolvedKey}); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), labelValue(inference.ProviderDisplay(normalized), "key saved"))
			return nil
		},
	}
	cmd.Flags().StringVar(&provider, "provider", "", "inference provider: anthropic or openai")
	cmd.Flags().StringVar(&key, "key", "", "API key to store; omit to paste interactively")
	return cmd
}

func inferencePromptInput(in io.Reader) io.Reader {
	if file, ok := in.(*os.File); ok && term.IsTerminal(file.Fd()) {
		return in
	}
	return bufio.NewReader(in)
}

func promptInferenceProvider(in io.Reader, out io.Writer) (string, error) {
	fmt.Fprintln(out, section("Inference provider"))
	fmt.Fprintf(out, "  %s  Anthropic\n", command("1"))
	fmt.Fprintf(out, "  %s  OpenAI\n", command("2"))
	fmt.Fprintf(out, "\n%s ", muted("Choose provider [1-2]:"))

	raw, err := readInferencePromptLine(in)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	choice := strings.ToLower(strings.TrimSpace(raw))
	if choice == "" {
		return "", fmt.Errorf("provider is required")
	}
	if index, err := strconv.Atoi(choice); err == nil {
		switch index {
		case 1:
			return inference.ProviderAnthropic, nil
		case 2:
			return inference.ProviderOpenAI, nil
		default:
			return "", fmt.Errorf("selection %d out of range", index)
		}
	}
	return inference.NormalizeProvider(choice)
}

func promptInferenceKey(in io.Reader, out io.Writer, provider string) (string, error) {
	fmt.Fprintf(out, "%s ", muted(fmt.Sprintf("Paste %s API key:", inference.ProviderDisplay(provider))))
	if file, ok := in.(*os.File); ok && term.IsTerminal(file.Fd()) {
		data, err := term.ReadPassword(file.Fd())
		fmt.Fprintln(out)
		if err != nil {
			return "", err
		}
		key := strings.TrimSpace(string(data))
		if key == "" {
			return "", fmt.Errorf("api key is required")
		}
		return key, nil
	}
	raw, err := readInferencePromptLine(in)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	key := strings.TrimSpace(raw)
	if key == "" {
		return "", fmt.Errorf("api key is required")
	}
	return key, nil
}

type stringLineReader interface {
	ReadString(delim byte) (string, error)
}

func readInferencePromptLine(in io.Reader) (string, error) {
	if reader, ok := in.(stringLineReader); ok {
		return reader.ReadString('\n')
	}
	return bufio.NewReader(in).ReadString('\n')
}
