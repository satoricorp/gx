package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/inference"
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
func value(text string) string   { return termstyle.Value(text) }
func labelWarningValue(label, value string) string {
	return termstyle.LabelWarning(label, value)
}
func commandLine(invocation string, accentCommand bool) string {
	return termstyle.CommandLine(invocation, accentCommand)
}

const (
	groupSetup = "setup"
	groupWork  = "work"
	groupHelp  = "help"
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

	commandBlock := `{{commandSections .}}`

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
	fmt.Fprintln(out, rootInferenceStatusLine())
	fmt.Fprintln(out)
}

func rootAuthStatusLine() string {
	creds, err := cloud.LoadCloudCredentials()
	if err == nil && authKindForCredentials(creds) != "none" && strings.TrimSpace(creds.Login) != "" {
		return success("●") + " " + value("Signed in as "+strings.TrimSpace(creds.Login))
	}
	return danger("●") + " " + muted("Not signed in") + "  " + logoText("gx auth login")
}

func rootInferenceStatusLine() string {
	creds, ok := inference.Resolve()
	if !ok {
		return danger("●") + " " + muted("API Key required. Run `gx set key` to set an API Key.")
	}
	return success("●") + " " + value("Using "+creds.Provider+": "+maskedAPIKey(creds.APIKey))
}

func maskedAPIKey(key string) string {
	key = strings.TrimSpace(key)
	if len(key) < 2 {
		return "..."
	}
	visible := min(7, len(key)-1)
	return key[:visible] + "..."
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
		"commandSections":         commandSections,
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

func commandDisplayPadded(cmd *cobra.Command, commands []*cobra.Command) string {
	return logoText(padRight(commandDisplay(cmd), commandDisplayPadding(commands)))
}

func commandSections(cmd *cobra.Command) string {
	commands := availableHelpCommands(cmd.Commands())
	if len(commands) == 0 {
		return ""
	}
	grouped := groupedHelpCommands(cmd, commands)
	if len(grouped) == 0 {
		return renderCommandSection("Available Commands:", commands)
	}
	var out bytes.Buffer
	for _, group := range cmd.Groups() {
		sectionCommands := orderedHelpCommands(grouped[group.ID], group.ID)
		if len(sectionCommands) == 0 {
			continue
		}
		renderCommandSectionTo(&out, group.Title, sectionCommands)
	}
	var ungrouped []*cobra.Command
	for _, sub := range commands {
		if sub.Name() == "help" && cmd.ContainsGroup(groupHelp) {
			continue
		}
		if strings.TrimSpace(sub.GroupID) == "" {
			ungrouped = append(ungrouped, sub)
		}
	}
	if len(ungrouped) > 0 {
		title := "Other:"
		if len(ungrouped) == 1 && ungrouped[0].Name() == "help" {
			title = "Help:"
		}
		renderCommandSectionTo(&out, title, ungrouped)
	}
	return strings.TrimRight(out.String(), "\n")
}

func orderedHelpCommands(commands []*cobra.Command, groupID string) []*cobra.Command {
	order := rootHelpCommandOrder(groupID)
	if len(order) == 0 || len(commands) < 2 {
		return commands
	}
	byName := make(map[string]*cobra.Command, len(commands))
	for _, cmd := range commands {
		byName[cmd.Name()] = cmd
	}
	var ordered []*cobra.Command
	seen := make(map[string]bool, len(commands))
	for _, name := range order {
		cmd := byName[name]
		if cmd == nil {
			continue
		}
		ordered = append(ordered, cmd)
		seen[name] = true
	}
	for _, cmd := range commands {
		if seen[cmd.Name()] {
			continue
		}
		ordered = append(ordered, cmd)
	}
	return ordered
}

func rootHelpCommandOrder(groupID string) []string {
	switch groupID {
	case groupSetup:
		return []string{"init", "auth"}
	case groupWork:
		return []string{"review"}
	case groupHelp:
		return []string{"doctor", "version", "help"}
	default:
		return nil
	}
}

func availableHelpCommands(commands []*cobra.Command) []*cobra.Command {
	var available []*cobra.Command
	for _, cmd := range commands {
		if !cmd.IsAvailableCommand() && cmd.Name() != "help" {
			continue
		}
		available = append(available, cmd)
	}
	return available
}

func groupedHelpCommands(cmd *cobra.Command, commands []*cobra.Command) map[string][]*cobra.Command {
	if len(cmd.Groups()) == 0 {
		return nil
	}
	grouped := make(map[string][]*cobra.Command)
	for _, sub := range commands {
		groupID := strings.TrimSpace(sub.GroupID)
		if groupID == "" && sub.Name() == "help" && cmd.ContainsGroup(groupHelp) {
			groupID = groupHelp
		}
		if groupID == "" {
			continue
		}
		grouped[groupID] = append(grouped[groupID], sub)
	}
	if len(grouped) == 0 {
		return nil
	}
	return grouped
}

func renderCommandSection(title string, commands []*cobra.Command) string {
	var out bytes.Buffer
	renderCommandSectionTo(&out, title, commands)
	return strings.TrimRight(out.String(), "\n")
}

func renderCommandSectionTo(out io.Writer, title string, commands []*cobra.Command) {
	if len(commands) == 0 {
		return
	}
	fmt.Fprintln(out, section(title))
	for _, cmd := range commands {
		fmt.Fprintf(out, "  %s %s\n", commandDisplayPadded(cmd, commands), helpDescription(commandShortDescription(cmd)))
	}
	fmt.Fprintln(out)
}

func commandShortDescription(cmd *cobra.Command) string {
	if cmd.Name() == "help" {
		return "Help screen"
	}
	return cmd.Short
}

func helpDescription(text string) string {
	return muted(text)
}

func commandDisplay(cmd *cobra.Command) string {
	name := cmd.Name()
	if len(cmd.Aliases) == 0 {
		return name
	}
	return name + " (" + strings.Join(cmd.Aliases, ", ") + ")"
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
