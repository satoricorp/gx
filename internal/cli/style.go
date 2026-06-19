package cli

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/termstyle"
)

func enableColor() bool {
	return termstyle.Enabled()
}

func danger(text string) string  { return termstyle.Danger(text) }
func success(text string) string { return termstyle.Success(text) }
func section(text string) string { return termstyle.Section(text) }
func command(text string) string { return termstyle.Command(text) }
func muted(text string) string   { return termstyle.Muted(text) }
func accent(text string) string  { return termstyle.Accent(text) }
func hint(text string) string    { return termstyle.Hint(text) }
func mint(text string) string    { return termstyle.Mint(text) }
func value(text string) string   { return termstyle.Value(text) }
func divider(width int) string   { return termstyle.Divider(width) }
func labelWarningValue(label, value string) string {
	return termstyle.LabelWarning(label, value)
}
func labelMintValue(label, value string) string { return termstyle.LabelMint(label, value) }
func labelStatus(label, value string) string    { return termstyle.LabelStatus(label, value) }
func progressBar(percent float64) string {
	return termstyle.ProgressBar(percent)
}
func commandLine(invocation string, accentCommand bool) string {
	return termstyle.CommandLine(invocation, accentCommand)
}

func hyperlink(url, text string) string {
	return termstyle.Hyperlink(url, text)
}

const (
	groupSetup    = "setup"
	groupWork     = "work"
	groupShip     = "ship"
	groupAdvanced = "advanced"
)

func installHelpStyling(root *cobra.Command) {
	cobra.AddTemplateFunc("section", section)
	cobra.AddTemplateFunc("command", command)
	cobra.AddTemplateFunc("commandPadded", commandPadded)
	cobra.AddTemplateFunc("commandDisplayPadded", commandDisplayPadded)
	cobra.AddTemplateFunc("helpDescription", helpDescription)
	cobra.AddTemplateFunc("muted", muted)

	subcommandTemplate := helpBodyTemplate()
	rootHelpBodyTemplate := helpBodyTemplate()

	walkCommands(root, func(cmd *cobra.Command) {
		if cmd == root {
			cmd.SetHelpTemplate(rootHelpBodyTemplate)
			return
		}
		cmd.SetHelpTemplate(subcommandTemplate)
	})
}

func helpBodyTemplate() string {
	intro := `{{if not (eq .CommandPath "gx")}}{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}{{end}}

{{end}}`

	commandBlock := `{{$commands := .Commands}}{{section "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{commandDisplayPadded . $commands}} {{helpDescription .Short}}{{end}}{{end}}`

	return intro + `{{if or .Runnable .HasSubCommands}}{{section "Usage:"}}
  {{if .Runnable}}{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}{{if .Runnable}}
  {{end}}{{.CommandPath}} [command]{{end}}{{end}}{{if gt (len .Aliases) 0}}

{{section "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{section "Examples:"}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

` + commandBlock + `{{end}}{{if .HasAvailableLocalFlags}}

{{section "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{section "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

{{section "Additional Help Topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{commandPadded .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

{{if eq .CommandPath "gx"}}{{muted (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}{{end}}
`
}

func printRootIntro(out io.Writer) {
	fmt.Fprintln(out, renderStaticLogo())
	fmt.Fprintln(out, muted(gxTagline))
	fmt.Fprintln(out, rootAuthStatusLine())
	fmt.Fprintln(out)
}

func rootAuthStatusLine() string {
	creds, err := cloud.LoadCloudCredentials()
	if err == nil && creds != nil && strings.TrimSpace(creds.Login) != "" {
		return success("●") + " " + value("Signed in as "+strings.TrimSpace(creds.Login))
	}
	return danger("●") + " " + muted("Not signed in") + "  " + logoText("gx auth login")
}

func printRootHelp(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()
	printRootIntro(out)
	return renderHelpTemplate(cmd, out)
}

func renderHelpTemplate(cmd *cobra.Command, out io.Writer) error {
	invocation := helpInvocation(cmd)
	if invocation != "" {
		fmt.Fprintln(out, commandLine(invocation, cmd.Root() != cmd))
		fmt.Fprintln(out)
	}
	t := template.New("help").Funcs(helpTemplateFuncs())
	parsed, err := t.Parse(cmd.HelpTemplate())
	if err != nil {
		return err
	}
	return parsed.Execute(out, cmd)
}

func helpInvocation(cmd *cobra.Command) string {
	path := strings.TrimSpace(cmd.CommandPath())
	if path == "" {
		return ""
	}
	if cmd.Root() == cmd {
		return "gx help"
	}
	if cmd.Name() == "help" {
		return path
	}
	return path + " -h"
}

func helpTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"trim":                    strings.TrimSpace,
		"trimRightSpace":          trimRightSpace,
		"trimTrailingWhitespaces": trimRightSpace,
		"rpad":                    helpRpad,
		"gt":                      cobra.Gt,
		"eq":                      cobra.Eq,
		"section":                 section,
		"command":                 command,
		"commandPadded":           commandPadded,
		"commandDisplayPadded":    commandDisplayPadded,
		"helpDescription":         helpDescription,
		"muted":                   muted,
	}
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func helpRpad(s string, padding int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", padding), s)
}

func walkCommands(root *cobra.Command, fn func(*cobra.Command)) {
	fn(root)
	for _, cmd := range root.Commands() {
		walkCommands(cmd, fn)
	}
}

func labelValue(label, value string) string {
	return termstyle.LabelValue(label, value)
}

func commandPadded(text string, width int) string {
	return command(padRight(text, width))
}

const helpNameAnnotation = "gx.helpName"

func setHelpName(cmd *cobra.Command, name string) {
	if strings.TrimSpace(name) == "" {
		return
	}
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[helpNameAnnotation] = name
}

func commandDisplayPadded(cmd *cobra.Command, commands []*cobra.Command) string {
	return logoText(padRight(commandDisplay(cmd), commandDisplayPadding(commands)))
}

func helpDescription(text string) string {
	return muted(text)
}

func commandDisplay(cmd *cobra.Command) string {
	return cmd.Name()
}

func commandDisplayPadding(commands []*cobra.Command) int {
	width := 10
	for _, cmd := range commands {
		if !cmd.IsAvailableCommand() && cmd.Name() != "help" {
			continue
		}
		width = max(width, runewidth.StringWidth(commandDisplay(cmd)))
	}
	return width
}

func padRight(value string, width int) string {
	padding := width - runewidth.StringWidth(value)
	if padding <= 0 {
		return value
	}
	return value + strings.Repeat(" ", padding)
}
