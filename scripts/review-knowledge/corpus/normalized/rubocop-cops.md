# RuboCop

Role models are important.

## Overview

**RuboCop** is a Ruby static code analyzer (a.k.a. linter) and code
formatter. Out of the box it will enforce many of the guidelines
outlined in the community Ruby Style Guide.

RuboCop packs a lot of features on top of what you’d normally expect from a linter:

-
Works with every major Ruby implementation
-
Autocorrection of many of the code offenses it detects
-
Robust code formatting capabilities
-
Multiple result formatters for both interactive use and for feeding data into other tools
-
Ability to have different configuration for different parts of your codebase
-
Ability to disable certain cops only for specific files or parts of files
-
Extremely flexible configuration that allows you to adapt RuboCop to pretty much every style and preference
-
It’s easy to extend RuboCop with custom cops and formatters
-
A vast number of ready-made plugins (e.g. `rubocop-rails`,`rubocop-rspec`,`rubocop-performance`and`rubocop-minitest`)
-
Built-in Language Server Protocol (LSP) and Model Context Protocol (MCP) support
-
Wide editor/IDE support via LSP
-
Many online services use RuboCop internally (e.g. CodeClimate and Codacy)
-
Best logo/stickers ever

The project is closely tied to several efforts to document and promote the best practices of the Ruby community:

A long-term goal of RuboCop (and its core extensions) is to cover with cops all the guidelines from the community style guides.

## Philosophy

Early on RuboCop aimed to be an opinionated linter/formatter that adhered very closely to the Ruby Style Guide (think `gofmt` and the like).
In those days cops supported just a single style and you couldn’t even turn individual cops off. Eventually, we realized
that in the Ruby community there were so many competing styles and preferences that it was going to be really
challenging to find one set of defaults that makes everyone happy. Part of this was Ruby’s own culture and philosophy,
part was the lack of common standards for almost 20 years. It’s hard to undo any of those, but it’s also not really necessary.

The early feedback we got led us to adopt a philosophy of (extreme) configurability and flexibility, and trying to account for every *common* style
of programming in Ruby. While we still believe that there’s a lot of merit to just sticking to the community
style guides, we acknowledge that Ruby is all about diversity and doing things the way that makes you happy. Whatever
style preferences you have RuboCop is there for you. That’s our promise and our guarantee. Within the subjective limits of sanity that is.

## Next Steps

So, what to do next? While you can peruse the documentation in whatever way you’d like, here are a few recommendations:

-
Start with the "Getting Started" guide to get up and running quickly.
-
See "CLI Reference" for the full command-line reference.
-
Adjust RuboCop to your style/preferences. RuboCop is an extremely flexible tool and most aspects of its behavior can be tweaked via various configuration options. See "Configuration" for more details.
-
See "Versioning" for information about RuboCop versioning, updates, and the process of introducing new cops.
-
Explore the available plugins.

# RuboCop

Role models are important.

## Overview

**RuboCop** is a Ruby static code analyzer (a.k.a. linter) and code
formatter. Out of the box it will enforce many of the guidelines
outlined in the community Ruby Style Guide.

RuboCop packs a lot of features on top of what you’d normally expect from a linter:

-
Works with every major Ruby implementation
-
Autocorrection of many of the code offenses it detects
-
Robust code formatting capabilities
-
Multiple result formatters for both interactive use and for feeding data into other tools
-
Ability to have different configuration for different parts of your codebase
-
Ability to disable certain cops only for specific files or parts of files
-
Extremely flexible configuration that allows you to adapt RuboCop to pretty much every style and preference
-
It’s easy to extend RuboCop with custom cops and formatters
-
A vast number of ready-made plugins (e.g. `rubocop-rails`,`rubocop-rspec`,`rubocop-performance`and`rubocop-minitest`)
-
Built-in Language Server Protocol (LSP) and Model Context Protocol (MCP) support
-
Wide editor/IDE support via LSP
-
Many online services use RuboCop internally (e.g. CodeClimate and Codacy)
-
Best logo/stickers ever

The project is closely tied to several efforts to document and promote the best practices of the Ruby community:

A long-term goal of RuboCop (and its core extensions) is to cover with cops all the guidelines from the community style guides.

## Philosophy

Early on RuboCop aimed to be an opinionated linter/formatter that adhered very closely to the Ruby Style Guide (think `gofmt` and the like).
In those days cops supported just a single style and you couldn’t even turn individual cops off. Eventually, we realized
that in the Ruby community there were so many competing styles and preferences that it was going to be really
challenging to find one set of defaults that makes everyone happy. Part of this was Ruby’s own culture and philosophy,
part was the lack of common standards for almost 20 years. It’s hard to undo any of those, but it’s also not really necessary.

The early feedback we got led us to adopt a philosophy of (extreme) configurability and flexibility, and trying to account for every *common* style
of programming in Ruby. While we still believe that there’s a lot of merit to just sticking to the community
style guides, we acknowledge that Ruby is all about diversity and doing things the way that makes you happy. Whatever
style preferences you have RuboCop is there for you. That’s our promise and our guarantee. Within the subjective limits of sanity that is.

## Next Steps

So, what to do next? While you can peruse the documentation in whatever way you’d like, here are a few recommendations:

-
Start with the "Getting Started" guide to get up and running quickly.
-
See "CLI Reference" for the full command-line reference.
-
Adjust RuboCop to your style/preferences. RuboCop is an extremely flexible tool and most aspects of its behavior can be tweaked via various configuration options. See "Configuration" for more details.
-
See "Versioning" for information about RuboCop versioning, updates, and the process of introducing new cops.
-
Explore the available plugins.

# Installation

RuboCop’s installation is pretty standard:

`$ gem install rubocop`If you’d rather install RuboCop using `bundler`, don’t require it in your `Gemfile`:

`gem 'rubocop', require: false`RuboCop is stable between major versions, both in terms of API and cop configuration.
We aim to ease the maintenance of RuboCop extensions and the upgrades between RuboCop
releases. All big changes are reserved for major releases.
To prevent an unwanted RuboCop update you might want to use a conservative version lock
in your `Gemfile`:

`gem 'rubocop', '~> 1.89', require: false`

# Compatibility

RuboCop targets Ruby 2.0+ code analysis.[1]

RuboCop officially supports MRI (a.k.a. CRuby) and JRuby at runtime.

-
MRI 2.7+
-
JRuby 9.4+

The oldest supported JRuby version is derived from the oldest compatible MRI version.

| RuboCop might be working with other Ruby implementations as well, but it’s tested only on MRI and JRuby. |

## Support Matrix

RuboCop generally aims to follow MRI’s own support policy - meaning RuboCop would support all officially supported MRI releases.[2] To give people extra time for a smooth transition, we’ve customarily provided support for about one year after EOL of an MRI version.[3]

The following table is the runtime support matrix.

| Supported runtime Ruby version | Last supported RuboCop version |
|---|---|
| 1.9 | 0.41 |
| 2.0 | 0.50 |
| 2.1 | 0.57 |
| 2.2 | 0.68 |
| 2.3 | 0.81 |
| 2.4 | 1.12 |
| 2.5 | 1.28 |
| 2.6 | 1.50 |
| 2.7 | - |
| 3.0 | - |
| 3.1 | - |
| 3.2 | - |
| 3.3 | - |
| 3.4 | - |
| 4.0 | - |
| 4.1 (experimental) | - |

| The table above is about runtimesupport (which Ruby can run RuboCop itself).Code analysissupport is broader — RuboCop can analyze code targeting Ruby 2.0+ regardless of the Ruby version it runs on. See`TargetRubyVersion`for details. |

## Parser engines

RuboCop supports two parser backends, configured via
`ParserEngine`:

Since RuboCop 1.75, `parser_prism` is used by default when `TargetRubyVersion` is 3.4 or higher.

# Getting Started

This guide will get you up and running with RuboCop in a few minutes. RuboCop serves three primary roles:

-
**Code style checker**(a.k.a. linter) — enforces style conventions from the Ruby Style Guide
-
**Lint tool**— catches bugs and suspicious code, like a smarter`ruby -w`
-
**Code formatter**— automatically fixes layout and formatting

## First Run

After installing RuboCop, just run it from your project’s root directory:

`$ rubocop`RuboCop will recursively check all Ruby files and report any offenses it finds:

```
Inspecting 5 files
.W.C.
Offenses:
lib/foo.rb:2:3: C: Style/IfUnlessModifier: Favor modifier if usage when having a single-line body.
 if something
 ^^
lib/bar.rb:5:5: W: Lint/UselessAssignment: Useless assignment to variable - x.
 x = 42
 ^
5 files inspected, 2 offenses detected
```
Each offense shows the file, line, column, severity (`C` for convention, `W` for warning, `E` for error), the cop name, and a description. The letter on the progress line (`.` = clean, `C`/`W`/`E` = offense found) gives you a quick overview.

## Setting Up Your Project

For an existing project with many offenses, the easiest way to get started is to generate a TODO file:

`$ rubocop --auto-gen-config`This creates two files:

-
`.rubocop_todo.yml`— temporarily disables all current offenses
-
`.rubocop.yml`— your configuration file (with`inherit_from: .rubocop_todo.yml`added)

Now `rubocop` will pass cleanly, and you can work through the TODO entries at your own pace. See Auto-generating Configuration for details.

For a new project, you can generate a starter config file instead:

`$ rubocop --init`## Customizing RuboCop

All configuration lives in `.rubocop.yml`. Here are a few common adjustments:

```
# Increase the maximum line length
Layout/LineLength:
 Max: 120
# Disable a cop entirely
Style/Documentation:
 Enabled: false
# Restrict a cop to specific files
Rails/HasAndBelongsToMany:
 Include:
 - app/models/**/*.rb
```
See Configuration for the full reference.

## Fixing Offenses Automatically

Many cops can autocorrect the problems they find. Use `-a` for safe corrections only, or `-A` to include unsafe ones:

```
# Safe autocorrections only
$ rubocop -a
# All autocorrections (safe and unsafe)
$ rubocop -A
# Fix only formatting (layout) offenses
$ rubocop -x
```
| Review the diff after running autocorrect, especially with `-A`. See Autocorrect for more details on safe vs. unsafe corrections. |

## Editor Integration

For the best experience, set up RuboCop in your editor so you get real-time feedback as you type. RuboCop has a built-in LSP server that works with any editor that supports the Language Server Protocol.

The recommended setup for popular editors:

See Integration with Other Tools for more editors and alternative approaches.

## Next Steps

-
CLI Reference — full command-line reference with all flags
-
Configuration — everything about `.rubocop.yml`
-
Source Code Directives — disable/enable cops inline with comments
-
Cops — browse all available cops by department
-
Plugins — add cops for Rails, RSpec, Performance, and more

# CLI Reference

Full reference for the `rubocop` command-line interface. For a tutorial introduction, see Getting Started.

## Command-line flags

For more details check the available command-line options:

`$ rubocop -h`To specify multiple cops for one flag, separate cops with commas and no spaces:

`$ rubocop --only Rails/Blank,Layout/HeredocIndentation,Naming/FileName`| Command flag | Description |
|---|---|
|
 | Autocorrect offenses (only when it’s safe). See Autocorrect. |
|
 | Deprecated alias of |
|
 | Autocorrect offenses (safe and unsafe). See Autocorrect. |
|
 | Deprecated alias of |
|
 | Generate a configuration file acting as a TODO list. |
|
 | Force color output on or off. |
|
 | Run with specified config file. |
|
 | Store and reuse results for faster operation. |
|
 | Displays some extra debug output. |
|
 | Run with all cops disabled by default, except |
|
 | Run without pending cops. |
|
 | Used with --autocorrect to annotate any offenses that do not support autocorrect with |
|
 | Displays cop names in offense messages. Default is true. |
|
 | Display elapsed time in seconds. |
|
 | Only output offense messages at the specified |
|
 | Only output correctable offense messages. |
|
 | Only output safe correctable offense messages. |
|
 | Run with all cops enabled, including those disabled by default. Overrides |
|
 | Run with pending cops. |
|
 | Run all cops enabled by configuration except the specified cop(s) and/or departments. |
|
 | Limit how many individual files |
|
 | Displays extra details in offense messages. |
|
 | Choose a formatter, see Formatters. |
|
 | Inspect files in order of modification time and stop after the first file with offenses. |
|
 | Minimum severity for exit with error code. Full severity name or upper case initial can be given. Normally, autocorrected offenses are ignored. Use |
|
 | Force excluding files specified in the configuration |
|
 | Inspect files given on the command line only if they are listed in |
|
 | Print usage information. |
|
 | Report offenses even if they have been manually disabled with a |
|
 | Ignores all Exclude: settings from all .rubocop.yml files present in parent folders. This is useful when you are importing submodules when you want to test them without being affected by the parent module’s rubocop settings. |
|
 | Ignore unrecognized cops or departments in the config. |
|
 | Generate a .rubocop.yml file in the current directory. |
|
 | Run only lint cops. |
|
 | List all files RuboCop will inspect. |
|
 | Generate only |
|
 | Include the date and time when |
|
 | Show offense counts in config file generated by |
|
 | Run only the specified cop(s) and/or cops in the specified departments. |
|
 | Write output to a file instead of STDOUT. |
|
 | Use available CPUs to execute inspection in parallel. Default is parallel. |
|
 | Raise cop-related errors with cause and location. This is used to prevent cops from failing silently. Default is false. |
|
 | Plugin RuboCop extension file (see Loading Plugin Extensions). |
|
 | Require Ruby file (see Loading Legacy Extensions). |
|
 | Regenerate the TODO list using the same options as the last time it was generated with |
|
 | Run only safe cops. |
|
 | Deprecated alias of |
|
 | Shows available cops and their configuration. You can use |
|
 | Shows URLs for documentation pages of supplied cops. |
|
 | Write all output to stderr except for the autocorrected source. This is especially useful when combined with |
|
 | Pipe source from STDIN. This is useful for editor integration. Takes one argument, a path, relative to the root of the project. RuboCop will use this path to determine which cops are enabled (via eg. Include/Exclude), and so that certain cops like Naming/FileName can be checked. |
|
 | Optimize real-time feedback in editors, adjusting behaviors for editing experience. Editors that run RuboCop directly (e.g., by shelling out) encounter the same issues as with |
|
 | Display style guide URLs in offense messages. |
|
 | Autocorrect only code layout (formatting) offenses. |
|
 | Displays the current version and exits. |
|
 | Displays the current version plus the version of Parser and Ruby. |

Default command-line options are loaded from `.rubocop` and `RUBOCOP_OPTS` and are combined with command-line options that are explicitly passed to `rubocop`.
Thus, the options have the following order of precedence (from highest to lowest):

-
Explicit command-line options
-
Options from `RUBOCOP_OPTS`environment variable
-
Options from `.rubocop`file.

## Exit codes

RuboCop exits with the following status codes:

-
`0`if no offenses are found or if the severity of all offenses are less than`--fail-level`. (By default, if you use`--autocorrect`, offenses which are autocorrected do not cause RuboCop to fail.)
-
`1`if one or more offenses equal or greater to`--fail-level`are found. (By default, this is any offense which is not autocorrected.)
-
`2`if RuboCop terminates abnormally due to invalid configuration, invalid CLI options, or an internal error.

## Parallel Processing

RuboCop enables the `--parallel` option by default. This option allows for
parallel processing using the number of available CPUs based on `Etc.nprocessors`.
If the default number of parallel processes causes excessive consumption of CPU and memory,
you can adjust the number of parallel processes using the `$PARALLEL_PROCESSOR_COUNT` environment variable.
Here is how to set the parallel processing to use two CPUs:

`$ PARALLEL_PROCESSOR_COUNT=2 bundle exec rubocop`Alternatively, specifying `--no-parallel` will disable parallel processing.

# Auto-generating Configuration

If you have a code base with an overwhelming amount of offenses, it can
be a good idea to use `rubocop --auto-gen-config`, which creates
`.rubocop_todo.yml` and adds `inherit_from: .rubocop_todo.yml` in your
`.rubocop.yml`. The generated file `.rubocop_todo.yml` contains
configuration to disable cops that currently detect an offense in the
code by changing the configuration for the cop, excluding the offending
files, or disabling the cop altogether once a file count limit has been
reached.

By adding the option `--exclude-limit COUNT`, e.g., ```
rubocop
--auto-gen-config --exclude-limit 5
```
, you can change how many files are
excluded before the cop is entirely disabled. The default COUNT is 15.
If you don’t want the cop to be entirely disabled regardless of the
number of files, use the `--no-exclude-limit` option, e.g.,
`rubocop --auto-gen-config --no-exclude-limit`.

## Working through the TODO

The next step is to cut and paste configuration from `.rubocop_todo.yml`
into `.rubocop.yml` for everything that you think is in line with your
(organization’s) code style and not a good fit for a todo list.

| Pay attention to the comments above each entry in `.rubocop_todo.yml`.
They can reveal configuration parameters such as`EnforcedStyle`, which can
be used to modify the behavior of a cop instead of disabling it completely. |

Then you can start removing the entries in the generated
`.rubocop_todo.yml` file one by one as you work through all the offenses
in the code. You can also regenerate your `.rubocop_todo.yml` using
the same options by running `rubocop --regenerate-todo`.

Another way of silencing offense reports, aside from configuration, is through source code directives. These can be added manually or automatically.

## Metrics cops

The cops in the `Metrics` department will by default get `Max` parameters
generated in `.rubocop_todo.yml`. The value of these will be just high enough
so that no offenses are reported the next time you run `rubocop`. If you
prefer to exclude files, like for other cops, add `--auto-gen-only-exclude`
when running with `--auto-gen-config`. It will still change the maximum if the
number of excluded files is higher than the exclude limit.

## EnforcedStyle

Some cops have a configurable option named `EnforcedStyle`.
By default, when generating the `.rubocop_todo.yml`, if one style is used
for all files, these cops will add the settings for the style being used.
If you want to exclude on a file-by-file basis,
add the `--no-auto-gen-enforced-style` option along with `--auto-gen-config`.

# Autocorrect

In autocorrect mode, RuboCop will try to automatically fix offenses.

There are three autocorrect modes:

| Flag | Description |
|---|---|
|
 | Safe corrections only — won’t change code semantics. |
|
 | All corrections, including unsafe ones that may change semantics. |
|
 | Layout (formatting) corrections only — never changes logic. |

| Always review the diff and run your test suite after autocorrecting,
especially when using `-A`. |

## Safe vs. unsafe

RuboCop distinguishes between safe and unsafe cops and autocorrections:

-
**Safe**(`true/false`) — indicates whether the cop can yield false positives by design.
-
**SafeAutoCorrect**(`true/false`) — indicates whether the autocorrection preserves code semantics. If a cop itself is unsafe, its autocorrect is automatically considered unsafe as well.

When you run `rubocop -a`, cops and autocorrections marked as unsafe are
skipped. When you run `rubocop -A`, everything is applied.

| Some cops may not yet be annotated for safety. Eventually, the safety of each cop will be specified in the default configuration. |

### Example: unsafe cop

```
class Miner
 def dig(how_deep)
 # ...
 end
end
Miner.new.dig(42) # => Style/SingleArgumentDig
 # => Use Miner.new[] instead of dig
```
This is the wrong diagnostic; this (contrived) use of `dig` is not an issue,
and there might not be an alternative. This cop is marked as `Safe: false`.

### Example: safe cop with unsafe autocorrect

```
# example.rb:
str = 'hello' # => Missing magic comment `# frozen_string_literal: true`
str << 'world'
# autocorrects to:
# frozen_string_literal: true
str = 'hello'
str << 'world' # => now fails because `str` is frozen
# must be manually corrected to:
# frozen_string_literal: true
str = +'hello' # => We want an unfrozen string literal here...
str << 'world' # => ok
```
This diagnostic is valid since the magic comment is indeed missing (thus `Safe: true`),
but the autocorrection is not; some string literals need to be prefixed with `+` to avoid
having them frozen.

## Disabling uncorrectable offenses

`$ rubocop --autocorrect --disable-uncorrectable`The `--disable-uncorrectable` flag generates `# rubocop:todo` comments in the
code to suppress reporting of offenses that could not be corrected
automatically. This is useful for gradually adopting new cops.

## Configuring autocorrect per cop

Individual cops can have their autocorrect behavior configured in `.rubocop.yml`.
See AutoCorrect for the
available settings: `always`, `contextual`, and `disabled`.

# Caching

Large projects containing hundreds or even thousands of files can take a really long time to inspect, but RuboCop has functionality to mitigate this problem. There’s a caching mechanism that stores information about offenses found in inspected files.

## Cache Validity

Later runs will be able to retrieve this information and present the stored information instead of inspecting the file again. This will be done if the cache for the file is still valid, which it is if there are no changes in:

-
the contents of the inspected file
-
RuboCop configuration for the file
-
the options given to `rubocop`, with some exceptions that have no bearing on which offenses are reported
-
the Ruby version used to invoke `rubocop`
-
version of the `rubocop`program (or to be precise, anything in the source code of the invoked`rubocop`program)

In rare cases, you may need to invalidate the cache when changing
external dependencies. This happens when your cop depends on external
configuration. In this case, override the method
`RuboCop::Cop::Base#external_dependency_checksum`.

## Enabling and Disabling the Cache

The caching functionality is enabled if the configuration parameter
`AllCops: UseCache` is `true`, which it is by default. The command
line option `--cache false` can be used to turn off caching, thus
overriding the configuration parameter. If `AllCops: UseCache` is set
to `false` in the local `.rubocop.yml`, then it’s `--cache true` that
overrides the setting.

## Cache Path

By default, the cache is stored in either
`$XDG_CACHE_HOME/$UID/rubocop_cache` if `$XDG_CACHE_HOME` is set or in
`$HOME/.cache/rubocop_cache/` if it’s not.

The root can be set to a different path in a number of ways (from
**highest** precedence to **lowest**):

-
the `--cache-root`command line option
-
the `$RUBOCOP_CACHE_ROOT`environment variable
-
the `AllCops: CacheRootDirectory`configuration parameter

One reason to set the cache root could be that there’s a network disk where users on different machines want to have a common RuboCop cache. Another could be that a Continuous Integration system allows directories, but not a temporary directory, to be saved between runs, or that the system caches certain folders by default.

## Cache Pruning

Each time a file has changed, its offenses will be stored under a new
key in the cache. This means that the cache will continue to grow
until we do something to stop it. The configuration parameter
`AllCops: MaxFilesInCache` sets a limit, and when the number of files
in the cache exceeds that limit, the oldest files will be automatically
removed from the cache.

To disable cache pruning entirely, set `MaxFilesInCache` to `false`:

```
AllCops:
 MaxFilesInCache: false
```
This can be useful in CI environments where cache cleanup is wasteful, as the cache directory is typically discarded after each build anyway.

# LSP (Language Server Protocol)

| The built-in language server was introduced in RuboCop 1.53. |

The Language Server Protocol is the modern standard for providing cross-editor support for various programming languages.

This feature enables extremely fast interactions through the LSP.

Offense detection and autocorrection are performed in real-time by editors and IDEs using the language server. The Server Mode is primarily used to speed up RuboCop runs in the terminal. Therefore, if you want real-time feedback from RuboCop in your editor or IDE, opting to use this language server instead of the server mode will not only provide a fast and efficient solution, but also offer a straightforward setup for integration.

## Examples of LSP Clients

Here are examples of LSP client configurations.

### VS Code

vscode-rubocop integrates RuboCop into VS Code.

You can install this VS Code extension from the Visual Studio Marketplace.

For VS Code-based IDEs like VSCodium or Eclipse Theia, the extension can be installed from the Open VSX Registry.

### Emacs (Eglot)

Eglot is a client for Language Server Protocol servers on Emacs.

Add the following to your Emacs configuration file (e.g. `~/.emacs.d/init.el`):

```
(require 'eglot)
(add-to-list 'eglot-server-programs '(ruby-mode . ("bundle" "exec" "rubocop" "--lsp")))
(add-hook 'ruby-mode-hook 'eglot-ensure)
```
Below is an example of additional setting for autocorrecting on save:

`(add-hook 'ruby-mode-hook (lambda () (add-hook 'before-save-hook 'eglot-format-buffer nil 'local)))`If you run into problems, first use "M-x eglot-reconnect" to reconnect to the language server.

See Eglot’s official documentation for more information.

### Emacs (LSP Mode)

LSP Mode is an Emacs client/library for the Language Server Protocol.

You can get the new `lsp-mode` package from MELPA.

See LSP Mode official documentation for more information: https://emacs-lsp.github.io/lsp-mode/page/lsp-rubocop/

### Vim and Neovim (coc.nvim)

coc.nvim is an extension host for Vim and Neovim, powered by Node.js. It allows the loading of extensions similar to VSCode and provides hosting for language servers.

Add the following to your coc.nvim configuration file. For example, in Vim, it would be `~/.vim/coc-settings.json`,
and in Neovim, it would be `~/.config/nvim/coc-settings.json`:

```
{
 "languageserver": {
 "rubocop": {
 "command": "bundle",
 "args" : ["exec", "rubocop", "--lsp"],
 "filetypes": ["ruby"],
 "rootPatterns": [".git", "Gemfile"],
 "requireRootPattern": true
 }
 }
}
```
Below is an example of additional setting for autocorrecting on save:

```
{
 "coc.preferences.formatOnSave": true
}
```
See coc.nvim’s official documentation for more information.

### Neovim (nvim-lspconfig)

nvim-lspconfig provides quickstart configs for Neovim’s LSP.

Add the following to your nvim-lspconfig configuration file (e.g. `~/.config/nvim/init.lua`):

```
vim.opt.signcolumn = "yes"
vim.api.nvim_create_autocmd("FileType", {
 pattern = "ruby",
 callback = function()
 vim.lsp.start {
 name = "rubocop",
 cmd = { "bundle", "exec", "rubocop", "--lsp" },
 }
 end,
})
```
Below is an example of additional setting for autocorrecting on save:

```
vim.api.nvim_create_autocmd("BufWritePre", {
 pattern = "*.rb",
 callback = function()
 vim.lsp.buf.format()
 end,
})
```
See nvim-lspconfig’s official documentation for more information.

### Helix

Helix is a post-modern modal text editor with built-in language server support.

Add the following to your Helix language configuration file (e.g. `~/.config/helix/languages.toml`):

Helix 23.10 or later:

```
[language-server.rubocop]
command = "bundle"
args = ["exec", "rubocop", "--lsp"]
[[language]]
name = "ruby"
auto-format = true
language-servers = [
 { name = "rubocop" }
]
```
Helix 23.10 or earlier:

```
[[language]]
name = "ruby"
language-server = { command = "bundle", args = ["exec", "rubocop", "--lsp"] }
auto-format = true
```
See Helix’s official documentation for more information: https://docs.helix-editor.com/languages.html

### Sublime Text

For Sublime Text LSP support is available through the Sublime-LSP plugin.
Add the following to its settings (accessible via `Preferences → Package Settings → LSP → Settings`) to enable RuboCop:

```
{
 "clients": {
 "rubocop": {
 "enabled": true,
 "command": ["bundle", "exec", "rubocop", "--lsp"],
 "selector": "source.ruby | text.html.ruby | text.html.rails",
 },
 },
}
```
## Autocorrection

The language server supports the `textDocument/formatting` method and is autocorrectable. The autocorrection is safe by default (`rubocop -a`).

An LSP client can switch to unsafe autocorrection (`rubocop -A`) by passing the following `safeAutocorrect` parameter in the `initialize` request.

```
{
 "jsonrpc": "2.0",
 "id": 42,
 "method": "initialize",
 "params": {
 "initializationOptions": {
 "safeAutocorrect": false
 }
 }
}
```
For detailed instructions on setting the parameter, please refer to the configuration methods of your LSP client.

| The `safeAutocorrect`parameter was introduced in RuboCop 1.54. |

As execute commands in the `workspace/executeCommand` parameters, it provides `rubocop.formatAutocorrects` for safe autocorrections (`rubocop -a`) and
`rubocop.formatAutocorrectsAll` for unsafe autocorrections (`rubocop -A`).
These parameters take precedence over the `initializationOptions:safeAutocorrect` value set in the `initialize` parameter.

| The `rubocop.formatAutocorrectsAll`execute command was introduced in RuboCop 1.56. |

## Lint Mode

An LSP client can run lint cops by passing the following `lintMode` parameter in the `initialize` request
if you only want to enable the feature as a linter like `ruby -w`:

```
{
 "jsonrpc": "2.0",
 "id": 42,
 "method": "initialize",
 "params": {
 "initializationOptions": {
 "lintMode": true
 }
 }
}
```
Furthermore, enabling autocorrect in an LSP client at the time of saving is equivalent to the `rubocop -l` option.

For detailed instructions on setting the parameter, please refer to the configuration methods of your LSP client.

| The `lintMode`parameter was introduced in RuboCop 1.55. |

## Layout Mode

An LSP client can run layout cops by passing the following `layoutMode` parameter in the `initialize` request
if you only want to enable the feature as a formatter:

```
{
 "jsonrpc": "2.0",
 "id": 42,
 "method": "initialize",
 "params": {
 "initializationOptions": {
 "layoutMode": true
 }
 }
}
```
Furthermore, enabling autocorrect in an LSP client at the time of saving is equivalent to the `rubocop -x` option.

For detailed instructions on setting the parameter, please refer to the configuration methods of your LSP client.

| The `layoutMode`parameter was introduced in RuboCop 1.55. |

## Enable YJIT

YJIT, a Ruby JIT compiler, has been supported since Ruby 3.1.
In an LSP client, you can enable YJIT by launching `rubocop --lsp` with the `RUBY_YJIT_ENABLE=1` environment variable using the `env` command:

`env RUBY_YJIT_ENABLE=1 bundle exec rubocop --lsp`Below is an example for Emacs’s Eglot:

`(add-to-list 'eglot-server-programs '(ruby-mode . ("env" "RUBY_YJIT_ENABLE=1" "bundle" "exec" "rubocop" "--lsp")))`The console of the LSP client will display `+YJIT`:

`RuboCop 1.63.4 language server +YJIT initialized, PID 13501`For more details, please refer to the respective LSP configuration documentation. In some cases, like with vscode-rubocop, it may be available as a built-in option: https://github.com/rubocop/vscode-rubocop#rubocopyjitenabled

## Run as a Language Server

Run `rubocop --lsp` command from LSP client.

When the language server is started, the command displays the language server’s PID:

```
$ ps aux | grep 'rubocop --lsp'
user 17414 0.0 0.2 5557716 144376 ?? Ss 4:48PM 0:02.13 rubocop --lsp /Users/user/src/github.com/rubocop/rubocop
```
| `rubocop --lsp`starts the LSP server and is intended to be invoked by an LSP client, not run manually by users. |

## Language Server Development

RuboCop provides APIs for developers of custom language servers or tools analogous to LSP, using RuboCop as the backend instead of RuboCop’s built-in LSP.

-
`RuboCop::LSP.enable`enables LSP mode, customizing for LSP-specific features such as autocorrection and short offense messages.
-
`RuboCop::LSP.disable`disables LSP mode, which can be particularly useful for testing. Intentional autocorrection by the user can be specified via e.g.`workspace/executeCommand`and`textDocument/codeAction`LSP methods.

When implementing custom cops, `RuboCop::LSP.enabled?` can be used to achieve behavior that considers these states.

# MCP (Model Context Protocol)

| The built-in MCP server was introduced in RuboCop 1.85. |

| This feature is experimental and should not be considered stable. Changes to its behavior or interface may occur. |

Model Context Protocol is an open-source standard for connecting AI applications to external systems.

This feature enables interactions through MCP clients that communicate with RuboCop via the Model Context Protocol.

Through MCP, RuboCop operations are initiated by an MCP client, with decisions made by an LLM.

The MCP server runs as a long-lived process over stdio, allowing an MCP client to request analysis and autocorrection without spawning a new RuboCop process for each request. This is based on the same principles as Server Mode and LSP.

## Setup

The `mcp` gem is required but is not included as a dependency of RuboCop. You must install it separately.

If you use Bundler, add the following to your `Gemfile`:

`gem 'mcp', '~> 0.6'`Then run:

`$ bundle install`If you do not use Bundler:

`$ gem install mcp`## Available Tools

The MCP server exposes two tools:

- `rubocop_inspection`
-
Inspect Ruby code for offenses. Accepts a `path`to check files on disk or`source_code`to check inline code. Returns detected offenses as JSON.
- `rubocop_autocorrection`
-
Autocorrect RuboCop offenses in Ruby code. Accepts `path`or`source_code`like the inspection tool, plus a`safety`boolean (defaults to`true`). When`safety`is`true`, only safe corrections are applied; set it to`false`to include unsafe corrections.

## Client Configuration

For a list of tools that support MCP, see MCP Clients.

### JSON Configuration

Many MCP clients (e.g. VS Code, Cursor, Windsurf) accept a JSON configuration file. The exact file path varies by client — consult your client’s documentation.

```
{
 "mcpServers": {
 "rubocop": {
 "type": "stdio",
 "command": "bundle",
 "args": [
 "exec",
 "rubocop",
 "--mcp"
 ],
 "cwd": "/path/to/your/project"
 }
 }
}
```
| `rubocop --mcp`starts the MCP server and is intended to be invoked by an MCP client, not run manually by users. |

# Server Mode

| The server mode was introduced in RuboCop 1.31. If you’re using an older RuboCop version you can check out the rubocop-daemon project that served as the inspiration for RuboCop’s built-in functionality. |

You can reduce the RuboCop boot time significantly (something like 850x faster) by using the `--server` command-line option.

The `--server` option speeds up the launch of the `rubocop` command by utilizing
a standalone server process that loads the RuboCop runtime production files (i.e. `require 'rubocop'`).

Normally RuboCop starts somewhat slowly because it needs to `require` a ton of files and that’s fairly
slow. With the RuboCop server we sidestep this nasty issue and make it much more pleasant to
interact with RuboCop from text editors and IDEs.

| The feature cannot be used on JRuby and Windows, as they do not support the `fork`system call. |

## Run with Server

There are two ways to enable server:

-
`rubocop --server`: If server process has not started yet, start server process and execute inspection with server.
-
`rubocop --start-server`: Just start server process.

When the server is started, it outputs the host and port.

```
$ rubocop --start-server
RuboCop server starting on 127.0.0.1:55772.
```
| The `rubocop`command is executed using the server process if a server is started.
Whenever a server process is not running, it will load the RuboCop runtime files and execute.
(same behavior as with RuboCop 1.30 and lower) |

If a server is already running, the command only displays the server’s PID. A new server will not be started.

```
$ rubocop --start-server
RuboCop server (16060) is already running.
```
The server process name is basically `rubocop --server` and the project directory path:

```
$ ps aux | grep 'rubocop --server'
user 16060 0.0 0.0 5078568 2264 ?? S 7:54AM 0:00.00 rubocop --server /Users/user/src/github.com/rubocop/rubocop
user 16337 0.0 0.0 5331560 2396 ?? S 23:51PM 0:00.00 rubocop --server /Users/user/src/github.com/rubocop/rubocop-rails
```
When you run `bundle update` or update a local config file (e.g., `.rubocop.yml` or `.rubocop_todo.yml`), and then run `rubocop`, the server process will automatically restart.

```
$ rubocop --server
RuboCop version incompatibility found, RuboCop server restarting...
RuboCop server starting on 127.0.0.1:60665.
```
| Detection of incompatibility changes in the local configuration also includes changes to local file paths specified by `inherit_from`and`require`in`.rubocop.yml`.
Changes involving remote files or those considered to be searched on`$LOAD_PATH`are not detected. |

If you would like to start the server in the foreground, which can be useful when running within Docker, you can pass the `--no-detach` option.

`$ rubocop --start-server --no-detach`## Restart Server

The started server does not reload the configuration file. You will need to restart the server when you upgrade RuboCop or change the RuboCop configuration.

```
$ rubocop --restart-server
RuboCop server starting on 127.0.0.1:55822.
```
## Command Line Options

These are the command-line options for server operations:

| Command flag | Description |
|---|---|
|
 | If a server process has not been started yet, start the server process and execute inspection with server. |
|
 | If a server process has been started, stop the server process and execute inspection without the server. |
|
 | Restart server process. |
|
 | Start server process. |
|
 | Stop server process. |
|
 | Show server status. |
|
 | Run the server process in the foreground. |

| You can specify the server host and port with the $RUBOCOP_SERVER_HOST and the $RUBOCOP_SERVER_PORT environment variables. |

If `RUBOCOP_OPTS` environment variable or `.rubocop` file contains `--server` option, `rubocop` command defaults to server mode.
Other server options such as `stop-server`, `restart-server` specified on the command line will take precedence over them.

# Source Code Directives

RuboCop can be controlled within your source code using special comments.
These directives let you disable or enable cops for specific sections of a file,
without changing your `.rubocop.yml` configuration.

## Disabling Cops within Source Code

One or more individual cops can be disabled locally in a section of a file by adding a comment such as

```
# rubocop:disable Layout/LineLength, Style/StringLiterals
[...]
# rubocop:enable Layout/LineLength, Style/StringLiterals
```
You can also disable entire departments by giving a department name in the comment.

```
# rubocop:disable Metrics, Layout/LineLength
[...]
# rubocop:enable Metrics, Layout/LineLength
```
You can also disable *all* cops with

```
# rubocop:disable all
[...]
# rubocop:enable all
```
In cases where you want to differentiate intentionally-disabled cops vs. cops
you’d like to revisit later, you can use `rubocop:todo` as an alias of
`rubocop:disable`.

```
# rubocop:todo Layout/LineLength, Style/StringLiterals
[...]
# rubocop:enable Layout/LineLength, Style/StringLiterals
```
One or more cops can be disabled on a single line with an end-of-line comment.

`for x in (0..19) # rubocop:disable Style/For`If you want to disable a cop that inspects comments, you can do so by adding an "inner comment" on the comment line.

`# coding: utf-8 # rubocop:disable Style/Encoding`Running `rubocop --autocorrect --disable-uncorrectable` will
create comments to disable all offenses that can’t be automatically
corrected.

You can add a comment to the disabling/enabling directive by prefixing it with `--`. For example:

`# rubocop:disable Layout/LineLength -- A comment explaining why the cop is disabled`The syntax of directives can be checked using the cop `Lint/CopDirectiveSyntax`.

## Temporarily Enabling Cops in Source Code

In a similar way to disabling cops within source code, you can also temporarily enable specific cops if you want to enforce specific rules for part of a file.

Let’s use the cop `Style/AsciiComments`, which is by default `Enabled: false`. Say you want a
specific file to have ASCII-only comments to be compatible with some post-processing tool:

```
# rubocop:enable Style/AsciiComments
# If applicable, leave a comment to others explaining the rationale:
# We need the comments to remain ASCII only for compatibility with lib/post_processor.rb
class Restaurant
 # This comment has to be ASCII-only because of the rubocop:enable directive
 def menu
 return dishes.map(&:humanize)
 end
end
```
You can also enforce the same for part of a file by disabling the cop afterwards

```
class Dish
 def humanize
 return [
 "Delicious #{self.name}"
 *ingredients
 ].join("\n")
 end
end
# rubocop:enable Style/AsciiComments
# If applicable, leave a comment to others explaining the rationale:
# We need the comments to remain ASCII only for compatibility with lib/post_processor.rb
class Restaurant
 # This comment has to be ASCII-only because of the rubocop:enable directive
 def menu
 return dishes.map(&:humanize)
 end
end
# rubocop:disable Style/AsciiComments
class Ingredient
 # Notice how the comment below is non-ASCII
 # Gets rid of odd characters like 😀,
 def sanitize
 self.name.gsub(/[^a-z]/, '')
 end
end
```
| If a file is excluded via configuration (e.g., in `.rubocop.yml`or`.rubocop_todo.yml`),`rubocop:enable`comments within that file will have no effect. Configuration-based exclusions take
precedence over in-source opt-in directives. |

## Scoped Disabling with Push/Pop Directives

When you want to temporarily change cop settings for a specific section of code
and then automatically restore the previous state, you can use `rubocop:push` and
`rubocop:pop` directives. This is particularly useful when you need to disable
cops for a block of code without affecting the rest of the file.

### Basic Push/Pop Usage

The `push` directive saves the current state of all cop settings, and `pop`
restores them:

```
def process_data(input)
 result = input.upcase
 # rubocop:push
 # rubocop:disable Style/GuardClause
 if result.present?
 return result.strip
 end
 # rubocop:pop
 nil
end
```
After `pop`, the `Style/GuardClause` cop is automatically re-enabled, returning
to its state before `push`.

### Inline Push Arguments

For convenience, you can combine `push` with enable/disable operations using
inline arguments. Use `-` to disable a cop and `+` to enable a cop:

```
def process_data(input)
 result = input.upcase
 # rubocop:push -Style/GuardClause
 if result.present?
 return result.strip
 end
 # rubocop:pop
 nil
end
```
You can specify multiple cops with different operations:

```
# rubocop:disable Style/For
for x in [1, 2, 3]
 puts x
end
# rubocop:push +Style/For -Style/GuardClause
for y in [4, 5, 6] # Style/For is re-enabled here
 if y > 0
 return y # Style/GuardClause is disabled here
 end
end
# rubocop:pop
# Back to original state: Style/For disabled, Style/GuardClause enabled
```
### Nested Push/Pop

Push/pop directives can be nested for complex scenarios:

```
# rubocop:disable Metrics/MethodLength
def complex_method
 step1
 # rubocop:push
 # rubocop:enable Metrics/MethodLength
 # rubocop:disable Style/GuardClause
 def helper_method
 # rubocop:push
 # rubocop:enable Style/GuardClause
 if condition
 return value
 end
 # rubocop:pop
 other_code
 end
 # rubocop:pop
 step2
end
```
Each `pop` restores the state to what it was at the corresponding `push`.

# Project Index

| The project index feature was introduced in RuboCop 1.87. |

| This feature is experimental and should not be considered stable. Changes to its behavior or interface may occur. |

RuboCop can optionally use Rubydex to build a project-wide index of declarations and references. When enabled, cops that opt in can consult the index to detect issues that span multiple files.

This integration is **opt-in and experimental**. The default behavior of RuboCop is unchanged.

## Enabling

-
Add `rubydex`to your`Gemfile`and`bundle install`:`gem 'rubydex', require: false`
-
Set the flag in your `.rubocop.yml`:`AllCops: UseProjectIndex: true`

If `UseProjectIndex` is `true` but the `rubydex` gem is not installed, or the running Ruby
is older than the version `rubydex` supports, RuboCop prints a warning and falls back to
its standard file-local behavior.

The integration requires Ruby 3.2 or later. On Ruby 3.1 and older, `AllCops/UseProjectIndex`
has no effect even if set to `true`.

### Indexing gem sources

By default the index covers only the project’s own files, so ancestry chains and members
that live in gems are unresolvable, and index-aware cops fall back to their conservative
behavior whenever one is involved (e.g. a class inheriting from a framework base class).
Setting `AllCops/ProjectIndexIncludesGems: true` additionally indexes the sources of every
gem in the project’s bundle:

```
AllCops:
 UseProjectIndex: true
 ProjectIndexIncludesGems: true
```
This makes ancestry-based reasoning conclusive for most real-world classes (on RuboCop’s own repository it raises the share of classes with a fully resolvable ancestry from about 20% to about 97%) at the cost of extra memory proportional to the bundle’s size and a slightly longer index build. The option requires RuboCop to run inside Bundler; outside a bundle it silently degrades to project-only indexing.

## What it enables

`Lint/ConstantReassignment` reports reassignments whose previous definition lives in another file.
For example:

```
# a.rb
CROSS_FILE_CONST = :first
# b.rb
CROSS_FILE_CONST = :second
```
With `UseProjectIndex: true`, RuboCop reports a `Lint/ConstantReassignment` offense referencing the other file.

`Lint/DuplicateMethods` reports methods whose duplicate definition lives in another file.
For example:

```
# a.rb
class Foo
 def bar; end
end
# b.rb
class Foo
 def bar; end
end
```
With `UseProjectIndex: true`, RuboCop reports a `Lint/DuplicateMethods` offense in each file,
referencing the definition in the other one. Redefining a method from another file on purpose
(e.g. a monkey patch) can be signaled with the self-alias trick (`alias bar bar` right before
the redefinition), which also suppresses Ruby’s method redefinition warning.

`Lint/DeprecatedReference` (pending) is entirely powered by the index: it reports calls to methods
and references to constants documented with a YARD `@deprecated` tag anywhere in the project.

`Lint/UnusedPrivateMethod` (disabled by default) uses the index for project-wide dead-code
detection: it reports private instance methods whose names are never referenced anywhere in
the indexed project. Since symbol-based references from other files (e.g. Rails callbacks
declared in a concern) cannot be detected, it is best suited for occasional dead-code sweeps.

`Style/MissingRespondToMissing` accepts a `respond_to_missing?` defined in another definition
of the same class or module (e.g. a reopening in another file).

`Lint/InheritException` also reports classes that inherit from `Exception` indirectly, through
a parent class defined elsewhere in the project.

`Style/StaticClass` does not report classes that are subclassed anywhere in the project.

`Naming/PredicatePrefix` and `Naming/AccessorMethodName` do not suggest renaming methods
that override a method defined by an ancestor elsewhere in the project.

`Style/ClassAndModuleChildren` uses the index to make its (unsafe) autocorrection more reliable:
when nesting a compact definition (`class Foo::Bar`), the namespace wrapper’s keyword (`class` or
`module`) is resolved from the actual definition of `Foo` instead of guessed, and compacting a
nested definition is skipped when the outer namespace is not defined elsewhere (the compact form
would raise `NameError` at load time).

`Lint/MissingSuper` skips the constructor offense when the class' entire ancestry is resolvable
in the index and no ancestor defines `initialize` - in that case `super` would only reach the
no-op `Object#initialize`. This removes the need to list project-local abstract base classes in
`AllowedParentClasses`.

`Style/RedundantConstantBase` also reports a leading `::` inside a namespace when the
constant provably resolves identically without it.

`Lint/ConstantResolution` reports only genuinely ambiguous constants - those that resolve to a
different declaration through the surrounding nesting than they would fully qualified - which
makes the cop practical to enable without `Only`/`Ignore` lists.

`Style/Documentation` accepts a reopened class or module when any of its other definition
sites carries a documentation comment, instead of requiring a comment at every reopening.

`Lint/NameTypo` (pending) is entirely powered by the index: it reports qualified constant
references and constant-receiver method calls that do not resolve anywhere in the project
when a similarly named alternative exists, and suggests it. `CheckConstants` and
`CheckMethods` toggle the two checks independently.

Index-aware cops automatically pick up the index whenever it is built; no per-cop opt-in is required.

## Notes

-
`rubydex`requires Ruby 3.2 or newer and ships native (Rust) extensions.
-
Parallel inspection is currently disabled on Windows when `UseProjectIndex`is on; use serial inspection (omit`--parallel`).
-
The index always covers the whole project (rooted at the directory containing `Gemfile`or`gems.rb`, falling back to the current directory), regardless of which files a particular run inspects. Inspecting a single file therefore reports the same cross-file offenses as a full run.
-
The index is rebuilt once per `rubocop`invocation; no on-disk index is shared between runs. Indexing is fast (roughly 100ms for a couple thousand files), but very large monorepos will notice the per-run cost on single-file runs.

# Profiling

When RuboCop feels slow, profiling helps you find the bottleneck — whether it’s a particular cop, a large file, or excessive memory allocation.

| Profiling requires MRI (CRuby). It is not available on JRuby or Windows. |

## CPU profiling with `--profile`

The `--profile` flag uses the stackprof gem
to record wall-time samples while RuboCop runs. Add it to your `Gemfile` first:

`gem 'stackprof', require: false`Then run:

```
$ rubocop --profile
# ...
Profile report generated at tmp/rubocop-stackprof.dump
```
The dump file is written to `tmp/rubocop-stackprof.dump` in your project root.
You can inspect it with the `stackprof` CLI:

```
# Top methods by wall time
$ stackprof tmp/rubocop-stackprof.dump --limit 20
# Call tree for a specific method
$ stackprof tmp/rubocop-stackprof.dump --method 'RuboCop::Cop::Style::MethodCallWithArgsParentheses#on_send'
# Flamegraph (open in a browser)
$ stackprof --flamegraph tmp/rubocop-stackprof.dump > tmp/flamegraph.json
$ stackprof --flamegraph-viewer tmp/flamegraph.json
```
| To profile a single cop, combine `--profile`with`--only`:`rubocop --profile --only Style/StringLiterals`. |

## Memory profiling with `--memory`

The `--memory` flag (used together with `--profile`) additionally tracks memory
allocations using the
memory_profiler gem. Add both
gems to your `Gemfile`:

```
gem 'stackprof', require: false
gem 'memory_profiler', require: false
```
Then run:

```
$ rubocop --profile --memory
# ...
Profile report generated at tmp/rubocop-stackprof.dump
Building memory report...
Memory report generated at tmp/rubocop-memory_profiler.txt
```
The memory report is a human-readable text file showing allocations grouped by gem, file, location, and class.

| Memory profiling is expensive. On a large codebase it can take a very long time and consume significant memory. Limit the scope to a subset of files or cops: |

```
$ rubocop --profile --memory app/models/
$ rubocop --profile --memory --only Layout/LineLength lib/
```
## Tips

-
**Caching is disabled automatically**when profiling, so results reflect actual cop execution rather than cache hits.
-
**Narrow the scope**— profiling the entire codebase produces noisy results. Focus on a single directory or cop to get actionable data.
-
**Compare before and after**— if you’re optimizing a cop, profile it in isolation (`--only`) before and after your change to measure the improvement.

## Reporting performance problems

If you discover a performance bottleneck in RuboCop, please open an issue or pull request with the profiling data. See Contributing for details.

# Configuration

The behavior of RuboCop can be controlled via the .rubocop.yml configuration file. It makes it possible to enable/disable certain cops (checks) and to alter their behavior if they accept any parameters. The file can be placed in your home directory, XDG config directory, or in some project directory.

The file has the following format:

```
inherit_from: ../.rubocop.yml
Style/Encoding:
 Enabled: false
Layout/LineLength:
 Max: 99
```
This page covers general configuration topics. For more specific areas, see:

-
Inheritance — `inherit_from`,`inherit_gem`, and`inherit_mode`
-
Including and Excluding Files — controlling which files RuboCop inspects
-
Configuring Cops — `Enabled`,`Severity`,`AutoCorrect`,`AllowedMethods`,`AllowedPatterns`

## Config file locations

RuboCop will start looking for the configuration file in the directory where the inspected file is and continue its way up to the root directory.

If it cannot be found until reaching the project’s root directory, then it will be searched for in the .config directory of the project root and the user’s global config locations. The user’s global config locations consist of a dotfile or a config file inside the XDG Base Directory specification.

-
`.config/.rubocop.yml`or`.config/rubocop/config.yml`at the project root
-
`~/.rubocop.yml`
-
`$XDG_CONFIG_HOME/rubocop/config.yml`(expands to`~/.config/rubocop/config.yml`if`$XDG_CONFIG_HOME`is not set)

If both files exist, the dotfile will be selected.

As an example, if RuboCop is invoked from inside `/path/to/project/lib/utils`,
then RuboCop will use the config as specified inside the first of the following
files:

-
`/path/to/project/lib/utils/.rubocop.yml`
-
`/path/to/project/lib/.rubocop.yml`
-
`/path/to/project/.rubocop.yml`
-
`/path/to/project/.config/.rubocop.yml`
-
`/path/to/project/.config/rubocop/config.yml`
-
`~/.rubocop.yml`
-
`~/.config/rubocop/config.yml`

| All the previous logic does not apply if a specific configuration file is passed
on the command line through the `--config`flag. In that case, the resolved
configuration file will be the one passed to the CLI. |

## Pre-processing

Configuration files are pre-processed using the ERB templating mechanism. This makes it possible to add dynamic content that will be evaluated when the configuration file is read. For example, you could let RuboCop ignore all files ignored by Git.

```
AllCops:
 Exclude:
 <% `git status --ignored --porcelain`.lines.grep(/^!! /).each do |path| %>
 - <%= path.sub(/^!! /, '').sub(/\/$/, '/**/*') %>
 <% end %>
```
## Defaults

The file config/default.yml under the RuboCop home directory contains the
default settings that all configurations inherit from. Project and personal
`.rubocop.yml` files need only make settings that are different from the
default ones. If there is no `.rubocop.yml` file in the project, home or XDG
directories, `config/default.yml` will be used.

## Setting the target Ruby version

Some checks are dependent on the version of the Ruby interpreter which the
inspected code must run on. For example, enforcing using Ruby 2.6+ endless
ranges `foo[n..]` rather than `foo[n..-1]` can help make your code shorter and
more consistent… *unless* it must run on e.g. Ruby 2.5.

Users may let RuboCop know the oldest version of Ruby which your project supports with:

```
AllCops:
 TargetRubyVersion: 2.5
```
If a `TargetRubyVersion` is not specified in your config, then RuboCop will
check your project for a series of other files where the Ruby version may be
specified already. The files that will be checked are (in this order):
`*.gemspec`, `.ruby-version`, `mise.toml`, `.tool-versions`, and `Gemfile.lock`.

The target ruby version may also be specified by setting the
`RUBOCOP_TARGET_RUBY_VERSION` environment variable to the desired version: for
example, running `RUBOCOP_TARGET_RUBY_VERSION=3.3 rubocop` will
run rubocop with a target ruby version of 3.3. Using this environment variable
will override all other sources of version information, including
`.rubocop.yml`.

If a target Ruby version cannot be found via any of the above sources, then a default target Ruby version will be used.

### Finding target Ruby in a `*.gemspec` file

In order for RuboCop to parse a `*.gemspec` file’s `required_ruby_version`, the
Ruby version must be specified using one of these syntaxes:

-
a string range, e.g. `'~> 3.2.0'`or`'>= 3.2.2'`
-
an array of strings, e.g. `['>= 3.0.0', '< 3.4.0']`
-
a `Gem::Requirement`, e.g.`Gem::Requirement.new('>= 3.1.2')`

If a `*.gemspec` file specifies a range of supported Ruby versions via any of
these means, then the greater of the following Ruby versions will be used:

-
the lowest Ruby version that is compatible with your specified range
-
the lowest version of Ruby that is still supported by your version of RuboCop

If a `*.gemspec` file defines its `required_ruby_version` dynamically (e.g. by
reading from a `.ruby-version` file, via an environment variable, referencing a
constant or local variable, etc), then RuboCop will *not* detect that Ruby
version, and will instead try to find a target Ruby version elsewhere.

## Setting the parser engine

| The parser engine configuration was introduced in RuboCop 1.62. Since RuboCop 1.75, RuboCop chooses the parser engine automatically, so you don’t need to configure it yourself. |

RuboCop allows switching the backend parser by specifying either
`parser_whitequark` or `parser_prism` as the value for the `ParserEngine`.

Here are the parsers used as backends for each value:

-
`ParserEngine: default`
-
`ParserEngine: parser_whitequark`… https://github.com/whitequark/parser
-
`ParserEngine: parser_prism`… https://github.com/ruby/prism (`Prism::Translation::Parser`)

`parser_whitequark` can analyze source code from Ruby 2.0 until Ruby 3.4:

```
AllCops:
 ParserEngine: parser_whitequark
```
`parser_prism` can analyze source code from Ruby 3.3 and above:

```
AllCops:
 ParserEngine: parser_prism
 TargetRubyVersion: 3.3
```
| `parser_prism`tends to perform analysis faster than`parser_whitequark`. |

## Setting the style guide URL

You can specify the base URL of the style guide using `StyleGuideBaseURL`.
If specified under `AllCops`, all cops are targeted.

```
AllCops:
 StyleGuideBaseURL: https://rubystyle.guide
```
`StyleGuideBaseURL` is combined with `StyleGuide` specified to the cop.

```
Lint/UselessAssignment:
 StyleGuide: '#underscore-unused-vars'
```
The style guide URL is https://rubystyle.guide#underscore-unused-vars.

If specified under a specific department, it takes precedence over `AllCops`.
The following is an example of specifying `Rails` department.

```
Rails:
 StyleGuideBaseURL: https://rails.rubystyle.guide
```
```
Rails/TimeZone:
 StyleGuide: '#time'
```
The style guide URL is https://rails.rubystyle.guide#time.

## Setting the documentation URL

You can specify the base URL of the documentation using `DocumentationBaseURL`.
If specified under `AllCops`, all cops are targeted.

```
AllCops:
 DocumentationBaseURL: https://docs.rubocop.org/rubocop
```
If specified under a specific department, it takes precedence over `AllCops`.
The following is an example of specifying `Rails` department.

```
Rails:
 DocumentationBaseURL: https://docs.rubocop.org/rubocop-rails
```
By default, documentation is expected to be served as HTML but if you prefer
to use something else like markdown you can set `DocumentationExtension`.

With markdown as the documentation format you are able to host it directly through
GitHub without having to own a domain or using GitHub Pages. The `rubocop-sorbet`
extension is an example of this, its docs are available
here.

```
Sorbet:
 DocumentationBaseURL: https://github.com/Shopify/rubocop-sorbet/blob/main/manual
 DocumentationExtension: .md
```
## Setting the version tracking metadata for cops

This configuration is particularly useful when custom cops are distributed as a gem.

Each cop can have the following additional metadata:

-
`VersionAdded`- the RuboCop version in which it was added
-
`VersionChanged`(optional) - the latest RuboCop version in which it was changed in a user-impacting way (new config, updated defaults, etc)

```
Style/HashSyntax:
 VersionAdded: '0.9'
 VersionChanged: '1.67'
```
| These values do not include patch versions. |

Those will be pretty useful for the documentation (so the manual generation has to be enhanced to include them) and keeping track of changes.

## Enable checking Active Support extensions

Some cops for checking specified methods (e.g. `Style/HashExcept`) support Active Support extensions.
This is off by default, but can be enabled by the `ActiveSupportExtensionsEnabled` option.

```
AllCops:
 ActiveSupportExtensionsEnabled: true
```
## Opting into globally frozen string literals

Ruby continues to move into the direction of having all string literals frozen by default.
Ruby 3.4, for example, will show a warning if a non-frozen string literal from a file without
the frozen string literal magic comment gets modified. By starting Ruby with the environment
variable `RUBYOPT` set to `--enable=frozen-string-literal` you can opt into that behaviour today.
For RuboCop to provide accurate analysis you must also configure the `StringLiteralsFrozenByDefault`
option.

```
AllCops:
 StringLiteralsFrozenByDefault: true
```

# Inheritance

All configuration inherits from RuboCop’s default configuration (see "Defaults").

RuboCop also supports inheritance in user’s configuration files. The most common
example would be the `.rubocop_todo.yml` file (see Auto-generating Configuration).

Settings in the child file (that which inherits) override those in the parent (that which is inherited), with the following caveats.

## Inheritance of hashes vs. other types

Configuration parameters that are hashes, for example `PreferredMethods` in
`Style/CollectionMethods`, are merged with the same parameter in the parent
configuration. This means that any key-value pairs given in child configuration
override the same keys in parent configuration. Giving `~`, YAML’s
representation of `nil`, as a value cancels the setting of the corresponding
key in the parent configuration. For example:

```
Style/CollectionMethods:
 Enabled: true
 PreferredMethods:
 # No preference for collect, keep all others from default config.
 collect: ~
```
Other types, such as `AllCops` / `Include` (an array), are overridden by the
child setting.

Arrays override because if they were merged, there would be no way to remove elements in child files.

However, advanced users can still merge arrays using the `inherit_mode` setting.
See Merging arrays using inherit_mode below.

## Inheriting from another configuration file in the project

The optional `inherit_from` directive is used to include configuration
from one or more files. This makes it possible to have the common
project settings in the `.rubocop.yml` file at the project root, and
then only the deviations from those rules in the subdirectories. The
files can be given with absolute paths or paths relative to the file
where they are referenced. The settings after an `inherit_from`
directive override any settings in the file(s) inherited from. When
multiple files are included, the first file in the list has the lowest
precedence and the last one has the highest. The format for multiple
inheritance is:

```
inherit_from:
 - ../.rubocop.yml
 - ../conf/.rubocop.yml
```
`inherit_from` also accepts a glob, for example:

```
inherit_from:
 - packages/*/.rubocop_todo.yml
```
The example above is one potential use-case: allowing components within your repo to organize their own `.rubocop_todo.yml` files.

## Inheriting configuration from a remote URL

The optional `inherit_from` directive can contain a full URL to a remote
file. This makes it possible to have common project settings stored on an HTTP
server and shared between many projects.

The remote config file is cached locally and is only updated if:

-
The file does not exist.
-
The file has not been updated in the last 24 hours.
-
The remote copy has a newer modification time than the local copy.

This local cache is stored in the configured cache path, maintaining the remote filename prefix.

| The root of relative paths specified in `Include`or`Exclude`changes
depending on the configuration file name, as explained in
"Path relativity".
This behavior was introduced in RuboCop 1.84. |

You can inherit from both remote and local files in the same config and the same inheritance rules apply to remote URLs and inheriting from local files where the first file in the list has the lowest precedence and the last one has the highest. The format for multiple inheritance using URLs is:

```
inherit_from:
 - http://www.example.com/rubocop.yml
 - ../.rubocop.yml
```
You can inherit from a repo with basic auth that is authorized to access the repo as follows:

```
inherit_from:
 - http://<user_name>:<password>@raw.githubusercontent.com/example/rubocop.yml
```
A GitHub personal access token can also be configured as follows:

```
inherit_from:
 - http://<personal_access_token>@raw.githubusercontent.com/example/rubocop.yml
```
## Inheriting configuration from a dependency gem

The optional `inherit_gem` directive is used to include configuration from
one or more gems external to the current project. This makes it possible to
inherit a shared dependency’s RuboCop configuration that can be used from
multiple disparate projects.

Configurations inherited in this way will be essentially *prepended* to the
`inherit_from` directive, such that the `inherit_gem` configurations will be
loaded first, then the `inherit_from` relative file paths will be loaded
(overriding the configurations from the gems), and finally the remaining
directives in the configuration file will supersede any of the inherited
configurations. This means the configurations inherited from one or more gems
have the lowest precedence of inheritance.

The directive should be formatted as a YAML Hash using the gem name as the key and the relative path within the gem as the value:

```
inherit_gem:
 my-shared-gem: .rubocop.yml
 cucumber: conf/rubocop.yml
```
An array can also be used as the value to include multiple configuration files from a single gem:

```
inherit_gem:
 my-shared-gem:
 - default.yml
 - strict.yml
```
| If the shared dependency is declared using a Bundler
Gemfile and the gem was installed using `bundle install`, it would be
necessary to also invoke RuboCop using Bundler in order to find the
dependency’s installation path at runtime: |

`$ bundle exec rubocop <options...>`## Merging arrays using inherit_mode

The optional directive `inherit_mode` specifies which configuration keys that
have array values should be merged together instead of overriding the inherited
value.

This applies to explicit inheritance using `inherit_from` as well as implicit
inheritance from the default configuration.

Given the following config:

```
# .rubocop.yml
inherit_from:
 - shared.yml
inherit_mode:
 merge:
 - Exclude
AllCops:
 Exclude:
 - 'generated/**/*.rb'
Style/For:
 Exclude:
 - bar.rb
```
```
# .shared.yml
Style/For:
 Exclude:
 - foo.rb
```
The list of `Exclude`s for the `Style/For` cop in this example will be
`['foo.rb', 'bar.rb']`. Similarly, the `AllCops:Exclude` list will contain all
the default patterns plus the `generated/**/*.rb` entry that was added locally.

The directive can also be used on individual cop configurations to override the global setting.

```
inherit_from:
 - shared.yml
inherit_mode:
 merge:
 - Exclude
Style/For:
 inherit_mode:
 override:
 - Exclude
 Exclude:
 - bar.rb
```
In this example the `Exclude` would only include `bar.rb`.

# Including and Excluding Files

RuboCop does a recursive file search starting from the directory it is
run in, or directories given as command line arguments. Files that
match any pattern listed under `AllCops`/`Include` and extensionless
files with a hash-bang (`#!`) declaration containing one of the known
ruby interpreters listed under `AllCops`/`RubyInterpreters` are
inspected, unless the file also matches a pattern in
`AllCops`/`Exclude`. Hidden directories (i.e., directories whose names
start with a dot) are not searched by default.

Here is an example that might be used for a Rails project:

```
AllCops:
 Exclude:
 - 'db/**/*'
 - 'config/**/*'
 - 'script/**/*'
 - 'bin/{rails,rake}'
 - !ruby/regexp /old_and_unused\.rb$/
# other configuration
# ...
```
| When inspecting a certain directory (or file)
given as RuboCop’s command line arguments,
patterns listed under `AllCops`/`Exclude`are also inspected.
If you want to apply`AllCops`/`Exclude`rules in this circumstance,
add`--force-exclusion`to the command line argument. |

Here is an example:

```
# .rubocop.yml
AllCops:
 Exclude:
 - foo.rb
```
If `foo.rb` is specified as a RuboCop’s command line argument, the result is:

```
# RuboCop inspects foo.rb.
$ bundle exec rubocop foo.rb
# RuboCop does not inspect foo.rb.
$ bundle exec rubocop --force-exclusion foo.rb
```
## Path relativity

In `.rubocop.yml` and any other configuration file beginning with `.rubocop`,
files, and directories are specified relative to the directory where the
configuration file is. In configuration files that don’t begin with `.rubocop`,
e.g. `our_company_defaults.yml`, paths are relative to the directory where
`rubocop` is run.

This affects cops that have customisable paths: if the default is `db/migrate/*.rb`,
and the cop is enabled in `db/migrate/.rubocop.yml`, the path will need to be
explicitly set as `*.rb`, as the default will look for `db/migrate/db/migrate/*.rb`.
This is unlikely to be what you wanted.

## Unusual files, that would not be included by default

RuboCop comes with a comprehensive list of common ruby file names and
extensions. But, if you’d like RuboCop to check files that are not included by
default, you’ll need to pass them in on the command line, or add entries for
them under `AllCops`/`Include`.

| Your configuration files override
RuboCop’s defaults.
In the following example, we want to include `foo.unusual_extension`, but we also
must copy any other patterns we need from the overridden`default.yml`. |

```
AllCops:
 Include:
 - foo.unusual_extension
 - '**/*.rb'
 - '**/*.gemfile'
 - '**/*.gemspec'
 - '**/*.rake'
 - '**/*.ru'
 - '**/Gemfile'
 - '**/Rakefile'
```
This behavior of `Include` (overriding `default.yml`) was introduced in
0.56.0
via #5882. This change allows
people to include/exclude precisely what they need to, without the defaults
getting in the way.

## Deprecated patterns

Patterns that are just a file name, e.g. `Rakefile`, will match
that file name in any directory, but this pattern style is deprecated. The
correct way to match the file in any directory, including the current, is
`**/Rakefile`.

The pattern `config/**` will match any file recursively under
`config`, but this pattern style is deprecated and should be replaced by
`config/**/*`.

`Include` and `Exclude` are relative to their directory

The `Include` and `Exclude` parameters are special. They are
valid for the directory tree starting where they are defined. They are not
shadowed by the setting of `Include` and `Exclude` in other `.rubocop.yml`
files in subdirectories. This is different from all other parameters, which
follow RuboCop’s general principle that configuration for an inspected file
is taken from the nearest `.rubocop.yml`, searching upwards.

| This behavior
will be overridden if you specify the `--ignore-parent-exclusion`command line
argument. |

## Cop-specific `Include` and `Exclude`

Cops can be run only on specific sets of files when that’s needed (for
instance you might want to run some Rails model checks only on files whose
paths match `app/models/*.rb`). All cops support the
`Include` param.

```
Rails/HasAndBelongsToMany:
 Include:
 - app/models/*.rb
```
Cops can also exclude only specific sets of files when that’s needed (for
instance you might want to run some cop only on a specific file). All cops support the
`Exclude` param.

```
Rails/HasAndBelongsToMany:
 Exclude:
 - app/models/problematic.rb
```

# Configuring Cops

Every cop supports a set of common configuration parameters, and many cops have additional parameters specific to their behavior.

| Qualifying cop name with its type, e.g., `Style`, is recommended,
but not necessary as long as the cop name is unique across all types. |

## Enabled

Specific cops can be disabled by setting `Enabled` to `false` for that specific cop.

```
Layout/LineLength:
 Enabled: false
```
Most cops are enabled by default. Cops, introduced or significantly updated
between major versions, are in a special pending status (read more in
"Versioning"). Some cops, configured with `Enabled: false`
in config/default.yml,
are disabled by default.

The cop enabling process can be altered by setting `DisabledByDefault` or
`EnabledByDefault` (but not both) to `true`. These settings override the default for **all**
cops to disabled or enabled, except `Lint/Syntax` which is always enabled,
regardless of the cops' default values (whether enabled, disabled or pending).

```
AllCops:
 DisabledByDefault: true
```
All cops except `Lint/Syntax` are then disabled by default. Only cops appearing in user
configuration files with `Enabled: true` will be enabled; every other cop will
be disabled without having to explicitly disable them in configuration. It is
also possible to enable entire departments by adding for example

```
Style:
 Enabled: true
```
All cops in the `Style` department are then enabled. In this case, only the cops
in the `Style` department that are enabled by default will be enabled.
The cops in the `Style` department that are disabled by default will remain disabled.

The same effects can be obtained from the command line using `--disable-all-cops`
and `--enable-all-cops`. These options take precedence over the configuration
files, which is useful when you want to inspect the source with a different cop
set without modifying `.rubocop.yml`. Note that `--disable-all-cops` does not
disable `Lint/Syntax`, which is always enabled.

If a department is disabled, cops in that department can still be individually enabled, and that setting overrides the setting for its department in the same configuration file and in any inherited file.

```
inherit_from: config_that_disables_the_metrics_department.yml
Metrics/MethodLength:
 Enabled: true
Style:
 Enabled: false
Style/Alias:
 Enabled: true
```
## Severity

Each cop has a default severity level based on which department it belongs
to. The level is normally `warning` for `Lint` and `convention` for all the
others, but this can be changed in user configuration. Cops can customize their
severity level. Allowed values are `info`, `refactor`, `convention`, `warning`, `error`
and `fatal`.

Cops with severity `info` will be reported but will not cause `rubocop` to return
a non-zero value.

There is one exception from the general rule above and that is `Lint/Syntax`, a
special cop that checks for syntax errors before the other cops are invoked. It
cannot be disabled and its severity (`fatal`) cannot be changed in
configuration.

```
Lint:
 Severity: error
Metrics/CyclomaticComplexity:
 Severity: warning
```
## Details

Individual cops can be embellished with extra details in offense messages:

```
Layout/LineLength:
 Details: >-
 If lines are too short, text becomes hard to read because you must
 constantly jump from one line to the next while reading. If lines are too
 long, the line jumping becomes too hard because you "lose the line" while
 going back to the start of the next line. 80 characters is a good
 compromise.
```
These details will only be seen when RuboCop is run with the `--extra-details` flag or if `ExtraDetails` is set to true in your global RuboCop configuration.

## AutoCorrect

Cops that support the `--autocorrect` option offer flexible settings for autocorrection.
These settings can be specified in the configuration file as follows:

-
`always`
-
`contextual`
-
`disabled`

`always (Default)`

This setting enables autocorrection always by default. For backward compatibility, `true` is treated the same as `always`.

```
Style/PerlBackrefs:
 AutoCorrect: always # or true
```
`contextual`

This setting enables autocorrection when launched from the `rubocop` command, but it is not available through LSP.
e.g., `rubocop --lsp`, `rubocop --editor-mode`, or a program where `RuboCop::LSP.enable` has been applied.

Inspections via the command line are treated as code that has been finalized.

```
Style/PerlBackrefs:
 AutoCorrect: contextual
```
This setting prevents autocorrection during editing in the editor, e.g., with the `textDocument/formatting` LSP method.
However, the `workspace/executeCommand` LSP method, which is triggered by intentional user actions, respects the user’s intention for autocorrection.

Additionally, for cases like `Metrics` cops where the highlight range extends over the entire body of classes, modules, methods, or blocks,
the offending range will be confined to only the name. This helps avoid redundant and noisy offenses in editor display.

## AllowedMethods

Many cops can be configured to exclude specific methods from inspection.
`AllowedMethods` accepts a list of method names as strings. A cop will skip
any method whose name exactly matches one of the entries:

```
Metrics/BlockLength:
 AllowedMethods:
 - refine
 - class_methods
 - instance_methods
```
Only bare method names are supported — you cannot use qualified names like `Foo.bar`
or `foo.bar`. Class and module names are not supported either. If you need
more flexibility, use `AllowedPatterns` instead.

| The `IgnoredMethods`and`ExcludedMethods`parameters are deprecated aliases for`AllowedMethods`. They still work, but you should migrate your configuration to use`AllowedMethods`. |

## AllowedPatterns

`AllowedPatterns` provides regex-based exclusions for cops that support it.
Each entry is treated as a regular expression — strings are automatically
converted via `Regexp.new`, so the `!ruby/regexp` YAML tag is optional.

```
Metrics/BlockLength:
 AllowedPatterns:
 # These two forms are equivalent:
 - !ruby/regexp /\b(class|instance)_methods\b/
 - '\b(class|instance)_methods\b'
```
What the pattern matches against depends on the cop:

-
**Metrics cops**(e.g.`BlockLength`,`MethodLength`) match against the**method name**.
-
**Layout/LineLength**matches against the**entire source line**.
-
**Naming cops**(e.g.`MethodName`,`VariableName`) match against the**identifier name**.

For example, to allow all methods starting with `test_` in `Metrics/MethodLength`:

```
Metrics/MethodLength:
 AllowedPatterns:
 - ^test_
```
To allow long lines containing URLs in `Layout/LineLength`:

```
Layout/LineLength:
 AllowedPatterns:
 - '^\s*#\s*https?://'
```
| The `IgnoredPatterns`parameter is a deprecated alias for`AllowedPatterns`.
It still works, but you should migrate your configuration to use`AllowedPatterns`. |

# Cops

In RuboCop lingo the various checks performed on the code are called cops. Each cop is responsible for detecting one particular offense. There are several cop departments, grouping the cops by class of offense.

Many of the Style and Layout cops have configuration options, allowing them to
enforce different coding conventions. See Configuring Cops
for details on `Enabled`, `Severity`, `AutoCorrect`, `AllowedMethods`, and other settings.

You can also load custom cops.

## Departments

### Style

Style cops check for stylistic consistency of your code. Many of them are
based on the Ruby Style Guide. This is the largest
department and covers everything from preferred hash syntax to method naming
conventions. Most Style cops support multiple `EnforcedStyle` options so you
can adapt them to your team’s preferences.

### Layout

Layout cops inspect your code for consistent use of indentation, alignment,
and white space. They deal exclusively with formatting — changing Layout cops
never affects program behavior. Running `rubocop -x` applies only Layout
autocorrections.

### Lint

Lint cops check for ambiguities and possible errors in your code.

RuboCop implements, in a portable way, all built-in MRI lint checks
(`ruby -wc`) and adds a lot of extra lint checks of its own.

You can run only the Lint cops like this:

`$ rubocop -l`The `-l`/`--lint` option can be used together with `--only` to run all the
enabled Lint cops plus a selection of other cops.

Disabling Lint cops is generally a bad idea.

### Metrics

Metrics cops deal with properties of the source code that can be measured,
such as class length, method length, and cyclomatic complexity. They have a
configuration parameter called `Max` and when running
`rubocop --auto-gen-config`, this parameter will be set to the highest value
found for the inspected code.

### Naming

Naming cops check for naming issues in your code, such as method names, constant names, file names, and predicate prefixes. Like Style cops, many support multiple enforced styles.

## Cop metadata

Each cop in config/default.yml has metadata that helps you understand its status:

-
`VersionAdded`— the RuboCop version that introduced the cop
-
`VersionChanged`— the last version that changed the cop’s behavior or defaults
-
`Enabled`— whether the cop is enabled by default
-
`SafeAutoCorrect`— whether the autocorrection is safe (won’t change behavior)

Cop-related errors are silenced by default but can be surfaced using the
`--raise-cop-error` option, which is useful for debugging.

# Formatters

You can change the output format of RuboCop by specifying formatters with the `-f/--format` option.
RuboCop ships with several built-in formatters, and also you can create your custom formatter.

Additionally the output can be redirected to a file instead of `$stdout` with the `-o/--out` option.

Some of the built-in formatters produce **machine-parsable** output
and they are considered public APIs.
The rest of the formatters are for humans, so parsing their outputs is discouraged.

You can enable multiple formatters at the same time by specifying `-f/--format` multiple times.
The `-o/--out` option applies to the previously specified `-f/--format`,
or the default `progress` format if no `-f/--format` is specified before the `-o/--out` option.

### Quick reference

| Formatter | Flag | Use case |
|---|---|---|
| Progress (default) |
 | Interactive terminal use |
| Clang |
 | Detailed offenses with source context |
| Simple |
 | Compact human-readable output |
| Quiet |
 | Only show output when there are offenses |
| Fuubar |
 | Progress bar with real-time offense reporting |
| Pacman |
 | Fun progress indicator |
| Emacs |
 | Emacs-compatible output (machine-parsable) |
| JSON |
 | Machine consumption / custom tooling (machine-parsable) |
| JUnit |
 | CI systems like Jenkins (machine-parsable) |
| TAP |
 | Test Anything Protocol consumers (machine-parsable) |
| GitHub Actions |
 | GitHub Actions annotations |
| HTML |
 | HTML report for sharing / CI artifacts |
| Markdown |
 | Markdown report for PR comments |
| File List |
 | Pipe file names to other tools (machine-parsable) |
| Offense Count |
 | Overview of which cops trigger the most |
| Worst Offenders |
 | Find files needing the most attention |

```
# Simple format to $stdout.
$ rubocop --format simple
# Progress (default) format to the file result.txt.
$ rubocop --out result.txt
# Both progress and offense count formats to $stdout.
# The offense count formatter outputs only the final summary,
# so you'll mostly see the outputs from the progress formatter,
# and at the end the offense count summary will be output.
$ rubocop --format progress --format offenses
# Progress format to $stdout and JSON format to the file rubocop.json.
$ rubocop --format progress --format json --out rubocop.json
# ~~~~~~~~~~~~~~~~~ ~~~~~~~~~~~~~ ~~~~~~~~~~~~~~~~~~
# |_______________| |
# default format $stdout
# Progress format to result.txt and simple format to $stdout.
$ rubocop --out result.txt --format simple
# ~~~~~~~~~~~~~~~~ ~~~~~~~~~~~~~~~
# | |
# $stdout default format
```
You can also load custom formatters.

## Progress Formatter (default)

The default `progress` formatter outputs a character for each inspected file,
and at the end it displays all detected offenses in the `clang` format.
A `.` represents a clean file, and each of the capital letters means
the severest offense (convention, warning, error, or fatal) found in a file.

```
$ rubocop
Inspecting 26 files
..W.C....C..CWCW.C...WC.CC
Offenses:
lib/foo.rb:6:5: C: Style/Documentation: Missing top-level class documentation comment.
 class Foo
 ^^^^^
...
26 files inspected, 46 offenses detected
```
## Auto Gen Formatter

Behaves like Progress Formatter except that it will not show any offenses.

```
$ rubocop --format autogenconf
Inspecting 26 files
..W.C....C..CWCW.C...WC.CC
26 files inspected, 46 offenses detected
```
## Clang Style Formatter

The `clang` formatter displays the offenses in a manner similar to `clang`:

```
$ rubocop --format clang test.rb
Inspecting 1 file
W
Offenses:
test.rb:1:5: C: Naming/MethodName: Use snake_case for method names.
def badName
 ^^^^^^^
test.rb:2:3: C: Style/GuardClause: Use a guard clause instead of wrapping the code inside a conditional expression.
 if something
 ^^
test.rb:2:3: C: Style/IfUnlessModifier: Favor modifier if usage when having a single-line body. Another good alternative is the usage of control flow &&/||.
 if something
 ^^
test.rb:4:5: W: Layout/DefEndAlignment: end at 4, 4 is not aligned with if at 2, 2
 end
 ^^^
1 file inspected, 4 offenses detected
```
## Fuubar Style Formatter

The `fuubar` style formatter displays a progress bar
and shows details of offenses in the `clang` format as soon as they are detected.
This is inspired by the Fuubar formatter for RSpec.

```
$ rubocop --format fuubar
lib/foo.rb.rb:1:1: C: Naming/MethodName: Use snake_case for method names.
def badName
 ^^^^^^^
lib/bar.rb:13:14: W: Lint/UnreachableCode: Unreachable code detected.
 x + 1
 ^^^^^
 22/53 files |======== 43 ========> | ETA: 00:00:02
```
## Pacman Style Formatter

The `pacman` style formatter prints a PACDOT per every file to be analyzed. Pacman will "eat" one PACDOT per file when no offense is detected. Otherwise it will print a Ghost.
This is inspired by the Pacman formatter for RSpec.

```
$ rubocop --format pacman
Eating 31 files
src/foo.rb:1:1: C: Style/FrozenStringLiteralComment: Missing magic comment # frozen_string_literal: true.
src/bar.rb:14:15: C: Style/MutableConstant: Freeze mutable objects assigned to constants.
 GHOST = 'ᗣ'
 ^^^
....ᗣ...ᗣ...ᗧ••••••••••••••••••
31 examples, 2 failures
```
## Emacs Style Formatter

**Machine-parsable**

The `emacs` formatter displays the offenses in a format suitable for consumption by `Emacs` (and possibly other tools).

```
$ rubocop --format emacs test.rb
/Users/bozhidar/projects/test.rb:1:1: C: Naming/MethodName: Use snake_case for method names.
/Users/bozhidar/projects/test.rb:2:3: C: Style/IfUnlessModifier: Favor modifier if/unless usage when you have a single-line body. Another good alternative is the usage of control flow &&/||.
/Users/bozhidar/projects/test.rb:4:5: W: Layout/DefEndAlignment: end at 4, 4 is not aligned with if at 2, 2
```
## Simple Formatter

The name of the formatter says it all :-)

```
$ rubocop --format simple test.rb
== test.rb ==
C: 1: 5: Naming/MethodName: Use snake_case for method names.
C: 2: 3: Style/GuardClause: Use a guard clause instead of wrapping the code inside a conditional expression.
C: 2: 3: Style/IfUnlessModifier: Favor modifier if usage when having a single-line body. Another good alternative is the usage of control flow &&/||.
W: 4: 5: Layout/DefEndAlignment: end at 4, 4 is not aligned with if at 2, 2
1 file inspected, 4 offenses detected
```
## Quiet Formatter

Behaves like Simple Formatter if there are offenses. Completely quiet otherwise:

`$ rubocop --format quiet`## File List Formatter

**Machine-parsable**

Sometimes you might want to just open all files with offenses in your favorite editor. This formatter outputs just the names of the files with offenses in them and makes it possible to do something like:

`$ rubocop --format files | xargs vim`## JSON Formatter

**Machine-parsable**

You can get RuboCop’s inspection result in JSON format by passing `--format json` option in command line.
The JSON structure is like the following example:

```
{
 "metadata": {
 "rubocop_version": "1.84.2",
 "ruby_engine": "ruby",
 "ruby_version": "3.4.8",
 "ruby_patchlevel": "72",
 "ruby_platform": "arm64-darwin24"
 },
 "files": [{
 "path": "lib/foo.rb",
 "offenses": []
 }, {
 "path": "lib/bar.rb",
 "offenses": [{
 "severity": "convention",
 "message": "Line is too long. [81/80]",
 "cop_name": "Layout/LineLength",
 "corrected": true,
 "correctable": true,
 "location": {
 "start_line": 546,
 "start_column": 80,
 "last_line": 546,
 "last_column": 84,
 "length": 4,
 "line": 546,
 "column": 80
 }
 }, {
 "severity": "warning",
 "message": "Unreachable code detected.",
 "cop_name": "Lint/UnreachableCode",
 "corrected": false,
 "correctable": false,
 "location": {
 "start_line": 15,
 "start_column": 9,
 "last_line": 15,
 "last_column": 19,
 "length": 10,
 "line": 15,
 "column": 9
 }
 }
 ]
 }
 ],
 "summary": {
 "offense_count": 2,
 "target_file_count": 2,
 "inspected_file_count": 2
 }
}
```
## JUnit Style Formatter

**Machine-parsable**

The `junit` style formatter provides the JUnit formatting.
This formatter is based on the rubocop-junit-formatter gem.

```
$ rubocop --format junit
<?xml version='1.0'?>
<testsuites>
 <testsuite name='rubocop' tests='2' failures='2'>
 <testcase classname='example' name='Style/FrozenStringLiteralComment'>
 <failure type='Style/FrozenStringLiteralComment' message='Style/FrozenStringLiteralComment: Missing frozen string literal comment.'>
 /tmp/src/example.rb:1:1
 </failure>
 </testcase>
 <testcase classname='example' name='Naming/MethodName'>
 <failure type='Naming/MethodName' message='Naming/MethodName: Use snake_case for method names.'>
 /tmp/src/example.rb:1:5
 </failure>
 </testcase>
 <testcase classname='example' name='Lint/DeprecatedClassMethods'>
 <failure type='Lint/DeprecatedClassMethods' message='Lint/DeprecatedClassMethods: `File.exists?` is deprecated in favor of `File.exist?`.'>
 /tmp/src/example.rb:2:8
 </failure>
 </testcase>
 </testsuite>
</testsuites>
```
The `junit` style formatter is very useful for continuous integration systems
such as Jenkins, most of which support junit formatting when parsing test
results. A typical invocation in this type of scenario might look like:

`$ rubocop --format junit --out test-reports/junit.xml`Since there is one XML node for each cop for each file, the size of the resulting
XML can get quite large. If it is too large for you, you can restrict the output
to just failures by adding the `--display-only-failed` option.

## Offense Count Formatter

Sometimes when first applying RuboCop to a codebase, it’s nice to be able to see where most of your style cleanup is going to be spent.

With this in mind, you can use the offense count formatter to outline the offended cops and the number of offenses found for each by running:

```
$ rubocop --format offenses
36 Layout/LineLength [Safe Correctable]
18 Style/StringLiterals [Safe Correctable]
13 Style/Documentation
10 Style/ExpandPathArguments [Safe Correctable]
8 Style/EmptyMethod [Safe Correctable]
6 Layout/IndentationConsistency [Safe Correctable]
4 Lint/SuppressedException
3 Layout/EmptyLinesAroundAccessModifier [Safe Correctable]
2 Layout/ExtraSpacing [Safe Correctable]
1 Layout/AccessModifierIndentation [Safe Correctable]
1 Style/ClassAndModuleChildren [Unsafe Correctable]
--
102 Total in 31 files
```
## Worst Offenders Formatter

Similar to the Offense Count formatter, but lists the files which need the most attention:

```
$ rubocop --format worst
89 this/file/is/really/bad.rb
2 much/better.rb
--
91 Total in 2 files
```
## HTML Formatter

Useful for CI environments. It will create an HTML report like this.

`$ rubocop --format html -o rubocop.html`## Markdown Formatter

Useful for CI environments, especially if posting comments back to pull requests. It will create a markdown report like this.

`$ rubocop --format markdown -o rubocop.md`## TAP Formatter

**Machine-parsable**

Useful for CI environments, it will format report following the Test Anything Protocol.

```
$ rubocop --format tap
1..3
not ok 1 - lib/rubocop.rb
# lib/rubocop.rb:2:3: C: foo
# This is line 2.
# ^
ok 2 - spec/spec_helper.rb
not ok 3 - exe/rubocop
# exe/rubocop:5:2: E: bar
# This is line 5.
# ^
# exe/rubocop:6:1: C: foo
# This is line 6.
# ^
3 files inspected, 3 offenses detected
```
## GitHub Actions Formatter

Useful for GitHub Actions. Formats offenses as workflow commands to create annotations in GitHub UI.

The formatter uses fail_level to determine which GitHub level to use for each annotation. By default, all RuboCop severities are errors. If fail level is set and severity is below fail level, a warning will be created instead.

```
$ rubocop --format github
::error file=lib/foo.rb,line=6,col=5::Style/Documentation: Missing top-level class documentation comment.
```

# Bundler

## Bundler/DuplicatedGem

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.46 | 1.40 |

A Gem’s requirements should be listed only once in a Gemfile.

### Examples

```
# bad
gem 'rubocop'
gem 'rubocop'
# bad
group :development do
 gem 'rubocop'
end
group :test do
 gem 'rubocop'
end
# good
group :development, :test do
 gem 'rubocop'
end
# good
gem 'rubocop', groups: [:development, :test]
# good - conditional declaration
if Dir.exist?(local)
 gem 'rubocop', path: local
elsif ENV['RUBOCOP_VERSION'] == 'master'
 gem 'rubocop', git: 'https://github.com/rubocop/rubocop.git'
else
 gem 'rubocop', '~> 0.90.0'
end
```
## Bundler/DuplicatedGroup

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 1.56 | - |

A Gem group, or a set of groups, should be listed only once in a Gemfile.

For example, if the values of `source`, `git`, `platforms`, or `path`
surrounding `group` are different, no offense will be registered:

```
platforms :ruby do
 group :default do
 gem 'openssl'
 end
end
platforms :jruby do
 group :default do
 gem 'jruby-openssl'
 end
end
```
### Examples

```
# bad
group :development do
 gem 'rubocop'
end
group :development do
 gem 'rubocop-rails'
end
# bad (same set of groups declared twice)
group :development, :test do
 gem 'rubocop'
end
group :test, :development do
 gem 'rspec'
end
# good
group :development do
 gem 'rubocop'
end
group :development, :test do
 gem 'rspec'
end
# good
gem 'rubocop', groups: [:development, :test]
gem 'rspec', groups: [:development, :test]
```
## Bundler/GemComment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.59 | 1.88 |

Each gem in the Gemfile should have a comment explaining its purpose in the project, or the reason for its version or source.

The optional "OnlyFor" configuration array can be used to only register offenses when the gems use certain options or have version specifiers.

When "version_specifiers" is included, a comment will be enforced if the gem has any version specifier.

When "restrictive_version_specifiers" is included, a comment will be enforced if the gem has a version specifier that holds back the version of the gem.

For any other value in the array, a comment will be enforced for a gem if an option by the same name is present. A useful use case is to enforce a comment when using options that change the source of a gem:

-
`bitbucket`
-
`gist`
-
`git`
-
`github`
-
`source`

For a full list of options supported by bundler, see https://bundler.io/man/gemfile.5.html .

## Bundler/GemFilename

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 1.20 | - |

Verifies that a project contains Gemfile or gems.rb file and correct associated lock file based on the configuration.

## Bundler/GemVersion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 1.14 | - |

Enforce that Gem version specifications or a commit reference (branch, ref, or tag) are either required or forbidden.

### Examples

#### EnforcedStyle: required (default)

```
# bad
gem 'rubocop'
# good
gem 'rubocop', '~> 1.12'
# good
gem 'rubocop', '>= 1.10.0'
# good
gem 'rubocop', '>= 1.5.0', '< 1.10.0'
# good
gem 'rubocop', branch: 'feature-branch'
# good
gem 'rubocop', ref: '74b5bfbb2c4b6fd6cdbbc7254bd7084b36e0c85b'
# good
gem 'rubocop', tag: 'v1.17.0'
```
#### EnforcedStyle: forbidden

```
# good
gem 'rubocop'
# bad
gem 'rubocop', '~> 1.12'
# bad
gem 'rubocop', '>= 1.10.0'
# bad
gem 'rubocop', '>= 1.5.0', '< 1.10.0'
# bad
gem 'rubocop', branch: 'feature-branch'
# bad
gem 'rubocop', ref: '74b5bfbb2c4b6fd6cdbbc7254bd7084b36e0c85b'
# bad
gem 'rubocop', tag: 'v1.17.0'
```
## Bundler/InsecureProtocolSource

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.50 | 1.40 |

Passing symbol arguments to `source` (e.g. `source :rubygems`) is
deprecated because they default to using HTTP requests. Instead, specify
`'https://rubygems.org'` if possible, or `'http://rubygems.org'` if not.

When autocorrecting, this cop will replace symbol arguments with
`'https://rubygems.org'`.

This cop will not replace existing sources that use `http://`. This may
be necessary where HTTPS is not available. For example, where using an
internal gem server via an intranet, or where HTTPS is prohibited.
However, you should strongly prefer `https://` where possible, as it is
more secure.

If you don’t allow `http://`, please set `false` to `AllowHttpProtocol`.
This option is `true` by default for safe autocorrection.

## Bundler/OrderedGems

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.46 | 0.47 |

Gems should be alphabetically sorted within groups.

# Gemspec

## Gemspec/AddRuntimeDependency

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.65 | - |

Prefer `add_dependency` over `add_runtime_dependency` as the latter is
considered soft-deprecated.

## Gemspec/AttributeAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.77 | - |

Use consistent style for Gemspec attributes assignment.

### Examples

```
# bad
# This example uses two styles for assignment of metadata attribute.
Gem::Specification.new do |spec|
 spec.metadata = { 'key' => 'value' }
 spec.metadata['another-key'] = 'another-value'
end
# good
Gem::Specification.new do |spec|
 spec.metadata['key'] = 'value'
 spec.metadata['another-key'] = 'another-value'
end
# good
Gem::Specification.new do |spec|
 spec.metadata = { 'key' => 'value', 'another-key' => 'another-value' }
end
# bad
# This example uses two styles for assignment of authors attribute.
Gem::Specification.new do |spec|
 spec.authors = %w[author-0 author-1]
 spec.authors[2] = 'author-2'
end
# good
Gem::Specification.new do |spec|
 spec.authors = %w[author-0 author-1 author-2]
end
# good
Gem::Specification.new do |spec|
 spec.authors[0] = 'author-0'
 spec.authors[1] = 'author-1'
 spec.authors[2] = 'author-2'
end
# good
# This example uses consistent assignment per attribute,
# even though two different styles are used overall.
Gem::Specification.new do |spec|
 spec.metadata = { 'key' => 'value' }
 spec.authors[0] = 'author-0'
 spec.authors[1] = 'author-1'
 spec.authors[2] = 'author-2'
end
```
## Gemspec/DependencyVersion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 1.29 | - |

Enforce that gem dependency version specifications or a commit reference (branch, ref, or tag) are either required or forbidden.

### Examples

#### EnforcedStyle: required (default)

```
# bad
Gem::Specification.new do |spec|
 spec.add_dependency 'parser'
end
# bad
Gem::Specification.new do |spec|
 spec.add_development_dependency 'parser'
end
# good
Gem::Specification.new do |spec|
 spec.add_dependency 'parser', '>= 2.3.3.1', '< 3.0'
end
# good
Gem::Specification.new do |spec|
 spec.add_development_dependency 'parser', '>= 2.3.3.1', '< 3.0'
end
```
#### EnforcedStyle: forbidden

```
# bad
Gem::Specification.new do |spec|
 spec.add_dependency 'parser', '>= 2.3.3.1', '< 3.0'
end
# bad
Gem::Specification.new do |spec|
 spec.add_development_dependency 'parser', '>= 2.3.3.1', '< 3.0'
end
# good
Gem::Specification.new do |spec|
 spec.add_dependency 'parser'
end
# good
Gem::Specification.new do |spec|
 spec.add_development_dependency 'parser'
end
```
## Gemspec/DeprecatedAttributeAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.30 | 1.40 |

Checks that deprecated attributes are not set in a gemspec file. Removing deprecated attributes allows the user to receive smaller packed gems.

### Examples

```
# bad
Gem::Specification.new do |spec|
 spec.name = 'your_cool_gem_name'
 spec.test_files = Dir.glob('test/**/*')
end
# bad
Gem::Specification.new do |spec|
 spec.name = 'your_cool_gem_name'
 spec.test_files += Dir.glob('test/**/*')
end
# good
Gem::Specification.new do |spec|
 spec.name = 'your_cool_gem_name'
end
```
## Gemspec/DevelopmentDependencies

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.44 | - |

Enforce that development dependencies for a gem are specified in
`Gemfile`, rather than in the `gemspec` using
`add_development_dependency`. Alternatively, using ```
EnforcedStyle:
gemspec
```
, enforce that all dependencies are specified in `gemspec`,
rather than in `Gemfile`.

### Examples

#### EnforcedStyle: Gemfile (default)

```
# Specify runtime dependencies in your gemspec,
# but all other dependencies in your Gemfile.
# bad
# example.gemspec
s.add_development_dependency "foo"
# good
# Gemfile
gem "foo"
# good
# gems.rb
gem "foo"
# good (with AllowedGems: ["bar"])
# example.gemspec
s.add_development_dependency "bar"
```
#### EnforcedStyle: gems.rb

```
# Specify runtime dependencies in your gemspec,
# but all other dependencies in your Gemfile.
#
# Identical to `EnforcedStyle: Gemfile`, but with a different error message.
# Rely on Bundler/GemFilename to enforce the use of `Gemfile` vs `gems.rb`.
# bad
# example.gemspec
s.add_development_dependency "foo"
# good
# Gemfile
gem "foo"
# good
# gems.rb
gem "foo"
# good (with AllowedGems: ["bar"])
# example.gemspec
s.add_development_dependency "bar"
```
## Gemspec/DuplicatedAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.52 | 1.40 |

An attribute assignment method call should be listed only once in a gemspec.

Assigning to an attribute with the same name using `spec.foo =` or
`spec.attribute#[]=` will be an unintended usage. On the other hand,
duplication of methods such as `spec.requirements`,
`spec.add_runtime_dependency`, and others are permitted because it is
the intended use of appending values.

### Examples

```
# bad
Gem::Specification.new do |spec|
 spec.name = 'rubocop'
 spec.name = 'rubocop2'
end
# good
Gem::Specification.new do |spec|
 spec.name = 'rubocop'
end
# good
Gem::Specification.new do |spec|
 spec.requirements << 'libmagick, v6.0'
 spec.requirements << 'A good graphics card'
end
# good
Gem::Specification.new do |spec|
 spec.add_dependency('parallel', '~> 1.10')
 spec.add_dependency('parser', '>= 2.3.3.1', '< 3.0')
end
# bad
Gem::Specification.new do |spec|
 spec.metadata["key"] = "value"
 spec.metadata["key"] = "value"
end
# good
Gem::Specification.new do |spec|
 spec.metadata["key"] = "value"
end
```
## Gemspec/OrderedDependencies

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.51 | - |

Dependencies in the gemspec should be alphabetically sorted.

### Examples

```
# bad
spec.add_dependency 'rubocop'
spec.add_dependency 'rspec'
# good
spec.add_dependency 'rspec'
spec.add_dependency 'rubocop'
# good
spec.add_dependency 'rubocop'
spec.add_dependency 'rspec'
# bad
spec.add_development_dependency 'rubocop'
spec.add_development_dependency 'rspec'
# good
spec.add_development_dependency 'rspec'
spec.add_development_dependency 'rubocop'
# good
spec.add_development_dependency 'rubocop'
spec.add_development_dependency 'rspec'
# bad
spec.add_runtime_dependency 'rubocop'
spec.add_runtime_dependency 'rspec'
# good
spec.add_runtime_dependency 'rspec'
spec.add_runtime_dependency 'rubocop'
# good
spec.add_runtime_dependency 'rubocop'
spec.add_runtime_dependency 'rspec'
```
## Gemspec/RequireMFA

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.23 | 1.40 |

Requires a gemspec to have `rubygems_mfa_required` metadata set.

This setting tells RubyGems that MFA (Multi-Factor Authentication) is required for accounts to be able to perform privileged operations, such as (see RubyGems' documentation for the full list of privileged operations):

-
`gem push`
-
`gem yank`
-
`gem owner --add/remove`
-
adding or removing owners using gem ownership page

This helps make your gem more secure, as users can be more confident that gem updates were pushed by maintainers.

### Examples

```
# bad
Gem::Specification.new do |spec|
 # no `rubygems_mfa_required` metadata specified
end
# good
Gem::Specification.new do |spec|
 spec.metadata = {
 'rubygems_mfa_required' => 'true'
 }
end
# good
Gem::Specification.new do |spec|
 spec.metadata['rubygems_mfa_required'] = 'true'
end
# bad
Gem::Specification.new do |spec|
 spec.metadata = {
 'rubygems_mfa_required' => 'false'
 }
end
# good
Gem::Specification.new do |spec|
 spec.metadata = {
 'rubygems_mfa_required' => 'true'
 }
end
# bad
Gem::Specification.new do |spec|
 spec.metadata['rubygems_mfa_required'] = 'false'
end
# good
Gem::Specification.new do |spec|
 spec.metadata['rubygems_mfa_required'] = 'true'
end
```
## Gemspec/RequiredRubyVersion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.52 | 1.40 |

Checks that `required_ruby_version` in a gemspec file is set to a valid
value (non-blank) and matches `TargetRubyVersion` as set in RuboCop’s
configuration for the gem.

This ensures that RuboCop is using the same Ruby version as the gem.

### Examples

```
# When `TargetRubyVersion` of .rubocop.yml is `2.5`.
# bad
Gem::Specification.new do |spec|
 # no `required_ruby_version` specified
end
# bad
Gem::Specification.new do |spec|
 spec.required_ruby_version = '>= 2.4.0'
end
# bad
Gem::Specification.new do |spec|
 spec.required_ruby_version = '>= 2.6.0'
end
# bad
Gem::Specification.new do |spec|
 spec.required_ruby_version = ''
end
# good
Gem::Specification.new do |spec|
 spec.required_ruby_version = '>= 2.5.0'
end
# good
Gem::Specification.new do |spec|
 spec.required_ruby_version = '>= 2.5'
end
# accepted but not recommended
Gem::Specification.new do |spec|
 spec.required_ruby_version = ['>= 2.5.0', '< 2.7.0']
end
# accepted but not recommended, since
# Ruby does not really follow semantic versioning
Gem::Specification.new do |spec|
 spec.required_ruby_version = '~> 2.5'
end
```
## Gemspec/RubyVersionGlobalsUsage

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.72 | 1.40 |

Checks that `RUBY_VERSION` and `Ruby::VERSION` constants are not used in gemspec.
Using `RUBY_VERSION` and `Ruby::VERSION` is dangerous because the value of the
constant is determined by `rake release`.
It’s possible to have a dependency based on the Ruby version used
to execute `rake release` and not the user’s Ruby version.

### Examples

```
# bad
Gem::Specification.new do |spec|
 if RUBY_VERSION >= '3.0'
 spec.add_dependency 'gem_a'
 else
 spec.add_dependency 'gem_b'
 end
end
# good
Gem::Specification.new do |spec|
 spec.add_dependency 'gem_a'
end
```

# Layout

## Layout/AccessModifierIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Bare access modifiers (those not applying to specific methods) should be
indented as deep as method definitions, or as deep as the `class`/`module`
keyword, depending on configuration.

### Examples

## Layout/ArgumentAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.68 | 0.77 |

Checks that the arguments on a multi-line method call are aligned.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| IndentationWidth |
 | Integer |

## Layout/ArrayAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.77 |

Checks that the elements of a multi-line array literal are aligned.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| IndentationWidth |
 | Integer |

## Layout/AssignmentIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.45 |

Checks the indentation of the first line of the right-hand-side of a multi-line assignment.

The indentation of the remaining lines can be corrected with
other cops such as `Layout/IndentationConsistency` and `Layout/EndAlignment`.

## Layout/BeginEndAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.91 | - |

Checks whether the end keyword of `begin` is aligned properly.

Two modes are supported through the `EnforcedStyleAlignWith` configuration
parameter. If it’s set to `start_of_line` (which is the default), the
`end` shall be aligned with the start of the line where the `begin`
keyword is. If it’s set to `begin`, the `end` shall be aligned with the
`begin` keyword.

`Layout/EndAlignment` cop aligns with keywords (e.g. `if`, `while`, `case`)
by default. On the other hand, `||= begin` that this cop targets tends to
align with the start of the line, it defaults to `EnforcedStyleAlignWith: start_of_line`.
These styles can be configured by each cop.

## Layout/BlockAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

Checks whether the end keywords are aligned properly for do end blocks.

Three modes are supported through the `EnforcedStyleAlignWith`
configuration parameter:

`start_of_block` : the `end` shall be aligned with the
start of the line where the `do` appeared.

`start_of_line` : the `end` shall be aligned with the
start of the line where the expression started.

`either` (which is the default) : the `end` is allowed to be in either
location. The autocorrect will default to `start_of_line`.

When the `do` or `{` appears on a continuation line of multiline
method arguments, the start of the line where the method is called
is used as the alignment target instead of that continuation line.

## Layout/CaseIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.16 |

Checks how the `when` and `in`s of a `case` expression
are indented in relation to its `case` or `end` keyword.

It will register a separate offense for each misaligned `when` and `in`.

### Examples

```
# If Layout/EndAlignment is set to keyword style (default)
# *case* and *end* should always be aligned to same depth,
# and therefore *when* should always be aligned to both -
# regardless of configuration.
# bad for all styles
case n
 when 0
 x * 2
 else
 y / 3
end
case n
 in pattern
 x * 2
 else
 y / 3
end
# good for all styles
case n
when 0
 x * 2
else
 y / 3
end
case n
in pattern
 x * 2
else
 y / 3
end
```
#### EnforcedStyle: case (default)

```
# if EndAlignment is set to other style such as
# start_of_line (as shown below), then *when* alignment
# configuration does have an effect.
# bad
a = case n
when 0
 x * 2
else
 y / 3
end
a = case n
in pattern
 x * 2
else
 y / 3
end
# good
a = case n
 when 0
 x * 2
 else
 y / 3
end
a = case n
 in pattern
 x * 2
 else
 y / 3
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| IndentOneStep |
 | Boolean |
| IndentationWidth |
 | Integer |

## Layout/ClassStructure

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always (Unsafe) | 0.52 | 1.89 |

Checks if the code style follows the `ExpectedOrder` configuration:

`Categories` allows us to map macro names into a category.

Consider an example of code style that covers the following order:

-
Module inclusion ( `include`,`prepend`,`extend`)
-
Constants
-
Associations ( `has_one`,`has_many`)
-
Public attribute macros ( `attr_accessor`,`attr_writer`,`attr_reader`)
-
Other macros ( `validates`,`validate`)
-
Public class methods
-
Initializer
-
Public instance methods
-
Protected attribute macros ( `attr_accessor`,`attr_writer`,`attr_reader`)
-
Protected instance methods
-
Private attribute macros ( `attr_accessor`,`attr_writer`,`attr_reader`)
-
Private instance methods

| Simply enabling the cop with `Enabled: true`will not use
the example order shown below.
To enforce the order of macros like`attr_reader`,
you must define both`ExpectedOrder`and`Categories`. |

You can configure the following order:

```
 Layout/ClassStructure:
 ExpectedOrder:
 - module_inclusion
 - constants
 - association
 - public_attribute_macros
 - public_delegate
 - macros
 - public_class_methods
 - initializer
 - public_methods
 - protected_attribute_macros
 - protected_methods
 - private_attribute_macros
 - private_delegate
 - private_methods
```
Instead of putting all literals in the expected order, it is also possible to group categories of macros. Visibility levels are handled automatically.

```
 Layout/ClassStructure:
 Categories:
 association:
 - has_many
 - has_one
 attribute_macros:
 - attr_accessor
 - attr_reader
 - attr_writer
 macros:
 - validates
 - validate
 module_inclusion:
 - include
 - prepend
 - extend
```
If you only set `ExpectedOrder`
without defining `Categories`,
macros such as `attr_reader` or `has_many`
will not be recognized as part of a category, and their order will not be validated.
For example, the following will NOT raise any offenses, even if the order is incorrect:

```
Layout/ClassStructure:
 Enabled: true
 ExpectedOrder:
 - public_attribute_macros
 - initializer
```
To make it work as expected, you must also specify `Categories` like this:

```
Layout/ClassStructure:
 ExpectedOrder:
 - public_attribute_macros
 - initializer
 Categories:
 attribute_macros:
 - attr_reader
 - attr_writer
 - attr_accessor
```
### Safety

Autocorrection is unsafe because class methods and module inclusion can behave differently, based on which methods or constants have already been defined.

Constants will only be moved when they are assigned with literals.

### Examples

```
# bad
# Expect extend be before constant
class Person < ApplicationRecord
 has_many :orders
 ANSWER = 42
 extend SomeModule
 include AnotherModule
end
# good
class Person
 # extend and include go first
 extend SomeModule
 include AnotherModule
 # inner classes
 CustomError = Class.new(StandardError)
 # constants are next
 SOME_CONSTANT = 20
 # afterwards we have public attribute macros
 attr_reader :name
 # followed by other macros (if any)
 validates :name
 # then we have public delegate macros
 delegate :to_s, to: :name
 # public class methods are next in line
 def self.some_method
 end
 # initialization goes between class methods and instance methods
 def initialize
 end
 # followed by other public instance methods
 def some_method
 end
 # protected attribute macros and methods go next
 protected
 attr_reader :protected_name
 def some_protected_method
 end
 # private attribute macros, delegate macros and methods
 # are grouped near the end
 private
 attr_reader :private_name
 delegate :some_private_delegate, to: :name
 def some_private_method
 end
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| Categories |
 | |
| ExpectedOrder |
 | Array |

## Layout/ClosingHeredocIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.57 | - |

Checks the indentation of here document closings.

## Layout/ClosingParenthesisIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.86 |

Checks the indentation of hanging closing parentheses in
method calls, method definitions, and grouped expressions. A hanging
closing parenthesis means `)` preceded by a line break.

### Examples

```
# bad
some_method(
 a,
 b
 )
some_method(
 a, b
 )
some_method(a, b, c
 )
some_method(a,
 b,
 c
 )
some_method(a,
 x: 1,
 y: 2
 )
# Scenario 1: When First Parameter Is On Its Own Line
# good: when first param is on a new line, right paren is *always*
# outdented by IndentationWidth
some_method(
 a,
 b
)
# good
some_method(
 a, b
)
# Scenario 2: When First Parameter Is On The Same Line
# good: when all other params are also on the same line, outdent
# right paren by IndentationWidth
some_method(a, b, c
 )
# good: when all other params are on multiple lines, but are lined
# up, align right paren with left paren
some_method(a,
 b,
 c
 )
# good: when other params are not lined up on multiple lines, outdent
# right paren by IndentationWidth
some_method(a,
 x: 1,
 y: 2
)
```
## Layout/CommentIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.86 |

Checks the indentation of comments.

### Examples

```
# bad
 # comment here
def method_name
end
 # comment here
a = 'hello'
# yet another comment
 if true
 true
 end
# good
# comment here
def method_name
end
# comment here
a = 'hello'
# yet another comment
if true
 true
end
```
## Layout/ConditionPosition

## Layout/DefEndAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

Checks whether the end keywords of method definitions are aligned properly.

Two modes are supported through the EnforcedStyleAlignWith configuration
parameter. If it’s set to `start_of_line` (which is the default), the
`end` shall be aligned with the start of the line where the `def`
keyword is. If it’s set to `def`, the `end` shall be aligned with the
`def` keyword.

## Layout/DotPosition

## Layout/ElseAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks the alignment of else keywords. Normally they should be aligned with an if/unless/while/until/begin/def/rescue keyword, but there are special cases when they should follow the same rules as the alignment of end.

## Layout/EmptyComment

## Layout/EmptyLineAfterGuardClause

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.56 | 0.59 |

Enforces empty line after guard clause.

This cop allows a SimpleCov directive comment after guard clause because
SimpleCov excludes code from the coverage report by wrapping it in such directives.
Both the legacy `# :nocov:` comment and the newer `# simplecov:disable` /
`# simplecov:enable` comments are recognized:

```
def foo
 # :nocov:
 return if condition
 # :nocov:
 bar
end
def foo
 # simplecov:disable
 return if condition
 # simplecov:enable
 bar
end
```
Refer to SimpleCov’s documentation for more details: https://github.com/simplecov-ruby/simplecov#ignoringskipping-code

## Layout/EmptyLineAfterMagicComment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for a newline after the final magic comment.

## Layout/EmptyLineAfterMultilineCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.90 | - |

Enforces empty line after multiline condition.

### Examples

```
# bad
if multiline &&
 condition
 do_something
end
# good
if multiline &&
 condition
 do_something
end
# bad
case x
when foo,
 bar
 do_something
end
# good
case x
when foo,
 bar
 do_something
end
# bad
begin
 do_something
rescue FooError,
 BarError
 handle_error
end
# good
begin
 do_something
rescue FooError,
 BarError
 handle_error
end
```
## Layout/EmptyLineBetweenDefs

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.23 |

Checks whether class/module/method definitions are separated by one or more empty lines.

`NumberOfEmptyLines` can be an integer (default is 1) or
an array (e.g. [1, 2]) to specify a minimum and maximum
number of empty lines permitted.

`AllowAdjacentOneLineDefs` configures whether adjacent
one-line definitions are considered an offense.

### Examples

#### EmptyLineBetweenMethodDefs: true (default)

```
# checks for empty lines between method definitions.
# bad
def a
end
def b
end
# good
def a
end
def b
end
```
#### EmptyLineBetweenClassDefs: true (default)

```
# checks for empty lines between class definitions.
# bad
class A
end
class B
end
def b
end
# good
class A
end
class B
end
def b
end
```
#### EmptyLineBetweenModuleDefs: true (default)

```
# checks for empty lines between module definitions.
# bad
module A
end
module B
end
def b
end
# good
module A
end
module B
end
def b
end
```
## Layout/EmptyLines

## Layout/EmptyLinesAfterModuleInclusion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.79 | - |

Checks for an empty line after a module inclusion method (`extend`,
`include` and `prepend`), or a group of them.

## Layout/EmptyLinesAroundAccessModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Access modifiers should be surrounded by blank lines.

### Examples

## Layout/EmptyLinesAroundArguments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks if empty lines exist around the arguments of a method invocation.

## Layout/EmptyLinesAroundAttributeAccessor

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.83 | 0.84 |

Checks for a newline after an attribute accessor or a group of them.
`alias` syntax and `alias_method`, `public`, `protected`, and `private` methods are allowed
by default. These are customizable with `AllowAliasSyntax` and `AllowedMethods` options.

### Examples

```
# bad
attr_accessor :foo
def do_something
end
# good
attr_accessor :foo
def do_something
end
# good
attr_accessor :foo
attr_reader :bar
attr_writer :baz
attr :qux
def do_something
end
```
## Layout/EmptyLinesAroundBeginBody

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks if empty lines exist around the bodies of begin-end blocks.

## Layout/EmptyLinesAroundBlockBody

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks if empty lines around the bodies of blocks match the configuration.

## Layout/EmptyLinesAroundClassBody

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.53 |

Checks if empty lines around the bodies of classes match the configuration.

### Examples

#### EnforcedStyle: no_empty_lines (default)

```
# bad
class Foo
 def bar
 # ...
 end
end
# good
class Foo
 def bar
 # ...
 end
end
```
## Layout/EmptyLinesAroundExceptionHandlingKeywords

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks if empty lines exist around the bodies of `begin`
sections. This cop doesn’t check empty lines at `begin` body
beginning/end and around method definition body.
`Layout/EmptyLinesAroundBeginBody` or `Layout/EmptyLinesAroundMethodBody`
can be used for this purpose.

## Layout/EmptyLinesAroundMethodBody

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks if empty lines exist around the bodies of methods.

## Layout/EmptyLinesAroundModuleBody

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks if empty lines around the bodies of modules match the configuration.

### Examples

#### EnforcedStyle: no_empty_lines (default)

```
# bad
module Foo
 def bar
 # ...
 end
end
# good
module Foo
 def bar
 # ...
 end
end
```
## Layout/EndAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

Checks whether the end keywords are aligned properly.

Three modes are supported through the `EnforcedStyleAlignWith`
configuration parameter:

If it’s set to `keyword` (which is the default), the `end`
shall be aligned with the start of the keyword (if, class, etc.).

If it’s set to `variable` the `end` shall be aligned with the
left-hand-side of the variable assignment, if there is one.

If it’s set to `start_of_line`, the `end` shall be aligned with the
start of the line where the matching keyword appears.

This `Layout/EndAlignment` cop aligns with keywords (e.g. `if`, `while`, `case`)
by default. On the other hand, `Layout/BeginEndAlignment` cop aligns with
`EnforcedStyleAlignWith: start_of_line` by default because `||= begin` tends
to align with the start of the line. `Layout/DefEndAlignment` cop also aligns with
`EnforcedStyleAlignWith: start_of_line` by default.
These styles can be configured by each cop.

## Layout/EndOfLine

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.49 | - |

Checks for Windows-style line endings in the source code.

### Examples

#### EnforcedStyle: native (default)

```
# The `native` style means that CR+LF (Carriage Return + Line Feed) is
# enforced on Windows, and LF is enforced on other platforms.
# bad
puts 'Hello' # Return character is LF on Windows.
puts 'Hello' # Return character is CR+LF on other than Windows.
# good
puts 'Hello' # Return character is CR+LF on Windows.
puts 'Hello' # Return character is LF on other than Windows.
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Layout/ExtraSpacing

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for extra/unnecessary whitespace.

### Examples

```
# good if AllowForAlignment is true
name = "RuboCop"
# Some comment and an empty line
website += "/rubocop/rubocop" unless cond
puts "rubocop" if debug
# bad for any configuration
set_app("RuboCop")
website = "https://github.com/rubocop/rubocop"
# good only if AllowBeforeTrailingComments is true
object.method(arg) # this is a comment
# good even if AllowBeforeTrailingComments is false or not set
object.method(arg) # this is a comment
# good with either AllowBeforeTrailingComments or AllowForAlignment
object.method(arg) # this is a comment
another_object.method(arg) # this is another comment
some_object.method(arg) # this is some comment
```
## Layout/FirstArgumentIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.68 | 0.77 |

Checks the indentation of the first argument in a method call.
Arguments after the first one are checked by `Layout/ArgumentAlignment`,
not by this cop.

For indenting the first parameter of method *definitions*, check out
`Layout/FirstParameterIndentation`.

This cop will respect `Layout/ArgumentAlignment` and will not work when
`EnforcedStyle: with_fixed_indentation` is specified for `Layout/ArgumentAlignment`.

### Examples

```
# bad
some_method(
first_param,
second_param)
foo = some_method(
first_param,
second_param)
foo = some_method(nested_call(
nested_first_param),
second_param)
foo = some_method(
nested_call(
nested_first_param),
second_param)
some_method nested_call(
nested_first_param),
second_param
```
#### EnforcedStyle: special_for_inner_method_call_in_parentheses (default)

```
# Same as `special_for_inner_method_call` except that the special rule
# only applies if the outer method call encloses its arguments in
# parentheses.
# good
some_method(
 first_param,
second_param)
foo = some_method(
 first_param,
second_param)
foo = some_method(nested_call(
 nested_first_param),
second_param)
foo = some_method(
 nested_call(
 nested_first_param),
second_param)
some_method nested_call(
 nested_first_param),
second_param
```
#### EnforcedStyle: consistent

```
# The first argument should always be indented one step more than the
# preceding line.
# good
some_method(
 first_param,
second_param)
foo = some_method(
 first_param,
second_param)
foo = some_method(nested_call(
 nested_first_param),
second_param)
foo = some_method(
 nested_call(
 nested_first_param),
second_param)
some_method nested_call(
 nested_first_param),
second_param
```
#### EnforcedStyle: consistent_relative_to_receiver

```
# The first argument should always be indented one level relative to
# the parent that is receiving the argument
# good
some_method(
 first_param,
second_param)
foo = some_method(
 first_param,
second_param)
foo = some_method(nested_call(
 nested_first_param),
second_param)
foo = some_method(
 nested_call(
 nested_first_param),
second_param)
some_method nested_call(
 nested_first_param),
second_params
```
#### EnforcedStyle: special_for_inner_method_call

```
# The first argument should normally be indented one step more than
# the preceding line, but if it's an argument for a method call that
# is itself an argument in a method call, then the inner argument
# should be indented relative to the inner method.
# good
some_method(
 first_param,
second_param)
foo = some_method(
 first_param,
second_param)
foo = some_method(nested_call(
 nested_first_param),
second_param)
foo = some_method(
 nested_call(
 nested_first_param),
second_param)
some_method nested_call(
 nested_first_param),
second_param
```
## Layout/FirstArrayElementIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.68 | 0.77 |

Checks the indentation of the first element in an array literal
where the opening bracket and the first element are on separate lines.
The other elements' indentations are handled by `Layout/ArrayAlignment` cop.

This cop will respect `Layout/ArrayAlignment` and will not work when
`EnforcedStyle: with_fixed_indentation` is specified for `Layout/ArrayAlignment`.

By default, array literals that are arguments in a method call with parentheses, and where the opening square bracket of the array is on the same line as the opening parenthesis of the method call, shall have their first element indented one step (two spaces) more than the position inside the opening parenthesis.

Other array literals shall have their first element indented one step more than the start of the line where the opening square bracket is.

This default style is called 'special_inside_parentheses'. Alternative styles are 'consistent' and 'align_brackets'. Here are examples:

### Examples

#### EnforcedStyle: special_inside_parentheses (default)

```
# The `special_inside_parentheses` style enforces that the first
# element in an array literal where the opening bracket and first
# element are on separate lines is indented one step (two spaces) more
# than the position inside the opening parenthesis.
# bad
array = [
 :value
]
and_in_a_method_call([
 :no_difference
 ])
# good
array = [
 :value
]
but_in_a_method_call([
 :its_like_this
 ])
```
#### EnforcedStyle: consistent

```
# The `consistent` style enforces that the first element in an array
# literal where the opening bracket and the first element are on
# separate lines is indented the same as an array literal which is not
# defined inside a method call.
# bad
array = [
 :value
]
but_in_a_method_call([
 :its_like_this
])
# good
array = [
 :value
]
and_in_a_method_call([
 :no_difference
])
```
## Layout/FirstArrayElementLineBreak

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.49 | - |

Checks for a line break before the first element in a multi-line array.

## Layout/FirstHashElementIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.68 | 0.77 |

Checks the indentation of the first key in a hash literal where the opening brace and the first key are on separate lines. The other keys' indentations are handled by the HashAlignment cop.

By default, `Hash` literals that are arguments in a method call with
parentheses, and where the opening curly brace of the hash is on the
same line as the opening parenthesis of the method call, shall have
their first key indented one step (two spaces) more than the position
inside the opening parenthesis.

Other hash literals shall have their first key indented one step more than the start of the line where the opening curly brace is.

This default style is called 'special_inside_parentheses'. Alternative styles are 'consistent' and 'align_braces'. Here are examples:

### Examples

#### EnforcedStyle: special_inside_parentheses (default)

```
# The `special_inside_parentheses` style enforces that the first key
# in a hash literal where the opening brace and the first key are on
# separate lines is indented one step (two spaces) more than the
# position inside the opening parentheses.
# bad
hash = {
 key: :value
 }
in_a_method_call({
 foo: :bar
})
takes_multi_pairs_hash(x: {
 a: 1,
 b: 2
},
 y: {
 c: 1,
 d: 2
 })
# good
hash = {
 key: :value
}
in_a_method_call({
 foo: :bar
 })
takes_multi_pairs_hash(x: {
 a: 1,
 b: 2
 },
 y: {
 c: 1,
 d: 2
 })
```
#### EnforcedStyle: consistent

```
# The `consistent` style enforces that the first key in a hash
# literal where the opening brace and the first key are on
# separate lines is indented the same as a hash literal which is not
# defined inside a method call.
# bad
hash = {
 key: :value
 }
in_a_method_call({
 foo: :bar
 })
# good
hash = {
 key: :value
}
in_a_method_call({
 foo: :bar
})
```
#### EnforcedStyle: align_braces

```
# The `align_brackets` style enforces that the opening and closing
# braces are indented to the same position.
# bad
and_now_for_something = {
 completely: :different
}
takes_multi_pairs_hash(x: {
 a: 1,
 b: 2
},
 y: {
 c: 1,
 d: 2
 })
# good
and_now_for_something = {
 completely: :different
 }
takes_multi_pairs_hash(x: {
 a: 1,
 b: 2
 },
 y: {
 c: 1,
 d: 2
 })
```
## Layout/FirstHashElementLineBreak

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.49 | - |

Checks for a line break before the first element in a multi-line hash.

## Layout/FirstMethodArgumentLineBreak

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.49 | - |

Checks for a line break before the first argument in a multi-line method call.

### Examples

```
# bad
method(foo, bar,
 baz)
# good
method(
 foo, bar,
 baz)
 # ignored
 method foo, bar,
 baz
```
#### AllowMultilineFinalElement: false (default)

```
# bad
method(foo, bar, {
 baz: "a",
 qux: "b",
})
# good
method(
 foo, bar, {
 baz: "a",
 qux: "b",
})
```
## Layout/FirstMethodParameterLineBreak

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.49 | - |

Checks for a line break before the first parameter in a multi-line method parameter definition.

### Examples

```
# bad
def method(foo, bar,
 baz)
 do_something
end
# good
def method(
 foo, bar,
 baz)
 do_something
end
# ignored
def method foo,
 bar
 do_something
end
```
## Layout/FirstParameterIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.77 |

Checks the indentation of the first parameter in a method
definition. Parameters after the first one are checked by
`Layout/ParameterAlignment`, not by this cop.

For indenting the first argument of method *calls*, check out
`Layout/FirstArgumentIndentation`, which supports options related to
nesting that are irrelevant for method *definitions*.

## Layout/HashAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.16 |

Checks that the keys, separators, and values of a multi-line hash literal are aligned according to configuration. The configuration options are:

-
key (left align keys, one space before hash rockets and values)
-
separator (align hash rockets and colons, right align keys)
-
table (left align keys, hash rockets, and values)

The treatment of hashes passed as the last argument to a method call can also be configured. The options are:

-
always_inspect
-
always_ignore
-
ignore_implicit (without curly braces)

Alternatively you can specify multiple allowed styles. That’s done by passing a list of styles to EnforcedHashRocketStyle and EnforcedColonStyle.

### Examples

#### EnforcedHashRocketStyle: key (default)

```
# bad
{
 :foo => bar,
 :ba => baz
}
{
 :foo => bar,
 :ba => baz
}
# good
{
 :foo => bar,
 :ba => baz
}
```
#### EnforcedHashRocketStyle: separator

```
# bad
{
 :foo => bar,
 :ba => baz
}
{
 :foo => bar,
 :ba => baz
}
# good
{
 :foo => bar,
 :ba => baz
}
```
#### EnforcedColonStyle: key (default)

```
# bad
{
 foo: bar,
 ba: baz
}
{
 foo: bar,
 ba: baz
}
# good
{
 foo: bar,
 ba: baz
}
```
#### EnforcedLastArgumentHashStyle: always_inspect (default)

```
# Inspect both implicit and explicit hashes.
# bad
do_something(foo: 1,
 bar: 2)
# bad
do_something({foo: 1,
 bar: 2})
# good
do_something(foo: 1,
 bar: 2)
# good
do_something(
 foo: 1,
 bar: 2
)
# good
do_something({foo: 1,
 bar: 2})
# good
do_something({
 foo: 1,
 bar: 2
})
```
#### EnforcedLastArgumentHashStyle: always_ignore

```
# Ignore both implicit and explicit hashes.
# good
do_something(foo: 1,
 bar: 2)
# good
do_something({foo: 1,
 bar: 2})
```
## Layout/HeredocArgumentClosingParenthesis

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.68 | - |

Checks for the placement of the closing parenthesis in a method call that passes a HEREDOC string as an argument. It should be placed at the end of the line containing the opening HEREDOC tag.

## Layout/HeredocIndentation

| Requires Ruby version 2.3 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.85 |

Checks the indentation of the here document bodies. The bodies are indented one step.

| When `Layout/LineLength`'s`AllowHeredoc`is false (not default),
 this cop does not add any offenses for long here documents to
 avoid`Layout/LineLength`'s offenses. |

## Layout/IndentationConsistency

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for inconsistent indentation.

The difference between `indented_internal_methods` and `normal` is
that the `indented_internal_methods` style prescribes that in
classes and modules the `protected` and `private` modifier keywords
shall be indented the same as public methods and that protected and
private members shall be indented one step more than the modifiers.
Other than that, both styles mean that entities on the same logical
depth shall have the same indentation.

### Examples

#### EnforcedStyle: normal (default)

```
# bad
class A
 def test
 puts 'hello'
 puts 'world'
 end
end
# bad
class A
 def test
 puts 'hello'
 puts 'world'
 end
 protected
 def foo
 end
 private
 def bar
 end
end
# good
class A
 def test
 puts 'hello'
 puts 'world'
 end
end
# good
class A
 def test
 puts 'hello'
 puts 'world'
 end
 protected
 def foo
 end
 private
 def bar
 end
end
```
#### EnforcedStyle: indented_internal_methods

```
# bad
class A
 def test
 puts 'hello'
 puts 'world'
 end
end
# bad
class A
 def test
 puts 'hello'
 puts 'world'
 end
 protected
 def foo
 end
 private
 def bar
 end
end
# good
class A
 def test
 puts 'hello'
 puts 'world'
 end
end
# good
class A
 def test
 puts 'hello'
 puts 'world'
 end
 protected
 def foo
 end
 private
 def bar
 end
end
```
## Layout/IndentationStyle

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.82 |

Checks that the indentation method is consistent. Either tabs only or spaces only are used for indentation.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| IndentationWidth |
 | Integer |
| EnforcedStyle |
 |
 |

## Layout/IndentationWidth

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for indentation that doesn’t use the specified number of spaces.
The indentation width can be configured using the `Width` setting. The default width is 2.
The block body indentation for method chain blocks can be configured using the
`EnforcedStyleAlignWith` setting.

See also the `Layout/IndentationConsistency` cop which is the companion to this one.

### Examples

#### Width: 2 (default)

```
# bad
class A
 def test
 puts 'hello'
 end
end
# good
class A
 def test
 puts 'hello'
 end
end
```
```
# bad
value = (
foo - bar
)
# good
value = (
 foo - bar
)
```
#### AllowedPatterns: ['^\s*module']

```
# bad
module A
class B
 def test
 puts 'hello'
 end
end
end
# good
module A
class B
 def test
 puts 'hello'
 end
end
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| Width |
 | Integer |
| EnforcedStyleAlignWith |
 |
 |
| AllowedPatterns |
 | Array |

## Layout/LeadingCommentSpace

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.73 |

Checks whether comments have a leading space after the
```
 denoting the start of the comment. The leading space is not
required for some RDoc special syntax, like
```
`++`, `--`,
`:nodoc`, `=begin`- and `=end` comments, "shebang" directives,
or rackup options.

### Examples

```
# bad
#Some comment
# good
# Some comment
```
#### AllowRBSInlineAnnotation: false (default)

```
# bad
include Enumerable #[Integer]
attr_reader :name #: String
attr_reader :age #: Integer?
#: (
#| Integer,
#| String
#| ) -> void
def foo; end
```
#### AllowRBSInlineAnnotation: true

```
# good
include Enumerable #[Integer]
attr_reader :name #: String
attr_reader :age #: Integer?
#: (
#| Integer,
#| String
#| ) -> void
def foo; end
```
#### AllowSteepAnnotation: false (default)

```
# bad
[1, 2, 3].each_with_object([]) do |n, list| #$ Array[Integer]
 list << n
end
name = 'John' #: String
```
#### AllowSteepAnnotation: true

```
# good
[1, 2, 3].each_with_object([]) do |n, list| #$ Array[Integer]
 list << n
end
name = 'John' #: String
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowDoxygenCommentStyle |
 | Boolean |
| AllowGemfileRubyComment |
 | Boolean |
| AllowRBSInlineAnnotation |
 | Boolean |
| AllowSteepAnnotation |
 | Boolean |
| AllowYARDCommentBlockSeparator |
 | Boolean |

## Layout/LeadingEmptyLines

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.57 | 0.77 |

Checks for unnecessary leading blank lines at the beginning of a file.

## Layout/LineContinuationLeadingSpace

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.31 | 1.45 |

Checks that strings broken over multiple lines (by a backslash) contain trailing spaces instead of leading spaces (default) or leading spaces instead of trailing spaces.

## Layout/LineContinuationSpacing

## Layout/LineEndStringConcatenationIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.18 | - |

Checks the indentation of the next line after a line that ends with a string literal and a backslash.

If `EnforcedStyle: aligned` is set, the concatenated string parts shall be aligned with the
first part. There are some exceptions, such as implicit return values, where the
concatenated string parts shall be indented regardless of `EnforcedStyle` configuration.

If `EnforcedStyle: indented` is set, it’s the second line that shall be indented one step
more than the first line. Lines 3 and forward shall be aligned with line 2.

### Examples

```
# bad
def some_method
 'x' \
 'y' \
 'z'
end
my_hash = {
 first: 'a message' \
 'in two parts'
}
# good
def some_method
 'x' \
 'y' \
 'z'
end
```
## Layout/LineLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.25 | 1.82 |

Checks the length of lines in the source code.
The maximum length is configurable.
The tab size is configured in the `IndentationWidth`
of the `Layout/IndentationStyle` cop.
It also ignores a shebang line by default.

This cop has some autocorrection capabilities. It can programmatically shorten certain long lines by inserting line breaks into expressions that can be safely split across lines. These include arrays, hashes, and method calls with argument lists.

If autocorrection is enabled, the following cops are recommended to further format the broken lines. (Many of these are enabled by default.)

-
`Layout/ArgumentAlignment`
-
`Layout/ArrayAlignment`
-
`Layout/BlockAlignment`
-
`Layout/BlockEndNewline`
-
`Layout/ClosingParenthesisIndentation`
-
`Layout/FirstArgumentIndentation`
-
`Layout/FirstArrayElementIndentation`
-
`Layout/FirstHashElementIndentation`
-
`Layout/FirstParameterIndentation`
-
`Layout/HashAlignment`
-
`Layout/IndentationWidth`
-
`Layout/MultilineArrayLineBreaks`
-
`Layout/MultilineBlockLayout`
-
`Layout/MultilineHashBraceLayout`
-
`Layout/MultilineHashKeyLineBreaks`
-
`Layout/MultilineMethodArgumentLineBreaks`
-
`Layout/MultilineMethodParameterLineBreaks`
-
`Layout/ParameterAlignment`
-
`Style/BlockDelimiters`

Together, these cops will pretty print hashes, arrays, method calls, etc. For example, let’s say the max columns is 25:

### Examples

```
# bad
{foo: "0000000000", bar: "0000000000", baz: "0000000000"}
# good
{foo: "0000000000",
bar: "0000000000", baz: "0000000000"}
# good (with recommended cops enabled)
{
 foo: "0000000000",
 bar: "0000000000",
 baz: "0000000000",
}
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| Max |
 | Integer |
| AllowHeredoc |
 | Boolean |
| AllowURI |
 | Boolean |
| AllowQualifiedName |
 | Boolean |
| URISchemes |
 | Array |
| AllowRBSInlineAnnotation |
 | Boolean |
| AllowCopDirectives |
 | Boolean |
| AllowedPatterns |
 | Array |
| SplitStrings |
 | Boolean |

## Layout/MultilineArrayBraceLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that the closing brace in an array literal is either on the same line as the last array element or on a new line.

When using the `symmetrical` (default) style:

If an array’s opening brace is on the same line as the first element of the array, then the closing brace should be on the same line as the last element of the array.

If an array’s opening brace is on the line above the first element of the array, then the closing brace should be on the line below the last element of the array.

When using the `new_line` style:

The closing brace of a multi-line array literal must be on the line after the last element of the array.

When using the `same_line` style:

The closing brace of a multi-line array literal must be on the same line as the last element of the array.

## Layout/MultilineArrayLineBreaks

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.67 | - |

Ensures that each item in a multi-line array starts on a separate line.

## Layout/MultilineAssignmentLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.49 | - |

Checks whether the multiline assignments have a newline after the assignment operator.

### Examples

#### EnforcedStyle: new_line (default)

```
# bad
foo = if expression
 'bar'
end
# good
foo =
 if expression
 'bar'
 end
# good
foo =
 begin
 compute
 rescue => e
 nil
 end
```
## Layout/MultilineBlockLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks whether the multiline do end blocks have a newline after the start of the block. Additionally, it checks whether the block arguments, if any, are on the same line as the start of the block. Putting block arguments on separate lines, because the whole line would otherwise be too long, is accepted.

## Layout/MultilineHashBraceLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that the closing brace in a hash literal is either on the same line as the last hash element, or a new line.

When using the `symmetrical` (default) style:

If a hash’s opening brace is on the same line as the first element of the hash, then the closing brace should be on the same line as the last element of the hash.

If a hash’s opening brace is on the line above the first element of the hash, then the closing brace should be on the line below the last element of the hash.

When using the `new_line` style:

The closing brace of a multi-line hash literal must be on the line after the last element of the hash.

When using the `same_line` style:

The closing brace of a multi-line hash literal must be on the same line as the last element of the hash.

## Layout/MultilineHashKeyLineBreaks

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.67 | - |

Ensures that each key in a multi-line hash starts on a separate line.

## Layout/MultilineMethodArgumentLineBreaks

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.67 | - |

Ensures that each argument in a multi-line method call starts on a separate line.

| This cop does not move the first argument, if you want that to
be on a separate line, see `Layout/FirstMethodArgumentLineBreak`. |

### Examples

```
# bad
foo(a, b,
 c
)
# bad
foo(a, b, {
 foo: "bar",
})
# good
foo(
 a,
 b,
 c
)
# good
foo(a, b, c)
```
## Layout/MultilineMethodCallBraceLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that the closing brace in a method call is either on the same line as the last method argument, or a new line.

When using the `symmetrical` (default) style:

If a method call’s opening brace is on the same line as the first argument of the call, then the closing brace should be on the same line as the last argument of the call.

If a method call’s opening brace is on the line above the first argument of the call, then the closing brace should be on the line below the last argument of the call.

When using the `new_line` style:

The closing brace of a multi-line method call must be on the line after the last argument of the call.

When using the `same_line` style:

The closing brace of a multi-line method call must be on the same line as the last argument of the call.

## Layout/MultilineMethodCallIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks the indentation of the method name part in method calls that span more than one line.

## Layout/MultilineMethodDefinitionBraceLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that the closing brace in a method definition is either on the same line as the last method parameter, or a new line.

When using the `symmetrical` (default) style:

If a method definition’s opening brace is on the same line as the first parameter of the definition, then the closing brace should be on the same line as the last parameter of the definition.

If a method definition’s opening brace is on the line above the first parameter of the definition, then the closing brace should be on the line below the last parameter of the definition.

When using the `new_line` style:

The closing brace of a multi-line method definition must be on the line after the last parameter of the definition.

When using the `same_line` style:

The closing brace of a multi-line method definition must be on the same line as the last parameter of the definition.

## Layout/MultilineMethodParameterLineBreaks

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 1.32 | - |

Ensures that each parameter in a multi-line method definition starts on a separate line.

| This cop does not move the first argument, if you want that to
be on a separate line, see `Layout/FirstMethodParameterLineBreak`. |

### Examples

```
# bad
def foo(a, b,
 c
)
end
# good
def foo(
 a,
 b,
 c
)
end
# good
def foo(
 a,
 b = {
 foo: "bar",
 }
)
end
# good
def foo(a, b, c)
end
```
## Layout/MultilineOperationIndentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks the indentation of the right hand side operand in binary operations that span more than one line.

The `aligned` style checks that operators are aligned if they are part of an `if` or `while`
condition, an explicit `return` statement, etc. In other contexts, the second operand should
be indented regardless of enforced style.

In both styles, operators should be aligned when an assignment begins on the next line.

## Layout/ParameterAlignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.77 |

Checks that the parameters on a multi-line method call or definition are aligned.

To set the alignment of the first argument, use the
`Layout/FirstParameterIndentation` cop.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| IndentationWidth |
 | Integer |

## Layout/RedundantLineBreak

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 1.13 | - |

Checks whether certain expressions, e.g. method calls, that could fit completely on a single line, are broken up into multiple lines unnecessarily.

## Layout/SingleLineBlockChain

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 1.14 | - |

Checks if method calls are chained onto single line blocks. It considers that a line break before the dot improves the readability of the code.

## Layout/SpaceAfterColon

## Layout/SpaceAfterComma

## Layout/SpaceAfterMethodName

## Layout/SpaceAfterNot

## Layout/SpaceAfterSemicolon

## Layout/SpaceAroundBlockParameters

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks the spacing inside and after block parameters pipes. Line breaks
inside parameter pipes are checked by `Layout/MultilineBlockLayout` and
not by this cop.

## Layout/SpaceAroundEqualsInParameterDefault

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that the equals signs in parameter default assignments have or don’t have surrounding space depending on configuration.

### Examples

## Layout/SpaceAroundKeyword

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks the spacing around the keywords.

## Layout/SpaceAroundMethodCallOperator

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.82 | - |

Checks method call operators to not have spaces around them.

## Layout/SpaceAroundOperators

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that operators have space around them, except for ** which should or shouldn’t have surrounding space depending on configuration. It allows vertical alignment consisting of one or more whitespace around operators.

This cop has `AllowForAlignment` option. When `true`, allows most
uses of extra spacing if the intent is to align with an operator on
the previous or next line, not counting empty lines or comment lines.

### Examples

```
# bad
total = 3*4
"apple"+"juice"
my_number = 38/4
# good
total = 3 * 4
"apple" + "juice"
my_number = 38 / 4
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowForAlignment |
 | Boolean |
| EnforcedStyleForExponentOperator |
 |
 |
| EnforcedStyleForRationalLiterals |
 |
 |

## Layout/SpaceBeforeBlockBraces

## Layout/SpaceBeforeBrackets

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.7 | - |

Checks for space between the name of a receiver and a left bracket.

## Layout/SpaceBeforeFirstArg

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that exactly one space is used between a method name and the first argument for method calls without parentheses.

Alternatively, extra spaces can be added to align the argument with something on a preceding or following line, if the AllowForAlignment config parameter is true.

## Layout/SpaceInLambdaLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for spaces between `→` and opening parameter
parenthesis (`(`) in lambda literals.

## Layout/SpaceInsideArrayLiteralBrackets

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks that brackets used for array literals have or don’t have surrounding space depending on configuration.

Array pattern matching is handled in the same way.

### Examples

#### EnforcedStyle: no_space (default)

```
# The `no_space` style enforces that array literals have
# no surrounding space.
# bad
array = [ a, b, c, d ]
array = [ a, [ b, c ]]
# good
array = [a, b, c, d]
array = [a, [b, c]]
```
#### EnforcedStyle: space

```
# The `space` style enforces that array literals have
# surrounding space.
# bad
array = [a, b, c, d]
array = [ a, [ b, c ]]
# good
array = [ a, b, c, d ]
array = [ a, [ b, c ] ]
```
#### EnforcedStyle: compact

```
# The `compact` style normally requires a space inside
# array brackets, with the exception that successive left
# or right brackets are collapsed together in nested arrays.
# bad
array = [a, b, c, d]
array = [ a, [ b, c ] ]
array = [
 [ a ],
 [ b, c ]
]
# good
array = [ a, b, c, d ]
array = [ a, [ b, c ]]
array = [[ a ],
 [ b, c ]]
```
## Layout/SpaceInsideArrayPercentLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for unnecessary additional spaces inside array percent literals (i.e. %i/%w).

Note that blank percent literals (e.g. `%i( )`) are checked by
`Layout/SpaceInsidePercentLiteralDelimiters`.

## Layout/SpaceInsideBlockBraces

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that block braces have or don’t have surrounding space inside them on configuration. For blocks taking parameters, it checks that the left brace has or doesn’t have trailing space depending on configuration.

### Examples

#### EnforcedStyle: space (default)

```
# The `space` style enforces that block braces have
# surrounding space.
# bad
some_array.each {puts e}
# good
some_array.each { puts e }
```
#### EnforcedStyle: no_space

```
# The `no_space` style enforces that block braces don't
# have surrounding space.
# bad
some_array.each { puts e }
# good
some_array.each {puts e}
```
#### EnforcedStyleForEmptyBraces: no_space (default)

```
# The `no_space` EnforcedStyleForEmptyBraces style enforces that
# block braces don't have a space in between when empty.
# bad
some_array.each { }
some_array.each { }
some_array.each { }
# good
some_array.each {}
```
#### EnforcedStyleForEmptyBraces: space

```
# The `space` EnforcedStyleForEmptyBraces style enforces that
# block braces have at least a space in between when empty.
# bad
some_array.each {}
# good
some_array.each { }
some_array.each { }
some_array.each { }
```
## Layout/SpaceInsideHashLiteralBraces

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks that braces used for hash literals have or don’t have surrounding space depending on configuration.

Hash pattern matching is handled in the same way.

### Examples

#### EnforcedStyle: space (default)

```
# The `space` style enforces that hash literals have
# surrounding space.
# bad
h = {a: 1, b: 2}
foo = {{ a: 1 } => { b: { c: 2 }}}
# good
h = { a: 1, b: 2 }
foo = { { a: 1 } => { b: { c: 2 } } }
```
#### EnforcedStyle: no_space

```
# The `no_space` style enforces that hash literals have
# no surrounding space.
# bad
h = { a: 1, b: 2 }
foo = {{ a: 1 } => { b: { c: 2 }}}
# good
h = {a: 1, b: 2}
foo = {{a: 1} => {b: {c: 2}}}
```
#### EnforcedStyle: compact

```
# The `compact` style normally requires a space inside
# hash braces, with the exception that successive left
# braces or right braces are collapsed together in nested hashes.
# bad
h = { a: { b: 2 } }
foo = { { a: 1 } => { b: { c: 2 } } }
# good
h = { a: { b: 2 }}
foo = {{ a: 1 } => { b: { c: 2 }}}
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| EnforcedStyleForEmptyBraces |
 |
 |

## Layout/SpaceInsideParens

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.22 |

Checks for spaces inside ordinary round parentheses.

### Examples

#### EnforcedStyle: no_space (default)

```
# The `no_space` style enforces that parentheses do not have spaces.
# bad
f( 3)
g = (a + 3 )
f( )
# good
f(3)
g = (a + 3)
f()
```
#### EnforcedStyle: space

```
# The `space` style enforces that parentheses have a space at the
# beginning and end.
# Note: Empty parentheses should not have spaces.
# bad
f(3)
g = (a + 3)
y( )
# good
f( 3 )
g = ( a + 3 )
y()
```
#### EnforcedStyle: compact

```
# The `compact` style enforces that parentheses have a space at the
# beginning with the exception that successive parentheses are allowed.
# Note: Empty parentheses should not have spaces.
# bad
f(3)
g = (a + 3)
y( )
g( f( x ) )
g( f( x( 3 ) ), 5 )
g( ( ( 3 + 5 ) * f) ** x, 5 )
# good
f( 3 )
g = ( a + 3 )
y()
g( f( x ))
g( f( x( 3 )), 5 )
g((( 3 + 5 ) * f ) ** x, 5 )
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Layout/SpaceInsidePercentLiteralDelimiters

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for unnecessary additional spaces inside the delimiters of %i/%w/%x literals.

## Layout/SpaceInsideRangeLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for spaces inside range literals.

## Layout/SpaceInsideReferenceBrackets

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | 0.53 |

Checks that reference brackets have or don’t have surrounding space depending on configuration.

### Examples

#### EnforcedStyle: no_space (default)

```
# The `no_space` style enforces that reference brackets have
# no surrounding space.
# bad
hash[ :key ]
array[ index ]
# good
hash[:key]
array[index]
```
#### EnforcedStyle: space

```
# The `space` style enforces that reference brackets have
# surrounding space.
# bad
hash[:key]
array[index]
# good
hash[ :key ]
array[ index ]
```
## Layout/SpaceInsideStringInterpolation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for whitespace within string interpolations.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Layout/TrailingEmptyLines

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.77 |

Looks for trailing blank lines and a final newline in the source code.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Layout/TrailingWhitespace

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.0 |

Looks for trailing whitespace in the source code.

### Examples

```
# The line in this example contains spaces after the 0.
# bad
x = 0
# The line in this example ends directly after the 0.
# good
x = 0
```
## AllowMultilineFinalElement

`AllowMultilineFinalElement` is a boolean option (`false` by default) that
is available on the following Layout cops:

-
FirstArrayElementLineBreak
-
FirstHashElementLineBreak
-
FirstMethodArgumentLineBreak
-
FirstMethodParameterLineBreak
-
MultilineArrayLineBreaks
-
MultilineHashKeyLineBreaks
-
MultilineMethodArgumentLineBreaks
-
MultilineMethodParameterLineBreaks

Those cops ignore their respective expressions if all of the elements
of the expression are on the same line. If `AllowMultilineFinalElement` is
set to `true`, the cop will also ignore multiline expressions if the last
element starts on the same line, but ends on a different one.

This works well with `Layout/LineLength` to present elements on
individual lines when wrapping, while not affecting expressions
that are already short enough, and only considered multiline
because of their last element.

Each cop can be configured independently allowing for more fine-grained control over what is considered ok in the codebase.

### Examples

Here are some examples of real world expressions that get wrapped
by their respective cops, but that are considered ok when setting
`AllowMultilineFinalElement` to `true` on those same cops.

```
# good
# FirstArrayElementLineBreak and MultilineArrayLineBreaks
# Array of error containing a single error
errors = [{
 error: "Something went wrong",
 error_code: error_code,
}]
# Array of flags, with last flag computed
flags = [:a, :b, foo(
 bar,
 baz
)]
# FirstHashElementLineBreak and MultilineHashKeyLineBreaks
hash = { foo: 1, bar: 2, baz: {
 c: 1,
 d: 2
}}
# FirstMethodArgumentLineBreak and MultilineMethodArgumentLineBreaks
single_argument_hash_method_call({
 a: 1,
 b: 2,
 c: 3
})
# Call some method, with a long last argument, that is a hash
write_log(:error, {
 "job_class" => job.class.name,
 "resource" => resource.id,
 "message" => "Something wrong happened here",
})
# Rails before action with long last argument
before_action :load_something, only: [
 :show,
 :list,
 :some_other_long_action_name,
]
# Rails validation with inline callback
validate :name, presence: true, on: [:create, :activate], if: -> {
 active? && some_relationship.any?
}
# Rails after commit hook, with some Sorbet bindings
after_commit :geolocate, unless: -> {
 T.bind(self, Address)
 archived? || invalid?
}
# FirstMethodParameterLineBreak and MultilineMethodParameterLineBreaks
# Method with a long last parameter default value
def foo(foo, bar, baz = {
 a: 1,
 b: 2,
 c: 3
})
 do_something
end
```

# Lint

## Lint/AmbiguousBlockAssociation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.48 | 1.89 |

Checks for ambiguous block association with method when param passed without parentheses.

This cop also detects `do…end` blocks that are likely intended for
an enumerable method in the arguments but actually bind to the outer
method call. For example, in `render json: data.map do |x| x end`,
Ruby parses the `do…end` block as belonging to `render`, not `map`.

This cop can customize allowed methods with `AllowedMethods`.
By default, there are no allowed methods.

### Examples

```
# bad
some_method a { |val| puts val }
# good
# With parentheses, there's no ambiguity.
some_method(a { |val| puts val })
# or (different meaning)
some_method(a) { |val| puts val }
# bad
render json: data.map do |item|
 item.to_h
end
# good
render json: data.map { |item| item.to_h }
# good
mapped = data.map { |item| item.to_h }
render json: mapped
# good
# Operator methods require no disambiguation
foo == bar { |b| b.baz }
# good
# Lambda arguments require no disambiguation
foo = ->(bar) { bar.baz }
```
## Lint/AmbiguousOperator

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.17 | 0.83 |

Checks for ambiguous operators in the first argument of a method invocation without parentheses.

## Lint/AmbiguousOperatorPrecedence

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.21 | - |

Looks for expressions containing multiple binary operators
where precedence is ambiguous due to lack of parentheses. For example,
in `1 + 2 * 3`, the multiplication will happen before the addition, but
lexically it appears that the addition will happen first.

The cop does not consider unary operators (ie. `!a` or `-b`) or comparison
operators (ie. `a =~ b`) because those are not ambiguous.

| Ranges are handled by `Lint/AmbiguousRange`. |

## Lint/AmbiguousRange

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.19 | - |

Checks for ambiguous ranges.

Ranges have quite low precedence, which leads to unexpected behavior when using a range with other operators. This cop avoids that by making ranges explicit by requiring parenthesis around complex range boundaries (anything that is not a literal: numerics, strings, symbols, etc.).

This cop can be configured with `RequireParenthesesForMethodChains` in order to
specify whether method chains (including `self.foo`) should be wrapped in parens
by this cop.

| Regardless of this configuration, if a method receiver is a basic literal
value, it will be wrapped in order to prevent the ambiguity of `1..2.to_a`. |

### Safety

The cop autocorrects by wrapping the entire boundary in parentheses, which makes the outcome more explicit but is possible to not be the intention of the programmer. For this reason, this cop’s autocorrect is unsafe (it will not change the behavior of the code, but will not necessarily match the intent of the program).

## Lint/AmbiguousRegexpLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.17 | 0.83 |

Checks for ambiguous regexp literals in the first argument of a method invocation without parentheses.

## Lint/ArrayLiteralInRegexp

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.71 | - |

Checks for an array literal interpolated inside a regexp.

When interpolating an array literal, it is converted to a string. This means
that when inside a regexp, it acts as a character class but with additional
quotes, spaces and commas that are likely not intended. For example,
`/#{%w[a b c]}/` parses as `/["a", "b", "c"]/` (or `/["a, bc]/` without
repeated characters).

The cop can autocorrect to a character class (if all items in the array are a single character) or alternation (if the array contains longer items).

| This only considers interpolated arrays that contain only strings, symbols, integers, and floats. Any other type is not easily convertible to a character class or regexp alternation. |

## Lint/AssignmentInCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.9 | 1.45 |

Checks for assignments in the conditions of if/while/until.

`AllowSafeAssignment` option for safe assignment.
By safe assignment we mean putting parentheses around
an assignment to indicate "I know I’m using an assignment
as a condition. It’s not a mistake."

### Safety

This cop’s autocorrection is unsafe because it assumes that the author meant to use an assignment result as a condition.

## Lint/BigDecimalNew

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

`BigDecimal.new()` is deprecated since BigDecimal 1.3.3.
This cop identifies places where `BigDecimal.new()`
can be replaced by `BigDecimal()`.

## Lint/BinaryOperatorWithIdenticalOperands

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | No | 0.89 | 1.69 |

Checks for places where binary operator has identical operands.

It covers comparison operators: `==`, `===`, `=~`, `>`, `>=`, `<`, `⇐`;
bitwise operators: `|`, `^`, `&`;
boolean operators: `&&`, `||`
and "spaceship" operator - `<⇒`.

Simple arithmetic operations are allowed by this cop: `+`, `, ``*`, `<<` and `>>`.
Although these can be rewritten in a different way, it should not be necessary to
do so. Operations such as `-` or `/` where the result will always be the same
(`x - x` will always be 0; `x / x` will always be 1) are offenses, but these
are covered by `Lint/NumericOperationWithConstantResult` instead.

## Lint/BooleanSymbol

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.50 | 1.22 |

Checks for `:true` and `:false` symbols.
In most cases it would be a typo.

## Lint/CircularArgumentReference

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.33 | - |

Checks for circular argument references in optional keyword arguments and optional ordinal arguments.

| This syntax was made invalid on Ruby 2.7 - Ruby 3.3 but is allowed again since Ruby 3.4. |

### Examples

```
# bad
def bake(pie: pie)
 pie.heat_up
end
# good
def bake(pie:)
 pie.refrigerate
end
# good
def bake(pie: self.pie)
 pie.feed_to(user)
end
# bad
def cook(dry_ingredients = dry_ingredients)
 dry_ingredients.reduce(&:+)
end
# good
def cook(dry_ingredients = self.dry_ingredients)
 dry_ingredients.combine
end
# bad
def foo(pie = pie = pie)
 pie.heat_up
end
# good
def foo(pie)
 pie.heat_up
end
# bad
def foo(pie = cake = pie)
 [pie, cake].each(&:heat_up)
end
# good
def foo(cake = pie)
 [pie, cake].each(&:heat_up)
end
```
## Lint/ConstantDefinitionInBlock

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.91 | 1.3 |

Do not define constants within a block, since the block’s scope does not isolate or namespace the constant in any way.

If you are trying to define that constant once, define it outside of the block instead, or use a variable or method if defining the constant in the outer scope would be problematic.

For meta-programming, use `const_set`.

### Examples

```
# bad
task :lint do
 FILES_TO_LINT = Dir['lib/*.rb']
end
# bad
describe 'making a request' do
 class TestRequest; end
end
# bad
module M
 extend ActiveSupport::Concern
 included do
 LIST = []
 end
end
# good
task :lint do
 files_to_lint = Dir['lib/*.rb']
end
# good
describe 'making a request' do
 let(:test_request) { Class.new }
 # see also `stub_const` for RSpec
end
# good
module M
 extend ActiveSupport::Concern
 included do
 const_set(:LIST, [])
 end
end
```
## Lint/ConstantOverwrittenInRescue

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.31 | - |

Checks for overwriting an exception with an exception result by using `rescue ⇒`.

You intended to write as `rescue StandardError`.
However, you have written `rescue ⇒ StandardError`.
In that case, the result of `rescue` will overwrite `StandardError`.

## Lint/ConstantReassignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.70 | 1.87 |

Checks for constant reassignments.

Emulates Ruby’s runtime warning "already initialized constant X" when a constant is reassigned in the same file and namespace.

The cop tracks constants defined via `NAME = value` syntax as well as
class/module keyword definitions. It detects reassignment when a constant
is first defined one way and then redefined using the `NAME = value` syntax.

The cop cannot catch all offenses, like, for example, when using metaprogramming
(`Module#const_set`).

By default the cop also cannot detect reassignment across files.
When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed,
the cop additionally consults the project-wide index and reports reassignments
whose previous definition lives in another file.

The cop only takes into account constants assigned in a "simple" way: directly inside class/module definition, or within another constant. Other type of assignments (e.g., inside a conditional) are disregarded.

The cop also tracks constant removal using `Module#remove_const` with symbol
or string argument.

### Examples

```
# bad
X = :foo
X = :bar
# bad
class A
 X = :foo
 X = :bar
end
# bad
module A
 X = :foo
 X = :bar
end
# bad
class FooError < StandardError; end
FooError = Class.new(RuntimeError)
# bad
module M; end
M = 1
# good - keep only one assignment
X = :bar
class A
 X = :bar
end
module A
 X = :bar
end
# good - use OR assignment
X = :foo
X ||= :bar
# good - use conditional assignment
X = :foo
X = :bar unless defined?(X)
# good - remove the assigned constant first
class A
 X = :foo
 remove_const :X
 X = :bar
end
```
## Lint/ConstantResolution

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.86 | 1.89 |

Checks that certain constants are fully qualified.

This is not enabled by default because it would mark a lot of offenses unnecessarily.

Generally, gems should fully qualify all constants to avoid conflicts with
the code that uses the gem. Enable this cop without using `Only`/`Ignore`

Large projects will over time end up with one or two constant names that
are problematic because of a conflict with a library or just internally
using the same name for a namespace and a class. To avoid too many unnecessary
offenses, enable this cop with `Only: [The, Constant, Names, Causing, Issues]`

| `Style/RedundantConstantBase`cop is disabled if this cop is enabled,
to prevent conflicting rules. This is because it respects user configurations
that want to enable this cop which is disabled by default. |

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, only genuinely ambiguous constants are reported: those
that resolve to a different declaration through the surrounding
nesting than they would fully qualified. This makes the cop practical
to enable without `Only`/`Ignore` lists.

## Lint/CopDirectiveSyntax

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.72 | - |

Checks that `# rubocop:enable …` and `# rubocop:disable …` statements
are strictly formatted.

A comment can be added to the directive by prefixing it with `--`.

### Examples

```
# bad
# rubocop:disable Layout/LineLength Style/Encoding
# good
# rubocop:disable Layout/LineLength, Style/Encoding
# bad
# rubocop:disable
# good
# rubocop:disable all
# bad
# rubocop:disable Layout/LineLength # rubocop:disable Style/Encoding
# good
# rubocop:disable Layout/LineLength
# rubocop:disable Style/Encoding
# bad
# rubocop:wrongmode Layout/LineLength
# good
# rubocop:disable Layout/LineLength
# bad
# rubocop:disable Layout/LineLength comment
# good
# rubocop:disable Layout/LineLength -- comment
```
## Lint/DataDefineOverride

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.85 | - |

Checks unexpected overrides of the `Data` built-in methods
via `Data.define`.

### Examples

```
# bad
Bad = Data.define(:members, :clone, :to_s)
b = Bad.new(members: [], clone: true, to_s: 'bad')
b.members #=> [] (overriding `Data#members`)
b.clone #=> true (overriding `Object#clone`)
b.to_s #=> "bad" (overriding `Data#to_s`)
# good
Good = Data.define(:id, :name)
g = Good.new(id: 1, name: "foo")
g.members #=> [:id, :name]
g.clone #=> #<data Good id=1, name="foo">
```
## Lint/Debugger

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.14 | 1.63 |

Checks for debug calls (such as `debugger` or `binding.pry`) that should
not be kept for production code.

The cop can be configured using `DebuggerMethods`. By default, a number of gems
debug entrypoints are configured (`Kernel`, `Byebug`, `Capybara`, `debug.rb`,
`Pry`, `Rails`, `RubyJard`, and `WebConsole`). Additional methods can be added.

Specific default groups can be disabled if necessary:

```
Lint/Debugger:
 DebuggerMethods:
 WebConsole: ~
```
You can also add your own methods by adding a new category:

```
Lint/Debugger:
 DebuggerMethods:
 MyDebugger:
 MyDebugger.debug_this
```
Some gems also ship files that will start a debugging session when required,
for example `require 'debug/start'` from `ruby/debug`. These requires can
be configured through `DebuggerRequires`. It has the same structure as
`DebuggerMethods`, which you can read about above.

### Examples

```
# bad (ok during development)
# using pry
def some_method
 binding.pry
 do_something
end
# bad (ok during development)
# using byebug
def some_method
 byebug
 do_something
end
# good
def some_method
 do_something
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| DebuggerMethods |
 | |
| DebuggerRequires |
 |

## Lint/DeprecatedClassMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.19 | - |

Checks for uses of the deprecated class method usages.

### Examples

```
# bad
File.exists?(some_path)
Dir.exists?(some_path)
iterator?
attr :name, true
attr :name, false
ENV.freeze # Calling `Env.freeze` raises `TypeError` since Ruby 2.7.
ENV.clone
ENV.dup # Calling `Env.dup` raises `TypeError` since Ruby 3.1.
Socket.gethostbyname(host)
Socket.gethostbyaddr(host)
# good
File.exist?(some_path)
Dir.exist?(some_path)
block_given?
attr_accessor :name
attr_reader :name
ENV # `ENV.freeze` cannot prohibit changes to environment variables.
ENV.to_h
ENV.to_h # `ENV.dup` cannot dup `ENV`, use `ENV.to_h` to get a copy of `ENV` as a hash.
Addrinfo.getaddrinfo(nodename, service)
Addrinfo.tcp(host, port).getnameinfo
```
## Lint/DeprecatedConstants

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.8 | 1.40 |

Checks for deprecated constants.

It has `DeprecatedConstants` config. If there is an alternative method, you can set
alternative value as `Alternative`. And you can set the deprecated version as
`DeprecatedVersion`. These options can be omitted if they are not needed.

```
DeprecatedConstants:
 'DEPRECATED_CONSTANT':
 Alternative: 'alternative_value'
 DeprecatedVersion: 'deprecated_version'
```
By default, `NIL`, `TRUE`, `FALSE`, `Net::HTTPServerException`, `Random::DEFAULT`,
`Struct::Group`, and `Struct::Passwd` are configured.

### Examples

```
# bad
NIL
TRUE
FALSE
Net::HTTPServerException
Random::DEFAULT # Return value of Ruby 2 is `Random` instance, Ruby 3.0 is `Random` class.
Struct::Group
Struct::Passwd
# good
nil
true
false
Net::HTTPClientException
Random.new # `::DEFAULT` has been deprecated in Ruby 3, `.new` is compatible with Ruby 2.
Etc::Group
Etc::Passwd
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| DeprecatedConstants |
 |

## Lint/DeprecatedOpenSSLConstant

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.84 | - |

Algorithmic constants for `OpenSSL::Cipher` and `OpenSSL::Digest`
deprecated since OpenSSL version 2.2.0. Prefer passing a string
instead.

## Lint/DeprecatedReference

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.89 | - |

Checks for calls to methods and references to constants that are documented
as deprecated with a YARD `@deprecated` tag.

The check is powered by the project-wide index, so it only runs when
`AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed.
Without the index the cop does nothing.

Only references that can be resolved without type inference are checked:
constants, method calls without an explicit receiver (or with `self`),
which are looked up in the enclosing class or module and its ancestry,
and calls whose receiver is a constant, which are looked up in that
namespace’s singleton class. Calls on arbitrary objects are not checked.

References made from a definition that is itself deprecated are allowed, so deprecated implementations can keep calling each other.

### Examples

```
# Given a deprecated method and constant:
#
# class Api
# # @deprecated Use `#new_method` instead.
# def old_method
# end
#
# # @deprecated
# OLD_TIMEOUT = 10
# end
# bad
class Client < Api
 def call
 old_method
 end
 def timeout
 OLD_TIMEOUT
 end
end
# good
class Client < Api
 def call
 new_method
 end
 def timeout
 NEW_TIMEOUT
 end
end
```
## Lint/DisjunctiveAssignmentInConstructor

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.62 | 0.88 |

Checks constructors for disjunctive assignments (`||=`) that should
be plain assignments.

So far, this cop is only concerned with disjunctive assignment of instance variables.

In ruby, an instance variable is nil until a value is assigned, so the disjunction is unnecessary. A plain assignment has the same effect.

### Safety

This cop is unsafe because it can register a false positive when a method is redefined in a subclass that calls super. For example:

```
class Base
 def initialize
 @config ||= 'base'
 end
end
class Derived < Base
 def initialize
 @config = 'derived'
 super
 end
end
```
Without the disjunctive assignment, `Derived` will be unable to override
the value for `@config`.

## Lint/DuplicateBranch

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.3 | 1.7 |

Checks that there are no repeated bodies
within `if/unless`, `case-when`, `case-in` and `rescue` constructs.

With `IgnoreLiteralBranches: true`, branches are not registered
as offenses if they return a basic literal value (string, symbol,
integer, float, rational, complex, `true`, `false`, or `nil`), or
return an array, hash, regexp or range that only contains one of
the above basic literal values.

With `IgnoreConstantBranches: true`, branches are not registered
as offenses if they return a constant value.

With `IgnoreDuplicateElseBranch: true`, in conditionals with multiple branches,
duplicate 'else' branches are not registered as offenses.

### Examples

```
# bad
if foo
 do_foo
 do_something_else
elsif bar
 do_foo
 do_something_else
end
# good
if foo || bar
 do_foo
 do_something_else
end
# bad
case x
when foo
 do_foo
when bar
 do_foo
else
 do_something_else
end
# good
case x
when foo, bar
 do_foo
else
 do_something_else
end
# bad
begin
 do_something
rescue FooError
 handle_error
rescue BarError
 handle_error
end
# good
begin
 do_something
rescue FooError, BarError
 handle_error
end
```
#### IgnoreLiteralBranches: true

```
# good
case size
when "small" then 100
when "medium" then 250
when "large" then 1000
else 250
end
```
## Lint/DuplicateCaseCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.45 | - |

Checks that there are no repeated conditions used in case 'when' expressions.

## Lint/DuplicateElsifCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.88 | - |

Checks that there are no repeated conditions used in if 'elsif'.

## Lint/DuplicateHashKey

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.34 | 0.77 |

Checks for duplicated keys in hash literals. This cop considers both primitive types and constants for the hash keys.

This cop mirrors a warning in Ruby 2.2.

## Lint/DuplicateMagicComment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.37 | - |

Checks for duplicated magic comments.

## Lint/DuplicateMatchPattern

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.50 | - |

Checks that there are no repeated patterns used in `in` keywords.

### Examples

```
# bad
case x
in 'first'
 do_something
in 'first'
 do_something_else
end
# good
case x
in 'first'
 do_something
in 'second'
 do_something_else
end
# bad - repeated alternate patterns with the same conditions don't depend on the order
case x
in 0 | 1
 first_method
in 1 | 0
 second_method
end
# good
case x
in 0 | 1
 first_method
in 2 | 3
 second_method
end
# bad - repeated hash patterns with the same conditions don't depend on the order
case x
in foo: a, bar: b
 first_method
in bar: b, foo: a
 second_method
end
# good
case x
in foo: a, bar: b
 first_method
in bar: b, baz: c
 second_method
end
# bad - repeated array patterns with elements in the same order
case x
in [foo, bar]
 first_method
in [foo, bar]
 second_method
end
# good
case x
in [foo, bar]
 first_method
in [bar, foo]
 second_method
end
# bad - repeated the same patterns and guard conditions
case x
in foo if bar
 first_method
in foo if bar
 second_method
end
# good
case x
in foo if bar
 first_method
in foo if baz
 second_method
end
```
## Lint/DuplicateMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.29 | 1.89 |

Checks for duplicated instance (or singleton) method definitions.

| Aliasing a method to itself is allowed, as it indicates that the developer intends to suppress Ruby’s method redefinition warnings. See https://bugs.ruby-lang.org/issues/13574. |

By default the cop can only detect duplicates within a single file.
When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed,
the cop additionally consults the project-wide index and reports methods
whose duplicate definition lives in another file.

| The project index does not record whether a definition in another file is wrapped in a conditional, so a platform-specific redefinition in another file may still be reported. Aliasing the method to itself (see above) before redefining marks the redefinition as intentional and is respected across files. |

### Examples

```
# bad
def foo
 1
end
def foo
 2
end
# bad
def foo
 1
end
alias foo bar
# good
def foo
 1
end
def bar
 2
end
# good
def foo
 1
end
alias bar foo
# good
alias foo foo
def foo
 1
end
# good
alias_method :foo, :foo
def foo
 1
end
# bad
class MyClass
 extend Forwardable
 # or with: `def_instance_delegator`, `def_delegators`, `def_instance_delegators`
 def_delegator :delegation_target, :delegated_method_name
 def delegated_method_name
 end
end
# good
class MyClass
 extend Forwardable
 def_delegator :delegation_target, :delegated_method_name
 def non_duplicated_delegated_method_name
 end
end
```
#### AllCops:ActiveSupportExtensionsEnabled: false (default)

```
# good
def foo
 1
end
delegate :foo, to: :bar
```
#### AllCops:ActiveSupportExtensionsEnabled: true

```
# bad
def foo
 1
end
delegate :foo, to: :bar
# good
def foo
 1
end
delegate :baz, to: :bar
# good - delegate with splat arguments is ignored
def foo
 1
end
delegate :foo, **options
# good - delegate inside a condition is ignored
def foo
 1
end
if cond
 delegate :foo, to: :bar
end
```
## Lint/DuplicateRequire

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.90 | 1.28 |

Checks for duplicate `require`s and `require_relative`s.

## Lint/DuplicateRescueException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.89 | - |

Checks that there are no repeated exceptions used in 'rescue' expressions.

## Lint/DuplicateSetElement

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.67 | - |

Checks for duplicate literal, constant, or variable elements in `Set` and `SortedSet`.

### Examples

```
# bad
Set[:foo, :bar, :foo]
# good
Set[:foo, :bar]
# bad
Set.new([:foo, :bar, :foo])
# good
Set.new([:foo, :bar])
# bad
[:foo, :bar, :foo].to_set
# good
[:foo, :bar].to_set
# bad
SortedSet[:foo, :bar, :foo]
# good
SortedSet[:foo, :bar]
# bad
SortedSet.new([:foo, :bar, :foo])
# good
SortedSet.new([:foo, :bar])
```
## Lint/EachWithObjectArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.31 | - |

Checks if each_with_object is called with an immutable argument. Since the argument is the object that the given block shall make calls on to build something based on the enumerable that each_with_object iterates over, an immutable argument makes no sense. It’s definitely a bug.

## Lint/ElseLayout

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.17 | 1.2 |

Checks for odd `else` block layout - like
having an expression on the same line as the `else` keyword,
which is usually a mistake.

Its autocorrection tweaks layout to keep the syntax. So, this autocorrection
is compatible correction for bad case syntax, but if your code makes a mistake
with `elsif` and `else`, you will have to correct it manually.

### Examples

```
# bad
if something
 # ...
else do_this
 do_that
end
# good
# This code is compatible with the bad case. It will be autocorrected like this.
if something
 # ...
else
 do_this
 do_that
end
# This code is incompatible with the bad case.
# If `do_this` is a condition, `elsif` should be used instead of `else`.
if something
 # ...
elsif do_this
 do_that
end
# bad
# For single-line conditionals using `then` the layout is disallowed
# when the `else` body is multiline because it is treated as a lint offense.
if something then on_the_same_line_as_then
else first_line
 second_line
end
# good
# For single-line conditional using `then` the layout is allowed
# when `else` body is a single-line because it is treated as intentional.
if something then on_the_same_line_as_then
else single_line
end
```
## Lint/EmptyBlock

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.1 | 1.15 |

Checks for blocks without a body. Such empty blocks are typically an oversight or we should provide a comment to clarify what we’re aiming for.

Empty lambdas and procs are ignored by default.

| For backwards compatibility, the configuration that allows/disallows
empty lambdas and procs is called `AllowEmptyLambdas`, even though it also
applies to procs. |

### Examples

```
# bad
items.each { |item| }
# good
items.each { |item| puts item }
```
#### AllowComments: true (default)

```
# good
items.each do |item|
 # TODO: implement later (inner comment)
end
items.each { |item| } # TODO: implement later (inline comment)
```
#### AllowComments: false

```
# bad
items.each do |item|
 # TODO: implement later (inner comment)
end
items.each { |item| } # TODO: implement later (inline comment)
```
## Lint/EmptyClass

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.3 | - |

Checks for classes and metaclasses without a body. Such empty classes and metaclasses are typically an oversight or we should provide a comment to be clearer what we’re aiming for.

### Examples

```
# bad
class Foo
end
class Bar
 class << self
 end
end
class << obj
end
# good
class Foo
 def do_something
 # ... code
 end
end
class Bar
 class << self
 attr_reader :bar
 end
end
class << obj
 attr_reader :bar
end
```
## Lint/EmptyConditionalBody

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.89 | 1.73 |

Checks for the presence of `if`, `elsif` and `unless` branches without a body.

| empty `else`branches are handled by`Style/EmptyElse`. |

### Examples

```
# bad
if condition
end
# bad
unless condition
end
# bad
if condition
 do_something
elsif other_condition
end
# good
if condition
 do_something
end
# good
unless condition
 do_something
end
# good
if condition
 do_something
elsif other_condition
 nil
end
# good
if condition
 do_something
elsif other_condition
 do_something_else
end
```
## Lint/EmptyEnsure

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.10 | 1.61 |

Checks for empty `ensure` blocks.

## Lint/EmptyFile

## Lint/EmptyInPattern

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.16 | - |

Checks for the presence of `in` pattern branches without a body.

## Lint/EmptyWhen

## Lint/EnsureReturn

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.9 | 0.83 |

Checks for `return` from an `ensure` block.
`return` from an ensure block is a dangerous code smell as it
will take precedence over any exception being raised,
and the exception will be silently thrown away as if it were rescued.

If you want to rescue some (or all) exceptions, best to do it explicitly

### Examples

```
# bad
def foo
 do_something
ensure
 cleanup
 return self
end
# good
def foo
 do_something
 self
ensure
 cleanup
end
# good
def foo
 begin
 do_something
 rescue SomeException
 # Let's ignore this exception
 end
 self
ensure
 cleanup
end
```
## Lint/ErbNewArguments

| Requires Ruby version 2.6 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.56 | - |

Emulates the following Ruby warnings in Ruby 2.6.

```
$ cat example.rb
ERB.new('hi', nil, '-', '@output_buffer')
$ ruby -rerb example.rb
example.rb:1: warning: Passing safe_level with the 2nd argument of ERB.new is
deprecated. Do not use it, and specify other arguments as keyword arguments.
example.rb:1: warning: Passing trim_mode with the 3rd argument of ERB.new is
deprecated. Use keyword argument like ERB.new(str, trim_mode:...) instead.
example.rb:1: warning: Passing eoutvar with the 4th argument of ERB.new is
deprecated. Use keyword argument like ERB.new(str, eoutvar: ...) instead.
```
Now non-keyword arguments other than first one are softly deprecated
and will be removed when Ruby 2.5 becomes EOL.
`ERB.new` with non-keyword arguments is deprecated since ERB 2.2.0.
Use `:trim_mode` and `:eoutvar` keyword arguments to `ERB.new`.
This cop identifies places where `ERB.new(str, trim_mode, eoutvar)` can
be replaced by `ERB.new(str, trim_mode: trim_mode, eoutvar: eoutvar)`.

### Examples

```
# Target codes supports Ruby 2.6 and higher only
# bad
ERB.new(str, nil, '-', '@output_buffer')
# good
ERB.new(str, trim_mode: '-', eoutvar: '@output_buffer')
# Target codes supports Ruby 2.5 and lower only
# good
ERB.new(str, nil, '-', '@output_buffer')
# Target codes supports Ruby 2.6, 2.5 and lower
# bad
ERB.new(str, nil, '-', '@output_buffer')
# good
# Ruby standard library style
# https://github.com/ruby/ruby/commit/3406c5d
if ERB.instance_method(:initialize).parameters.assoc(:key) # Ruby 2.6+
 ERB.new(str, trim_mode: '-', eoutvar: '@output_buffer')
else
 ERB.new(str, nil, '-', '@output_buffer')
end
# good
# Use `RUBY_VERSION` style
if RUBY_VERSION >= '2.6'
 ERB.new(str, trim_mode: '-', eoutvar: '@output_buffer')
else
 ERB.new(str, nil, '-', '@output_buffer')
end
```
## Lint/FlipFlop

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.16 | - |

Looks for uses of flip-flop operator based on the Ruby Style Guide.

Here is the history of flip-flops in Ruby. flip-flop operator is deprecated in Ruby 2.6.0 and the deprecation has been reverted by Ruby 2.7.0 and backported to Ruby 2.6. See: https://bugs.ruby-lang.org/issues/5400

### Examples

```
# bad
(1..20).each do |x|
 puts x if (x == 5) .. (x == 10)
end
# good
(1..20).each do |x|
 puts x if (x >= 5) && (x <= 10)
end
```
## Lint/FloatComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.89 | - |

Checks for the presence of precise comparison of floating point numbers.

Floating point values are inherently inaccurate, and comparing them for exact equality
is almost never the desired semantics. Comparison via the `==/!=` operators checks
floating-point value representation to be exactly the same, which is very unlikely
if you perform any arithmetic operations involving precision loss.

### Examples

```
# bad
x == 0.1
x != 0.1
# bad
case value
when 1.0
 foo
when 2.0
 bar
end
# good - using BigDecimal
x.to_d == 0.1.to_d
# good - comparing against zero
x == 0.0
x != 0.0
# good
(x - 0.1).abs < Float::EPSILON
# good
tolerance = 0.0001
(x - 0.1).abs < tolerance
# good - comparing against nil
Float(x, exception: false) == nil
# good - using epsilon comparison in case expression
case
when (value - 1.0).abs < Float::EPSILON
 foo
when (value - 2.0).abs < Float::EPSILON
 bar
end
# Or some other epsilon based type of comparison:
# https://www.embeddeduse.com/2019/08/26/qt-compare-two-floats/
```
## Lint/FloatOutOfRange

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.36 | - |

Identifies `Float` literals which are, like, really really really
really really really really really big. Too big. No-one needs Floats
that big. If you need a float that big, something is wrong with you.

## Lint/FormatParameterMismatch

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.33 | - |

This lint sees if there is a mismatch between the number of expected fields for format/sprintf/#% and what is actually passed as arguments.

In addition, it checks whether different formats are used in the same format string. Do not mix numbered, unnumbered, and named formats in the same format string.

## Lint/HashCompareByIdentity

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | No | 0.93 | - |

Prefer using `Hash#compare_by_identity` rather than using `object_id`
for hash keys.

This cop looks for hashes being keyed by objects' `object_id`, using
one of these methods: `key?`, `has_key?`, `fetch`, `[]` and `[]=`.

### Safety

This cop is unsafe. Although unlikely, the hash could store both object ids and other values that need be compared by value, and thus could be a false positive.

Furthermore, this cop cannot guarantee that the receiver of one of the
methods (`key?`, etc.) is actually a hash.

### Examples

```
# bad
hash = {}
hash[foo.object_id] = :bar
hash.key?(baz.object_id)
# good
hash = {}.compare_by_identity
hash[foo] = :bar
hash.key?(baz)
```
## Lint/HashNewWithKeywordArgumentsAsDefault

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.69 | - |

Checks for the deprecated use of keyword arguments as a default in `Hash.new`.

This usage raises a warning in Ruby 3.3 and results in an error in Ruby 3.4. In Ruby 3.4, keyword arguments will instead be used to change the behavior of a hash. For example, the capacity option can be passed to create a hash with a certain size if you know it in advance, for better performance.

| The following corner case may result in a false negative when upgrading from Ruby 3.3 or earlier, but it is intentionally not detected to respect the expected usage in Ruby 3.4. |

`Hash.new(capacity: 42)`## Lint/HeredocMethodCallPosition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.68 | - |

Checks for the ordering of a method call where the receiver of the call is a HEREDOC.

### Examples

```
# bad
<<-SQL
 bar
SQL
.strip_indent
<<-SQL
 bar
SQL
.strip_indent
.trim
# good
<<~SQL
 bar
SQL
<<~SQL.trim
 bar
SQL
```
## Lint/IdentityComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.91 | - |

Prefer `equal?` over `==` when comparing `object_id`.

`Object#equal?` is provided to compare objects for identity, and in contrast
`Object#==` is provided for the purpose of doing value comparison.

### Examples

```
# bad
foo.object_id == bar.object_id
foo.object_id != baz.object_id
# good
foo.equal?(bar)
!foo.equal?(baz)
```
## Lint/ImplicitStringConcatenation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.36 | - |

Checks for implicit string concatenation of string literals which are on the same line.

## Lint/IncompatibleIoSelectWithFiberScheduler

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.21 | 1.24 |

Checks for `IO.select` that is incompatible with Fiber Scheduler since Ruby 3.0.

When an array of IO objects waiting for an exception (the third argument of `IO.select`)
is used as an argument, there is no alternative API, so offenses are not registered.

## Lint/IneffectiveAccessModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.36 | - |

Checks for `private` or `protected` access modifiers which are
applied to a singleton method. These access modifiers do not make
singleton methods private/protected. `private_class_method` can be
used for that.

## Lint/InheritException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.41 | 1.89 |

Looks for error classes inheriting from `Exception`.
It is configurable to suggest using either `StandardError` (default) or
`RuntimeError` instead.

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, indirect inheritance is also detected: a class whose parent
(defined anywhere in the project) ultimately inherits from `Exception`
is reported, without autocorrection.

### Safety

This cop’s autocorrection is unsafe because `rescue` that omit
exception class handle `StandardError` and its subclasses,
but not `Exception` and its subclasses.

## Lint/InterpolationCheck

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.50 | 1.40 |

Checks for interpolation in a single quoted string.

## Lint/ItWithoutArgumentsInBlock

| Requires Ruby version ⇐ 3.3 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.59 | - |

Emulates the following Ruby warning in Ruby 3.3.

```
$ ruby -e '0.times { it }'
-e:1: warning: `it` calls without arguments will refer to the first block param in Ruby 3.4;
use it() or self.it
```
`it` calls without arguments will refer to the first block param in Ruby 3.4.
So use `it()` or `self.it` to ensure compatibility.

## Lint/LambdaWithoutLiteralBlock

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.8 | - |

Checks uses of lambda without a literal block. It emulates the following warning in Ruby 3.0:

```
$ ruby -vwe 'lambda(&proc {})'
ruby 3.0.0p0 (2020-12-25 revision 95aff21468) [x86_64-darwin19]
-e:1: warning: lambda without a literal block is deprecated; use the proc without
lambda instead
```
This way, proc object is never converted to lambda. Autocorrection replaces with compatible proc argument.

## Lint/LiteralAsCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.51 | - |

Checks for literals used as the conditions or as operands in and/or expressions serving as the conditions of if/while/until/case-when/case-in.

| Literals in `case-in`condition where the match variable is used in`in`are accepted as a pattern matching. |

### Examples

```
# bad
if 20
 do_something
end
# bad
# We're only interested in the left hand side being a truthy literal,
# because it affects the evaluation of the &&, whereas the right hand
# side will be conditionally executed/called and can be a literal.
if true && some_var
 do_something
end
# good
if some_var
 do_something
end
# good
# When using a boolean value for an infinite loop.
while true
 break if condition
end
```
## Lint/LiteralAssignmentInCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.58 | - |

Checks for literal assignments in the conditions of `if`, `while`, and `until`.
It emulates the following Ruby warning:

```
$ ruby -we 'if x = true; end'
-e:1: warning: found `= literal' in conditional, should be ==
```
As a lint cop, it cannot be determined if `==` is appropriate as intended,
therefore this cop does not provide autocorrection.

## Lint/LiteralInInterpolation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.19 | 0.32 |

Checks for interpolated literals.

| Array literals interpolated in regexps are not handled by this cop, but
by `Lint/ArrayLiteralInRegexp`instead. |

## Lint/Loop

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.9 | 1.3 |

Checks for uses of `begin…end while/until something`.

### Safety

The cop is unsafe because behavior can change in some cases, including
if a local variable inside the loop body is accessed outside of it, or if the
loop body raises a `StopIteration` exception (which `Kernel#loop` rescues).

### Examples

```
# bad
# using while
begin
 do_something
end while some_condition
# good
# while replacement
loop do
 do_something
 break unless some_condition
end
# bad
# using until
begin
 do_something
end until some_condition
# good
# until replacement
loop do
 do_something
 break if some_condition
end
```
## Lint/MissingCopEnableDirective

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.52 | 1.88 |

Checks that there is an `# rubocop:enable …` statement
after a `# rubocop:disable …` statement. This will prevent leaving
cop disables on wide ranges of code, that later contributors to
a file wouldn’t be aware of.

You can set `MaxRangeSize` to define the maximum number of
consecutive lines a cop can be disabled for.

-
`.inf`any size (default)
-
`0`allows only single-line disables
-
`1`means the maximum allowed is as follows:

```
# rubocop:disable SomeCop
a = 1
# rubocop:enable SomeCop
```
### Examples

#### MaxRangeSize: .inf (default)

```
# good
# rubocop:disable Layout/SpaceAroundOperators
x= 0
# rubocop:enable Layout/SpaceAroundOperators
# y = 1
# EOF
# bad
# rubocop:disable Layout/SpaceAroundOperators
x= 0
# EOF
```
#### MaxRangeSize: 2

```
# good
# rubocop:disable Layout/SpaceAroundOperators
x= 0
# With the previous, there are 2 lines on which cop is disabled.
# rubocop:enable Layout/SpaceAroundOperators
# bad
# rubocop:disable Layout/SpaceAroundOperators
x= 0
x += 1
# Including this, that's 3 lines on which the cop is disabled.
# rubocop:enable Layout/SpaceAroundOperators
```
## Lint/MissingSuper

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.89 | 1.89 |

Checks for the presence of constructors and lifecycle callbacks
without calls to `super`.

This cop does not consider `method_missing` (and `respond_to_missing?`)
because in some cases it makes sense to overtake what is considered a
missing method. In other cases, the theoretical ideal handling could be
challenging or verbose for no actual gain.

Autocorrection is not supported because the position of `super` cannot be
determined automatically.

`Object` and `BasicObject` are allowed by this cop because of their
stateless nature. However, sometimes you might want to allow other parent
classes from this cop, for example in the case of an abstract class that is
not meant to be called with `super`. In those cases, you can use the
`AllowedParentClasses` option to specify which classes should be allowed
**in addition to** `Object` and `BasicObject`.

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed,
the constructor check additionally consults the project-wide index: if the
class' entire ancestry is resolvable and no ancestor defines `initialize`,
no offense is registered, since `super` would only reach the no-op
`Object#initialize`. Classes whose ancestry contains an unresolvable
superclass or mixin (e.g. one defined in a gem) are still reported.

### Examples

```
# bad
class Employee < Person
 def initialize(name, salary)
 @salary = salary
 end
end
# good
class Employee < Person
 def initialize(name, salary)
 super(name)
 @salary = salary
 end
end
# bad
Employee = Class.new(Person) do
 def initialize(name, salary)
 @salary = salary
 end
end
# good
Employee = Class.new(Person) do
 def initialize(name, salary)
 super(name)
 @salary = salary
 end
end
# bad
class Parent
 def self.inherited(base)
 do_something
 end
end
# good
class Parent
 def self.inherited(base)
 super
 do_something
 end
end
# good
class ClassWithNoParent
 def initialize
 do_something
 end
end
```
## Lint/MixedCaseRange

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.53 | - |

Checks for mixed-case character ranges since they include likely unintended characters.

Offenses are registered for regexp character classes like `/[A-z]/`
as well as range objects like `('A'..'z')`.

| `Range`objects cannot be autocorrected. |

### Safety

The cop autocorrects regexp character classes
by replacing one character range with two: `A-z` becomes `A-Za-z`.
In most cases this is probably what was originally intended
but it changes the regexp to no longer match symbols it used to include.
For this reason, this cop’s autocorrect is unsafe (it will
change the behavior of the code).

## Lint/MixedRegexpCaptureTypes

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.85 | - |

Do not mix named captures and numbered captures in a `Regexp` literal
because numbered capture is ignored if they’re mixed.
Replace numbered captures with non-capturing groupings or
named captures.

## Lint/MultipleComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.47 | 1.1 |

In math and Python, we can use `x < y < z` style comparison to compare
multiple values. However, we can’t use the comparison in Ruby. However,
the comparison is not a syntax error. This cop checks the bad usage of
comparison operators.

## Lint/NameTypo

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.89 | - |

Checks for probable typos in constant and method names: a name that does not resolve anywhere in the project, used in a namespace the project does define, with a close-named sibling to suggest instead.

The check is powered by the project-wide index, so it only runs when
`AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed.
Without the index the cop does nothing.

Constants are checked only in qualified references (`Foo::Bar`) whose
namespace resolves in the index; bare names cannot be distinguished
from constants provided by gems or the standard library. Methods are
checked only in calls on constant receivers (`Foo.bar`) whose entire
indexed ancestry is resolved, so methods gained through gem classes or
dynamic definitions never produce offenses. In both cases an offense
requires a similarly named alternative to exist — an unknown name
alone is not reported, since the index does not see gems or the
standard library.

Names that appear as symbols or inside string literals in the same
file are never reported, since they usually belong to runtime
definitions the index cannot see (`stub_const`, `const_set`,
`define_method`, and the like).

## Lint/NestedMethodDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.32 | - |

Checks for nested method definitions.

### Examples

```
# bad
# `bar` definition actually produces methods in the same scope
# as the outer `foo` method. Furthermore, the `bar` method
# will be redefined every time `foo` is invoked.
def foo
 def bar
 end
end
# good
def foo
 bar = -> { puts 'hello' }
 bar.call
end
# good
# `class_eval`, `instance_eval`, `module_eval`, `class_exec`, `instance_exec`, and
# `module_exec` blocks are allowed by default.
def foo
 self.class.class_eval do
 def bar
 end
 end
end
def foo
 self.class.module_exec do
 def bar
 end
 end
end
# good
def foo
 class << self
 def bar
 end
 end
end
```
#### AllowedMethods: [] (default)

```
# bad
def do_something
 has_many :articles do
 def find_or_create_by_name(name)
 end
 end
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowedMethods |
 | Array |
| AllowedPatterns |
 | Array |

## Lint/NestedPercentLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.52 | - |

Checks for nested percent literals.

### Examples

```
# bad
# The percent literal for nested_attributes is parsed as four tokens,
# yielding the array [:name, :content, :"%i[incorrectly", :"nested]"].
attributes = {
 valid_attributes: %i[name content],
 nested_attributes: %i[name content %i[incorrectly nested]]
}
# good
# Neither is incompatible with the bad case, but probably the intended code.
attributes = {
 valid_attributes: %i[name content],
 nested_attributes: [:name, :content, %i[incorrectly nested]]
}
attributes = {
 valid_attributes: %i[name content],
 nested_attributes: [:name, :content, [:incorrectly, :nested]]
}
```
## Lint/NextWithoutAccumulator

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.36 | - |

Don’t omit the accumulator when calling `next` in a `reduce` block.

## Lint/NoReturnInBeginEndBlocks

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.2 | - |

Checks for the presence of a `return` inside a `begin..end` block
in assignment contexts.
In this situation, the `return` will result in an exit from the current
method, possibly leading to unexpected behavior.

### Examples

```
# bad
@some_variable ||= begin
 return some_value if some_condition_is_met
 do_something
end
# good
@some_variable ||= begin
 if some_condition_is_met
 some_value
 else
 do_something
 end
end
# good
some_variable = if some_condition_is_met
 return if another_condition_is_met
 some_value
 else
 do_something
 end
```
## Lint/NonAtomicFileOperation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.31 | - |

Checks for non-atomic file operation. And then replace it with a nearly equivalent and atomic method.

These can cause problems that are difficult to reproduce, especially in cases of frequent file operations in parallel, such as test runs with parallel_rspec.

For examples: creating a directory if there is none, has the following problems

An exception occurs when the directory didn’t exist at the time of `exist?`,
but someone else created it before `mkdir` was executed.

Subsequent processes are executed without the directory that should be there
when the directory existed at the time of `exist?`,
but someone else deleted it shortly afterwards.

### Safety

This cop is unsafe, because autocorrection change to atomic processing. The atomic processing of the replacement destination is not guaranteed to be strictly equivalent to that before the replacement.

### Examples

```
# bad - race condition with another process may result in an error in `mkdir`
unless Dir.exist?(path)
 FileUtils.mkdir(path)
end
# good - atomic and idempotent creation
FileUtils.mkdir_p(path)
# bad - race condition with another process may result in an error in `remove`
if File.exist?(path)
 FileUtils.remove(path)
end
# good - atomic and idempotent removal
FileUtils.rm_f(path)
```
## Lint/NonDeterministicRequireOrder

| Requires Ruby version ⇐ 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.78 | - |

`Dir[…]` and `Dir.glob(…)` do not make any guarantees about
the order in which files are returned. The final order is
determined by the operating system and file system.
This means that using them in cases where the order matters,
such as requiring files, can lead to intermittent failures
that are hard to debug. To ensure this doesn’t happen,
always sort the list.

`Dir.glob` and `Dir[]` sort globbed results by default in Ruby 3.0.
So all bad cases are acceptable when Ruby 3.0 or higher is used.

| This cop will be deprecated and removed when supporting only Ruby 3.0 and higher. |

### Examples

```
# bad
Dir["./lib/**/*.rb"].each do |file|
 require file
end
# good
Dir["./lib/**/*.rb"].sort.each do |file|
 require file
end
# bad
Dir.glob(Rails.root.join(__dir__, 'test', '*.rb')) do |file|
 require file
end
# good
Dir.glob(Rails.root.join(__dir__, 'test', '*.rb')).sort.each do |file|
 require file
end
# bad
Dir['./lib/**/*.rb'].each(&method(:require))
# good
Dir['./lib/**/*.rb'].sort.each(&method(:require))
# bad
Dir.glob(Rails.root.join('test', '*.rb'), &method(:require))
# good
Dir.glob(Rails.root.join('test', '*.rb')).sort.each(&method(:require))
# good - Respect intent if `sort` keyword option is specified in Ruby 3.0 or higher.
Dir.glob(Rails.root.join(__dir__, 'test', '*.rb'), sort: false).each(&method(:require))
```
## Lint/NonLocalExitFromIterator

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.30 | - |

Checks for non-local exits from iterators without a return value. It registers an offense under these conditions:

-
No value is returned,
-
the block is preceded by a method chain,
-
the block has arguments,
-
the method which receives the block is not `define_method`or`define_singleton_method`,
-
the return is not contained in an inner scope, e.g. a lambda or a method definition.

### Examples

```
class ItemApi
 rescue_from ValidationError do |e| # non-iteration block with arg
 return { message: 'validation error' } unless e.errors # allowed
 error_array = e.errors.map do |error| # block with method chain
 return if error.suppress? # warned
 return "#{error.param}: invalid" unless error.message # allowed
 "#{error.param}: #{error.message}"
 end
 { message: 'validation error', errors: error_array }
 end
 def update_items
 transaction do # block without arguments
 return unless update_necessary? # allowed
 find_each do |item| # block without method chain
 return if item.stock == 0 # false-negative...
 item.update!(foobar: true)
 end
 end
 end
end
```
## Lint/NumberConversion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always (Unsafe) | 0.53 | 1.88 |

Warns against the usage of unsafe number conversions. Unsafe number conversion can cause an unexpected error if auto type conversion fails. The cop prefers parsing with a number class instead.

Conversion with `Integer`, `Float`, etc. will raise an `ArgumentError`
if given input that is not numeric (eg. an empty string), whereas
`to_i`, etc. will try to convert regardless of input (`''.to_i ⇒ 0`).
As such, this cop is disabled by default because it’s not necessarily
always correct to raise if a value is not numeric.

| Some values cannot be converted properly using one of the `Kernel`methods (for instance,`Time`and`DateTime`values are allowed by this
cop by default). Similarly, Rails' duration methods do not work well
with`Integer()`and can be allowed with`AllowedMethods`. By default,
there are no allowed methods. |

### Safety

Autocorrection is unsafe because it is not guaranteed that the
replacement `Kernel` methods are able to properly handle the
input if it is not a standard class.

## Lint/NumberedParameterAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.9 | - |

Checks for uses of numbered parameter assignment. It emulates the following warning in Ruby 2.7:

$ ruby -ve '_1 = :value' ruby 2.7.2p137 (2020-10-01 revision 5445e04352) [x86_64-darwin19] -e:1: warning: `_1' is reserved for numbered parameter; consider another name

Assigning to a numbered parameter (from `_1` to `_9`) causes an error in Ruby 3.0.

$ ruby -ve '_1 = :value' ruby 3.0.0p0 (2020-12-25 revision 95aff21468) [x86_64-darwin19] -e:1: _1 is reserved for numbered parameter

| The numbered parameters are from `_1`to`_9`. This cop checks`_0`, and over`_10`as well to prevent confusion. |

## Lint/NumericOperationWithConstantResult

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.69 | 1.88 |

Certain numeric operations have a constant result, usually 0 or 1.
Multiplying a number by 0 will always return 0.
Dividing a number by itself or raising it to the power of 0 will always return 1.
As such, they can be replaced with that result.
These are probably leftover from debugging, or are mistakes.
Other numeric operations that are similarly leftover from debugging or mistakes
are handled by `Lint/UselessNumericOperation`.

| This cop doesn’t detect offenses for the `-`and`%`operator because it
can’t determine the type of`x`. If`x`is an`Array`or`String`, it doesn’t perform
a numeric operation. |

### Safety

This cop is unsafe because the autocorrection drops the operands, which
discards any side effects of evaluating them and can change behavior when
the result is not actually constant. For example, `x / x` raises
`ZeroDivisionError` when `x` is `0`, and returns `Float::NAN` (not `1`)
when `x` is `0.0`; replacing it with `1` silences that.

## Lint/OrAssignmentToConstant

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.9 | - |

Checks for unintended or-assignment to a constant.

Constants should always be assigned in the same location. And its value
should always be the same. If constants are assigned in multiple
locations, the result may vary depending on the order of `require`.

## Lint/OrderedMagicComments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.53 | 1.37 |

Checks the proper ordering of magic comments and whether a magic comment is not placed before a shebang.

### Examples

```
# bad
# frozen_string_literal: true
# encoding: ascii
p [''.frozen?, ''.encoding] #=> [true, #<Encoding:UTF-8>]
# good
# encoding: ascii
# frozen_string_literal: true
p [''.frozen?, ''.encoding] #=> [true, #<Encoding:US-ASCII>]
# good
#!/usr/bin/env ruby
# encoding: ascii
# frozen_string_literal: true
p [''.frozen?, ''.encoding] #=> [true, #<Encoding:US-ASCII>]
```
## Lint/OutOfRangeRegexpRef

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | No | 0.89 | - |

Looks for references of `Regexp` captures that are out of range
and thus always returns nil.

### Safety

This cop is unsafe because it is naive in how it determines what references are available based on the last encountered regexp, but it cannot handle some cases, such as conditional regexp matches, which leads to false positives, such as:

```
foo ? /(c)(b)/ =~ str : /(b)/ =~ str
do_something if $2
# $2 is defined for the first condition but not the second, however
# the cop will mark this as an offense.
```
This might be a good indication of code that should be refactored, however.

## Lint/ParenthesesAsGroupedExpression

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.12 | 0.83 |

Checks for space between the name of a called method and a left parenthesis.

### Examples

```
# bad
do_something (foo)
# good
do_something(foo)
do_something (2 + 3) * 4
do_something (foo * bar).baz
```
## Lint/PercentStringArray

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.41 | - |

Checks for quotes and commas in %w, e.g. `%w('foo', "bar")`

It is more likely that the additional characters are unintended (for example, mistranslating an array of literals to percent string notation) rather than meant to be part of the resulting strings.

## Lint/PercentSymbolArray

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.41 | - |

Checks for colons and commas in %i, e.g. `%i(:foo, :bar)`

It is more likely that the additional characters are unintended (for example, mistranslating an array of literals to percent string notation) rather than meant to be part of the resulting symbols.

## Lint/RaiseException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.81 | 0.86 |

Checks for `raise` or `fail` statements which raise `Exception` or
`Exception.new`. Use `StandardError` or a specific exception class instead.

If you have defined your own namespaced `Exception` class, it is possible
to configure the cop to allow it by setting `AllowedImplicitNamespaces` to
an array with the names of the namespaces to allow. By default, this is set to
`['Gem']`, which allows `Gem::Exception` to be raised without an explicit namespace.
If not allowed, a false positive may be registered if `raise Exception` is called
within the namespace.

Alternatively, use a fully qualified name with `raise`/`fail`
(eg. `raise Namespace::Exception`).

### Safety

This cop is unsafe because it will change the exception class being raised, which is a change in behavior.

### Examples

```
# bad
raise Exception, 'Error message here'
raise Exception.new('Error message here')
# good
raise StandardError, 'Error message here'
raise MyError.new, 'Error message here'
```
#### AllowedImplicitNamespaces: ['Gem'] (default)

```
# bad - `Foo` is not an allowed implicit namespace
module Foo
 def self.foo
 raise Exception # This is qualified to `Foo::Exception`.
 end
end
# good
module Gem
 def self.foo
 raise Exception # This is qualified to `Gem::Exception`.
 end
end
# good
module Foo
 def self.foo
 raise Foo::Exception
 end
end
```
## Lint/RedundantCopDisableDirective

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.76 | - |

Detects instances of rubocop:disable comments that can be removed without causing any offenses to be reported. It’s implemented as a cop in that it inherits from the Cop base class and calls add_offense. The unusual part of its implementation is that it doesn’t have any on_* methods or an investigate method. This means that it doesn’t take part in the investigation phase when the other cops do their work. Instead, it waits until it’s called in a later stage of the execution. The reason it can’t be implemented as a normal cop is that it depends on the results of all other cops to do its work.

## Lint/RedundantCopEnableDirective

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.76 | - |

Detects instances of rubocop:enable comments that can be removed.

When a comment enables all cops at once `rubocop:enable all`
the cop checks whether any cop was actually enabled.

## Lint/RedundantDirGlobSort

| Requires Ruby version 3.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.8 | 1.26 |

Sort globbed results by default in Ruby 3.0.
This cop checks for redundant `sort` method to `Dir.glob` and `Dir[]`.

## Lint/RedundantRegexpQuantifiers

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.53 | - |

Checks for redundant quantifiers inside `Regexp` literals.

It is always allowed when interpolation is used in a regexp literal, because it’s unknown what kind of string will be expanded as a result:

`/(?:a*#{interpolation})?/x`## Lint/RedundantRequireStatement

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.76 | 1.73 |

Checks for unnecessary `require` statement.

The following features are unnecessary `require` statement because
they are already loaded. e.g. Ruby 2.2:

```
ruby -ve 'p $LOADED_FEATURES.reject { |feature| %r|/| =~ feature }'
ruby 2.2.8p477 (2017-09-14 revision 59906) [x86_64-darwin13]
["enumerator.so", "rational.so", "complex.so", "thread.rb"]
```
Below are the features that each `TargetRubyVersion` targets.

-
2.0+ … `enumerator`
-
2.1+ … `thread`
-
2.2+ … Add `rational`and`complex`above
-
2.7+ … Add `ruby2_keywords`above
-
3.1+ … Add `fiber`above
-
3.2+ … Add `set`above
-
4.0+ … Add `pathname`above

This cop target those features.

## Lint/RedundantSafeNavigation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.93 | 1.79 |

Checks for redundant safe navigation calls.
Use cases where a constant, named in camel case for classes and modules is `nil` are rare,
and an offense is not detected when the receiver is a constant. The detection also applies
to `self`, and to literal receivers, except for `nil`.

For all receivers, the `instance_of?`, `kind_of?`, `is_a?`, `eql?`, `respond_to?`,
and `equal?` methods are checked by default.
These are customizable with `AllowedMethods` option.

The `AllowedMethods` option specifies nil-safe methods,
in other words, it is a method that is allowed to skip safe navigation.
Note that the `AllowedMethod` option is not an option that specifies methods
for which to suppress (allow) this cop’s check.

In the example below, the safe navigation operator (`&.`) is unnecessary
because `NilClass` has methods like `respond_to?` and `is_a?`.

The `InferNonNilReceiver` option specifies whether to look into previous code
paths to infer if the receiver can’t be nil. This check is unsafe because the receiver
can be redefined between the safe navigation call and previous regular method call.
It does the inference only in the current scope, e.g. within the same method definition etc.

The `AdditionalNilMethods` option specifies additional custom methods which are
defined on `NilClass`. When `InferNonNilReceiver` is set, they are used to determine
whether the receiver can be nil.

### Safety

This cop is unsafe, because autocorrection can change the return type of
the expression. An offending expression that previously could return `nil`
will be autocorrected to never return `nil`.

### Examples

```
# bad
CamelCaseConst&.do_something
# good
CamelCaseConst.do_something
# bad
foo.to_s&.strip
foo.to_i&.zero?
foo.to_f&.zero?
foo.to_a&.size
foo.to_h&.size
# good
foo.to_s.strip
foo.to_i.zero?
foo.to_f.zero?
foo.to_a.size
foo.to_h.size
# bad
do_something if attrs&.respond_to?(:[])
# good
do_something if attrs.respond_to?(:[])
# bad
foo&.bar ? foo&.bar.baz : qux
# good
foo&.bar ? foo.bar.baz : qux
# bad
if foo&.bar
 foo&.bar.baz
end
# good
if foo&.bar
 foo.bar.baz
end
# bad
while node&.is_a?(BeginNode)
 node = node.parent
end
# good
while node.is_a?(BeginNode)
 node = node.parent
end
# good - without `&.` this changes the return value for `nil`
foo&.respond_to?(:to_a)
foo&.respond_to?(:class)
# bad - for `nil`s conversion methods return default values for the type
foo&.to_h || {}
foo&.to_h { |k, v| [k, v] } || {}
foo&.to_a || []
foo&.to_i || 0
foo&.to_f || 0.0
foo&.to_s || ''
# good
foo.to_h
foo.to_h { |k, v| [k, v] }
foo.to_a
foo.to_i
foo.to_f
foo.to_s
# bad
self&.foo
# good
self.foo
```
#### AllowedMethods: [nil_safe_method]

```
# bad
do_something if attrs&.nil_safe_method(:[])
# good
do_something if attrs.nil_safe_method(:[])
do_something if attrs&.not_nil_safe_method(:[])
```
## Lint/RedundantSplatExpansion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.76 | 1.7 |

Checks for unneeded usages of splat expansion.

### Examples

```
# bad
a = *[1, 2, 3]
a = *'a'
a = *1
['a', 'b', *%w(c d e), 'f', 'g']
# good
c = [1, 2, 3]
a = *c
a, b = *c
a, *b = *c
a = *1..10
a = ['a']
['a', 'b', 'c', 'd', 'e', 'f', 'g']
# bad
do_something(*['foo', 'bar', 'baz'])
# good
do_something('foo', 'bar', 'baz')
# bad
begin
 foo
rescue *[StandardError, ApplicationError]
 bar
end
# good
begin
 foo
rescue StandardError, ApplicationError
 bar
end
# bad
case foo
when *[1, 2, 3]
 bar
else
 baz
end
# good
case foo
when 1, 2, 3
 bar
else
 baz
end
```
## Lint/RedundantStringCoercion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.19 | 0.77 |

Checks for string conversion in string interpolation, `print`, `puts`, and `warn` arguments,
which is redundant.

### Examples

```
# bad
"result is #{something.to_s}"
print something.to_s
puts something.to_s
warn something.to_s
# good
"result is #{something}"
print something
puts something
warn something
```
## Lint/RedundantTypeConversion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.72 | - |

Checks for redundant uses of `to_s`, `to_sym`, `to_i`, `to_f`, `to_d`, `to_r`, `to_c`,
`to_a`, `to_h`, and `to_set`.

When one of these methods is called on an object of the same type, that object
is returned, making the call unnecessary. The cop detects conversion methods called
on object literals, class constructors, class `[]` methods, and the `Kernel` methods
`String()`, `Integer()`, `Float()`, BigDecimal(), `Rational()`, `Complex()`, and `Array()`.

Specifically, these cases are detected for each conversion method:

-
`to_s`when called on a string literal, interpolated string, heredoc, or with`String.new`or`String()`.
-
`to_sym`when called on a symbol literal or interpolated symbol.
-
`to_i`when called on an integer literal or with`Integer()`.
-
`to_f`when called on a float literal or with`Float()`.
-
`to_d`when called with`BigDecimal()`.
-
`to_r`when called on a rational literal or with`Rational()`.
-
`to_c`when called on a complex literal or with`Complex()`.
-
`to_a`when called on an array literal, or with`Array.new`,`Array()`or`Array[]`.
-
`to_h`when called on a hash literal, or with`Hash.new`,`Hash()`or`Hash[]`.
-
`to_set`when called on`Set.new`or`Set[]`.

In all cases, chaining one of the same `to_*` conversion methods listed above is redundant.

The cop can also register an offense for chaining conversion methods on methods that are
expected to return a specific type regardless of receiver (eg. `foo.inspect.to_s` and
`foo.to_json.to_s`).

### Examples

```
# bad
"text".to_s
:sym.to_sym
42.to_i
8.5.to_f
12r.to_r
1i.to_c
[].to_a
{}.to_h
Set.new.to_set
# good
"text"
:sym
42
8.5
12r
1i
[]
{}
Set.new
# bad
Integer(var).to_i
# good
Integer(var)
# good - chaining to a type constructor with exceptions suppressed
# in this case, `Integer()` could return `nil`
Integer(var, exception: false).to_i
# bad
BigDecimal(var).to_d
# good
BigDecimal(var)
# bad - chaining the same conversion
foo.to_s.to_s
# good
foo.to_s
# bad - chaining a conversion to a method that is expected to return the same type
foo.inspect.to_s
foo.to_json.to_s
# good
foo.inspect
foo.to_json
```
## Lint/RedundantWithObject

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.51 | 1.88 |

Checks for redundant `with_object`.

## Lint/RefinementImportMethods

| Requires Ruby version 3.1 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.27 | 1.88 |

Checks if `include` or `prepend` is called in `refine` block.
These methods are deprecated and should be replaced with `Refinement#import_methods`.

It emulates deprecation warnings in Ruby 3.1. Functionality has been removed in Ruby 3.2.

## Lint/RegexpAsCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.51 | 0.86 |

Checks for regexp literals used as `match-current-line`.
If a regexp literal is in condition, the regexp matches `$_` implicitly.

## Lint/RequireParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.18 | - |

Checks for method calls with at least one argument where no parentheses
are used around the parameter list, and the call could be misread as an
operand of a boolean operator (`&&` or `||`). Two forms are flagged:

-
a predicate method whose last argument is a `&&`/`||`expression, and
-
any method whose first argument is a ternary expression with a `&&`/`||`condition.

The idea behind warning for these constructs is that the user might be under the impression that the return value from the method call is an operand of &&/||.

## Lint/RequireRangeParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.32 | - |

Checks that a range literal is enclosed in parentheses when the end of the range is at a line break.

| The following is maybe intended for `(42..)`. But, compatible is`42..do_something`.
So, this cop does not provide autocorrection because it is left to user. |

```
case condition
when 42..
 do_something
end
```
## Lint/RescueException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.9 | 0.27 |

Checks for `rescue` blocks targeting the `Exception` class.

### Examples

```
# bad
begin
 do_something
rescue Exception
 handle_exception
end
# good
begin
 do_something
rescue ArgumentError
 handle_exception
end
```
## Lint/RescueType

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | - |

Checks for arguments to `rescue` that will result in a `TypeError`
if an exception is raised.

## Lint/ReturnInVoidContext

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | - |

Checks for the use of a return with a value in a context where the value will be ignored. (initialize and setter methods)

## Lint/SafeNavigationChain

| Requires Ruby version 2.3 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.47 | 1.89 |

The safe navigation operator returns nil if the receiver is nil. If you chain an ordinary method call after a safe navigation operator, it raises NoMethodError. We should use a safe navigation operator after a safe navigation operator. This cop checks for the problem outlined above.

### Safety

This cop’s autocorrection is unsafe because it changes the behavior of
the code: `x&.foo.bar` raises `NoMethodError` when `x` or `x&.foo` is `nil`,
whereas the corrected `x&.foo&.bar` returns `nil` instead. The autocorrection
also assumes that extending safe navigation through the chain is intended,
while removing the first `&.` may be what was intended instead.

## Lint/SafeNavigationConsistency

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.55 | 1.85 |

Checks that if safe navigation is used in an `&&` or `||` condition,
consistent and appropriate safe navigation, without excess or deficiency,
is used for all method calls on the same object.

### Safety

Autocorrection is unsafe because if the receiver is not a local variable
but a method call, it may not be idempotent. For example, replacing
`foo&.bar` with `foo.bar` could raise `NoMethodError` if `foo` returns
`nil` on a subsequent call.

### Examples

```
# bad
foo&.bar && foo&.baz
# good
foo&.bar && foo.baz
# bad
foo.bar && foo&.baz
# good
foo.bar && foo.baz
# bad
foo&.bar || foo.baz
# good
foo&.bar || foo&.baz
# bad
foo.bar || foo&.baz
# good
foo.bar || foo.baz
# bad
foo&.bar && (foobar.baz || foo&.baz)
# good
foo&.bar && (foobar.baz || foo.baz)
```
## Lint/SafeNavigationWithEmpty

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.62 | 0.87 |

Checks to make sure safe navigation isn’t used with `empty?` in
a conditional.

While the safe navigation operator is generally a good idea, when
checking `foo&.empty?` in a conditional, `foo` being `nil` will actually
do the opposite of what the author intends.

## Lint/ScriptPermission

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 0.50 |

Checks if a file which has a shebang line as its first line is granted execute permission.

### Examples

```
# bad
# A file which has a shebang line as its first line is not
# granted execute permission.
#!/usr/bin/env ruby
puts 'hello, world'
# good
# A file which has a shebang line as its first line is
# granted execute permission.
#!/usr/bin/env ruby
puts 'hello, world'
# good
# A file which has not a shebang line as its first line is not
# granted execute permission.
puts 'hello, world'
```
## Lint/SelfAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.89 | - |

Checks for self-assignments.

### Examples

```
# bad
foo = foo
foo, bar = foo, bar
Foo = Foo
hash['foo'] = hash['foo']
obj.attr = obj.attr
# good
foo = bar
foo, bar = bar, foo
Foo = Bar
hash['foo'] = hash['bar']
obj.attr = obj.attr2
# good (method calls possibly can return different results)
hash[foo] = hash[foo]
```
## Lint/SendWithMixinArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.75 | - |

Checks for `send`, `public_send`, and *send*

`include` and `prepend` methods were private methods until Ruby 2.0,
they were mixed-in via `send` method. This cop uses Ruby 2.1 or
higher style that can be called by public methods.
And `extend` method that was originally a public method is also targeted
for style unification.

### Examples

```
# bad
Foo.send(:include, Bar)
Foo.send(:prepend, Bar)
Foo.send(:extend, Bar)
# bad
Foo.public_send(:include, Bar)
Foo.public_send(:prepend, Bar)
Foo.public_send(:extend, Bar)
# bad
Foo.__send__(:include, Bar)
Foo.__send__(:prepend, Bar)
Foo.__send__(:extend, Bar)
# good
Foo.include Bar
Foo.prepend Bar
Foo.extend Bar
```
## Lint/ShadowedArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.52 | - |

Checks for shadowed arguments.

This cop has `IgnoreImplicitReferences` configuration option.
It means argument shadowing is used in order to pass parameters
to zero arity `super` when `IgnoreImplicitReferences` is `true`.

### Examples

```
# bad
do_something do |foo|
 foo = 42
 puts foo
end
def do_something(foo)
 foo = 42
 puts foo
end
# good
do_something do |foo|
 foo = foo + 42
 puts foo
end
def do_something(foo)
 foo = foo + 42
 puts foo
end
def do_something(foo)
 puts foo
end
```
## Lint/ShadowedException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.41 | - |

Checks for a rescued exception that gets shadowed by a less specific exception being rescued before a more specific exception is rescued.

An exception is considered shadowed if it is rescued after its
ancestor is, or if it and its ancestor are both rescued in the
same `rescue` statement. In both cases, the more specific rescue is
unnecessary because it is covered by rescuing the less specific
exception. (ie. `rescue Exception, StandardError` has the same behavior
whether `StandardError` is included or not, because all `StandardError`s
are rescued by `rescue Exception`).

### Examples

```
# bad
begin
 something
rescue Exception
 handle_exception
rescue StandardError
 handle_standard_error
end
# bad
begin
 something
rescue Exception, StandardError
 handle_error
end
# good
begin
 something
rescue StandardError
 handle_standard_error
rescue Exception
 handle_exception
end
# good, however depending on runtime environment.
#
# This is a special case for system call errors.
# System dependent error code depends on runtime environment.
# For example, whether `Errno::EAGAIN` and `Errno::EWOULDBLOCK` are
# the same error code or different error code depends on environment.
# This good case is for `Errno::EAGAIN` and `Errno::EWOULDBLOCK` with
# the same error code.
begin
 something
rescue Errno::EAGAIN, Errno::EWOULDBLOCK
 handle_standard_error
end
```
## Lint/ShadowingOuterLocalVariable

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.9 | 1.76 |

Checks for the use of local variable names from an outer scope
in block arguments or block-local variables. This mirrors the warning
given by `ruby -cw` prior to Ruby 2.6:
"shadowing outer local variable - foo".

The cop is now disabled by default to match the upstream Ruby behavior. It’s useful, however, if you’d like to avoid shadowing variables from outer scopes, which some people consider an anti-pattern that makes it harder to keep track of what’s going on in a program.

| Shadowing of variables in block passed to `Ractor.new`is allowed
because`Ractor`should not access outer variables.
eg. following style is encouraged: |

```
worker_id, pipe = env
Ractor.new(worker_id, pipe) do |worker_id, pipe|
end
```
## Lint/SharedMutableDefault

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.70 | - |

Checks for `Hash` creation with a mutable default value.
Creating a `Hash` in such a way will share the default value
across all keys, causing unexpected behavior when modifying it.

For example, when the `Hash` was created with an `Array` as the argument,
calling `hash[:foo] << 'bar'` will also change the value of all
other keys that have not been explicitly assigned to.

### Examples

```
# bad
Hash.new([])
Hash.new({})
Hash.new(Array.new)
Hash.new(Hash.new)
# okay -- In rare cases that intentionally have this behavior,
# without disabling the cop, you can set the default explicitly.
h = Hash.new
h.default = []
h[:a] << 1
h[:b] << 2
h # => {:a => [1, 2], :b => [1, 2]}
# okay -- beware this will discard mutations and only remember assignments
Hash.new { Array.new }
Hash.new { Hash.new }
Hash.new { {} }
Hash.new { [] }
# good - frozen solution will raise an error when mutation attempted
Hash.new([].freeze)
Hash.new({}.freeze)
# good - using a proc will create a new object for each key
h = Hash.new
h.default_proc = ->(h, k) { [] }
h.default_proc = ->(h, k) { {} }
# good - using a block will create a new object for each key
Hash.new { |h, k| h[k] = [] }
Hash.new { |h, k| h[k] = {} }
```
## Lint/StructNewOverride

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.81 | - |

Checks unexpected overrides of the `Struct` built-in methods
via `Struct.new`.

### Examples

```
# bad
Bad = Struct.new(:members, :clone, :count)
b = Bad.new([], true, 1)
b.members #=> [] (overriding `Struct#members`)
b.clone #=> true (overriding `Object#clone`)
b.count #=> 1 (overriding `Enumerable#count`)
# good
Good = Struct.new(:id, :name)
g = Good.new(1, "foo")
g.members #=> [:id, :name]
g.clone #=> #<struct Good id=1, name="foo">
g.count #=> 2
```
## Lint/SuppressedException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.9 | 1.12 |

Checks for `rescue` blocks with no body.

### Examples

```
# bad
def some_method
 do_something
rescue
end
# bad
begin
 do_something
rescue
end
# good
def some_method
 do_something
rescue
 handle_exception
end
# good
begin
 do_something
rescue
 handle_exception
end
```
#### AllowComments: true (default)

```
# good
def some_method
 do_something
rescue
 # do nothing
end
# good
begin
 do_something
rescue
 # do nothing
end
```
#### AllowComments: false

```
# bad
def some_method
 do_something
rescue
 # do nothing
end
# bad
begin
 do_something
rescue
 # do nothing
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowComments |
 | Boolean |
| AllowNil |
 | Boolean |

## Lint/SuppressedExceptionInNumberConversion

| Requires Ruby version 2.6 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.72 | - |

Checks for cases where exceptions unrelated to the numeric constructors `Integer()`,
`Float()`, `BigDecimal()`, `Complex()`, and `Rational()` may be unintentionally swallowed.

### Safety

The cop is unsafe for autocorrection because unexpected errors occurring in the argument
passed to numeric constructor (e.g., `Integer()`) can lead to incompatible behavior.
For example, changing it to `Integer(potential_exception_method_call, exception: false)`
ensures that exceptions raised by `potential_exception_method_call` are not ignored.

`Integer(potential_exception_method_call) rescue nil`## Lint/SymbolConversion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.9 | 1.16 |

Checks for uses of literal strings converted to a symbol where a literal symbol could be used instead.

There are two possible styles for this cop.
`strict` (default) will register an offense for any incorrect usage.
`consistent` additionally requires hashes to use the same style for
every symbol key (ie. if any symbol key needs to be quoted it requires
all keys to be quoted).

### Examples

```
# bad
'string'.to_sym
:symbol.to_sym
'underscored_string'.to_sym
:'underscored_symbol'
'hyphenated-string'.to_sym
"string_#{interpolation}".to_sym
# good
:string
:symbol
:underscored_string
:underscored_symbol
:'hyphenated-string'
:"string_#{interpolation}"
```
## Lint/Syntax

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.9 | - |

Repacks Parser’s diagnostics/errors into RuboCop’s offenses.

## Lint/ToEnumArguments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.1 | - |

Ensures that `to_enum`/`enum_for`, called for the current method,
has correct arguments.

### Examples

```
# bad
def foo(x, y = 1)
 return to_enum(__callee__, x) # `y` is missing
end
# good
def foo(x, y = 1)
 # Alternatives to `__callee__` are `__method__` and `:foo`.
 return to_enum(__callee__, x, y)
end
# good
def foo(x, y = 1)
 # It is also allowed if it is wrapped in some method like Sorbet.
 return to_enum(T.must(__callee__), x, y)
end
```
## Lint/ToJSON

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.66 | - |

Checks to make sure `#to_json` includes an optional argument.
When overriding `#to_json`, callers may invoke JSON
generation via `JSON.generate(your_obj)`. Since `JSON#generate` allows
for an optional argument, your method should too.

## Lint/TopLevelReturnWithArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.89 | - |

Checks for top level return with arguments. If there is a top-level return statement with an argument, then the argument is always ignored. This is detected automatically since Ruby 2.7.

## Lint/TrailingCommaInAttributeDeclaration

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.90 | 1.61 |

Checks for trailing commas in attribute declarations, such as
`#attr_reader`. Leaving a trailing comma will nullify the next method
definition by overriding it with a getter method.

## Lint/TripleQuotes

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.9 | - |

Checks for "triple quotes" (strings delimited by any odd number of quotes greater than 1).

Ruby allows multiple strings to be implicitly concatenated by just
being adjacent in a statement (ie. `"foo""bar" == "foobar"`). This sometimes
gives the impression that there is something special about triple quotes, but
in fact it is just extra unnecessary quotes and produces the same string. Each
pair of quotes produces an additional concatenated empty string, so the result
is still only the "actual" string within the delimiters.

| Although this cop is called triple quotes, the same behavior is present for strings delimited by 5, 7, etc. quotation marks. |

## Lint/UnderscorePrefixedVariableName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.21 | - |

Checks for underscore-prefixed variables that are actually used.

Since block keyword arguments cannot be arbitrarily named at call
sites, the `AllowKeywordBlockArguments` will allow use of underscore-
prefixed block keyword arguments.

### Examples

## Lint/UnescapedBracketInRegexp

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.68 | - |

Checks for Regexps (both literals and via `Regexp.new` / `Regexp.compile`)
that contain unescaped `]` characters.

It emulates the following Ruby warning:

```
$ ruby -e '/abc]123/'
-e:1: warning: regular expression has ']' without escape: /abc]123/
```
## Lint/UnexpectedBlockArity

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | No | 1.5 | - |

Checks for a block that is known to need more positional
block arguments than are given (by default this is configured for
`Enumerable` methods needing 2 arguments). Optional arguments are allowed,
although they don’t generally make sense as the default value will
be used. Blocks that have no receiver, or take splatted arguments
(ie. `*args`) are always accepted.

Keyword arguments (including `**kwargs`) do not get counted towards
this, as they are not used by the methods in question.

Method names and their expected arity can be configured like this:

```
Methods:
 inject: 2
 reduce: 2
```
### Safety

This cop matches for method names only and hence cannot tell apart methods with same name in different classes, which may lead to a false positive.

## Lint/UnmodifiedReduceAccumulator

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.1 | 1.5 |

Looks for `reduce` or `inject` blocks where the value returned (implicitly or
explicitly) does not include the accumulator. A block is considered valid as
long as at least one return value includes the accumulator.

If the accumulator is not included in the return value, then the entire block will just return a transformation of the last element value, and could be rewritten as such without a loop.

Also catches instances where an index of the accumulator is returned, as this may change the type of object being retained.

| For the purpose of reducing false positives, this cop only flags
returns in `reduce`blocks where the element is the only variable in
the expression (since we will not be able to tell what other variables
relate to via static analysis). |

### Examples

```
# bad
(1..4).reduce(0) do |acc, el|
 el * 2
end
# bad, may raise a NoMethodError after the first iteration
%w(a b c).reduce({}) do |acc, letter|
 acc[letter] = true
end
# good
(1..4).reduce(0) do |acc, el|
 acc + el * 2
end
# good, element is returned but modified using the accumulator
values.reduce do |acc, el|
 el << acc
 el
end
# good, returns the accumulator instead of the index
%w(a b c).reduce({}) do |acc, letter|
 acc[letter] = true
 acc
end
# good, at least one branch returns the accumulator
values.reduce(nil) do |result, value|
 break result if something?
 value
end
# good, recursive
keys.reduce(self) { |result, key| result[key] }
# ignored as the return value cannot be determined
enum.reduce do |acc, el|
 x = foo(acc, el)
 bar(x)
end
```
## Lint/UnreachableCode

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.9 | - |

Checks for unreachable code.
The check is based on the presence of flow-of-control
statements in non-final position in `begin` (implicit) blocks.

## Lint/UnreachableLoop

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.89 | 1.7 |

Checks for loops that will have at most one iteration.

A loop that can never reach the second iteration is a possible error in the code.
In rare cases where only one iteration (or at most one iteration) is intended behavior,
the code should be refactored to use `if` conditionals.

| Block methods that are used with `Enumerable`s are considered to be loops. |

`AllowedPatterns` can be used to match against the block receiver in order to allow
code that would otherwise be registered as an offense (eg. `times` used not in an
`Enumerable` context).

### Examples

```
# bad
while node
 do_something(node)
 node = node.parent
 break
end
# good
while node
 do_something(node)
 node = node.parent
end
# bad
def verify_list(head)
 item = head
 begin
 if verify(item)
 return true
 else
 return false
 end
 end while(item)
end
# good
def verify_list(head)
 item = head
 begin
 if verify(item)
 item = item.next
 else
 return false
 end
 end while(item)
 true
end
# bad
def find_something(items)
 items.each do |item|
 if something?(item)
 return item
 else
 raise NotFoundError
 end
 end
end
# good
def find_something(items)
 items.each do |item|
 if something?(item)
 return item
 end
 end
 raise NotFoundError
end
# bad
2.times { raise ArgumentError }
```
## Lint/UnreachablePatternBranch

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.85 | - |

Checks for unreachable `in` pattern branches in `case…in` statements.

An `in` branch is unreachable when a previous branch uses an unguarded
catch-all pattern that matches any value unconditionally. Any `in` branches
(and `else`) that follow such a catch-all are dead code.

A catch-all pattern is one of:

-
A bare variable capture ( `in x`)
-
An underscore ( `in _`)
-
A pattern alias where the left side is a catch-all ( `in _ ⇒ y`)
-
An alternation pattern where at least one alternative is a catch-all ( `in _ | Integer`)

| A catch-all pattern with a guard clause (e.g., `in _ if condition`)
does NOT make subsequent branches unreachable because the guard might
not be satisfied. |

### Examples

```
# bad
case value
in Integer
 handle_integer
in x
 handle_other
in String
 handle_string
else
 handle_else
end
# good
case value
in Integer
 handle_integer
in String
 handle_string
in x
 handle_other
end
# bad - else is unreachable after catch-all
case value
in Integer
 handle_integer
in _
 handle_other
else
 handle_else
end
# good - guard clause means catch-all might not match
case value
in x if x.positive?
 handle_positive
in Integer
 handle_integer
else
 handle_other
end
```
## Lint/UnusedBlockArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.21 | 1.61 |

Checks for unused block arguments.

### Examples

```
# bad
do_something do |used, unused|
 puts used
end
do_something do |bar|
 puts :foo
end
define_method(:foo) do |bar|
 puts :baz
end
# good
do_something do |used, _unused|
 puts used
end
do_something do
 puts :foo
end
define_method(:foo) do |_bar|
 puts :baz
end
```
## Lint/UnusedMethodArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.21 | 1.69 |

Checks for unused method arguments.

### Examples

```
# bad
def some_method(used, unused, _unused_but_allowed)
 puts used
end
# good
def some_method(used, _unused, _unused_but_allowed)
 puts used
end
```
#### IgnoreNotImplementedMethods: true (default)

```
# with default value of `NotImplementedExceptions: ['NotImplementedError']`
# good
def do_something(unused)
 raise NotImplementedError
end
def do_something_else(unused)
 fail "TODO"
end
```
## Lint/UnusedPrivateMethod

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 1.89 | - |

Checks for private instance methods that are not referenced anywhere in the project.

The check is powered by the project-wide index, so it only runs when
`AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed.
Without the index the cop does nothing.

A method counts as referenced when a call with its name appears anywhere
in the indexed project (regardless of the receiver), when it is the source
of an `alias`, or when its name appears in the same file as a symbol or
inside a string literal (covering `send(:name)` and declarative DSLs like
`before_action :name`). Methods defined in classes or modules with
descendants are not checked, since they may be invoked through `super`
or inherited dispatch, and neither are methods whose names are built
dynamically (e.g. `send("do_#{action}")`).

The cop is disabled by default because symbol-based references from
**other** files (e.g. a Rails callback declared in a concern) cannot be
detected and would be reported as false positives. It is best suited
for occasional dead-code sweeps rather than permanent enforcement.

## Lint/UriEscapeUnescape

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | - |

Identifies places where `URI.escape` can be replaced by
`CGI.escape`, `URI.encode_www_form`, or `URI.encode_www_form_component`
depending on your specific use case.
Also this cop identifies places where `URI.unescape` can be replaced by
`CGI.unescape`, `URI.decode_www_form`,
or `URI.decode_www_form_component` depending on your specific use case.

### Examples

```
# bad
URI.escape('http://example.com')
URI.encode('http://example.com')
# good
CGI.escape('http://example.com')
URI.encode_uri_component(uri) # Since Ruby 3.1
URI.encode_www_form([['example', 'param'], ['lang', 'en']])
URI.encode_www_form(page: 10, locale: 'en')
URI.encode_www_form_component('http://example.com')
# bad
URI.unescape(enc_uri)
URI.decode(enc_uri)
# good
CGI.unescape(enc_uri)
URI.decode_uri_component(uri) # Since Ruby 3.1
URI.decode_www_form(enc_uri)
URI.decode_www_form_component(enc_uri)
```
## Lint/UriRegexp

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.50 | - |

Identifies places where `URI.regexp` is obsolete and should not be used.

For Ruby 3.3 or lower, use `URI::DEFAULT_PARSER.make_regexp`.
For Ruby 3.4 or higher, use `URI::RFC2396_PARSER.make_regexp`.

| If you need to support both Ruby 3.3 and lower as well as Ruby 3.4 and higher, consider manually changing the code as follows: |

`defined?(URI::RFC2396_PARSER) ? URI::RFC2396_PARSER : URI::DEFAULT_PARSER`## Lint/UselessAccessModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.20 | 1.61 |

Checks for redundant access modifiers, including those with no
code, those which are repeated, those which are on top-level, and
leading `public` modifiers in a class or module body.
Conditionally-defined methods are considered as always being defined,
and thus access modifiers guarding such methods are not redundant.

This cop has `ContextCreatingMethods` option. The default setting value
is an empty array that means no method is specified.
This setting is an array of methods which, when called, are known to
create its own context in the module’s current access context.

It also has `MethodCreatingMethods` option. The default setting value
is an empty array that means no method is specified.
This setting is an array of methods which, when called, are known to
create other methods in the module’s current access context.

### Examples

```
# bad
class Foo
 public # this is redundant (default access is public)
 def method
 end
end
# bad
class Foo
 # The following is redundant (methods defined on the class'
 # singleton class are not affected by the private modifier)
 private
 def self.method3
 end
end
# bad
class Foo
 protected
 define_method(:method2) do
 end
 protected # this is redundant (repeated from previous modifier)
 [1,2,3].each do |i|
 define_method("foo#{i}") do
 end
 end
end
# bad
class Foo
 private # this is redundant (no following methods are defined)
end
# bad
private # this is useless (access modifiers have no effect on top-level)
def method
end
# good
class Foo
 private # this is not redundant (a method is defined)
 def method2
 end
end
# good
class Foo
 # The following is not redundant (conditionally defined methods are
 # considered as always defining a method)
 private
 if condition?
 def method
 end
 end
end
# good
class Foo
 protected # this is not redundant (a method is defined)
 define_method(:method2) do
 end
end
```
#### ContextCreatingMethods: concerning

```
# Lint/UselessAccessModifier:
# ContextCreatingMethods:
# - concerning
# good
require 'active_support/concern'
class Foo
 concerning :Bar do
 def some_public_method
 end
 private
 def some_private_method
 end
 end
 # this is not redundant because `concerning` created its own context
 private
 def some_other_private_method
 end
end
```
## Lint/UselessAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.11 | 1.66 |

Checks for every useless assignment to local variable in every
scope.
The basic idea for this cop was from the warning of `ruby -cw`:

`assigned but unused variable - foo`Currently this cop has advanced logic that detects unreferenced reassignments and properly handles varied cases such as branch, loop, rescue, ensure, etc.

This cop’s autocorrection avoids cases like `a ||= 1` because removing assignment from
operator assignment can cause `NameError` if this assignment has been used to declare
a local variable. For example, replacing `a ||= 1` with `a || 1` may cause
"undefined local variable or method `a' for main:Object (NameError)".

| Given the assignment `foo = 1, bar = 2`, removing unused variables
can lead to a syntax error, so this case is not autocorrected. |

## Lint/UselessConstantScoping

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.72 | - |

Checks for useless constant scoping. Private constants must be defined using
`private_constant`. Even if `private` access modifier is used, it is public scope despite
its appearance.

It does not support autocorrection due to behavior change and multiple ways to fix it. Or a public constant may be intended.

## Lint/UselessDefaultValueArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.76 | - |

Checks for usage of method `fetch` or `Array.new` with default value argument
and block. In such cases, block will always be used as default value.

This cop emulates Ruby warning "block supersedes default value argument" which
applies to `Array.new`, `Array#fetch`, `Hash#fetch`, `ENV.fetch` and
`Thread#fetch`.

A `fetch` call without a receiver is considered a custom method and does not register
an offense.

### Safety

This cop is unsafe because the receiver could have nonstandard implementation
of `fetch`, or be a class other than the one listed above.

It is also unsafe because default value argument could have side effects:

```
def x(a) = puts "side effect"
Array.new(5, x(1)) { 2 }
```
so removing it would change behavior.

### Examples

```
# bad
x.fetch(key, default_value) { block_value }
Array.new(size, default_value) { block_value }
# good
x.fetch(key) { block_value }
Array.new(size) { block_value }
# also good - in case default value argument is desired instead
x.fetch(key, default_value)
Array.new(size, default_value)
# good - keyword arguments aren't registered as offenses
x.fetch(key, keyword: :arg) { block_value }
```
## Lint/UselessDefined

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.69 | - |

Checks for calls to `defined?` with strings or symbols as the argument.
Such calls will always return `'expression'`, you probably meant to
check for the existence of a constant, method, or variable instead.

`defined?` is part of the Ruby syntax and doesn’t behave like normal methods.
You can safely pass in what you are checking for directly, without encountering
a `NameError`.

When interpolation is used, oftentimes it is not possible to write the
code with `defined?`. In these cases, switch to one of the more specific methods:

-
`class_variable_defined?`
-
`const_defined?`
-
`method_defined?`
-
`instance_variable_defined?`
-
`binding.local_variable_defined?`

## Lint/UselessElseWithoutRescue

| Requires Ruby version ⇐ 2.5 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.17 | 1.31 |

Checks for useless `else` in `begin..end` without `rescue`.

| This syntax is no longer valid on Ruby 2.6 or higher. |

## Lint/UselessMethodDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Command-line only (Unsafe) | 0.90 | 1.61 |

Checks for useless method definitions, specifically: empty constructors
and methods just delegating to `super`.

## Lint/UselessNumericOperation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.66 | - |

Certain numeric operations have no impact, being: Adding or subtracting 0, multiplying or dividing by 1 or raising to the power of 1. These are probably leftover from debugging, or are mistakes.

## Lint/UselessOr

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Command-line only (Unsafe) | 1.76 | 1.82 |

Checks for useless OR (`||` and `or`) expressions.

Some methods always return a truthy value, even when called
on `nil` (e.g. `nil.to_i` evaluates to `0`). Therefore, OR expressions
appended after these methods will never evaluate.

### Safety

As shown in the examples below, there are generally two possible ways to correct the offense, but this cop’s autocorrection always chooses the option that preserves the current behavior. While this does not change how the code behaves, that option is not necessarily the appropriate fix in every situation. For this reason, the autocorrection provided by this cop is considered unsafe.

### Examples

```
# bad
x.to_a || fallback
x.to_c || fallback
x.to_d || fallback
x.to_i || fallback
x.to_f || fallback
x.to_h || fallback
x.to_r || fallback
x.to_s || fallback
x.to_sym || fallback
x.intern || fallback
x.inspect || fallback
x.hash || fallback
x.object_id || fallback
x.__id__ || fallback
x.to_s or fallback
# good - if fallback is same as return value of method called on nil
x.to_a # nil.to_a returns []
x.to_c # nil.to_c returns (0+0i)
x.to_d # nil.to_d returns 0.0
x.to_i # nil.to_i returns 0
x.to_f # nil.to_f returns 0.0
x.to_h # nil.to_h returns {}
x.to_r # nil.to_r returns (0/1)
x.to_s # nil.to_s returns ''
x.to_sym # nil.to_sym raises an error
x.intern # nil.intern raises an error
x.inspect # nil.inspect returns "nil"
x.hash # nil.hash returns an Integer
x.object_id # nil.object_id returns an Integer
x.__id__ # nil.object_id returns an Integer
# good - if the intention is not to call the method on nil
x&.to_a || fallback
x&.to_c || fallback
x&.to_d || fallback
x&.to_i || fallback
x&.to_f || fallback
x&.to_h || fallback
x&.to_r || fallback
x&.to_s || fallback
x&.to_sym || fallback
x&.intern || fallback
x&.inspect || fallback
x&.hash || fallback
x&.object_id || fallback
x&.__id__ || fallback
x&.to_s or fallback
```
## Lint/UselessRescue

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.43 | - |

Checks for useless `rescue`s, which only reraise rescued exceptions.

### Examples

```
# bad
def foo
 do_something
rescue
 raise
end
# bad
def foo
 do_something
rescue => e
 raise # or 'raise e', or 'raise $!', or 'raise $ERROR_INFO'
end
# good
def foo
 do_something
rescue
 do_cleanup
 raise
end
# bad (latest rescue)
def foo
 do_something
rescue ArgumentError
 # noop
rescue
 raise
end
# good (not the latest rescue)
def foo
 do_something
rescue ArgumentError
 raise
rescue
 # noop
end
```
## Lint/UselessRuby2Keywords

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.23 | - |

Looks for `ruby2_keywords` calls for methods that do not need it.

`ruby2_keywords` should only be called on methods that accept an argument splat
(`*args`) but do not have explicit keyword arguments (`k:` or `k: true`) or
a keyword splat (`**kwargs`).

### Examples

```
# good (splat argument without keyword arguments)
ruby2_keywords def foo(*args); end
# bad (no arguments)
ruby2_keywords def foo; end
# good
def foo; end
# bad (positional argument)
ruby2_keywords def foo(arg); end
# good
def foo(arg); end
# bad (double splatted argument)
ruby2_keywords def foo(**args); end
# good
def foo(**args); end
# bad (keyword arguments)
ruby2_keywords def foo(i:, j:); end
# good
def foo(i:, j:); end
# bad (splat argument with keyword arguments)
ruby2_keywords def foo(*args, i:, j:); end
# good
def foo(*args, i:, j:); end
# bad (splat argument with double splat)
ruby2_keywords def foo(*args, **kwargs); end
# good
def foo(*args, **kwargs); end
# bad (ruby2_keywords given a symbol)
def foo; end
ruby2_keywords :foo
# good
def foo; end
# bad (ruby2_keywords with dynamic method)
define_method(:foo) { |arg| }
ruby2_keywords :foo
# good
define_method(:foo) { |arg| }
```
## Lint/UselessSetterCall

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.13 | 1.2 |

Checks for setter call to local variable as the final expression of a function definition.

## Lint/UselessTimes

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Command-line only (Unsafe) | 0.91 | 1.61 |

Checks for uses of `Integer#times` that will never yield
(when the integer `⇐ 0`) or that will only ever yield once
(`1.times`).

## Lint/Void

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.9 | 1.61 |

Checks for operators, variables, literals, lambda, proc and nonmutating methods used in void context.

`each` blocks are allowed to prevent false positives.
For example, the expression inside the `each` block below.
It’s not void, especially when the receiver is an `Enumerator`:

```
enumerator = [1, 2, 3].filter
enumerator.each { |item| item >= 2 } #=> [2, 3]
```
| The last expression in an assignment method definition such as `def foo=(arg)`is not flagged. Ruby discards it (the method returns its argument), but the method can
still be called directly and its return value relied upon, so flagging it would be a
false positive for this lint. |

| A constant used in a void context is flagged but not autocorrected, since referencing a constant can trigger autoloading side effects (e.g. forcing a file to load before a monkey-patch), so removing it may change behavior. |

# Metrics

## Metrics/AbcSize

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.27 | 1.5 |

Checks that the ABC size of methods is not higher than the configured maximum. The ABC size is based on assignments, branches (method calls), and conditions. See http://c2.com/cgi/wiki?AbcMetric and https://en.wikipedia.org/wiki/ABC_Software_Metric.

Interpreting ABC size:

-
`⇐ 17`satisfactory
-
`18..30`unsatisfactory
-
`>`30 dangerous

You can have repeated "attributes" calls count as a single "branch".
For this purpose, attributes are any method with no argument; no attempt
is meant to distinguish actual `attr_reader` from other methods.

This cop also takes into account `AllowedMethods` (defaults to `[]`)
And `AllowedPatterns` (defaults to `[]`)

### Examples

#### CountRepeatedAttributes: false (default is true)

```
# `model` and `current_user`, referenced 3 times each,
# are each counted as only 1 branch each if
# `CountRepeatedAttributes` is set to 'false'
def search
 @posts = model.active.visible_by(current_user)
 .search(params[:q])
 @posts = model.some_process(@posts, current_user)
 @posts = model.another_process(@posts, current_user)
 render 'pages/search/page'
end
```
## Metrics/BlockLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.44 | 1.5 |

Checks if the length of a block exceeds some maximum value. Comment lines can optionally be ignored. The maximum allowed length is configurable. The cop can be configured to ignore blocks passed to certain methods.

You can set constructs you want to fold with `CountAsOne`.

Available are: 'array', 'hash', 'heredoc', and 'method_call'. Each construct will be counted as one line regardless of its actual size.

| This cop does not apply for `Struct`definitions. |

| The `ExcludedMethods`configuration is deprecated and only kept
for backwards compatibility. Please use`AllowedMethods`and`AllowedPatterns`instead. By default, there are no allowed methods. |

## Metrics/BlockNesting

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.25 | 1.65 |

Checks for excessive nesting of conditional and looping constructs. Deeply nested code is harder to read, understand, and maintain. Extracting nested logic into methods improves clarity.

You can configure if blocks are considered using the `CountBlocks` and `CountModifierForms`
options. When both are set to `false` (the default) blocks and modifier forms are not
counted towards the nesting level. Set them to `true` to include these in the nesting level
calculation as well.

The maximum level of nesting allowed is configurable.

## Metrics/ClassLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.25 | 0.87 |

Checks if the length of a class exceeds some maximum value. Comment lines can optionally be ignored. The maximum allowed length is configurable.

You can set constructs you want to fold with `CountAsOne`.

Available are: 'array', 'hash', 'heredoc', and 'method_call'. Each construct will be counted as one line regardless of its actual size.

| This cop also applies for `Struct`definitions. |

## Metrics/CollectionLiteralLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.47 | 1.88 |

Checks for literals with extremely many entries. This is indicative of configuration or data that may be better extracted somewhere else, like a database, fetched from an API, or read from a non-code file (CSV, JSON, YAML, etc.).

### Examples

```
# bad
# Huge Array literal
[1, 2, '...', 999_999_999]
# bad
# Huge Hash literal
{ 1 => 1, 2 => 2, '...' => '...', 999_999_999 => 999_999_999}
# bad
# Huge Set "literal"
Set[1, 2, '...', 999_999_999]
# good
# Reasonably sized Array literal
[1, 2, '...', 10]
# good
# Reading huge Array from external data source
# File.readlines('numbers.txt', chomp: true).map!(&:to_i)
# good
# Reasonably sized Hash literal
{ 1 => 1, 2 => 2, '...' => '...', 10 => 10}
# good
# Reading huge Hash from external data source
CSV.foreach('numbers.csv', headers: true).each_with_object({}) do |row, hash|
 hash[row["key"].to_i] = row["value"].to_i
end
# good
# Reasonably sized Set "literal"
Set[1, 2, '...', 10]
# good
# Reading huge Set from external data source
SomeFramework.config_for(:something)[:numbers].to_set
```
## Metrics/CyclomaticComplexity

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.25 | 0.81 |

Checks that the cyclomatic complexity of methods is not higher than the configured maximum. The cyclomatic complexity is the number of linearly independent paths through a method. The algorithm counts decision points and adds one.

An if statement (or unless or ?:) increases the complexity by one. An
else branch does not, since it doesn’t add a decision point. The &&
operator (or keyword and) can be converted to a nested if statement,
and ||/or is shorthand for a sequence of ifs, so they also add one.
Loops can be said to have an exit condition, so they add one.
Blocks that are calls to builtin iteration methods
(e.g. `ary.map{…}`) also add one, others are ignored.

### Examples

```
def each_child_node(*types) # count begins: 1
 unless block_given? # unless: +1
 return to_enum(__method__, *types)
 end
 children.each do |child| # each{}: +1
 next unless child.is_a?(Node) # unless: +1
 yield child if types.empty? || # if: +1, ||: +1
 types.include?(child.type)
 end
 self
end # total: 6
```
## Metrics/MethodLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.25 | 1.5 |

Checks if the length of a method exceeds some maximum value. Comment lines can optionally be allowed. The maximum allowed length is configurable.

You can set constructs you want to fold with `CountAsOne`.

Available are: 'array', 'hash', 'heredoc', and 'method_call'. Each construct will be counted as one line regardless of its actual size.

| The `ExcludedMethods`and`IgnoredMethods`configuration is
deprecated and only kept for backwards compatibility.
Please use`AllowedMethods`and`AllowedPatterns`instead.
By default, there are no allowed methods. |

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| CountComments |
 | Boolean |
| Max |
 | Integer |
| CountAsOne |
 | Array |
| AllowedMethods |
 | Array |
| AllowedPatterns |
 | Array |

## Metrics/ModuleLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.31 | 0.87 |

Checks if the length of a module exceeds some maximum value. Comment lines can optionally be ignored. The maximum allowed length is configurable.

You can set constructs you want to fold with `CountAsOne`.

Available are: 'array', 'hash', 'heredoc', and 'method_call'. Each construct will be counted as one line regardless of its actual size.

## Metrics/ParameterLists

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.25 | 1.5 |

Checks for methods with too many parameters.

The maximum number of parameters is configurable. Keyword arguments can optionally be excluded from the total count, as they add less complexity than positional or optional parameters.

Any number of arguments for `initialize` method inside a block of
`Struct.new` and `Data.define` like this is always allowed:

```
Struct.new(:one, :two, :three, :four, :five, keyword_init: true) do
 def initialize(one:, two:, three:, four:, five:)
 end
end
```
This is because checking the number of arguments of the `initialize` method
does not make sense.

| Explicit block argument `&block`is not counted to prevent
erroneous change that is avoided by making block argument implicit. |

This cop also checks for the maximum number of optional parameters.
This can be configured using the `MaxOptionalParameters` config option.

### Examples

#### CountKeywordArgs: true (default)

```
# counts keyword args towards the maximum
# bad (assuming Max is 3)
def foo(a, b, c, d: 1)
end
# good (assuming Max is 3)
def foo(a, b, c: 1)
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| Max |
 | Integer |
| CountKeywordArgs |
 | Boolean |
| MaxOptionalParameters |
 | Integer |

## Metrics/PerceivedComplexity

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.25 | 1.88 |

Tries to produce a complexity score that’s a measure of the
complexity the reader experiences when looking at a method. For that
reason it considers `when` nodes as something that doesn’t add as much
complexity as an `if` or a `&&`. Except if it’s one of those special
`case`/`when` constructs where there’s no expression after `case`. Then
the cop treats it as an `if`/`elsif`/`elsif`… and lets all the `when`
nodes count. In contrast to the CyclomaticComplexity cop, this cop
considers `else` nodes as adding complexity.

A `case`/`in` branch whose pattern is a simple literal (e.g. `in 1`, `in "red"`, `in 1..10`)
or a constant/type (e.g. `in Integer`) and has no guard is just as easy to read as a `when`
branch, so it is discounted the same way. Branches with structural patterns (e.g. array,
hash, or find patterns), bindings, alternatives, or a guard add the full complexity of
a decision point.

### Examples

```
def example_1 # 1
 if cond # 1
 case var # 2 (0.8 + 4 * 0.2, rounded)
 when 1 then func_one
 when 2 then func_two
 when 3 then func_three
 when 4..10 then func_other
 end
 else # 1
 do_something until a && b # 2
 end # ===
end # 7 complexity points
def example_2 # 1
 case color # 1 (3 * 0.2, rounded)
 in "red" then func_red
 in "blue" then func_blue
 in "green" then func_green
 end # ===
end # 2 complexity points
```

Migration Migration/DepartmentName Enabled by default Safe Supports autocorrection Version Added Version Changed Enabled Yes Always 0.75 - Checks that cop names in rubocop:disable comments are given with department name. Examples # bad # rubocop:disable AbcSize # good # rubocop:disable Metrics/AbcSize Metrics Naming

# Naming

## Naming/AccessorMethodName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 1.89 |

Avoid prefixing accessor method names with `get_` or `set_`.
Applies to both instance and class methods.

| Method names starting with `get_`or`set_`only register an offense
when the methods match the expected arity for getters and setters respectively.
Getters (`get_attribute`) must have no arguments to be registered,
and setters (`set_attribute(value)`) must have exactly one. |

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, methods that override a method defined by an ancestor
elsewhere in the project are not reported, since renaming an override
breaks the inherited contract.

## Naming/AsciiIdentifiers

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 0.87 |

Checks for non-ascii characters in identifier and constant names. Identifiers are always checked and whether constants are checked can be controlled using AsciiConstants config.

### Examples

```
# bad
def καλημερα # Greek alphabet (non-ascii)
end
# bad
def こんにちはと言う # Japanese character (non-ascii)
end
# bad
def hello_🍣 # Emoji (non-ascii)
end
# good
def say_hello
end
# bad
신장 = 10 # Hangul character (non-ascii)
# good
height = 10
# bad
params[:عرض_gteq] # Arabic character (non-ascii)
# good
params[:width_gteq]
```
## Naming/BinaryOperatorParameterName

## Naming/BlockForwarding

| Requires Ruby version 3.1 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.24 | - |

In Ruby 3.1, anonymous block forwarding has been added.

This cop identifies places where `do_something(&block)` can be replaced
by `do_something(&)`.

It also supports the opposite style by alternative `explicit` option.
You can specify the block variable name for autocorrection with `BlockForwardingName`.
The default variable name is `block`. If the name is already in use, it will not be
autocorrected.

| Because of a bug in Ruby 3.3.0, when a block is referenced inside of another block, no offense will be registered until Ruby 3.4: |

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| BlockForwardingName |
 | String |

## Naming/BlockParameterName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.53 | 0.77 |

Checks block parameter names for how descriptive they are. It is highly configurable.

The `MinNameLength` config option takes an integer. It represents
the minimum amount of characters the name must be. Its default is 1.
The `AllowNamesEndingInNumbers` config option takes a boolean. When
set to false, this cop will register offenses for names ending with
numbers. Its default is false. The `AllowedNames` config option
takes an array of permitted names that will never register an
offense. The `ForbiddenNames` config option takes an array of
restricted names that will always register an offense.

### Examples

```
# bad
bar do |varOne, varTwo|
 varOne + varTwo
end
# With `AllowNamesEndingInNumbers` set to false
foo { |num1, num2| num1 * num2 }
# With `MinNameLength` set to number greater than 1
baz { |a, b, c| do_stuff(a, b, c) }
# good
bar do |thud, fred|
 thud + fred
end
foo { |speed, distance| speed * distance }
baz { |age, height, gender| do_stuff(age, height, gender) }
```
## Naming/ClassAndModuleCamelCase

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 0.85 |

Checks for class and module names with an underscore in them.

`AllowedNames` config takes an array of permitted names.
Its default value is `['module_parent']`.
These names can be full class/module names or part of the name.
eg. Adding `my_class` to the `AllowedNames` config will allow names like
`my_class`, `my_class::User`, `App::my_class`, `App::my_class::User`, etc.

### Examples

```
# bad
class My_Class
end
module My_Module
end
# good
class MyClass
end
module MyModule
end
class module_parent::MyModule
end
```
## Naming/ConstantName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | - |

Checks whether constant names are written using SCREAMING_SNAKE_CASE.

To avoid false positives, it ignores cases in which we cannot know for certain the type of value that would be assigned to a constant.

## Naming/FileName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 1.23 |

Makes sure that Ruby source files have snake_case names. Ruby scripts (i.e. source files with a shebang in the first line) are ignored.

The cop also ignores `.gemspec` files, because Bundler
recommends using dashes to separate namespaces in nested gems
(i.e. `bundler-console` becomes `Bundler::Console`). As such, the
gemspec is supposed to be named `bundler-console.gemspec`.

When `ExpectMatchingDefinition` (default: `false`) is `true`, the cop requires
each file to have a class, module or `Struct` defined in it that matches
the filename. This can be further configured using
`CheckDefinitionPathHierarchy` (default: `true`) to determine whether the
path should match the namespace of the above definition.

When `IgnoreExecutableScripts` (default: `true`) is `true`, files that start
with a shebang line are not considered by the cop.

When `Regex` is set, the cop will flag any filename that does not match
the regular expression.

### Examples

```
# bad
lib/layoutManager.rb
anything/usingCamelCase
# good
lib/layout_manager.rb
anything/using_snake_case.rake
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| Exclude |
 | Array |
| ExpectMatchingDefinition |
 | Boolean |
| CheckDefinitionPathHierarchy |
 | Boolean |
| CheckDefinitionPathHierarchyRoots |
 | Array |
| Regex |
 | |
| IgnoreExecutableScripts |
 | Boolean |
| AllowedAcronyms |
 | Array |

## Naming/HeredocDelimiterCase

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.50 | 1.2 |

Checks that your heredocs are using the configured case. By default it is configured to enforce uppercase heredocs.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Naming/HeredocDelimiterNaming

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | - |

Checks that your heredocs are using meaningful delimiters.
By default it disallows `END` and `EO*`, and can be configured through
forbidden listing additional delimiters.

### Examples

```
# good
<<-SQL
 SELECT * FROM foo
SQL
# bad
<<-END
 SELECT * FROM foo
END
# bad
<<-EOS
 SELECT * FROM foo
EOS
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| ForbiddenDelimiters |
 | Array |

## Naming/InclusiveLanguage

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 1.18 | 1.49 |

Recommends the use of inclusive language instead of problematic terms. The cop can check the following locations for offenses:

-
identifiers
-
constants
-
variables
-
strings
-
symbols
-
comments
-
file paths

Each of these locations can be individually enabled/disabled via configuration, for example CheckIdentifiers = true/false.

Flagged terms are configurable for the cop. For each flagged term an optional
Regex can be specified to identify offenses. Suggestions for replacing a flagged term can
be configured and will be displayed as part of the offense message.
An AllowedRegex can be specified for a flagged term to exempt allowed uses of the term.
`WholeWord: true` can be set on a flagged term to indicate the cop should only match when
a term matches the whole word (partial matches will not be offenses).

The cop supports autocorrection when there is only one suggestion. When there are multiple suggestions, the best suggestion cannot be identified and will not be autocorrected.

### Examples

#### FlaggedTerms: { whitelist: { Suggestions: ['allowlist'] } }

```
# Suggest replacing identifier whitelist with allowlist
# bad
whitelist_users = %w(user1 user1)
# good
allowlist_users = %w(user1 user2)
```
#### FlaggedTerms: { master: { Suggestions: ['main', 'primary', 'leader'] } }

```
# Suggest replacing master in an instance variable name with main, primary, or leader
# bad
@master_node = 'node1.example.com'
# good
@primary_node = 'node1.example.com'
```
#### FlaggedTerms: { whitelist: { Regex: !ruby/regexp '/white[-_\s]?list' } }

```
# Identify problematic terms using a Regexp
# bad
white_list = %w(user1 user2)
# good
allow_list = %w(user1 user2)
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| CheckIdentifiers |
 | Boolean |
| CheckConstants |
 | Boolean |
| CheckVariables |
 | Boolean |
| CheckStrings |
 | Boolean |
| CheckSymbols |
 | Boolean |
| CheckComments |
 | Boolean |
| CheckFilepaths |
 | Boolean |
| FlaggedTerms |
 |

## Naming/MemoizedInstanceVariableName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.53 | 1.2 |

Checks for memoized methods whose instance variable name
does not match the method name. Applies to both regular methods
(defined with `def`) and dynamic methods (defined with
`define_method` or `define_singleton_method`).

This cop can be configured with the EnforcedStyleForLeadingUnderscores directive. It can be configured to allow for memoized instance variables prefixed with an underscore. Prefixing ivars with an underscore is a convention that is used to implicitly indicate that an ivar should not be set or referenced outside of the memoization method.

### Safety

This cop relies on the pattern `@instance_var ||= …`,
but this is sometimes used for other purposes than memoization
so this cop is considered unsafe. Also, its autocorrection is unsafe
because it may conflict with instance variable names already in use.

### Examples

#### EnforcedStyleForLeadingUnderscores: disallowed (default)

```
# bad
# Method foo is memoized using an instance variable that is
# not `@foo`. This can cause confusion and bugs.
def foo
 @something ||= calculate_expensive_thing
end
def foo
 return @something if defined?(@something)
 @something = calculate_expensive_thing
end
# good
def _foo
 @foo ||= calculate_expensive_thing
end
# good
def foo
 @foo ||= calculate_expensive_thing
end
# good
def foo
 @foo ||= begin
 calculate_expensive_thing
 end
end
# good
def foo
 helper_variable = something_we_need_to_calculate_foo
 @foo ||= calculate_expensive_thing(helper_variable)
end
# good
define_method(:foo) do
 @foo ||= calculate_expensive_thing
end
# good
define_method(:foo) do
 return @foo if defined?(@foo)
 @foo = calculate_expensive_thing
end
```
#### EnforcedStyleForLeadingUnderscores: required

```
# bad
def foo
 @something ||= calculate_expensive_thing
end
# bad
def foo
 @foo ||= calculate_expensive_thing
end
def foo
 return @foo if defined?(@foo)
 @foo = calculate_expensive_thing
end
# good
def foo
 @_foo ||= calculate_expensive_thing
end
# good
def _foo
 @_foo ||= calculate_expensive_thing
end
def foo
 return @_foo if defined?(@_foo)
 @_foo = calculate_expensive_thing
end
# good
define_method(:foo) do
 @_foo ||= calculate_expensive_thing
end
# good
define_method(:foo) do
 return @_foo if defined?(@_foo)
 @_foo = calculate_expensive_thing
end
```
#### EnforcedStyleForLeadingUnderscores: optional

```
# bad
def foo
 @something ||= calculate_expensive_thing
end
# good
def foo
 @foo ||= calculate_expensive_thing
end
# good
def foo
 @_foo ||= calculate_expensive_thing
end
# good
def _foo
 @_foo ||= calculate_expensive_thing
end
# good
def foo
 return @_foo if defined?(@_foo)
 @_foo = calculate_expensive_thing
end
# good
define_method(:foo) do
 @foo ||= calculate_expensive_thing
end
# good
define_method(:foo) do
 @_foo ||= calculate_expensive_thing
end
```
## Naming/MethodName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 1.75 |

Makes sure that all methods use the configured style, snake_case or camelCase, for their names.

Method names matching patterns are always allowed.

The cop can be configured with `AllowedPatterns` to allow certain regexp patterns:

```
Naming/MethodName:
 AllowedPatterns:
 - '\AonSelectionBulkChange\z'
 - '\AonSelectionCleared\z'
```
As well, you can also forbid specific method names or regexp patterns
using `ForbiddenIdentifiers` or `ForbiddenPatterns`:

```
Naming/MethodName:
 ForbiddenIdentifiers:
 - 'def'
 - 'super'
 ForbiddenPatterns:
 - '_v1\z'
 - '_gen1\z'
```
### Examples

#### EnforcedStyle: snake_case (default)

```
# bad
def fooBar; end
# good
def foo_bar; end
# bad
define_method :fooBar do
end
# good
define_method :foo_bar do
end
# bad
Struct.new(:fooBar)
# good
Struct.new(:foo_bar)
# bad
alias_method :fooBar, :some_method
# good
alias_method :foo_bar, :some_method
```
## Naming/MethodParameterName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.53 | 0.77 |

Checks method parameter names for how descriptive they are. It is highly configurable.

The `MinNameLength` config option takes an integer. It represents
the minimum amount of characters the name must be. Its default is 3.
The `AllowNamesEndingInNumbers` config option takes a boolean. When
set to false, this cop will register offenses for names ending with
numbers. Its default is false. The `AllowedNames` config option
takes an array of permitted names that will never register an
offense. The `ForbiddenNames` config option takes an array of
restricted names that will always register an offense.

### Examples

```
# bad
def bar(varOne, varTwo)
 varOne + varTwo
end
# With `AllowNamesEndingInNumbers` set to false
def foo(num1, num2)
 num1 * num2
end
# With `MinNameLength` set to number greater than 1
def baz(a, b, c)
 do_stuff(a, b, c)
end
# good
def bar(thud, fred)
 thud + fred
end
def foo(speed, distance)
 speed * distance
end
def baz(age_a, height_b, gender_c)
 do_stuff(age_a, height_b, gender_c)
end
```
## Naming/PredicateMethod

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.76 | 1.78 |

Checks that predicate methods end with `?` and non-predicate methods do not.

The names of predicate methods (methods that return a boolean value) should end in a question mark. Methods that don’t return a boolean shouldn’t end in a question mark.

The cop assesses a predicate method as one that returns boolean values. Likewise, a method that only returns literal values is assessed as non-predicate. Other predicate method calls are assumed to return boolean values. The cop does not make an assessment if the return type is unknown (non-predicate method calls, variables, etc.).

| The `initialize`method and operator methods (`def ==`, etc.) are ignored. |

By default, the cop runs in `conservative` mode, which allows a method to be named
with a question mark as long as at least one return value is boolean. In `aggressive`
mode, methods with a question mark will register an offense if any known non-boolean
return values are detected.

The cop also has `AllowedMethods` configuration in order to prevent the cop from
registering an offense from a method name that does not conform to the naming
guidelines. By default, `call` is allowed. The cop also has `AllowedPatterns`
configuration to allow method names by regular expression.

Although returning a call to another predicate method is treated as a boolean value,
certain method names can be known to not return a boolean, despite ending in a `?`
(for example, `Numeric#nonzero?` returns `self` or `nil`). These methods can be
configured using `NonBooleanPredicates`.

The cop can furthermore be configured to allow all bang methods (method names
ending with `!`), with `AllowBangMethods: true` (default false).

### Examples

#### Mode: conservative (default)

```
# bad
def foo
 bar == baz
end
# good
def foo?
 bar == baz
end
# bad
def foo?
 5
end
# good
def foo
 5
end
# bad
def foo
 x == y
end
# good
def foo?
 x == y
end
# bad
def foo
 !x
end
# good
def foo?
 !x
end
# bad - returns the value of another predicate method
def foo
 bar?
end
# good
def foo?
 bar?
end
# good - operator method
def ==(other)
 hash == other.hash
end
# good - at least one return value is boolean
def foo?
 return unless bar?
 true
end
# ok - return type is not known
def foo?
 bar
end
# ok - return type is not known
def foo
 bar?
end
```
## Naming/PredicatePrefix

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 1.89 |

Checks that predicate method names end with a question mark and do not start with a forbidden prefix.

A method is determined to be a predicate method if its name starts with
one of the prefixes listed in the `NamePrefix` configuration. The list
defaults to `is_`, `has_`, and `have_` but may be overridden.

Predicate methods must end with a question mark.

When `ForbiddenPrefixes` is also set (as it is by default), predicate
methods which begin with a forbidden prefix are not allowed, even if
they end with a `?`. These methods should be changed to remove the
prefix.

When `UseSorbetSigs` is set to true (optional), the cop will only report
offenses if the method has a Sorbet `sig` with a return type of
`T::Boolean`. Dynamic methods are not supported with this configuration.

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, methods that override a method defined by an ancestor
elsewhere in the project are not reported, since renaming an override
breaks the inherited contract.

### Examples

#### NamePrefix: ['is_', 'has_', 'have_'] (default)

```
# bad
def is_even(value)
end
# When ForbiddenPrefixes: ['is_', 'has_', 'have_'] (default)
# good
def even?(value)
end
# When ForbiddenPrefixes: []
# good
def is_even?(value)
end
```
#### NamePrefix: ['seems_to_be_']

```
# bad
def seems_to_be_even(value)
end
# When ForbiddenPrefixes: ['seems_to_be_']
# good
def even?(value)
end
# When ForbiddenPrefixes: []
# good
def seems_to_be_even?(value)
end
```
#### AllowedMethods: ['is_a?'] (default)

```
# Despite starting with the `is_` prefix, this method is allowed
# good
def is_a?(value)
end
```
#### UseSorbetSigs: false (default)

```
# bad
sig { returns(String) }
def is_this_thing_on
 "yes"
end
# good - Sorbet signature is not evaluated
sig { returns(String) }
def is_this_thing_on?
 "yes"
end
```
#### UseSorbetSigs: true

```
# bad
sig { returns(T::Boolean) }
def odd(value)
end
# good
sig { returns(T::Boolean) }
def odd?(value)
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| NamePrefix |
 | Array |
| ForbiddenPrefixes |
 | Array |
| AllowedMethods |
 | Array |
| MethodDefinitionMacros |
 | Array |
| UseSorbetSigs |
 | Boolean |
| Exclude |
 | Array |

## Naming/RescuedExceptionsVariableName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.67 | 0.68 |

Makes sure that rescued exception variables are named as expected.

The `PreferredName` config option takes a `String`. It represents
the required name of the variable. Its default is `e`.

| This cop does not consider nested rescues because it cannot guarantee that the variable from the outer rescue is not used within the inner rescue (in which case, changing the inner variable would shadow the outer variable). |

## Naming/VariableName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 1.73 |

Checks that the configured style (snake_case or camelCase) is used for all variable names. This includes local variables, instance variables, class variables, method arguments (positional, keyword, rest or block), and block arguments.

The cop can also be configured to forbid using specific names for variables, using
`ForbiddenIdentifiers` or `ForbiddenPatterns`. In addition to the above, this applies
to global variables as well.

Method definitions and method calls are not affected by this cop.

### Examples

## Naming/VariableNumber

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.50 | 1.4 |

Makes sure that all numbered variables use the configured style, snake_case, normalcase, or non_integer, for their numbering.

Additionally, `CheckMethodNames` and `CheckSymbols` configuration options
can be used to specify whether method names and symbols should be checked.
Both are enabled by default.

### Examples

#### EnforcedStyle: normalcase (default)

```
# bad
:some_sym_1
variable_1 = 1
def some_method_1; end
def some_method1(arg_1); end
# good
:some_sym1
variable1 = 1
def some_method1; end
def some_method1(arg1); end
```
#### EnforcedStyle: snake_case

```
# bad
:some_sym1
variable1 = 1
def some_method1; end
def some_method_1(arg1); end
# good
:some_sym_1
variable_1 = 1
def some_method_1; end
def some_method_1(arg_1); end
```
#### EnforcedStyle: non_integer

```
# bad
:some_sym1
:some_sym_1
variable1 = 1
variable_1 = 1
def some_method1; end
def some_method_1; end
def some_methodone(arg1); end
def some_methodone(arg_1); end
# good
:some_symone
:some_sym_one
variableone = 1
variable_one = 1
def some_methodone; end
def some_method_one; end
def some_methodone(argone); end
def some_methodone(arg_one); end
# In the following examples, we assume `EnforcedStyle: normalcase` (default).
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| CheckMethodNames |
 | Boolean |
| CheckSymbols |
 | Boolean |
| AllowedIdentifiers |
 | Array |
| AllowedPatterns |
 | Array |

# Security

## Security/CompoundHash

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | No | 1.28 | 1.51 |

Checks for implementations of the `hash` method which combine
values using custom logic instead of delegating to `Array#hash`.

Manually combining hashes is error prone and hard to follow, especially
when there are many values. Poor implementations may also introduce
performance or security concerns if they are prone to collisions.
Delegating to `Array#hash` is clearer and safer, although it might be slower
depending on the use case.

## Security/Eval

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.47 | - |

Checks for the use of `Kernel#eval` and `Binding#eval` with
dynamic strings as arguments. Evaluating non-literal strings
can enable code injection attacks and makes it difficult to
reason about what code will actually be executed.

Calls to `eval` with literal strings are not flagged by this cop,
as they do not pose the same injection risk.

## Security/IoMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.22 | - |

Checks for the first argument to `IO.read`, `IO.binread`, `IO.write`, `IO.binwrite`,
`IO.foreach`, and `IO.readlines`.

If argument starts with a pipe character (`'|'`) and the receiver is the `IO` class,
a subprocess is created in the same way as `Kernel#open`, and its output is returned.
`Kernel#open` may allow unintentional command injection, which is the reason these
`IO` methods are a security risk.
Consider using `File.read` to disable the behavior of subprocess invocation.

## Security/JSONLoad

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.43 | 1.22 |

Checks for the use of JSON class methods which have potential security issues.

`JSON.load` and similar methods allow deserialization of arbitrary ruby objects:

```
require 'json/add/string'
result = JSON.load('{ "json_class": "String", "raw": [72, 101, 108, 108, 111] }')
pp result # => "Hello"
```
Never use `JSON.load` for untrusted user input. Prefer `JSON.parse` unless you have
a concrete use-case for `JSON.load`.

| Starting with `json`gem version 2.8.0, triggering this behavior without explicitly
passing the`create_additions`keyword argument emits a deprecation warning, with the
goal of being secure by default in the next major version 3.0.0. |

### Safety

This cop’s autocorrection is unsafe because it’s potentially dangerous.
If using a stream, like `JSON.load(open('file'))`, you will need to call
`#read` manually, like `JSON.parse(open('file').read)`.
Other similar issues may apply.

## Security/MarshalLoad

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.47 | - |

Checks for the use of Marshal class methods which have potential security issues leading to remote code execution when loading from an untrusted source.

## Security/Open

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | No | 0.53 | 1.0 |

Checks for the use of `Kernel#open` and `URI.open` with dynamic
data.

`Kernel#open` and `URI.open` enable not only file access but also process
invocation by prefixing a pipe symbol (e.g., `open("| ls")`).
So, it may lead to a serious security risk by using variable input to
the argument of `Kernel#open` and `URI.open`. It would be better to use
`File.open`, `IO.popen` or `URI.parse#open` explicitly.

| `open`and`URI.open`with literal strings are not flagged by this
cop. |

## Security/YAMLLoad

| Requires Ruby version ⇐ 3.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.47 | - |

Checks for the use of YAML class methods which have potential security issues leading to remote code execution when loading from an untrusted source.

| Ruby 3.1+ (Psych 4) uses `Psych.load`as`Psych.safe_load`by default. |

### Safety

The behavior of the code might change depending on what was
in the YAML payload, since `YAML.safe_load` is more restrictive.

### Examples

```
# bad
YAML.load("--- !ruby/object:Foo {}") # Psych 3 is unsafe by default
# good
YAML.safe_load("--- !ruby/object:Foo {}", [Foo]) # Ruby 2.5 (Psych 3)
YAML.safe_load("--- !ruby/object:Foo {}", permitted_classes: [Foo]) # Ruby 3.0- (Psych 3)
YAML.load("--- !ruby/object:Foo {}", permitted_classes: [Foo]) # Ruby 3.1+ (Psych 4)
YAML.dump(foo)
```

# Style

## Style/AccessModifierDeclarations

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.57 | 1.70 |

Access modifiers should be declared to apply to a group of methods
or inline before each method, depending on configuration.
EnforcedStyle config covers only method definitions.
Applications of visibility methods to symbols can be controlled
using AllowModifiersOnSymbols config.
Also, the visibility of `attr*` methods can be controlled using
AllowModifiersOnAttrs config.

In Ruby 3.0, `attr*` methods now return an array of defined method names
as symbols. So we can write the modifier and `attr*` in inline style.
AllowModifiersOnAttrs config allows `attr*` methods to be written in
inline style without modifying applications that have been maintained
for a long time in group style. Furthermore, developers who are not very
familiar with Ruby may know that the modifier applies to `def`, but they
may not know that it also applies to `attr*` methods. It would be easier
to understand if we could write `attr*` methods in inline style.

### Safety

Autocorrection is not safe, because the visibility of dynamically defined methods can vary depending on the state determined by the group access modifier.

### Examples

#### EnforcedStyle: group (default)

```
# bad
class Foo
 private def bar; end
 private def baz; end
end
# good
class Foo
 private
 def bar; end
 def baz; end
end
```
#### EnforcedStyle: inline

```
# bad
class Foo
 private
 def bar; end
 def baz; end
end
# good
class Foo
 private def bar; end
 private def baz; end
end
```
#### AllowModifiersOnSymbols: true (default)

```
# good
class Foo
 private :bar, :baz
 private *%i[qux quux]
 private *METHOD_NAMES
 private *private_methods
end
```
#### AllowModifiersOnSymbols: false

```
# bad
class Foo
 private :bar, :baz
 private *%i[qux quux]
 private *METHOD_NAMES
 private *private_methods
end
```
#### AllowModifiersOnAttrs: true (default)

```
# good
class Foo
 public attr_reader :bar
 protected attr_writer :baz
 private attr_accessor :qux
 private attr :quux
 def public_method; end
 private
 def private_method; end
end
```
#### AllowModifiersOnAttrs: false

```
# bad
class Foo
 public attr_reader :bar
 protected attr_writer :baz
 private attr_accessor :qux
 private attr :quux
end
```
## Style/AccessorGrouping

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.87 | - |

Checks for grouping of accessors in `class` and `module` bodies.
By default it enforces accessors to be placed in grouped
declarations, reducing boilerplate. It can also be configured
to enforce separating them into individual declarations for
easier diffing and per-attribute documentation.

| If there is a method call before the accessor method it is always allowed as it might be intended like Sorbet. |

| If there is a RBS::Inline annotation comment just after the accessor method it is always allowed. |

### Examples

#### EnforcedStyle: grouped (default)

```
# bad
class Foo
 attr_reader :bar
 attr_reader :bax
 attr_reader :baz
end
# good
class Foo
 attr_reader :bar, :bax, :baz
end
# good
class Foo
 # may be intended comment for bar.
 attr_reader :bar
 sig { returns(String) }
 attr_reader :bax
 may_be_intended_annotation :baz
 attr_reader :baz
end
```
## Style/Alias

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.36 |

Enforces the use of either `#alias` or `#alias_method`
depending on configuration. Consistent use of one or the
other prevents confusion about their different semantics
(e.g., `alias` is resolved at parse time, while `alias_method`
is resolved at runtime).
It also flags uses of `alias :symbol` rather than `alias bareword`.

However, it will always enforce `alias_method` when `alias` is used
in an instance method definition and in a singleton method definition.
If used in a block, always enforce `alias_method`
unless it is an `instance_eval` block.

### Examples

## Style/AmbiguousEndlessMethodDefinition

| Requires Ruby version 3.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.68 | - |

Looks for endless methods inside operations of lower precedence (`and`, `or`, and
modifier forms of `if`, `unless`, `while`, `until`) that are ambiguous due to
lack of parentheses. This may lead to unexpected behavior as the code may appear
to use these keywords as part of the method but in fact they modify
the method definition itself.

In these cases, using a normal method definition is more clear.

## Style/AndOr

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.9 | 1.21 |

Checks for uses of `and` and `or`, and suggests using `&&` and
`||` instead. It can be configured to check only in conditions or in
all contexts.

### Safety

Autocorrection is unsafe because there is a different operator precedence
between logical operators (`&&` and `||`) and semantic operators (`and` and `or`),
and that might change the behavior.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/ArgumentsForwarding

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.1 | 1.58 |

In Ruby 2.7, arguments forwarding has been added.

This cop identifies places where `do_something(*args, &block)`
can be replaced by `do_something(…)`.

In Ruby 3.1, anonymous block forwarding has been added.

This cop identifies places where `do_something(&block)` can be replaced
by `do_something(&)`; if desired, this functionality can be disabled
by setting `UseAnonymousForwarding: false`.

In Ruby 3.2, anonymous args/kwargs forwarding has been added.

This cop also identifies places where `use_args(*args)`/`use_kwargs(**kwargs)` can be
replaced by `use_args(*)`/`use_kwargs(**)`; if desired, this functionality can be
disabled by setting `UseAnonymousForwarding: false`.

And this cop has `RedundantRestArgumentNames`, `RedundantKeywordRestArgumentNames`,
and `RedundantBlockArgumentNames` options. This configuration is a list of redundant names
that are sufficient for anonymizing meaningless naming.

Meaningless names that are commonly used can be anonymized by default:
e.g., `*args`, `**options`, `&block`, and so on.

Names not on this list are likely to be meaningful and are allowed by default.

This cop handles not only method forwarding but also forwarding to `super`.

| Because of a bug in Ruby 3.3.0, when a block is referenced inside of another block, no offense will be registered until Ruby 3.4: |

### Examples

```
# bad
def foo(*args, &block)
 bar(*args, &block)
end
# bad
def foo(*args, **kwargs, &block)
 bar(*args, **kwargs, &block)
end
# good
def foo(...)
 bar(...)
end
```
#### UseAnonymousForwarding: true (default, only relevant for Ruby >= 3.2)

```
# bad
def foo(*args, **kwargs, &block)
 args_only(*args)
 kwargs_only(**kwargs)
 block_only(&block)
end
# good
def foo(*, **, &)
 args_only(*)
 kwargs_only(**)
 block_only(&)
end
```
#### UseAnonymousForwarding: false (only relevant for Ruby >= 3.2)

```
# good
def foo(*args, **kwargs, &block)
 args_only(*args)
 kwargs_only(**kwargs)
 block_only(&block)
end
```
#### AllowOnlyRestArgument: true (default, only relevant for Ruby < 3.2)

```
# good
def foo(*args)
 bar(*args)
end
def foo(**kwargs)
 bar(**kwargs)
end
```
#### AllowOnlyRestArgument: false (only relevant for Ruby < 3.2)

```
# bad
# The following code can replace the arguments with `...`,
# but it will change the behavior. Because `...` forwards block also.
def foo(*args)
 bar(*args)
end
def foo(**kwargs)
 bar(**kwargs)
end
```
#### RedundantRestArgumentNames: ['args', 'arguments'] (default)

```
# bad
def foo(*args)
 bar(*args)
end
# good
def foo(*)
 bar(*)
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowOnlyRestArgument |
 | Boolean |
| UseAnonymousForwarding |
 | Boolean |
| RedundantRestArgumentNames |
 | Array |
| RedundantKeywordRestArgumentNames |
 | Array |
| RedundantBlockArgumentNames |
 | Array |

## Style/ArrayCoercion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 0.88 | - |

Enforces the use of `Array()` instead of explicit `Array` check or `[*var]`.

The cop is disabled by default due to safety concerns.

### Safety

This cop is unsafe because a false positive may occur if
the argument of `Array()` is (or could be) nil or depending
on how the argument is handled by `Array()` (which can be
different than just wrapping the argument in an array).

For example:

```
[nil] #=> [nil]
Array(nil) #=> []
[{a: 'b'}] #= [{a: 'b'}]
Array({a: 'b'}) #=> [[:a, 'b']]
[Time.now] #=> [#<Time ...>]
Array(Time.now) #=> [14, 16, 14, 16, 9, 2021, 4, 259, true, "EDT"]
```
### Examples

```
# bad
paths = [paths] unless paths.is_a?(Array)
paths.each { |path| do_something(path) }
# bad (always creates a new Array instance)
[*paths].each { |path| do_something(path) }
# good (and a bit more readable)
Array(paths).each { |path| do_something(path) }
```
## Style/ArrayFirstLast

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 1.58 | - |

Identifies usages of `arr[0]` and `arr[-1]` and suggests to change
them to use `arr.first` and `arr.last` instead.

The cop is disabled by default due to safety concerns.

## Style/ArrayIntersect

| Requires Ruby version 3.1 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.40 | - |

In Ruby 3.1, `Array#intersect?` has been added.

This cop identifies places where:

-
`(array1 & array2).any?`
-
`(array1.intersection(array2)).any?`
-
`array1.any? { |elem| array2.member?(elem) }`
-
`(array1 & array2).count > 0`
-
`(array1 & array2).size > 0`

can be replaced with `array1.intersect?(array2)`.

`array1.intersect?(array2)` is faster and more readable.

In cases like the following, compatibility is not ensured, so it will not be detected when using block argument.

```
([1] & [1,2]).any? { |x| false } # => false
[1].intersect?([1,2]) { |x| false } # => true
```
| Although `Array#intersection`can take zero or multiple arguments,
only cases where exactly one argument is provided can be replaced with`Array#intersect?`and are handled by this cop. |

| In the block form, `include?`is only detected when its receiver is
an array literal, because`include?`is defined with different semantics
on many non-array classes (e.g.`String#include?`checks for substrings).`member?`does not have this restriction. |

### Safety

This cop cannot guarantee that `array1` and `array2` are
actually arrays while method `intersect?` is for arrays only.

### Examples

```
# bad
(array1 & array2).any?
(array1 & array2).empty?
(array1 & array2).none?
# bad
array1.intersection(array2).any?
array1.intersection(array2).empty?
array1.intersection(array2).none?
# bad
array1.any? { |elem| array2.member?(elem) }
array1.none? { |elem| [1, 2].include?(elem) }
# good
array1.intersect?(array2)
!array1.intersect?(array2)
# bad
(array1 & array2).count > 0
(array1 & array2).count.positive?
(array1 & array2).count != 0
(array1 & array2).count == 0
(array1 & array2).count.zero?
# good
array1.intersect?(array2)
!array1.intersect?(array2)
```
## Style/ArrayJoin

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.20 | 0.31 |

Checks for uses of ` as a substitute for ``Array#join`.
Using `join` is clearer about intent and more readable than
overloading the ` operator for string conversion.`

Not all cases can be reliably checked, due to Ruby’s dynamic types, so we consider only cases when the first argument is an array literal or the second is a string literal.

## Style/AsciiComments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.9 | 1.21 |

Checks for non-ascii (non-English) characters
in comments. Non-ascii characters can cause issues with
portability and encoding across different environments
and editors. You could set an array of allowed non-ascii
chars in `AllowedChars` attribute (copyright notice "©"
by default).

## Style/Attr

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.12 |

Checks for uses of `Module#attr`. The `attr` method has confusing
behavior: with a single argument it creates a reader (like `attr_reader`),
but with a second boolean argument it creates an accessor (deprecated in
Ruby 1.9). Use `attr_reader` or `attr_accessor` to make intent explicit.

### Examples

```
# bad
attr :something, true
attr :one, :two, :three # behaves as attr_reader
# good
attr_accessor :something
attr_reader :one, :two, :three
```
## Style/AutoResourceCleanup

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.30 | - |

Checks for cases when you could use a block accepting version of a method that does automatic resource cleanup.

## Style/BarePercentLiterals

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.25 | - |

Checks if usage of `%()` or `%Q()` matches configuration.
Consistent use of one style makes the codebase easier
to read.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/BeginBlock

## Style/BisectedAttrAccessor

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.87 | - |

Checks for places where `attr_reader` and `attr_writer`
for the same method can be combined into single `attr_accessor`.

## Style/BitwisePredicate

| Requires Ruby version 2.5 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.68 | - |

Prefer bitwise predicate methods over direct comparison operations.

### Safety

This cop is unsafe, as it can produce false positives if the receiver
is not an `Integer` object.

## Style/BlockComments

## Style/BlockDelimiters

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.30 | 0.35 |

Checks for uses of braces or do/end around single line or multi-line blocks.

Methods that can be either procedural or functional and cannot be
categorised from their usage alone is ignored.
`lambda`, `proc`, and `it` are their defaults.
Additional methods can be added to the `AllowedMethods`.

### Examples

#### EnforcedStyle: line_count_based (default)

```
# bad - single line block
items.each do |item| item / 5 end
# good - single line block
items.each { |item| item / 5 }
# bad - multi-line block
things.map { |thing|
 something = thing.some_method
 process(something)
}
# good - multi-line block
things.map do |thing|
 something = thing.some_method
 process(something)
end
```
#### EnforcedStyle: semantic

```
# Prefer `do...end` over `{...}` for procedural blocks.
# return value is used/assigned
# bad
foo = map do |x|
 x
end
puts (map do |x|
 x
end)
# return value is not used out of scope
# good
map do |x|
 x
end
# Prefer `{...}` over `do...end` for functional blocks.
# return value is not used out of scope
# bad
each { |x|
 x
}
# return value is used/assigned
# good
foo = map { |x|
 x
}
map { |x|
 x
}.inspect
# The AllowBracesOnProceduralOneLiners option is allowed unless the
# EnforcedStyle is set to `semantic`. If so:
# If the AllowBracesOnProceduralOneLiners option is unspecified, or
# set to `false` or any other falsey value, then semantic purity is
# maintained, so one-line procedural blocks must use do-end, not
# braces.
# bad
collection.each { |element| puts element }
# good
collection.each do |element| puts element end
# If the AllowBracesOnProceduralOneLiners option is set to `true`, or
# any other truthy value, then one-line procedural blocks may use
# either style. (There is no setting for requiring braces on them.)
# good
collection.each { |element| puts element }
# also good
collection.each do |element| puts element end
```
#### EnforcedStyle: braces_for_chaining

```
# bad
words.each do |word|
 word.flip.flop
end.join("-")
# good
words.each { |word|
 word.flip.flop
}.join("-")
```
#### EnforcedStyle: always_braces

```
# bad
words.each do |word|
 word.flip.flop
end
# good
words.each { |word|
 word.flip.flop
}
```
#### BracesRequiredMethods: ['sig']

```
# Methods listed in the BracesRequiredMethods list, such as 'sig'
# in this example, will require `{...}` braces. This option takes
# precedence over all other configurations except AllowedMethods.
# bad
sig do
 params(
 foo: string,
 ).void
end
def bar(foo)
 puts foo
end
# good
sig {
 params(
 foo: string,
 ).void
}
def bar(foo)
 puts foo
end
```
#### AllowedMethods: ['lambda', 'proc', 'it' ] (default)

```
# good
foo = lambda do |x|
 puts "Hello, #{x}"
end
foo = lambda do |x|
 x * 100
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| ProceduralMethods |
 | Array |
| FunctionalMethods |
 | Array |
| AllowedMethods |
 | Array |
| AllowedPatterns |
 | Array |
| AllowBracesOnProceduralOneLiners |
 | Boolean |
| BracesRequiredMethods |
 | Array |

## Style/CaseEquality

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.89 |

Checks for uses of the case equality operator (`===`).
The `===` operator has different behavior depending on the
receiver and its use outside of `case`/`when` is confusing.
Prefer more explicit alternatives like `is_a?`, `include?`,
or `match?`.

If `AllowOnConstant` option is enabled, the cop will ignore violations when the receiver of
the case equality operator is a constant.

If `AllowOnSelfClass` option is enabled, the cop will ignore violations when the receiver of
the case equality operator is `self.class`. Note intermediate variables are not accepted.

| Regexp case equality ( `/regexp/ === var`) is allowed because changing it to`/regexp/.match?(var)`needs to take into account`Regexp.last_match?`,`$~`,`$1`, etc.
This potentially incompatible transformation is handled by`Performance/RegexpMatch`cop. |

## Style/CaseLikeIf

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.88 | 1.48 |

Identifies places where `if-elsif` constructions
can be replaced with `case-when`.

### Safety

This cop is unsafe. `case` statements use `===` for equality,
so if the original conditional used a different equality operator, the
behavior may be different.

### Examples

#### MinBranchesCount: 3 (default)

```
# bad
if status == :active
 perform_action
elsif status == :inactive || status == :hibernating
 check_timeout
elsif status == :invalid
 report_invalid
else
 final_action
end
# good
case status
when :active
 perform_action
when :inactive, :hibernating
 check_timeout
when :invalid
 report_invalid
else
 final_action
end
```
## Style/CharacterLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | - |

Checks for uses of the character literal ?x. Starting with Ruby 1.9 character literals are essentially one-character strings, so this syntax is mostly redundant at this point.

A `?` character literal can be used to express meta and control characters.
That’s a good use case of a `?` literal so it doesn’t count as an offense.

## Style/ClassAndModuleChildren

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.19 | 1.89 |

Checks that namespaced classes and modules are defined with a consistent style.

With `nested` style, classes and modules should be defined separately (one constant
on each line, without `::`). With `compact` style, classes and modules should be
defined with fully qualified names (using `::` for namespaces).

| The style chosen will affect `Module.nesting`for the class or module. Using`nested`style will result in each level being added, whereas`compact`style will
only include the fully qualified class or module name. |

By default, `EnforcedStyle` applies to both classes and modules. If desired, separate
styles can be defined for classes and modules by using `EnforcedStyleForClasses` and
`EnforcedStyleForModules` respectively. If not set, or set to nil, the `EnforcedStyle`
value will be used.

The compact style is only forced for classes/modules with one child.

### Safety

Autocorrection is unsafe.

Moving from `compact` to `nested` children requires knowledge of whether the
outer parent is a module or a class. Moving from `nested` to `compact` requires
verification that the outer parent is defined elsewhere. By default RuboCop does
not have the knowledge to perform either operation safely and thus requires
manual oversight.

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is installed,
the project-wide index is consulted to resolve whether the outer parent is a
class or a module, and compacting is skipped when the outer parent is not
defined elsewhere.

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| EnforcedStyleForClasses |
 |
 |
| EnforcedStyleForModules |
 |
 |

## Style/ClassCheck

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.24 | - |

Enforces consistent use of `Object#is_a?` or `Object#kind_of?`.

### Examples

## Style/ClassEqualityComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.93 | 1.57 |

Enforces the use of `Object#instance_of?` instead of class comparison
for equality.
`==`, `equal?`, and `eql?` custom method definitions are allowed by default.
These are customizable with `AllowedMethods` option.

### Safety

This cop’s autocorrection is unsafe because there is no guarantee that
the constant `Foo` exists when autocorrecting `var.class.name == 'Foo'` to
`var.instance_of?(Foo)`.

### Examples

```
# bad
var.class == Date
var.class.equal?(Date)
var.class.eql?(Date)
var.class.name == 'Date'
# good
var.instance_of?(Date)
```
#### AllowedMethods: ['==', 'equal?', 'eql?'] (default)

```
# good
def ==(other)
 self.class == other.class && name == other.name
end
def equal?(other)
 self.class.equal?(other.class) && name.equal?(other.name)
end
def eql?(other)
 self.class.eql?(other.class) && name.eql?(other.name)
end
```
## Style/ClassMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.20 |

Checks for uses of the class/module name instead of self, when defining class/module methods.

## Style/ClassMethodsDefinitions

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.89 | - |

Enforces using `def self.method_name` or `class << self` to define class methods.

### Examples

#### EnforcedStyle: def_self (default)

```
# bad
class SomeClass
 class << self
 attr_accessor :class_accessor
 def class_method
 # ...
 end
 end
end
# good
class SomeClass
 def self.class_method
 # ...
 end
 class << self
 attr_accessor :class_accessor
 end
end
# good - contains private method
class SomeClass
 class << self
 attr_accessor :class_accessor
 private
 def private_class_method
 # ...
 end
 end
end
```
## Style/ClassVars

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.13 | - |

Checks for uses of class variables. Offenses are signaled only on assignment to class variables to reduce the number of offenses that would be reported.

You have to be careful when setting a value for a class variable; if a class has been inherited, changing the value of a class variable also affects the inheriting classes. This means that it’s almost always better to use a class instance variable instead.

### Examples

```
# bad
class A
 @@test = 10
end
class A
 def self.test(name, value)
 class_variable_set("@@#{name}", value)
 end
end
class A; end
A.class_variable_set(:@@test, 10)
# good
class A
 @test = 10
end
class A
 def test
 @@test # you can access class variable without offense
 end
end
class A
 def self.test(name)
 class_variable_get("@@#{name}") # you can access without offense
 end
end
```
## Style/CollectionCompact

| Requires Ruby version 2.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.2 | 1.3 |

Checks for places where custom logic on rejection nils from arrays
and hashes can be replaced with `{Array,Hash}#{compact,compact!}`.

### Safety

It is unsafe by default because false positives may occur in the
`nil` check of block arguments to the receiver object. Additionally,
we can’t know the type of the receiver object for sure, which may
result in false positives as well.

For example, `[[1, 2], [3, nil]].reject { |first, second| second.nil? }`
and `[[1, 2], [3, nil]].compact` are not compatible. This will work fine
when the receiver is a hash object.

### Examples

```
# bad
array.reject(&:nil?)
array.reject { |e| e.nil? }
array.select { |e| !e.nil? }
array.filter { |e| !e.nil? }
array.grep_v(nil)
array.grep_v(NilClass)
# good
array.compact
# bad
hash.reject!(&:nil?)
hash.reject! { |k, v| v.nil? }
hash.select! { |k, v| !v.nil? }
hash.filter! { |k, v| !v.nil? }
# good
hash.compact!
```
## Style/CollectionMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 0.9 | 1.7 |

Enforces the use of consistent method names from the Enumerable module.

You can customize the mapping from undesired method to desired method.

e.g. to use `detect` over `find`:

```
Style/CollectionMethods:
 PreferredMethods:
 find: detect
```
### Safety

This cop is unsafe because it finds methods by name, without actually being able to determine if the receiver is an Enumerable or not, so this cop may register false positives.

### Examples

```
# These examples are based on the default mapping for `PreferredMethods`.
# bad
items.collect
items.collect!
items.collect_concat
items.inject
items.detect
items.find_all
items.member?
# good
items.map
items.map!
items.flat_map
items.reduce
items.find
items.select
items.include?
```
## Style/CollectionQuerying

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.77 | - |

Prefer `Enumerable` predicate methods over expressions with `count`.

The cop checks calls to `count` without arguments, or with a
block. It doesn’t register offenses for `count` with a positional
argument because its behavior differs from predicate methods (`count`
matches the argument using `==`, while `any?`, `none?` and `one?` use
`===`).

| This cop doesn’t check `length`and`size`methods because they
would yield false positives. For example,`String`implements`length`and`size`, but it doesn’t include`Enumerable`. |

### Safety

The cop is unsafe because receiver might not include `Enumerable`, or
it has nonstandard implementation of `count` or any replacement
methods.

It’s also unsafe because for collections with falsey values, expressions
with `count` without a block return a different result than methods `any?`,
`none?` and `one?`:

```
[nil, false].count.positive?
[nil].count == 1
# => true
[nil, false].any?
[nil].one?
# => false
[nil].count == 0
# => false
[nil].none?
# => true
```
Autocorrection is unsafe when replacement methods don’t iterate over every element in collection and the given block runs side effects:

```
x.count(&:method_with_side_effects).positive?
# calls `method_with_side_effects` on every element
x.any?(&:method_with_side_effects)
# calls `method_with_side_effects` until first element returns a truthy value
```
### Examples

```
# bad
x.count.positive?
x.count > 0
x.count != 0
x.count(&:foo?).positive?
x.count { |item| item.foo? }.positive?
# good
x.any?
x.any?(&:foo?)
x.any? { |item| item.foo? }
# bad
x.count.zero?
x.count == 0
# good
x.none?
# bad
x.count == 1
x.one?
```
## Style/ColonMethodCall

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | - |

Checks for methods invoked via the `::` operator instead
of the `.` operator (like `FileUtils::rmdir` instead of
`FileUtils.rmdir`). The `::` operator is conventionally used to
reference constants, so using it for method calls can be misleading.

### Examples

```
# bad
Timeout::timeout(500) { do_something }
FileUtils::rmdir(dir)
Marshal::dump(obj)
# good
Timeout.timeout(500) { do_something }
FileUtils.rmdir(dir)
Marshal.dump(obj)
```
## Style/ColonMethodDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for class methods that are defined using the `::`
operator instead of the `.` operator.

## Style/CombinableDefined

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.68 | - |

Checks for multiple `defined?` calls joined by `&&` that can be combined
into a single `defined?`.

When checking that a nested constant or chained method is defined, it is not necessary to check each ancestor or component of the chain.

## Style/CombinableLoops

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.90 | - |

Checks for places where multiple consecutive loops over the same data can be combined into a single loop. It is very likely that combining them will make the code more efficient and more concise.

| Autocorrection is not applied when the block variable names differ in separate loops, as it is impossible to determine which variable name should be prioritized. |

### Safety

The cop is unsafe, because the first loop might modify state that the second loop depends on; these two aren’t combinable.

### Examples

```
# bad
def method
 items.each do |item|
 do_something(item)
 end
 items.each do |item|
 do_something_else(item)
 end
end
# good
def method
 items.each do |item|
 do_something(item)
 do_something_else(item)
 end
end
# bad
def method
 for item in items do
 do_something(item)
 end
 for item in items do
 do_something_else(item)
 end
end
# good
def method
 for item in items do
 do_something(item)
 do_something_else(item)
 end
end
# good
def method
 each_slice(2) { |slice| do_something(slice) }
 each_slice(3) { |slice| do_something(slice) }
end
```
## Style/CommandLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.30 | - |

Enforces using `` or %x around command literals.

### Examples

#### EnforcedStyle: backticks (default)

```
# bad
folders = %x(find . -type d).split
# bad
%x(
 ln -s foo.example.yml foo.example
 ln -s bar.example.yml bar.example
)
# good
folders = `find . -type d`.split
# good
`
 ln -s foo.example.yml foo.example
 ln -s bar.example.yml bar.example
`
```
#### EnforcedStyle: mixed

```
# bad
folders = %x(find . -type d).split
# bad
`
 ln -s foo.example.yml foo.example
 ln -s bar.example.yml bar.example
`
# good
folders = `find . -type d`.split
# good
%x(
 ln -s foo.example.yml foo.example
 ln -s bar.example.yml bar.example
)
```
#### EnforcedStyle: percent_x

```
# bad
folders = `find . -type d`.split
# bad
`
 ln -s foo.example.yml foo.example
 ln -s bar.example.yml bar.example
`
# good
folders = %x(find . -type d).split
# good
%x(
 ln -s foo.example.yml foo.example
 ln -s bar.example.yml bar.example
)
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| AllowInnerBackticks |
 | Boolean |

## Style/CommentAnnotation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.10 | 1.20 |

Checks that comment annotation keywords are written according to guidelines.

Annotation keywords can be specified by overriding the cop’s `Keywords`
configuration. Keywords are allowed to be single words or phrases.

| With a multiline comment block (where each line is only a
comment), only the first line will be able to register an offense, even
if an annotation keyword starts another line. This is done to prevent
incorrect registering of keywords (eg. `review`) inside a paragraph as an
annotation. |

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| Keywords |
 | Array |
| RequireColon |
 | Boolean |

## Style/CommentedKeyword

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.51 | 1.19 |

Checks for comments put on the same line as some keywords.
These keywords are: `class`, `module`, `def`, `begin`, `end`.

Note that some comments
(`:nodoc:`, `:yields:`, `rubocop:disable` and `rubocop:todo`),
RBS::Inline annotation, and Steep annotation (`steep:ignore`) are allowed.

Autocorrection removes comments from `end` keyword and keeps comments
for `class`, `module`, `def` and `begin` above the keyword.

## Style/ComparableBetween

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.74 | 1.75 |

Checks for logical comparison which can be replaced with `Comparable#between?`.

| `Comparable#between?`is on average slightly slower than logical comparison,
although the difference generally isn’t observable. If you require maximum
performance, consider using logical comparison. |

## Style/ComparableClamp

| Requires Ruby version 2.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.44 | - |

Enforces the use of `Comparable#clamp` instead of comparison by minimum and maximum.

This cop supports autocorrection for `if/elsif/else` bad style only.
Because `ArgumentError` occurs if the minimum and maximum of `clamp` arguments are reversed.
When these are variables, it is not possible to determine which is the minimum and maximum:

```
[1, [2, 3].max].min # => 1
1.clamp(3, 1) # => min argument must be smaller than max argument (ArgumentError)
```
## Style/ConcatArrayLiterals

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.41 | - |

Enforces the use of `Array#push(item)` instead of `Array#concat([item])`
to avoid redundant array literals.

## Style/ConditionalAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.36 | 0.47 |

Checks for `if` and `case` statements where each branch is used for
both the assignment and comparison of the same variable
when using the return of the condition can be used instead.

### Examples

#### EnforcedStyle: assign_to_condition (default)

```
# bad
if foo
 bar = 1
else
 bar = 2
end
case foo
when 'a'
 bar += 1
else
 bar += 2
end
if foo
 some_method
 bar = 1
else
 some_other_method
 bar = 2
end
# good
bar = if foo
 1
 else
 2
 end
bar += case foo
 when 'a'
 1
 else
 2
 end
bar << if foo
 some_method
 1
 else
 some_other_method
 2
 end
```
#### EnforcedStyle: assign_inside_condition

```
# bad
bar = if foo
 1
 else
 2
 end
bar += case foo
 when 'a'
 1
 else
 2
 end
bar << if foo
 some_method
 1
 else
 some_other_method
 2
 end
# good
if foo
 bar = 1
else
 bar = 2
end
case foo
when 'a'
 bar += 1
else
 bar += 2
end
if foo
 some_method
 bar = 1
else
 some_other_method
 bar = 2
end
```
## Style/ConstantVisibility

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.66 | 1.10 |

Checks that constants defined in classes and modules have an explicit visibility declaration. By default, Ruby makes all class- and module constants public, which litters the public API of the class or module. Explicitly declaring a visibility makes intent more clear, and prevents outside actors from touching private state.

## Style/Copyright

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.30 | - |

Checks that a copyright notice was given in each source file.

The default regexp for an acceptable copyright notice can be found in config/default.yml. The default can be changed as follows:

```
Style/Copyright:
 Notice: '^Copyright (\(c\) )?2\d{3} Acme Inc'
```
This regex string is treated as an unanchored regex. For each file that RuboCop scans, a comment that matches this regex must be found or an offense is reported.

## Style/DataInheritance

| Requires Ruby version 3.2 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.49 | 1.51 |

Checks for inheritance from `Data.define` to avoid creating the anonymous parent class.
Inheriting from `Data.define` adds a superfluous level in inheritance tree.

### Safety

Autocorrection is unsafe because it will change the inheritance
tree (e.g. return value of `Module#ancestors`) of the constant.

## Style/DateTime

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always (Unsafe) | 0.51 | 0.92 |

Checks for consistent usage of the `Time` class over the
`DateTime` class. This cop is disabled by default since these classes,
although highly overlapping, have particularities that make them not
replaceable in certain situations when dealing with multiple timezones
and/or DST.

### Safety

Autocorrection is not safe, because `DateTime` and `Time` do not have
exactly the same behavior, although in most cases the autocorrection
will be fine.

### Examples

```
# bad - uses `DateTime` for current time
DateTime.now
# good - uses `Time` for current time
Time.now
# bad - uses `DateTime` for modern date
DateTime.iso8601('2016-06-29')
# good - uses `Time` for modern date
Time.iso8601('2016-06-29')
# good - uses `DateTime` with start argument for historical date
DateTime.iso8601('1751-04-23', Date::ENGLAND)
```
## Style/DefWithParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.12 |

Checks for parentheses in the definition of a method, that does not take any arguments. Both instance and class/singleton methods are checked.

### Examples

```
# bad
def foo()
 do_something
end
# good
def foo
 do_something
end
# bad
def foo() = do_something
# good
def foo = do_something
# good - without parentheses it's a syntax error
def foo() do_something end
def foo()=do_something
# bad
def Baz.foo()
 do_something
end
# good
def Baz.foo
 do_something
end
```
## Style/DigChain

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.69 | - |

Checks for chained `dig` calls that can be collapsed into a single `dig`.

## Style/Dir

| Requires Ruby version 2.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.50 | - |

Checks for places where the `#_\_dir\_\_` method can replace more
complex constructs to retrieve a canonicalized absolute path to the
current file.

## Style/DirEmpty

| Requires Ruby version 2.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.48 | - |

Prefer to use `Dir.empty?('path/to/dir')` when checking if a directory is empty.

## Style/DisableCopsWithinSourceCodeDirective

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.82 | 1.89 |

Detects comments to enable/disable RuboCop. This is useful if you want to make sure that every RuboCop error gets fixed and not quickly disabled with a comment.

Specific cops can be allowed with the `AllowedCops` configuration. Note that
if this configuration is set, `rubocop:disable all` is still disallowed.

Alternatively, specific cops can be disallowed with the `DisallowedCops`
configuration. When `DisallowedCops` is set, only directives for the listed
cops (and `all`) will be flagged. This is useful when you want
to protect a small set of critical cops from being disabled rather than
allowlisting all other cops. `AllowedCops` and `DisallowedCops` should not
both be set at the same time; if `DisallowedCops` is set, it takes precedence.

This cop cannot be disabled via directive comments when it is explicitly
enabled with `Enabled: true`. This prevents users from bypassing the cop
with `# rubocop:disable Style/DisableCopsWithinSourceCodeDirective`.

## Style/DocumentDynamicEvalDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.1 | 1.3 |

When using `class_eval` (or other `eval`) with string interpolation,
add a comment block showing its appearance if interpolated (a practice used in Rails code).

### Examples

```
# from activesupport/lib/active_support/core_ext/string/output_safety.rb
# bad
UNSAFE_STRING_METHODS.each do |unsafe_method|
 if 'String'.respond_to?(unsafe_method)
 class_eval <<-EOT, __FILE__, __LINE__ + 1
 def #{unsafe_method}(*params, &block)
 to_str.#{unsafe_method}(*params, &block)
 end
 def #{unsafe_method}!(*params)
 @dirty = true
 super
 end
 EOT
 end
end
# good, inline comments in heredoc
UNSAFE_STRING_METHODS.each do |unsafe_method|
 if 'String'.respond_to?(unsafe_method)
 class_eval <<-EOT, __FILE__, __LINE__ + 1
 def #{unsafe_method}(*params, &block) # def capitalize(*params, &block)
 to_str.#{unsafe_method}(*params, &block) # to_str.capitalize(*params, &block)
 end # end
 def #{unsafe_method}!(*params) # def capitalize!(*params)
 @dirty = true # @dirty = true
 super # super
 end # end
 EOT
 end
end
# good, block comments in heredoc
class_eval <<-EOT, __FILE__, __LINE__ + 1
 # def capitalize!(*params)
 # @dirty = true
 # super
 # end
 def #{unsafe_method}!(*params)
 @dirty = true
 super
 end
EOT
# good, block comments before heredoc
class_eval(
 # def capitalize!(*params)
 # @dirty = true
 # super
 # end
 <<-EOT, __FILE__, __LINE__ + 1
 def #{unsafe_method}!(*params)
 @dirty = true
 super
 end
 EOT
)
# bad - interpolated string without comment
class_eval("def #{unsafe_method}!(*params); end")
# good - with inline comment or replace it with block comment using heredoc
class_eval("def #{unsafe_method}!(*params); end # def capitalize!(*params); end")
```
## Style/Documentation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.9 | 1.89 |

Checks for missing top-level documentation of classes and modules. Classes with no body are exempt from the check and so are namespace modules - modules that have nothing in their bodies except classes, other modules, constant definitions or constant visibility declarations.

The documentation requirement is annulled if the class or module has
a `:nodoc:` comment next to it. Likewise, `:nodoc: all` does the
same for all its children.

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, a reopened class or module is not reported when any of its
other definition sites (in the same or another file) carries a
documentation comment.

### Examples

```
# bad
class Person
 # ...
end
module Math
end
# good
# Description/Explanation of Person class
class Person
 # ...
end
# allowed
# Class without body
class Person
end
# Namespace - A namespace can be a class or a module
# Containing a class
module Namespace
 # Description/Explanation of Person class
 class Person
 # ...
 end
end
# Containing constant visibility declaration
module Namespace
 class Private
 end
 private_constant :Private
end
# Containing constant definition
module Namespace
 Public = Class.new
end
# Macro calls
module Namespace
 extend Foo
end
```
## Style/DocumentationMethod

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.43 | - |

Checks for missing documentation comment for public methods. It can optionally be configured to also require documentation for non-public methods.

| This cop allows `initialize`method because`initialize`is
a special method called from`new`. In some programming languages
they are called constructor to distinguish it from method. |

### Examples

```
# bad
class Foo
 def bar
 puts baz
 end
end
module Foo
 def bar
 puts baz
 end
end
def foo.bar
 puts baz
end
# good
class Foo
 # Documentation
 def bar
 puts baz
 end
end
module Foo
 # Documentation
 def bar
 puts baz
 end
end
# Documentation
def foo.bar
 puts baz
end
```
#### RequireForNonPublicMethods: false (default)

```
# good
class Foo
 protected
 def do_something
 end
end
class Foo
 private
 def do_something
 end
end
```
## Style/DoubleCopDisableDirective

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.73 | - |

Detects double disable comments on one line. This is mostly to catch automatically generated comments that need to be regenerated.

## Style/DoubleNegation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.19 | 1.2 |

Checks for uses of double negation (`!!`) to convert something to a boolean value.

When using `EnforcedStyle: allowed_in_returns`, allow double negation in contexts
that use boolean as a return value. When using `EnforcedStyle: forbidden`, double negation
should be forbidden always.

| when `something`is a boolean value`!!something`and`!something.nil?`are not the same thing.
As you’re unlikely to write code that can accept values of any type
this is rarely a problem in practice. |

### Safety

Autocorrection is unsafe when the value is `false`, because the result
of the expression will change.

```
!!false #=> false
!false.nil? #=> true
```
### Examples

```
# bad
!!something
# good
!something.nil?
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/EachForSimpleLoop

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.41 | - |

Checks for loops which iterate a constant number of times,
using a `Range` literal and `#each`. This can be done more readably using
`Integer#times`.

This check only applies if the block takes no parameters.

## Style/EachWithObject

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.22 | 0.42 |

Looks for inject / reduce calls where the passed in object is returned at the end and so could be replaced by each_with_object without the need to return the object at the end.

However, we can’t replace with each_with_object if the accumulator parameter is assigned to within the block.

## Style/EmptyBlockParameter

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for pipes for empty block parameters. Pipes for empty block parameters do not cause syntax errors, but they are redundant.

## Style/EmptyCaseCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.40 | - |

Checks for case statements with an empty condition.

### Examples

```
# bad:
case
when x == 0
 puts 'x is 0'
when y == 0
 puts 'y is 0'
else
 puts 'neither is 0'
end
# good:
if x == 0
 puts 'x is 0'
elsif y == 0
 puts 'y is 0'
else
 puts 'neither is 0'
end
# good: (the case condition node is not empty)
case n
when 0
 puts 'zero'
when 1
 puts 'one'
else
 puts 'more'
end
```
## Style/EmptyClassDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.84 | 1.86 |

Enforces consistent style for empty class definitions.

This cop can enforce either a standard class definition or `Class.new`
for classes with no body.

The supported styles are:

-
class_keyword (default) - prefer standard class definition over `Class.new`
-
class_new - prefer `Class.new`over class definition

One difference between the two styles is that the `Class.new` form does not make
the subclass name available to the base class’s `inherited` callback.
For this reason, `EnforcedStyle: class_keyword` is set as the default style.
Class definitions without a superclass, which are not involved in inheritance,
are not detected. This ensures safe detection regardless of the applied style.
This avoids overlapping responsibilities with the `Lint/EmptyClass` cop.

Use `AllowedParentClasses` to permit both styles for specific parent classes.
For example, adding `StandardError` allows both `Error = Class.new(StandardError)`
and `class Error < StandardError; end` regardless of the enforced style.

### Examples

#### EnforcedStyle: class_keyword (default)

```
# bad
FooError = Class.new(StandardError)
# okish
class FooError < StandardError; end
# good
class FooError < StandardError
end
```
## Style/EmptyElse

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.28 | 1.61 |

Checks for empty else-clauses, possibly including comments and/or an
explicit `nil` depending on the EnforcedStyle.

### Examples

#### EnforcedStyle: both (default)

```
# warn on empty else and else with nil in it
# bad
if condition
 statement
else
 nil
end
# bad
if condition
 statement
else
end
# good
if condition
 statement
else
 statement
end
# good
if condition
 statement
end
```
#### EnforcedStyle: empty

```
# warn only on empty else
# bad
if condition
 statement
else
end
# good
if condition
 statement
else
 nil
end
# good
if condition
 statement
else
 statement
end
# good
if condition
 statement
end
```
#### EnforcedStyle: nil

```
# warn on else with nil in it
# bad
if condition
 statement
else
 nil
end
# good
if condition
 statement
else
end
# good
if condition
 statement
else
 statement
end
# good
if condition
 statement
end
```
## Style/EmptyHeredoc

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Command-line only | 1.32 | 1.61 |

Checks for using empty heredoc to reduce redundancy.

## Style/EmptyLambdaParameter

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for parentheses for empty lambda parameters. Parentheses for empty lambda parameters do not cause syntax errors, but they are redundant.

## Style/EmptyLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.12 |

Checks for the use of a method, the result of which would be a literal, like an empty array, hash, or string.

| When frozen string literals are enabled, `String.new`isn’t corrected to an empty string since the former is
mutable and the latter would be frozen. |

### Examples

```
# bad
a = Array.new
a = Array[]
h = Hash.new
h = Hash[]
s = String.new
# good
a = []
h = {}
s = ''
```
## Style/EmptyMethod

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Command-line only | 0.46 | 1.61 |

Checks for the formatting of empty method definitions.
By default it enforces empty method definitions to go on a single
line (compact style), but it can be configured to enforce the `end`
to go on its own line (expanded style).

| A method definition is not considered empty if it contains comments. |

| Autocorrection will not be applied for the `compact`style
if the resulting code is longer than the`Max`configuration for`Layout/LineLength`, but an offense will still be registered. |

### Examples

## Style/EmptyStringInsideInterpolation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.76 | - |

Checks for empty strings being assigned inside string interpolation.

Empty strings are a meaningless outcome inside of string interpolation, so we remove them. Alternatively, when configured to do so, we prioritise using empty strings.

While this cop would also apply to variables that are only going to be used as strings, RuboCop can’t detect that, so we only check inside of string interpolation.

### Examples

## Style/Encoding

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.50 |

Checks that source files have no utf-8 encoding comments. Since Ruby 2.0, UTF-8 is the default source encoding, so these comments are no longer necessary and just add noise.

### Examples

```
# bad
# encoding: UTF-8
# coding: UTF-8
# -*- coding: UTF-8 -*-
# good
# No encoding comment needed
```
## Style/EndBlock

## Style/EndlessMethod

| Requires Ruby version 3.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.8 | - |

Checks for endless methods.

It can enforce endless method definitions whenever possible or with single line methods. It can also disallow multiline endless method definitions or all endless definitions.

`require_single_line` style enforces endless method definitions for single line methods.
`require_always` style enforces endless method definitions for single statement methods.

Other method definition types are not considered by this cop.

The supported styles are:

-
allow_single_line (default) - only single line endless method definitions are allowed.
-
allow_always - all endless method definitions are allowed.
-
disallow - all endless method definitions are disallowed.
-
require_single_line - endless method definitions are required for single line methods.
-
require_always - all endless method definitions are required.

| Incorrect endless method definitions will always be corrected to a multi-line definition. |

### Examples

#### EnforcedStyle: allow_single_line (default)

```
# bad, multi-line endless method
def my_method = x.foo
 .bar
 .baz
# good
def my_method
 x
end
# good
def my_method = x
# good
def my_method
 x.foo
 .bar
 .baz
end
```
#### EnforcedStyle: allow_always

```
# good
def my_method
 x
end
# good
def my_method = x
# good
def my_method = x.foo
 .bar
 .baz
# good
def my_method
 x.foo
 .bar
 .baz
end
```
#### EnforcedStyle: disallow

```
# bad
def my_method = x
# bad
def my_method = x.foo
 .bar
 .baz
# good
def my_method
 x
end
# good
def my_method
 x.foo
 .bar
 .baz
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/EnvHome

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.29 | - |

Checks for consistent usage of `ENV['HOME']`. If `nil` is used as
the second argument of `ENV.fetch`, it is treated as a bad case like `ENV[]`.

## Style/EvalWithLocation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Ensures that eval methods (`eval`, `instance_eval`, `class_eval`
and `module_eval`) are given filename and line number values (`__FILE__`
and `__LINE__`). This data is used to ensure that any errors raised
within the evaluated code will be given the correct identification
in a backtrace.

The cop also checks that the line number given relative to `__LINE__` is
correct.

This cop will autocorrect incorrect or missing filename and line number
values. However, if `eval` is called without a binding argument, the cop
will not attempt to automatically add a binding, or add filename and
line values.

| This cop works only when a string literal is given as a code string. No offense is reported if a string variable is given as below: |

```
code = <<-RUBY
 def do_something
 end
RUBY
eval code # not checked.
```
## Style/EvenOdd

## Style/ExactRegexpMatch

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.51 | - |

Checks for exact regexp match inside `Regexp` literals.

## Style/ExpandPathArguments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

Checks for use of the `File.expand_path` arguments.
Likewise, it also checks for the `Pathname.new` argument.

Contrastive bad case and good case are alternately shown in the following examples.

### Examples

```
# bad
File.expand_path('..', __FILE__)
# good
File.expand_path(__dir__)
# bad
File.expand_path('../..', __FILE__)
# good
File.expand_path('..', __dir__)
# bad
File.expand_path('.', __FILE__)
# good
File.expand_path(__FILE__)
# bad
Pathname(__FILE__).parent.expand_path
# good
Pathname(__dir__).expand_path
# bad
Pathname.new(__FILE__).parent.expand_path
# good
Pathname.new(__dir__).expand_path
```
## Style/ExplicitBlockArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.89 | 1.8 |

Enforces the use of explicit block argument to avoid writing block literal that just passes its arguments to another block.

| This cop only registers an offense if the block args match the yield args exactly. |

### Examples

```
# bad
def with_tmp_dir
 Dir.mktmpdir do |tmp_dir|
 Dir.chdir(tmp_dir) { |dir| yield dir } # block just passes arguments
 end
end
# bad
def nine_times
 9.times { yield }
end
# good
def with_tmp_dir(&block)
 Dir.mktmpdir do |tmp_dir|
 Dir.chdir(tmp_dir, &block)
 end
end
with_tmp_dir do |dir|
 puts "dir is accessible as a parameter and pwd is set: #{dir}"
end
# good
def nine_times(&block)
 9.times(&block)
end
```
## Style/ExponentialNotation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.82 | - |

Enforces consistency when using exponential notation for numbers in the code (eg 1.2e4). Different styles are supported:

-
`scientific`which enforces a mantissa between 1 (inclusive) and 10 (exclusive).
-
`engineering`which enforces the exponent to be a multiple of 3 and the mantissa to be between 0.1 (inclusive) and 1000 (exclusive).
-
`integral`which enforces the mantissa to always be a whole number without trailing zeroes.

### Examples

#### EnforcedStyle: scientific (default)

```
# Enforces a mantissa between 1 (inclusive) and 10 (exclusive).
# bad
10e6
0.3e4
11.7e5
3.14e0
# good
1e7
3e3
1.17e6
3.14
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/FetchEnvVar

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.28 | 1.88 |

Suggests `ENV.fetch` for the replacement of `ENV[]`.
`ENV[]` silently fails and returns `nil` when the environment variable is unset,
which may cause unexpected behaviors when the developer forgets to set it.
On the other hand, `ENV.fetch` raises `KeyError` or returns the explicitly
specified default value.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowedVariables |
 | Array |
| DefaultToNil |
 | Boolean |

## Style/FileEmpty

| Requires Ruby version 2.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.48 | - |

Prefer to use `File.empty?('path/to/file')` when checking if a file is empty.

### Safety

This cop is unsafe, because `File.size`, `File.read`, and `File.binread`
raise `ENOENT` exception when there is no file corresponding to the path,
while `File.empty?` does not raise an exception.

### Examples

```
# bad
File.zero?('path/to/file')
File.size('path/to/file') == 0
File.size('path/to/file') >= 0
File.size('path/to/file').zero?
File.read('path/to/file').empty?
File.binread('path/to/file') == ''
FileTest.zero?('path/to/file')
# good
File.empty?('path/to/file')
FileTest.empty?('path/to/file')
```
## Style/FileNull

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.69 | - |

Use `File::NULL` instead of hardcoding the null device (`/dev/null` on Unix-like
OSes, `NUL` or `NUL:` on Windows), so that code is platform independent.
Only looks for full string matches, substrings within a longer string are not
considered.

However, only files that use the string `'/dev/null'` are targeted for detection.
This is because the string `'NUL'` is not limited to the null device.
This behavior results in false negatives when the `'/dev/null'` string is not used,
but it is a trade-off to avoid false positives. `NULL:`
Unlike `'NUL'`, `'NUL:'` is regarded as something like `C:` and is always detected.

| Uses inside arrays and hashes are ignored. |

## Style/FileOpen

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | No | 1.85 | - |

Checks for `File.open` without a block, which can leak file descriptors.

When `File.open` is called without a block, the caller is responsible
for closing the file descriptor. If it is not explicitly closed, it
will only be closed when the garbage collector runs, which may lead
to resource exhaustion. Using the block form ensures the file is
automatically closed when the block exits.

This cop only registers an offense when the result of `File.open` is
assigned to a variable or has a method chained on it, as those are the
clearest indicators that the block form should be used instead. When
`File.open` is used as a return value or passed as an argument, the
caller is likely managing the file descriptor intentionally.

### Safety

This cop is unsafe because it relies on syntax heuristics and cannot
verify whether the file descriptor is safely managed. For example, it
still flags intentional one-shot reads (`File.open("f").read`) where
the file descriptor is closed by the garbage collector.

### Examples

```
# bad
f = File.open('file')
# bad
File.open('file').read
# good
File.open('file') do |f|
 f.read
end
# good
File.open('file', &:read)
# good - pass an open file object to an API that manages its lifecycle
process(io: File.open('file'))
# good - return an open file object for the caller to manage
def json_key_io
 File.open('file')
end
# good - use File.read for one-shot reads
File.read('file')
```
## Style/FileRead

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.24 | - |

Favor `File.(bin)read` convenience methods.

### Examples

```
# bad - text mode
File.open(filename).read
File.open(filename, &:read)
File.open(filename) { |f| f.read }
File.open(filename) do |f|
 f.read
end
File.open(filename, 'r').read
File.open(filename, 'r', &:read)
File.open(filename, 'r') do |f|
 f.read
end
# good
File.read(filename)
# bad - binary mode
File.open(filename, 'rb').read
File.open(filename, 'rb', &:read)
File.open(filename, 'rb') do |f|
 f.read
end
# good
File.binread(filename)
```
## Style/FileTouch

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.69 | - |

Checks for usage of `File.open` in append mode with empty block.

Such a usage only creates a new file, but it doesn’t update timestamps for an existing file, which might have been the intention.

For example, for an existing file `foo.txt`:

```
ruby -e "puts File.mtime('foo.txt')"
# 2024-11-26 12:17:23 +0100
```
`ruby -e "File.open('foo.txt', 'a') {}"`
```
ruby -e "puts File.mtime('foo.txt')"
# 2024-11-26 12:17:23 +0100 -> unchanged
```
If the intention was to update timestamps, `FileUtils.touch('foo.txt')`
should be used instead.

## Style/FileWrite

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.24 | - |

Favor `File.(bin)write` convenience methods.

| There are different method signatures between `File.write`(class method)
and`File#write`(instance method). The following case will be allowed because
static analysis does not know the contents of the splat argument: |

```
File.open(filename, 'w') do |f|
 f.write(*objects)
end
```
### Examples

```
# bad - text mode
File.open(filename, 'w').write(content)
File.open(filename, 'w') do |f|
 f.write(content)
end
# good
File.write(filename, content)
# bad - binary mode
File.open(filename, 'wb').write(content)
File.open(filename, 'wb') do |f|
 f.write(content)
end
# good
File.binwrite(filename, content)
```
## Style/FloatDivision

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.72 | 1.9 |

Checks for division with integers coerced to floats.
It is recommended to either always use `fdiv` or coerce one side only.
This cop also provides other options for code consistency.

For `Regexp.last_match` and nth reference (e.g., `$1`), it assumes that the value
is a string matched by a regular expression, and allows conversion with `#to_f`.

### Safety

This cop is unsafe, because if the operand variable is a string object
then `#to_f` will be removed and an error will occur.

```
a = '1.2'
b = '3.4'
a.to_f / b.to_f # Both `to_f` calls are required here
```
## Style/For

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.13 | 1.26 |

Looks for uses of the `for` keyword or `each` method. The
preferred alternative is set in the EnforcedStyle configuration
parameter. An `each` call with a block on a single line is always
allowed.

| `each`is preferred in idiomatic Ruby because`for`leaks
its loop variable into the surrounding scope. |

### Safety

This cop’s autocorrection is unsafe because the scope of
variables is different between `each` and `for`.

### Examples

## Style/FormatString

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.19 | 0.49 |

Enforces the use of a single string formatting utility.
Valid options include `Kernel#format`, `Kernel#sprintf`, and `String#%`.

The detection of `String#%` cannot be implemented in a reliable
manner for all cases, so only two scenarios are considered -
if the first argument is a string literal and if the second
argument is an array literal.

Autocorrection will be applied when the argument is a literal or uses a known
built-in conversion method such as `to_d`, `to_f`, `to_h`, `to_i`, `to_r`, `to_s`,
and `to_sym` on variables, provided that their return value is not an array.
For example, when using `to_s`,
`'%s' % [1, 2, 3].to_s` can be autocorrected without any incompatibility:

```
'%s' % [1, 2, 3] #=> '1'
format('%s', [1, 2, 3]) #=> '[1, 2, 3]'
'%s' % [1, 2, 3].to_s #=> '[1, 2, 3]'
```
### Examples

#### EnforcedStyle: format (default)

```
# bad
puts sprintf('%10s', 'foo')
puts '%10s' % 'foo'
# good
puts format('%10s', 'foo')
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/FormatStringToken

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.74 |

Use a consistent style for tokens within a format string.

By default, all strings are evaluated. In some cases, this may be undesirable, as they could be used as arguments to a method that does not consider them to be tokens, but rather other identifiers or just part of the string.

`AllowedMethods` or `AllowedPatterns` can be configured with in order to mark specific
methods as always allowed, thereby avoiding an offense from the cop. By default, there
are no allowed methods.

Additionally, the cop can be made conservative by configuring it with
`Mode: conservative` (default `aggressive`). In this mode, tokens (regardless
of `EnforcedStyle`) are only considered if used in the format string argument to the
methods `printf`, `sprintf`, `format` and `%`.

| In `aggressive`mode, offenses are registered for all strings containing tokens,
but autocorrection is only applied when the string appears in a known formatting context
(`format`,`sprintf`,`printf`, or`%`). This is done in order to prevent false
autocorrections for strings that are not actually format strings. |

| Tokens in the `unannotated`style (eg.`%s`) are always treated as if
configured with`Conservative: true`. This is done in order to prevent false positives,
because this format is very similar to encoded URLs or Date/Time formatting strings. |

It is allowed to contain unannotated token
if the number of them is less than or equals to
`MaxUnannotatedPlaceholdersAllowed`.

### Examples

#### EnforcedStyle: annotated (default)

```
# bad
format('%{greeting}', greeting: 'Hello')
format('%s', 'Hello')
# good
format('%<greeting>s', greeting: 'Hello')
```
#### EnforcedStyle: template

```
# bad
format('%<greeting>s', greeting: 'Hello')
format('%s', 'Hello')
# good
format('%{greeting}', greeting: 'Hello')
```
#### EnforcedStyle: unannotated

```
# bad
format('%<greeting>s', greeting: 'Hello')
format('%{greeting}', greeting: 'Hello')
# good
format('%s', 'Hello')
```
#### MaxUnannotatedPlaceholdersAllowed: 0

```
# bad
format('%06d', 10)
format('%s %s.', 'Hello', 'world')
# good
format('%<number>06d', number: 10)
```
#### MaxUnannotatedPlaceholdersAllowed: 1 (default)

```
# bad
format('%s %s.', 'Hello', 'world')
# good
format('%06d', 10)
```
#### Mode: aggressive (default), EnforcedStyle: annotated

```
# bad
"%{greeting}"
foo("%{greeting}")
# bad
format("%{greeting}", greeting: 'Hello')
printf("%{greeting}", greeting: 'Hello')
sprintf("%{greeting}", greeting: 'Hello')
"%{greeting}" % { greeting: 'Hello' }
# good
format("%<greeting>s", greeting: 'Hello')
printf("%<greeting>s", greeting: 'Hello')
sprintf("%<greeting>s", greeting: 'Hello')
"%<greeting>s" % { greeting: 'Hello' }
```
#### Mode: conservative, EnforcedStyle: annotated

```
# good
"%{greeting}"
foo("%{greeting}")
# bad
format("%{greeting}", greeting: 'Hello')
printf("%{greeting}", greeting: 'Hello')
sprintf("%{greeting}", greeting: 'Hello')
"%{greeting}" % { greeting: 'Hello' }
# good
format("%<greeting>s", greeting: 'Hello')
printf("%<greeting>s", greeting: 'Hello')
sprintf("%<greeting>s", greeting: 'Hello')
"%<greeting>s" % { greeting: 'Hello' }
```
## Style/FrozenStringLiteralComment

| Requires Ruby version 2.3 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.36 | 0.79 |

Helps you transition from mutable string literals
to frozen string literals.
It will add the `# frozen_string_literal: true` magic comment to the top
of files to enable frozen string literals. Frozen string literals may be
default in future Ruby. The comment will be added below a shebang and
encoding comment. The frozen string literal comment is only valid in Ruby 2.3+.

Note that the cop will accept files where the comment exists but is set
to `false` instead of `true`.

To require a blank line after this comment, please see
`Layout/EmptyLineAfterMagicComment` cop.

### Safety

This cop’s autocorrection is unsafe since any strings mutations will
change from being accepted to raising `FrozenError`, as all strings
will become frozen by default, and will need to be manually refactored.

### Examples

#### EnforcedStyle: always (default)

```
# The `always` style will always add the frozen string literal comment
# to a file, regardless of the Ruby version or if `freeze` or `<<` are
# called on a string literal.
# bad
module Bar
 # ...
end
# good
# frozen_string_literal: true
module Bar
 # ...
end
# good
# frozen_string_literal: false
module Bar
 # ...
end
```
#### EnforcedStyle: never

```
# The `never` will enforce that the frozen string literal comment does
# not exist in a file.
# bad
# frozen_string_literal: true
module Baz
 # ...
end
# good
module Baz
 # ...
end
```
#### EnforcedStyle: always_true

```
# The `always_true` style enforces that the frozen string literal
# comment is set to `true`. This is a stricter option than `always`
# and forces projects to use frozen string literals.
# bad
# frozen_string_literal: false
module Baz
 # ...
end
# bad
module Baz
 # ...
end
# good
# frozen_string_literal: true
module Bar
 # ...
end
```
## Style/GlobalStdStream

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.89 | - |

Enforces the use of `$stdout/$stderr/$stdin` instead of `STDOUT/STDERR/STDIN`.
`STDOUT/STDERR/STDIN` are constants, and while you can actually
reassign (possibly to redirect some stream) constants in Ruby, you’ll get
an interpreter warning if you do so.

Additionally, `$stdout/$stderr/$stdin` can safely be accessed in a Ractor because they
are ractor-local, while `STDOUT/STDERR/STDIN` will raise `Ractor::IsolationError`.

### Safety

Autocorrection is unsafe because `STDOUT` and `$stdout` may point to different
objects, for example.

### Examples

```
# bad
STDOUT.puts('hello')
hash = { out: STDOUT, key: value }
def m(out = STDOUT)
 out.puts('hello')
end
# good
$stdout.puts('hello')
hash = { out: $stdout, key: value }
def m(out = $stdout)
 out.puts('hello')
end
```
## Style/GlobalVars

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.13 | - |

Looks for uses of global variables. Global variables introduce shared mutable state that makes code harder to test, debug, and reason about, since any part of the program can read or modify them.

It does not report offenses for built-in global variables. Built-in global variables are allowed by default. Additionally users can allow additional variables via the AllowedVariables option.

Note that backreferences like $1, $2, etc are not global variables.

## Style/GuardClause

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.20 | 1.31 |

Use a guard clause instead of wrapping the code inside a conditional expression

A condition with an `elsif` or `else` branch is allowed unless
one of `return`, `break`, `next`, `raise`, or `fail` is used
in the body of the conditional expression.

| Autocorrect works in most cases except with if-else statements
 that contain logical operators such as `foo || raise('exception')` |

### Examples

```
# bad
def test
 if something
 work
 end
end
# good
def test
 return unless something
 work
end
# also good
def test
 work if something
end
# bad
if something
 raise 'exception'
else
 ok
end
# good
raise 'exception' if something
ok
# bad
define_method(:test) do
 if something
 work
 end
end
# good
define_method(:test) do
 return unless something
 work
end
# also good
define_method(:test) do
 work if something
end
```
## Style/HashAsLastArrayItem

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.88 | - |

Checks for presence or absence of braces around hash literal as a last array item depending on configuration.

| This cop will ignore arrays where multiple items are all hashes,
regardless of `EnforcedStyle`. |

`[{ one: 1 }, { two: 2 }]`### Examples

## Style/HashConversion

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.10 | 1.55 |

Checks the usage of pre-2.1 `Hash[args]` method of converting enumerables and
sequences of values to hashes.

Correction code from splat argument (`Hash[*ary]`) is not simply determined. For example,
`Hash[*ary]` can be replaced with `ary.each_slice(2).to_h` but it will be complicated.
So, `AllowSplatArgument` option is true by default to allow splat argument for simple code.

### Safety

This cop’s autocorrection is unsafe because `ArgumentError` occurs
if the number of elements is odd:

```
Hash[[[1, 2], [3]]] #=> {1=>2, 3=>nil}
[[1, 2], [5]].to_h #=> wrong array length at 1 (expected 2, was 1) (ArgumentError)
```
## Style/HashEachMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.80 | 1.16 |

Checks for uses of `each_key` and `each_value` `Hash` methods.

| If you have an array of two-element arrays, you can put parentheses around the block arguments to indicate that you’re not working with a hash, and suppress RuboCop offenses. |

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
is a `Hash`. The `AllowedReceivers` configuration can mitigate,
but not fully resolve, this safety issue.

### Examples

```
# bad
hash.keys.each { |k| p k }
hash.each { |k, unused_value| p k }
# good
hash.each_key { |k| p k }
# bad
hash.values.each { |v| p v }
hash.each { |unused_key, v| p v }
# good
hash.each_value { |v| p v }
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowedReceivers |
 | Array |

## Style/HashExcept

| Requires Ruby version 3.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.7 | 1.39 |

Checks for usages of `Hash#reject`, `Hash#select`, and `Hash#filter` methods
that can be replaced with `Hash#except` method.

This cop should only be enabled on Ruby version 3.0 or higher.
(`Hash#except` was added in Ruby 3.0.)

For safe detection, it is limited to commonly used string and symbol comparisons
when using `==` or `!=`.

This cop doesn’t check for `Hash#delete_if` and `Hash#keep_if` because they
modify the receiver.

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
is a `Hash` or responds to the replacement method.

### Examples

```
# bad
{foo: 1, bar: 2, baz: 3}.reject {|k, v| k == :bar }
{foo: 1, bar: 2, baz: 3}.select {|k, v| k != :bar }
{foo: 1, bar: 2, baz: 3}.filter {|k, v| k != :bar }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| k.eql?(:bar) }
# bad
{foo: 1, bar: 2, baz: 3}.reject {|k, v| %i[bar].include?(k) }
{foo: 1, bar: 2, baz: 3}.select {|k, v| !%i[bar].include?(k) }
{foo: 1, bar: 2, baz: 3}.filter {|k, v| !%i[bar].include?(k) }
# good
{foo: 1, bar: 2, baz: 3}.except(:bar)
```
#### AllCops:ActiveSupportExtensionsEnabled: false (default)

```
# good
{foo: 1, bar: 2, baz: 3}.reject {|k, v| !%i[bar].exclude?(k) }
{foo: 1, bar: 2, baz: 3}.select {|k, v| %i[bar].exclude?(k) }
# good
{foo: 1, bar: 2, baz: 3}.reject {|k, v| k.in?(%i[bar]) }
{foo: 1, bar: 2, baz: 3}.select {|k, v| !k.in?(%i[bar]) }
```
#### AllCops:ActiveSupportExtensionsEnabled: true

```
# bad
{foo: 1, bar: 2, baz: 3}.reject {|k, v| !%i[bar].exclude?(k) }
{foo: 1, bar: 2, baz: 3}.select {|k, v| %i[bar].exclude?(k) }
# bad
{foo: 1, bar: 2, baz: 3}.reject {|k, v| k.in?(%i[bar]) }
{foo: 1, bar: 2, baz: 3}.select {|k, v| !k.in?(%i[bar]) }
# good
{foo: 1, bar: 2, baz: 3}.except(:bar)
```
## Style/HashFetchChain

| Requires Ruby version 2.3 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.75 | - |

Use `Hash#dig` instead of chaining potentially null `fetch` calls.

When `fetch(identifier, nil)` calls are chained on a hash, the expectation
is that each step in the chain returns either `nil` or another hash,
and in both cases, these can be simplified with a single call to `dig` with
multiple arguments.

If the 2nd parameter is `{}` or `Hash.new`, an offense will also be registered,
as long as the final call in the chain is a nil value. If a non-nil value is given,
the chain will not be registered as an offense, as the default value cannot be safely
given with `dig`.

| See `Style/DigChain`for replacing chains of`dig`calls with
a single method call. |

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
is a `Hash` or that `fetch` or `dig` have the expected standard implementation.

### Examples

```
# bad
hash.fetch('foo', nil)&.fetch('bar', nil)
# bad
# earlier members of the chain can return `{}` as long as the final `fetch`
# has `nil` as a default value
hash.fetch('foo', {}).fetch('bar', nil)
# good
hash.dig('foo', 'bar')
# ok - not handled by the cop since the final `fetch` value is non-nil
hash.fetch('foo', {}).fetch('bar', {})
```
## Style/HashLikeCase

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.88 | - |

Checks for places where `case-when` represents a simple 1:1
mapping and can be replaced with a hash lookup.

### Examples

#### MinBranchesCount: 3 (default)

```
# bad
case country
when 'europe'
 'http://eu.example.com'
when 'america'
 'http://us.example.com'
when 'australia'
 'http://au.example.com'
end
# good
SITES = {
 'europe' => 'http://eu.example.com',
 'america' => 'http://us.example.com',
 'australia' => 'http://au.example.com'
}
SITES[country]
```
## Style/HashLookupMethod

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 1.84 | - |

Enforces the use of either `Hash#[]` or `Hash#fetch` for hash lookup.

This cop can be configured to prefer either bracket-style (`[]`)
or fetch-style lookup. It is disabled by default.

When enforcing `fetch` style, only single-argument bracket access is flagged.
When enforcing `brackets` style, only `fetch` calls with a single key
argument are flagged (not those with default values or blocks).

### Safety

This cop is unsafe because `Hash#[]` and `Hash#fetch` have different
semantics. `Hash#[]` returns `nil` for missing keys, while `Hash#fetch`
raises a `KeyError`. Replacing one with the other can change program
behavior in cases where the key is missing.

Additionally, it cannot be guaranteed that the receiver is a `Hash`
or responds to the replacement method.

## Style/HashSlice

| Requires Ruby version 2.5 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.71 | - |

Checks for usages of `Hash#reject`, `Hash#select`, and `Hash#filter` methods
that can be replaced with `Hash#slice` method.

This cop should only be enabled on Ruby version 2.5 or higher.
(`Hash#slice` was added in Ruby 2.5.)

For safe detection, it is limited to commonly used string and symbol comparisons
when using `==` or `!=`.

This cop doesn’t check for `Hash#delete_if` and `Hash#keep_if` because they
modify the receiver.

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
is a `Hash` or responds to the replacement method.

Additionally, the replacement may change the order of the resulting
hash: `Hash#slice` returns entries in the order the keys are given,
whereas `select`, `filter`, and `reject` preserve the entry order of
the receiver.

For example:

```
hash = {foo: 1, bar: 2, baz: 3}
keys = %i[baz foo]
hash.select { |k, _v| keys.include?(k) } # => {foo: 1, baz: 3}
hash.slice(*keys) # => {baz: 3, foo: 1}
```
### Examples

```
# bad
{foo: 1, bar: 2, baz: 3}.select {|k, v| k == :bar }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| k != :bar }
{foo: 1, bar: 2, baz: 3}.filter {|k, v| k == :bar }
{foo: 1, bar: 2, baz: 3}.select {|k, v| k.eql?(:bar) }
# bad
{foo: 1, bar: 2, baz: 3}.select {|k, v| %i[bar].include?(k) }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| !%i[bar].include?(k) }
{foo: 1, bar: 2, baz: 3}.filter {|k, v| %i[bar].include?(k) }
# good
{foo: 1, bar: 2, baz: 3}.slice(:bar)
```
#### AllCops:ActiveSupportExtensionsEnabled: false (default)

```
# good
{foo: 1, bar: 2, baz: 3}.select {|k, v| !%i[bar].exclude?(k) }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| %i[bar].exclude?(k) }
# good
{foo: 1, bar: 2, baz: 3}.select {|k, v| k.in?(%i[bar]) }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| !k.in?(%i[bar]) }
```
#### AllCops:ActiveSupportExtensionsEnabled: true

```
# bad
{foo: 1, bar: 2, baz: 3}.select {|k, v| !%i[bar].exclude?(k) }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| %i[bar].exclude?(k) }
# bad
{foo: 1, bar: 2, baz: 3}.select {|k, v| k.in?(%i[bar]) }
{foo: 1, bar: 2, baz: 3}.reject {|k, v| !k.in?(%i[bar]) }
# good
{foo: 1, bar: 2, baz: 3}.slice(:bar)
```
## Style/HashSyntax

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 1.67 |

Checks hash literal syntax.

It can enforce either the use of the class hash rocket syntax or the use of the newer Ruby 1.9 syntax (when applicable).

A separate offense is registered for each problematic pair.

The supported styles are:

-
ruby19 - forces use of the 1.9 syntax (e.g. `{a: 1}`) when hashes have all symbols for keys
-
hash_rockets - forces use of hash rockets for all hashes
-
no_mixed_keys - simply checks for hashes with mixed syntaxes
-
ruby19_no_mixed_keys - forces use of ruby 1.9 syntax and forbids mixed syntax hashes

This cop has `EnforcedShorthandSyntax` option.
It can enforce either the use of the explicit hash value syntax or
the use of Ruby 3.1’s hash value shorthand syntax.

The supported styles are:

-
always - forces use of the 3.1 syntax (e.g. {foo:})
-
never - forces use of explicit hash literal value
-
either - accepts both shorthand and explicit use of hash literal value
-
consistent - forces use of the 3.1 syntax only if all values can be omitted in the hash
-
either_consistent - accepts both shorthand and explicit use of hash literal value, but they must be consistent

### Examples

#### EnforcedStyle: ruby19 (default)

```
# bad
{:a => 2}
{b: 1, :c => 2}
# good
{a: 2, b: 1}
{:c => 2, 'd' => 2} # acceptable since 'd' isn't a symbol
{d: 1, 'e' => 2} # technically not forbidden
```
#### EnforcedStyle: no_mixed_keys

```
# bad
{:a => 1, b: 2}
{c: 1, 'd' => 2}
# good
{:a => 1, :b => 2}
{c: 1, d: 2}
```
#### EnforcedStyle: ruby19_no_mixed_keys

```
# bad
{:a => 1, :b => 2}
{c: 2, 'd' => 3} # should just use hash rockets
# good
{a: 1, b: 2}
{:c => 3, 'd' => 4}
```
#### EnforcedShorthandSyntax: always

```
# bad
{foo: foo, bar: bar}
# good
{foo:, bar:}
# good - allowed to mix syntaxes
{foo:, bar: baz}
```
#### EnforcedShorthandSyntax: either (default)

```
# good
{foo: foo, bar: bar}
# good
{foo: foo, bar:}
# good
{foo:, bar:}
```
#### EnforcedShorthandSyntax: consistent

```
# bad - `foo` and `bar` values can be omitted
{foo: foo, bar: bar}
# bad - `bar` value can be omitted
{foo:, bar: bar}
# bad - mixed syntaxes
{foo:, bar: baz}
# good
{foo:, bar:}
# good - can't omit `baz`
{foo: foo, bar: baz}
```
#### EnforcedShorthandSyntax: either_consistent

```
# good - `foo` and `bar` values can be omitted, but they are consistent, so it's accepted
{foo: foo, bar: bar}
# bad - `bar` value can be omitted
{foo:, bar: bar}
# bad - mixed syntaxes
{foo:, bar: baz}
# good
{foo:, bar:}
# good - can't omit `baz`
{foo: foo, bar: baz}
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| EnforcedShorthandSyntax |
 |
 |
| UseHashRocketsWithSymbolValues |
 | Boolean |
| PreferHashRocketsForNonAlnumEndingSymbols |
 | Boolean |

## Style/HashTransformKeys

| Requires Ruby version 2.5 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.80 | 0.90 |

Looks for uses of `_.each_with_object({}) {...}`,
`_.map {...}.to_h`, and `Hash[_.map {...}]` that are actually just
transforming the keys of a hash, and tries to use a simpler & faster
call to `transform_keys` instead.
It should only be enabled on Ruby version 2.5 or newer.
(`transform_keys` was added in Ruby 2.5.)

### Safety

This cop identifies the receiver as a hash by checking for literal hash
syntax and common methods that are known to return hashes (e.g. `to_h`,
`merge`, `invert`, `group_by`, etc.). However, it is unsafe because it
is possible for a custom class to define one of these methods and return
something other than a hash.

### Examples

```
# bad
{a: 1, b: 2}.each_with_object({}) { |(k, v), h| h[foo(k)] = v }
Hash[{a: 1, b: 2}.collect { |k, v| [foo(k), v] }]
{a: 1, b: 2}.map { |k, v| [k.to_s, v] }.to_h
{a: 1, b: 2}.to_h { |k, v| [k.to_s, v] }
foo.to_h.each_with_object({}) { |(k, v), h| h[k.to_sym] = v }
foo.merge(bar).map { |k, v| [k.to_s, v] }.to_h
# good
{a: 1, b: 2}.transform_keys { |k| foo(k) }
{a: 1, b: 2}.transform_keys { |k| k.to_s }
foo.to_h.transform_keys { |k| k.to_sym }
foo.merge(bar).transform_keys { |k| k.to_s }
# Won't register an offense - receiver is not known to be a hash
foo.bar.each_with_object({}) { |(k, v), h| h[k.to_s] = v }
baz.map { |k, v| [k.to_s, v] }.to_h
```
## Style/HashTransformValues

| Requires Ruby version 2.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.80 | 0.90 |

Looks for uses of `_.each_with_object({}) {...}`,
`_.map {...}.to_h`, and `Hash[_.map {...}]` that are actually just
transforming the values of a hash, and tries to use a simpler & faster
call to `transform_values` instead.

### Safety

This cop identifies the receiver as a hash by checking for literal hash
syntax and common methods that are known to return hashes (e.g. `to_h`,
`merge`, `invert`, `group_by`, etc.). However, it is unsafe because it
is possible for a custom class to define one of these methods and return
something other than a hash.

### Examples

```
# bad
{a: 1, b: 2}.each_with_object({}) { |(k, v), h| h[k] = foo(v) }
Hash[{a: 1, b: 2}.collect { |k, v| [k, foo(v)] }]
{a: 1, b: 2}.map { |k, v| [k, v * v] }.to_h
{a: 1, b: 2}.to_h { |k, v| [k, v * v] }
foo.to_h.each_with_object({}) { |(k, v), h| h[k] = foo(v) }
foo.merge(bar).map { |k, v| [k, v.to_s] }.to_h
# good
{a: 1, b: 2}.transform_values { |v| foo(v) }
{a: 1, b: 2}.transform_values { |v| v * v }
foo.to_h.transform_values { |v| foo(v) }
foo.merge(bar).transform_values { |v| v.to_s }
# Won't register an offense - receiver is not known to be a hash
foo.bar.each_with_object({}) { |(k, v), h| h[k] = v.to_s }
baz.map { |k, v| [k, v.to_s] }.to_h
```
## Style/IdenticalConditionalBranches

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.36 | 1.19 |

Checks for identical expressions at the beginning or end of each branch of a conditional expression. Such expressions should normally be placed outside the conditional expression - before or after it.

| The cop is poorly named and some people might think that it actually checks for duplicated conditional branches. The name will probably be changed in a future major RuboCop release. |

### Safety

Autocorrection is unsafe because changing the order of method invocations may change the behavior of the code. For example:

```
if method_that_modifies_global_state # 1
 method_that_relies_on_global_state # 2
 foo # 3
else
 method_that_relies_on_global_state # 2
 bar # 3
end
```
In this example, `method_that_relies_on_global_state` will be moved before
`method_that_modifies_global_state`, which changes the behavior of the program.

### Examples

```
# bad
if condition
 do_x
 do_z
else
 do_y
 do_z
end
# good
if condition
 do_x
else
 do_y
end
do_z
# bad
if condition
 do_z
 do_x
else
 do_z
 do_y
end
# good
do_z
if condition
 do_x
else
 do_y
end
# bad
case foo
when 1
 do_x
when 2
 do_x
else
 do_x
end
# good
case foo
when 1
 do_x
 do_y
when 2
 # nothing
else
 do_x
 do_z
end
# bad
case foo
in 1
 do_x
in 2
 do_x
else
 do_x
end
# good
case foo
in 1
 do_x
 do_y
in 2
 # nothing
else
 do_x
 do_z
end
```
## Style/IfInsideElse

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.36 | 1.3 |

If the `else` branch of a conditional consists solely of an `if` node,
it can be combined with the `else` to become an `elsif`.
This helps to keep the nesting level from getting too deep.

### Examples

```
# bad
if condition_a
 action_a
else
 if condition_b
 action_b
 else
 action_c
 end
end
# good
if condition_a
 action_a
elsif condition_b
 action_b
else
 action_c
end
```
## Style/IfUnlessModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.30 |

Checks for `if` and `unless` statements that would fit on one line if
written as modifier `if`/`unless`. The cop also checks for modifier
`if`/`unless` lines that exceed the maximum line length.

The maximum line length is configured in the `Layout/LineLength`
cop. The tab size is configured in the `IndentationWidth` of the
`Layout/IndentationStyle` cop.

One-line pattern matching is always allowed. To ensure that there are few cases
where the match variable is not used, and to prevent oversights. The variable `x`
becomes undefined and raises `NameError` when the following example is changed to
the modifier form:

```
if [42] in [x]
 x # `x` is undefined when using modifier form.
end
```
The code `def method_name = body if condition` is considered a bad case by
`Style/AmbiguousEndlessMethodDefinition` cop. So, to respect the user’s intention to use
an endless method definition in the `if` body, the following code is allowed:

```
if condition
 def method_name = body
end
```
| It is allowed when `defined?`argument has an undefined value,
because using the modifier form causes the following incompatibility: |

```
unless defined?(undefined_foo)
 undefined_foo = 'default_value'
end
undefined_foo # => 'default_value'
undefined_bar = 'default_value' unless defined?(undefined_bar)
undefined_bar # => nil
```
### Examples

```
# bad
if condition
 do_stuff(bar)
end
unless qux.empty?
 Foo.do_something
end
do_something_with_a_long_name(arg) if long_condition_that_prevents_code_fit_on_single_line
# good
do_stuff(bar) if condition
Foo.do_something unless qux.empty?
if long_condition_that_prevents_code_fit_on_single_line
 do_something_with_a_long_name(arg)
end
if short_condition # a long comment that makes it too long if it were just a single line
 do_something
end
```
## Style/IfUnlessModifierOfIfUnless

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.39 | 0.87 |

Checks for if and unless statements used as modifiers of other if or unless statements.

## Style/IfWithBooleanLiteralBranches

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.9 | - |

Checks for redundant `if` with boolean literal branches.
It checks only conditions to return boolean value (`true` or `false`) for safe detection.
The conditions to be checked are comparison methods, predicate methods, and
double negation (!!).
`nonzero?` method is allowed by default.
These are customizable with `AllowedMethods` option.

This cop targets only `if`s with a single `elsif` or `else` branch. The following
code will be allowed, because it has two `elsif` branches:

```
if foo
 true
elsif bar > baz
 true
elsif qux > quux # Single `elsif` is warned, but two or more `elsif`s are not.
 true
else
 false
end
```
### Safety

Autocorrection is unsafe because there is no guarantee that all predicate methods
will return a boolean value. Those methods can be allowed with `AllowedMethods` config.

## Style/IfWithSemicolon

## Style/ImplicitRuntimeError

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.41 | - |

Checks for `raise` or `fail` statements which do not specify an
explicit exception class. (This raises a `RuntimeError`. Some projects
might prefer to use exception classes which more precisely identify the
nature of the error.)

## Style/InPatternThen

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.16 | - |

Checks for `in;` uses in `case` expressions.

## Style/InfiniteLoop

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.26 | 0.61 |

Use `Kernel#loop` for infinite loops.

### Safety

This cop is unsafe as the rule should not necessarily apply if the loop
body might raise a `StopIteration` exception; contrary to other infinite
loops, `Kernel#loop` silently rescues that and returns `nil`.

## Style/InlineComment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.23 | - |

Checks for trailing inline comments. Inline comments can make lines harder to read, especially when they are long. Placing comments on their own line above the code they describe is often clearer.

## Style/InverseMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.48 | - |

Checks for usages of not (`not` or `!`) called on a method
when an inverse of that method can be used instead.

Methods that can be inverted by a not (`not` or `!`) should be defined
in `InverseMethods`.

Methods that are inverted by inverting the return
of the block that is passed to the method should be defined in
`InverseBlocks`.

### Safety

This cop is unsafe because it cannot be guaranteed that the method and its inverse method are both defined on receiver, and also are actually inverse of each other.

### Examples

```
# bad
!foo.none?
!foo.any? { |f| f.even? }
!foo.blank?
!(foo == bar)
foo.select { |f| !f.even? }
foo.reject { |f| f != 7 }
# good
foo.none?
foo.blank?
foo.any? { |f| f.even? }
foo != bar
foo == bar
!!('foo' =~ /^\w+$/)
!(foo.class < Numeric) # Checking class hierarchy is allowed
# Blocks with guard clauses are ignored:
foo.select do |f|
 next if f.zero?
 f != 1
end
```
## Style/InvertibleUnlessCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 1.44 | 1.50 |

Checks for usages of `unless` which can be replaced by `if` with inverted condition.
Code without `unless` is easier to read, but that is subjective, so this cop
is disabled by default.

Methods that can be inverted should be defined in `InverseMethods`. Note that
the relationship of inverse methods needs to be defined in both directions.
For example,

```
InverseMethods:
 :!=: :==
 :even?: :odd?
 :odd?: :even?
```
will suggest both `even?` and `odd?` to be inverted, but only `!=` (and not `==`).

### Safety

This cop is unsafe because it cannot be guaranteed that the method and its inverse method are both defined on receiver, and also are actually inverse of each other.

### Examples

```
# bad (simple condition)
foo unless !bar
foo unless x != y
foo unless x >= 10
foo unless x.even?
foo unless odd?
# good
foo if bar
foo if x == y
foo if x < 10
foo if x.odd?
foo if even?
# bad (complex condition)
foo unless x != y || x.even?
# good
foo if x == y && x.odd?
# good (if)
foo if !condition
```
## Style/IpAddresses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.58 | 0.91 |

Checks for hardcoded IP addresses, which can make code brittle. IP addresses are likely to need to be changed when code is deployed to a different server or environment, which may break a deployment if forgotten. Prefer setting IP addresses in ENV or other configuration.

## Style/ItAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.70 | - |

Checks for local variables and method parameters named `it`,
where `it` can refer to the first anonymous parameter as of Ruby 3.4.
Use a meaningful variable name instead.

| Although Ruby allows reassigning `it`in these cases, it could
cause confusion if`it`is used as a block parameter elsewhere. |

### Examples

```
# bad
it = 5
# good
var = 5
# bad
def foo(it)
end
# good
def foo(arg)
end
# bad
def foo(it = 5)
end
# good
def foo(arg = 5)
end
# bad
def foo(*it)
end
# good
def foo(*args)
end
# bad
def foo(it:)
end
# good
def foo(arg:)
end
# bad
def foo(it: 5)
end
# good
def foo(arg: 5)
end
# bad
def foo(**it)
end
# good
def foo(**kwargs)
end
# bad
def foo(&it)
end
# good
def foo(&block)
end
```
## Style/ItBlockParameter

| Requires Ruby version 3.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.75 | 1.76 |

Checks for blocks with one argument where `it` block parameter can be used.

It provides four `EnforcedStyle` options:

-
`allow_single_line`(default) … Always uses the`it`block parameter in a single line.
-
`only_numbered_parameters`… Detects only numbered block parameters.
-
`always`… Always uses the`it`block parameter.
-
`disallow`… Disallows the`it`block parameter.

A single numbered parameter is detected when `allow_single_line`,
`only_numbered_parameters`, or `always`.

### Examples

#### EnforcedStyle: allow_single_line (default)

```
# bad
block do
 do_something(it)
end
block { do_something(_1) }
# good
block { do_something(it) }
block { |named_param| do_something(named_param) }
```
#### EnforcedStyle: only_numbered_parameters

```
# bad
block { do_something(_1) }
# good
block { do_something(it) }
block { |named_param| do_something(named_param) }
```
## Style/KeywordArgumentsMerging

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.68 | - |

When passing an existing hash as keyword arguments, provide additional arguments
directly rather than using `merge`.

Providing arguments directly is more performant than using `merge`, and
also leads to shorter and simpler code.

## Style/KeywordParametersOrder

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.90 | 1.7 |

Enforces that optional keyword parameters are placed at the end of the parameters list.

This improves readability, because when looking through the source, it is expected to find required parameters at the beginning of parameters list and optional parameters at the end.

### Examples

```
# bad
def some_method(first: false, second:, third: 10)
 # body omitted
end
# good
def some_method(second:, first: false, third: 10)
 # body omitted
end
# bad
do_something do |first: false, second:, third: 10|
 # body omitted
end
# good
do_something do |second:, first: false, third: 10|
 # body omitted
end
```
## Style/Lambda

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.40 |

(by default) checks for uses of the lambda literal syntax for single line lambdas, and the method call syntax for multiline lambdas. It is configurable to enforce one of the styles for both single line and multiline lambdas as well.

### Examples

#### EnforcedStyle: line_count_dependent (default)

```
# bad
f = lambda { |x| x }
f = ->(x) do
 x
 end
# good
f = ->(x) { x }
f = lambda do |x|
 x
 end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/LambdaCall

## Style/LineEndConcatenation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.18 | 0.64 |

Checks for string literal concatenation at the end of a line.

## Style/MagicCommentFormat

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.35 | - |

Ensures magic comments are written consistently throughout your code base.
Looks for discrepancies in separators (`-` vs `_`) and capitalization for
both magic comment directives and values.

Required capitalization can be set with the `DirectiveCapitalization` and
`ValueCapitalization` configuration keys.

| If one of these configurations is set to nil, any capitalization is allowed. |

### Examples

#### EnforcedStyle: snake_case (default)

```
# The `snake_case` style will enforce that the frozen string literal
# comment is written in snake case. (Words separated by underscores)
# bad
# frozen-string-literal: true
module Bar
 # ...
end
# good
# frozen_string_literal: false
module Bar
 # ...
end
```
#### EnforcedStyle: kebab_case

```
# The `kebab_case` style will enforce that the frozen string literal
# comment is written in kebab case. (Words separated by hyphens)
# bad
# frozen_string_literal: true
module Baz
 # ...
end
# good
# frozen-string-literal: true
module Baz
 # ...
end
```
#### DirectiveCapitalization: lowercase (default)

```
# bad
# FROZEN-STRING-LITERAL: true
# good
# frozen-string-literal: true
```
#### DirectiveCapitalization: uppercase

```
# bad
# frozen-string-literal: true
# good
# FROZEN-STRING-LITERAL: true
```
#### DirectiveCapitalization: nil

```
# any capitalization is accepted
# good
# frozen-string-literal: true
# good
# FROZEN-STRING-LITERAL: true
```
#### ValueCapitalization: nil (default)

```
# any capitalization is accepted
# good
# frozen-string-literal: true
# good
# frozen-string-literal: TRUE
```
## Style/MapCompactWithConditionalBlock

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.30 | 1.88 |

Prefer `select` or `reject` over `map { … }.compact`.
This cop also handles `filter_map { … }`, similar to `map { … }.compact`.

### Safety

This cop is unsafe because `compact` also removes `nil` elements that
were already present in the receiver, whereas `select`/`reject` keep
them. The result therefore differs when the collection contains `nil`:

```
[nil, 1].map { |e| e if e }.compact # => [1]
[nil, 1].select { |e| e } # => [nil, 1]
```
### Examples

```
# bad
array.map { |e| some_condition? ? e : next }.compact
# bad
array.filter_map { |e| some_condition? ? e : next }
# bad
array.map do |e|
 if some_condition?
 e
 else
 next
 end
end.compact
# bad
array.map do |e|
 next if some_condition?
 e
end.compact
# bad
array.map do |e|
 e if some_condition?
end.compact
# good
array.select { |e| some_condition? }
# good
array.reject { |e| some_condition? }
```
## Style/MapIntoArray

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.63 | 1.67 |

Checks for usages of `each` with `<<`, `push`, or `append` which
can be replaced by `map`.

If `PreferredMethods` is configured for `map` in `Style/CollectionMethods`,
this cop uses the specified method for replacement.

| The return value of `Enumerable#each`is`self`, whereas the
return value of`Enumerable#map`is an`Array`. They are not autocorrected
when a return value could be used because these types differ. |

| It only detects when the mapping destination is either:
* a local variable initialized as an empty array and referred to only by the
pushing operation;
* or, if it is the single block argument to a `[].tap`block.
This is because, if not, it’s challenging to statically guarantee that the
mapping destination variable remains an empty array: |

```
ret = []
src.each { |e| ret << e * 2 } # `<<` method may mutate `ret`
dest = []
src.each { |e| dest << transform(e, dest) } # `transform` method may mutate `dest`
```
### Safety

This cop is unsafe because not all objects that have an `each`
method also have a `map` method (e.g. `ENV`). Additionally, for calls
with a block, not all objects that have a `map` method return an array
(e.g. `Enumerator::Lazy`).

### Examples

```
# bad
dest = []
src.each { |e| dest << e * 2 }
dest
# good
dest = src.map { |e| e * 2 }
# bad
[].tap do |dest|
 src.each { |e| dest << e * 2 }
end
# good
dest = src.map { |e| e * 2 }
# good - contains another operation
dest = []
src.each { |e| dest << e * 2; puts e }
dest
```
## Style/MapJoin

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.85 | - |

Checks for `map { |x| x.to_s }.join` and similar calls where the
`map` is redundant because `Array#join` implicitly calls `#to_s` on
each element.

## Style/MapToHash

| Requires Ruby version 2.6 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.24 | - |

Looks for uses of `map.to_h` or `collect.to_h` that could be
written with just `to_h` in Ruby >= 2.6.

| `Style/HashTransformKeys`and`Style/HashTransformValues`will
also change this pattern if only hash keys or hash values are being
transformed. |

## Style/MapToSet

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.42 | - |

Looks for uses of `map.to_set` or `collect.to_set` that could be
written with just `to_set`.

## Style/MethodCallWithArgsParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.47 | 1.7 |

Enforces the presence (default) or absence of parentheses in method calls containing arguments.

In the default style (require_parentheses), macro methods are allowed.
Additional methods can be added to the `AllowedMethods` or
`AllowedPatterns` list. These options are valid only in the default
style. Macros can be included by either setting `IgnoreMacros` to false,
adding specific macros to the `IncludedMacros` list, or using
`IncludedMacroPatterns` for pattern-based matching.

Precedence of options is as follows:

-
`AllowedMethods`
-
`AllowedPatterns`
-
`IncludedMacros`
-
`IncludedMacroPatterns`

If a method is listed in both `IncludedMacros`/`IncludedMacroPatterns`
and `AllowedMethods`, then the latter takes precedence (that is, the
method is allowed).

In the alternative style (omit_parentheses), there are three additional options.

-
`AllowParenthesesInChaining`is`false`by default. Setting it to`true`allows the presence of parentheses in the last call during method chaining.
-
`AllowParenthesesInMultilineCall`is`false`by default. Setting it to`true`allows the presence of parentheses in multi-line method calls.
-
`AllowParenthesesInCamelCaseMethod`is`false`by default. This allows the presence of parentheses when calling a method whose name begins with a capital letter and which has no arguments. Setting it to`true`allows the presence of parentheses in such a method call even with arguments.

| The style of `omit_parentheses`allows parentheses in cases where
omitting them results in ambiguous or syntactically incorrect code. |

Non-exhaustive list of examples:

-
Parentheses are required allowed in method calls with arguments inside literals, logical operators, setting default values in position and keyword arguments, chaining and more.
-
Parentheses are allowed in method calls with arguments inside operators to avoid ambiguity. triple-dot syntax introduced in Ruby 2.7 as omitting them starts an endless range.
-
Parentheses are allowed when forwarding arguments with the triple-dot syntax introduced in Ruby 2.7 as omitting them starts an endless range.
-
Parentheses are required in calls with arguments when inside an endless method definition introduced in Ruby 3.0.
-
Ruby 3.1’s hash omission syntax allows parentheses if the method call is in conditionals and requires parentheses if the call is not the value-returning expression. See https://bugs.ruby-lang.org/issues/18396.
-
Parentheses are required in anonymous arguments, keyword arguments and block passing in Ruby 3.2.
-
Parentheses are required when the first argument is a beginless range or the last argument is an endless range.

### Examples

#### EnforcedStyle: require_parentheses (default)

```
# bad
array.delete e
# good
array.delete(e)
# good
# Operators don't need parens
foo == bar
# good
# Setter methods don't need parens
foo.bar = baz
# okay with `puts` listed in `AllowedMethods`
puts 'test'
# okay with `^assert` listed in `AllowedPatterns`
assert_equal 'test', x
```
#### EnforcedStyle: omit_parentheses

```
# bad
array.delete(e)
# good
array.delete e
# bad
action.enforce(strict: true)
# good
action.enforce strict: true
# good
# Parentheses are allowed for code that can be ambiguous without
# them.
action.enforce(condition) || other_condition
# good
# Parentheses are allowed for calls that won't produce valid Ruby
# without them.
yield path, File.basename(path)
# good
# Omitting the parentheses in Ruby 3.1 hash omission syntax can lead
# to ambiguous code. We allow them in conditionals and non-last
# expressions. See https://bugs.ruby-lang.org/issues/18396
if meets(criteria:, action:)
 safe_action(action) || dangerous_action(action)
end
```
#### AllowedMethods: ["puts", "print"]

```
# good
puts "Hello world"
print "Hello world"
# still enforces parentheses on other methods
array.delete(e)
```
#### AllowedPatterns: ["^assert"]

```
# good
assert_equal 'test', x
assert_match(/foo/, bar)
# still enforces parentheses on other methods
array.delete(e)
```
#### IncludedMacroPatterns: ["^assert", "^refute"]

```
# bad
assert_equal 'test', x
refute_nil value
# good
assert_equal('test', x)
refute_nil(value)
```
#### AllowParenthesesInMultilineCall: false (default)

```
# bad
foo.enforce(
 strict: true
)
# good
foo.enforce \
 strict: true
```
#### AllowParenthesesInMultilineCall: true

```
# good
foo.enforce(
 strict: true
)
# good
foo.enforce \
 strict: true
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| IgnoreMacros |
 | Boolean |
| AllowedMethods |
 | Array |
| AllowedPatterns |
 | Array |
| IncludedMacros |
 | Array |
| IncludedMacroPatterns |
 | Array |
| AllowParenthesesInMultilineCall |
 | Boolean |
| AllowParenthesesInChaining |
 | Boolean |
| AllowParenthesesInCamelCaseMethod |
 | Boolean |
| AllowParenthesesInStringInterpolation |
 | Boolean |
| EnforcedStyle |
 |
 |

## Style/MethodCallWithoutArgsParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.47 | 0.55 |

Checks for unwanted parentheses in parameterless method calls.

This cop’s allowed methods can be customized with `AllowedMethods`.
By default, there are no allowed methods.

| This cop allows the use of `it()`without arguments in blocks,
as in`0.times { it() }`, following`Lint/ItWithoutArgumentsInBlock`cop. |

## Style/MethodCalledOnDoEndBlock

## Style/MethodDefParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.16 | 1.7 |

Checks for parentheses around the arguments in method definitions. Both instance and class/singleton methods are checked.

Regardless of style, parentheses are necessary for:

-
Endless methods
-
Argument lists containing a `forward-arg`(`…`)
-
Argument lists containing an anonymous rest arguments forwarding ( `*`)
-
Argument lists containing an anonymous keyword rest arguments forwarding ( `**`)
-
Argument lists containing an anonymous block forwarding ( `&`)

Removing the parens would be a syntax error here.

### Examples

#### EnforcedStyle: require_parentheses (default)

```
# The `require_parentheses` style requires method definitions
# to always use parentheses
# bad
def bar num1, num2
 num1 + num2
end
def foo descriptive_var_name,
 another_descriptive_var_name,
 last_descriptive_var_name
 do_something
end
# good
def bar(num1, num2)
 num1 + num2
end
def foo(descriptive_var_name,
 another_descriptive_var_name,
 last_descriptive_var_name)
 do_something
end
```
#### EnforcedStyle: require_no_parentheses

```
# The `require_no_parentheses` style requires method definitions
# to never use parentheses
# bad
def bar(num1, num2)
 num1 + num2
end
def foo(descriptive_var_name,
 another_descriptive_var_name,
 last_descriptive_var_name)
 do_something
end
# good
def bar num1, num2
 num1 + num2
end
def foo descriptive_var_name,
 another_descriptive_var_name,
 last_descriptive_var_name
 do_something
end
```
#### EnforcedStyle: require_no_parentheses_except_multiline

```
# The `require_no_parentheses_except_multiline` style prefers no
# parentheses when method definition arguments fit on single line,
# but prefers parentheses when arguments span multiple lines.
# bad
def bar(num1, num2)
 num1 + num2
end
def foo descriptive_var_name,
 another_descriptive_var_name,
 last_descriptive_var_name
 do_something
end
# good
def bar num1, num2
 num1 + num2
end
def foo(descriptive_var_name,
 another_descriptive_var_name,
 last_descriptive_var_name)
 do_something
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/MinMaxComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.42 | - |

Enforces the use of `max` or `min` instead of comparison for greater or less.

| It can be used if you want to present limit or threshold in Ruby 2.7+.
It is slow though. So autocorrection will apply generic `max`or`min`: |

```
a.clamp(b..) # Same as `[a, b].max`
a.clamp(..b) # Same as `[a, b].min`
```
## Style/MissingElse

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.30 | 0.38 |

Checks for `if` expressions that do not have an `else` branch.

| Pattern matching is allowed to have no `else`branch because unlike`if`and`case`,
it raises`NoMatchingPatternError`if the pattern doesn’t match and without having`else`. |

Supported styles are: if, case, both.

### Examples

#### EnforcedStyle: both (default)

```
# warn when an `if` or `case` expression is missing an `else` branch.
# bad
if condition
 statement
end
# bad
case var
when condition
 statement
end
# good
if condition
 statement
else
 # the content of `else` branch will be determined by Style/EmptyElse
end
# good
case var
when condition
 statement
else
 # the content of `else` branch will be determined by Style/EmptyElse
end
```
#### EnforcedStyle: if

```
# warn when an `if` expression is missing an `else` branch.
# bad
if condition
 statement
end
# good
if condition
 statement
else
 # the content of `else` branch will be determined by Style/EmptyElse
end
# good
case var
when condition
 statement
end
# good
case var
when condition
 statement
else
 # the content of `else` branch will be determined by Style/EmptyElse
end
```
#### EnforcedStyle: case

```
# warn when a `case` expression is missing an `else` branch.
# bad
case var
when condition
 statement
end
# good
case var
when condition
 statement
else
 # the content of `else` branch will be determined by Style/EmptyElse
end
# good
if condition
 statement
end
# good
if condition
 statement
else
 # the content of `else` branch will be determined by Style/EmptyElse
end
```
## Style/MissingRespondToMissing

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.56 | 1.89 |

Checks for the presence of `method_missing` without also
defining `respond_to_missing?`.

Not defining `respond_to_missing?` will cause metaprogramming
methods like `respond_to?` to behave unexpectedly:

```
class StringDelegator
 def initialize(string)
 @string = string
 end
 def method_missing(name, *args)
 @string.send(name, *args)
 end
end
delegator = StringDelegator.new("foo")
# Claims to not respond to `upcase`.
delegator.respond_to?(:upcase) # => false
# But you can call it.
delegator.upcase # => FOO
```
When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, `respond_to_missing?` defined in another definition of the
same class or module (e.g. a reopening in another file) also
satisfies the check.

### Examples

```
# bad
def method_missing(name, *args)
 if @delegate.respond_to?(name)
 @delegate.send(name, *args)
 else
 super
 end
end
# good
def respond_to_missing?(name, include_private)
 @delegate.respond_to?(name) || super
end
def method_missing(name, *args)
 if @delegate.respond_to?(name)
 @delegate.send(name, *args)
 else
 super
 end
end
```
## Style/MixinGrouping

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.48 | 0.49 |

Checks for grouping of mixins in `class` and `module` bodies.
By default it enforces mixins to be placed in separate declarations,
but it can be configured to enforce grouping them in one declaration.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/MixinUsage

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.51 | - |

Checks that `include`, `extend` and `prepend` statements appear
inside classes and modules, not at the top level, so as to not affect
the behavior of `Object`.

## Style/ModuleFunction

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.11 | 0.65 |

Checks for use of `extend self` or `module_function` in a module.

Supported styles are: `module_function` (default), `extend_self` and `forbidden`.

A couple of things to keep in mind:

-
`forbidden`style prohibits the usage of both styles
-
in default mode ( `module_function`), the cop won’t be activated when the module contains any private methods

### Safety

Autocorrection is unsafe (and is disabled by default) because `extend self`
and `module_function` do not behave exactly the same.

### Examples

#### EnforcedStyle: module_function (default)

```
# bad
module Test
 extend self
 # ...
end
# good
module Test
 module_function
 # ...
end
# good
module Test
 extend self
 # ...
 private
 # ...
end
# good
module Test
 class << self
 # ...
 end
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| Autocorrect |
 | Boolean |

## Style/ModuleMemberExistenceCheck

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.82 | - |

Checks for usage of `Module` methods returning arrays that can be replaced
with equivalent predicates.

Calling a method returning an array then checking if an element is inside
it is much slower than using an equivalent predicate method. For example,
`instance_methods.include?` will return an array of all public and protected
instance methods in the module, then check if a given method is inside that
array, while `method_defined?` will do direct method lookup, which is much
faster and consumes less memory.

| `constants.include?`is not handled by this cop because`Module#const_defined?`has different lookup behavior than`Module#constants`-`const_defined?`searches up to`Object`(top-level constants like`String`,`Integer`, etc.) while`constants`does not, which can cause behavior changes after autocorrection. |

### Examples

```
# bad
Array.instance_methods.include?(:size)
Array.instance_methods.member?(:size)
Array.instance_methods(true).include?(:size)
Array.instance_methods(false).include?(:find)
# good
Array.method_defined?(:size)
Array.method_defined?(:find, false)
# bad
Array.class_variables.include?(:foo)
Array.private_instance_methods.include?(:foo)
Array.protected_instance_methods.include?(:foo)
Array.public_instance_methods.include?(:foo)
# good
Array.class_variable_defined?(:foo)
Array.private_method_defined?(:foo)
Array.protected_method_defined?(:foo)
Array.public_method_defined?(:foo)
```
## Style/MultilineBlockChain

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | No | 0.13 | - |

Checks for chaining of a block after another block that spans multiple lines.

### Examples

```
# bad
Thread.list.select do |t|
 t.alive?
end.map do |t|
 t.object_id
end
# good
alive_threads = Thread.list.select do |t|
 t.alive?
end
alive_threads.map do |t|
 t.object_id
end
```
## Style/MultilineIfModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.45 | - |

Checks for uses of if/unless modifiers with multiple-lines bodies.

## Style/MultilineIfThen

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.26 |

Checks for uses of the `then` keyword in multi-line `if` statements.
In multi-line `if` statements, `then` is redundant because the newline
already separates the condition from the body.

### Examples

```
# bad
# This is considered bad practice.
if cond then
end
# good
# If statements can contain `then` on the same line.
if cond then a
elsif cond then b
end
```
## Style/MultilineInPatternThen

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.16 | - |

Checks uses of the `then` keyword in multi-line `in` statement.

### Examples

```
# bad
case expression
in pattern then
end
# good
case expression
in pattern
end
# good
case expression
in pattern then do_something
end
# good
case expression
in pattern then do_something(arg1,
 arg2)
end
```
## Style/MultilineMemoization

## Style/MultilineTernaryOperator

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.86 |

Checks for multi-line ternary op expressions.

| `return if … else … end`is syntax error. If`return`is used before
multiline ternary operator expression, it will be autocorrected to single-line
ternary operator. The same is true for`break`,`next`, and method call. |

### Examples

```
# bad
a = cond ?
 b : c
a = cond ? b :
 c
a = cond ?
 b :
 c
return cond ?
 b :
 c
# good
a = cond ? b : c
a = if cond
 b
else
 c
end
return cond ? b : c
```
## Style/MultilineWhenThen

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.73 | - |

Checks uses of the `then` keyword
in multi-line when statements.

### Examples

```
# bad
case foo
when bar then
end
# good
case foo
when bar
end
# good
case foo
when bar then do_something
end
# good
case foo
when bar then do_something(arg1,
 arg2)
end
```
## Style/MultipleComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.49 | 1.1 |

Checks against comparing a variable with multiple items, where
`Array#include?`, `Set#include?` or a `case` could be used instead
to avoid code repetition.
It accepts comparisons of multiple method calls to avoid unnecessary method calls
by default. It can be configured by `AllowMethodComparison` option.

### Examples

```
# bad
a = 'a'
foo if a == 'a' || a == 'b' || a == 'c'
# good
a = 'a'
foo if ['a', 'b', 'c'].include?(a)
VALUES = Set['a', 'b', 'c'].freeze
# elsewhere...
foo if VALUES.include?(a)
case foo
when 'a', 'b', 'c' then foo
# ...
end
# accepted (but consider `case` as above)
foo if a == b.lightweight || a == b.heavyweight
```
## Style/MutableConstant

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.34 | 1.88 |

Checks whether some constant value isn’t a mutable literal (e.g. array or hash).

When the `Recursive` option is enabled, mutable literals nested inside
arrays and hashes are also frozen, so an offense on the outermost
unfrozen literal will autocorrect every nested mutable literal as well.
When the outer literal already has `.freeze` appended, the cop descends
into it and reports each outermost unfrozen literal underneath. The
option is disabled by default to preserve existing behavior; opt in to
get strict nested freezing.

Strict mode can be used to freeze all constants, rather than just literals. Strict mode is considered an experimental feature. It has not been updated with an exhaustive list of all methods that will produce frozen objects so there is a decent chance of getting some false positives. Luckily, there is no harm in freezing an already frozen object.

From Ruby 3.0, this cop honours the magic comment 'shareable_constant_value'. When this magic comment is set to any acceptable value other than none, it will suppress the offenses raised by this cop. It enforces frozen state.

| `Regexp`and`Range`literals are frozen objects since Ruby 3.0. |

| From Ruby 3.0, interpolated strings are not frozen when `# frozen-string-literal: true`is used, so this cop enforces explicit
freezing for such strings. |

| From Ruby 3.0, this cop allows explicit freezing of constants when
the `shareable_constant_value`directive is used. |

### Safety

This cop’s autocorrection is unsafe since any mutations on objects that
are made frozen will change from being accepted to raising `FrozenError`,
and will need to be manually refactored.

### Examples

#### EnforcedStyle: literals (default)

```
# bad
CONST = [1, 2, 3]
# good
CONST = [1, 2, 3].freeze
# good
CONST = <<~TESTING.freeze
 This is a heredoc
TESTING
# good
CONST = Something.new
```
#### Recursive: false (default)

```
# good - only the outer container needs to be frozen
CONST = [{ a: [], b: 'foo' }].freeze
```
#### Recursive: true

```
# bad - nested mutable literals must be frozen too
CONST = [{ a: [], b: 'foo' }].freeze
# good
CONST = [{ a: [].freeze, b: 'foo'.freeze }.freeze].freeze
```
#### EnforcedStyle: strict

```
# bad
CONST = Something.new
# bad
CONST = Struct.new do
 def foo
 puts 1
 end
end
# good
CONST = Something.new.freeze
# good
CONST = Struct.new do
 def foo
 puts 1
 end
end.freeze
# good - `Data.define` declares an immutable value type
CONST = Data.define(:foo, :bar)
```
```
# Magic comment - shareable_constant_value: literal
# bad
CONST = [1, 2, 3]
# good
# shareable_constant_value: literal
CONST = [1, 2, 3]
```
## Style/NegatedIf

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.20 | 0.48 |

Checks for uses of if with a negated condition. Only ifs without else are considered. There are three different styles:

-
both
-
prefix
-
postfix

### Examples

#### EnforcedStyle: both (default)

```
# enforces `unless` for `prefix` and `postfix` conditionals
# bad
if !foo
 bar
end
# good
unless foo
 bar
end
# bad
bar if !foo
# good
bar unless foo
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/NegatedIfElseCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.2 | - |

Checks for uses of `if-else` and ternary operators with a negated condition
which can be simplified by inverting condition and swapping branches.

## Style/NegatedUnless

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.69 | - |

Checks for uses of unless with a negated condition. Only unless without else are considered. There are three different styles:

-
both
-
prefix
-
postfix

### Examples

#### EnforcedStyle: both (default)

```
# enforces `if` for `prefix` and `postfix` conditionals
# bad
unless !foo
 bar
end
# good
if foo
 bar
end
# bad
bar unless !foo
# good
bar if foo
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/NegatedWhile

## Style/NegativeArrayIndex

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.84 | - |

Identifies usages of `arr[arr.length - n]`, `arr[arr.size - n]`, or
`arr[arr.count - n]` and suggests to change them to use `arr[-n]` instead.
Also handles range patterns like `arr[0..(arr.length - n)]`.

The cop recognizes preserving methods (`sort`, `reverse`, `shuffle`, `rotate`)
and their combinations, allowing safe replacement when the receiver matches.
It works with variables, instance variables, class variables, and constants.

## Style/NestedFileDirname

| Requires Ruby version 3.1 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.26 | - |

Checks for nested `File.dirname`.
It replaces nested `File.dirname` with the level argument introduced in Ruby 3.1.

## Style/NestedModifier

## Style/NestedParenthesizedCalls

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.36 | 0.77 |

Checks for unparenthesized method calls in the argument list
of a parenthesized method call.
`be`, `be_a`, `be_an`, `be_between`, `be_falsey`, `be_kind_of`, `be_instance_of`,
`be_truthy`, `be_within`, `eq`, `eql`, `end_with`, `include`, `match`, `raise_error`,
`respond_to`, and `start_with` methods are allowed by default.
These are customizable with `AllowedMethods` option.

## Style/NestedTernaryOperator

## Style/Next

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.22 | 1.75 |

Use `next` to skip iteration instead of a condition at the end.

### Examples

#### EnforcedStyle: skip_modifier_ifs (default)

```
# bad
[1, 2].each do |a|
 if a == 1
 puts a
 end
end
# good
[1, 2].each do |a|
 next unless a == 1
 puts a
end
# good
[1, 2].each do |a|
 puts a if a == 1
end
```
#### EnforcedStyle: always

```
# With `always` all conditions at the end of an iteration needs to be
# replaced by next - with `skip_modifier_ifs` the modifier if like
# this one are ignored: `[1, 2].each { |a| puts a if a == 1 }`
# bad
[1, 2].each do |a|
 puts a if a == 1
end
# bad
[1, 2].each do |a|
 if a == 1
 puts a
 end
end
# good
[1, 2].each do |a|
 next unless a == 1
 puts a
end
```
## Style/NilComparison

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.12 | 0.59 |

Checks for comparison of something with nil using `==` and
`nil?`. Enforcing a consistent style (either the `nil?`
predicate or `==` comparison) improves readability.

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |

## Style/NilLambda

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.3 | 1.15 |

Checks for lambdas and procs that always return nil, which can be replaced with an empty lambda or proc instead.

| A `proc`that returns nil via an explicit`return`is allowed,
because in a`proc``return`exits the enclosing method, so removing it
would change behavior. A lambda is still reported, since there`return`only exits the lambda itself. |

## Style/NonNilCheck

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.20 | 0.22 |

Checks for non-nil checks, which are usually redundant.

With `IncludeSemanticChanges` set to `false` by default, this cop
does not report offenses for `!x.nil?` and does no changes that might
change behavior.
Also `IncludeSemanticChanges` set to `false` with `EnforcedStyle: comparison` of
`Style/NilComparison` cop, this cop does not report offenses for `x != nil` and
does no changes to `!x.nil?` style.

With `IncludeSemanticChanges` set to `true`, this cop reports offenses
for `!x.nil?` and autocorrects that and `x != nil` to solely `x`, which
is **usually** OK, but might change behavior.

### Examples

```
# bad
if x != nil
end
# good
if x
end
# Non-nil checks are allowed if they are the final nodes of predicate.
# good
def signed_in?
 !current_user.nil?
end
```
## Style/Not

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.20 |

Checks for uses of the keyword `not` instead of `!`.
The `not` keyword has lower precedence than `!`, which can
lead to surprising behavior and often requires parentheses.

### Examples

```
# bad - parentheses are required because of op precedence
x = (not something)
# good
x = !something
```
## Style/NumberedParameters

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.22 | - |

Checks for numbered parameters.

It can either restrict the use of numbered parameters to single-lined blocks, or disallow completely numbered parameters.

## Style/NumberedParametersLimit

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.22 | - |

Detects use of an excessive amount of numbered parameters in a single block. Having too many numbered parameters can make code too cryptic and hard to read.

The cop defaults to registering an offense if there is more than 1 numbered
parameter but this maximum can be configured by setting `Max`.

## Style/NumericLiteralPrefix

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.41 | - |

Checks for octal, hex, binary, and decimal literals using uppercase prefixes and corrects them to lowercase prefix or no prefix (in case of decimals).

### Examples

## Style/NumericLiterals

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.48 |

Checks for big numeric literals without `_` between groups
of digits in them. Underscores make large numbers easier to
read by visually separating groups of digits.

Additional allowed patterns can be added by adding regexps to
the `AllowedPatterns` configuration. All regexps are treated
as anchored even if the patterns do not contain anchors (so
`\d{4}_\d{4}` will allow `1234_5678` but not `1234_5678_9012`).

| Even if `AllowedPatterns`are given, autocorrection will
only correct to the standard pattern of an`_`every 3 digits. |

## Style/NumericPredicate

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.42 | 0.59 |

Checks for usage of comparison operators (`==`,
`>`, `<`) to test numbers as zero, positive, or negative.
These can be replaced by their respective predicate methods.
This cop can also be configured to do the reverse.

This cop’s allowed methods can be customized with `AllowedMethods`.
By default, there are no allowed methods.

This cop disregards `#nonzero?` as its value is truthy or falsey,
but not `true` and `false`, and thus not always interchangeable with
`!= 0`.

This cop allows comparisons to global variables, since they are often
populated with objects which can be compared with integers, but are
not themselves `Integer` polymorphic.

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver defines the predicates or can be compared to a number, which may lead to a false positive for non-standard classes.

### Examples

#### EnforcedStyle: predicate (default)

```
# bad
foo == 0
0 > foo
bar.baz > 0
# good
foo.zero?
foo.negative?
bar.baz.positive?
```
#### EnforcedStyle: comparison

```
# bad
foo.zero?
foo.negative?
bar.baz.positive?
# good
foo == 0
0 > foo
bar.baz > 0
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| AllowedMethods |
 | Array |
| AllowedPatterns |
 | Array |
| Exclude |
 | Array |

## Style/ObjectThen

| Requires Ruby version 2.6 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.28 | - |

Enforces the use of consistent method names
`Object#yield_self` or `Object#then`.

## Style/OneClassPerFile

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.85 | 1.86 |

Checks that each source file defines at most one top-level class or module.

Keeping one class or module per file makes it easier to find and navigate code, and follows the convention used by most Ruby projects.

Classes and modules listed in `AllowedClasses` are not counted toward the
limit. This is useful for small ancillary classes like custom exception
classes that logically belong with the main class.

### Examples

```
# bad - Multiple top-level classes
class Foo
end
class Bar
end
# bad - Multiple top-level modules
module Foo
end
module Bar
end
# bad - A top-level class and a top-level module
class Foo
end
module Bar
end
# good - A single top-level class
class Foo
end
# good - A single top-level module
module Foo
end
# good - Nested classes within a single top-level class
class Foo
 class Bar
 end
end
# good - Multiple classes within a single top-level module
module Foo
 class Bar
 end
 class Baz
 end
end
```
## Style/OneLineConditional

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.90 |

Checks for uses of `if/then/else/end` constructs on a single line.
A ternary operator (`?:`) or multi-line `if` is more readable.
`AlwaysCorrectToMultiline` config option can be set to `true` to autocorrect all offenses to
multi-line constructs. When `AlwaysCorrectToMultiline` is `false` (default case) the
autocorrect will first try converting them to ternary operators.

### Examples

```
# bad
if foo then bar else baz end
# bad
unless foo then baz else bar end
# good
foo ? bar : baz
# good
bar if foo
# good
if foo then bar end
# good
if foo
 bar
else
 baz
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AlwaysCorrectToMultiline |
 | Boolean |

## Style/OpenStructUse

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | No | 1.23 | 1.51 |

Flags uses of `OpenStruct`, as it is now officially discouraged
to be used for performance, version compatibility, and potential security issues.

### Safety

Note that this cop may flag false positives; for instance, the following legal
use of a hand-rolled `OpenStruct` type would be considered an offense:

```
module MyNamespace
 class OpenStruct # not the OpenStruct we're looking for
 end
 def new_struct
 OpenStruct.new # resolves to MyNamespace::OpenStruct
 end
end
```
## Style/OperatorMethodCall

## Style/OptionHash

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 0.33 | 0.34 |

Checks for options hashes and discourages them if the current Ruby version supports keyword arguments.

### Examples

```
# bad
def fry(options = {})
 temperature = options.fetch(:temperature, 300)
 # ...
end
# good
def fry(temperature: 300)
 # ...
end
```
## Style/OptionalArguments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | No | 0.33 | 0.83 |

Checks for optional arguments to methods that do not come at the end of the argument list.

### Examples

```
# bad
def foo(a = 1, b, c)
end
# good
def baz(a, b, c = 1)
end
def foobar(a = 1, b = 2, c = 3)
end
```
## Style/OptionalBooleanParameter

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | No | 0.89 | - |

Checks for places where keyword arguments can be used instead of
boolean arguments when defining methods. `respond_to_missing?` method is allowed by default.
These are customizable with `AllowedMethods` option.

### Examples

```
# bad
def some_method(bar = false)
 puts bar
end
# bad - common hack before keyword args were introduced
def some_method(options = {})
 bar = options.fetch(:bar, false)
 puts bar
end
# good
def some_method(bar: false)
 puts bar
end
```
## Style/OrAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.50 | - |

Checks for potential usage of the `||=` operator.

## Style/ParallelAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.32 | - |

Checks for simple usages of parallel assignment. Parallel assignment is less readable than individual assignments and makes it harder to follow what each variable is being set to.

This will only complain when the number of variables being assigned matched the number of assigning variables.

### Examples

```
# bad
a, b, c = 1, 2, 3
a, b, c = [1, 2, 3]
# good
one, two = *foo
a, b = foo
a, b = b, a
a = 1
b = 2
c = 3
```
## Style/ParenthesesAroundCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.56 |

Checks for the presence of superfluous parentheses around the condition of if/unless/while/until.

`AllowSafeAssignment` option for safe assignment.
By safe assignment we mean putting parentheses around
an assignment to indicate "I know I’m using an assignment
as a condition. It’s not a mistake."

### Examples

```
# bad
x += 1 while (x < 10)
foo unless (bar || baz)
if (x > 10)
elsif (x < 3)
end
# good
x += 1 while x < 10
foo unless bar || baz
if x > 10
elsif x < 3
end
```
## Style/PartitionInsteadOfDoubleSelect

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.85 | - |

Checks for consecutive calls to `select`/`filter`/`find_all` and `reject`
on the same receiver with the same block body, where `partition` could be
used instead. Also detects two `select` or two `reject` calls where one
block negates the other with `!`. Using `partition` reduces two collection
traversals to one.

### Safety

This cop is unsafe because:

-
`Hash#select`and`Hash#reject`return hashes, but`Hash#partition`returns nested arrays.
-
When the receiver has side effects, calling it once (with `partition`) versus twice (with`select`+`reject`) may produce different results.
-
Custom classes may override `select`/`reject`without providing a compatible`partition`method.

### Examples

```
# bad
positives = array.select { |x| x > 0 }
negatives = array.reject { |x| x > 0 }
# bad
positives = array.filter { |x| x > 0 }
negatives = array.reject { |x| x > 0 }
# bad
negatives = array.reject { |x| x > 0 }
positives = array.select { |x| x > 0 }
# bad
positives = array.select(&:positive?)
negatives = array.reject(&:positive?)
# bad
positives = array.select(&:positive?)
negatives = array.reject { |x| x.positive? }
# bad
positives = array.select { |x| x.positive? }
non_positives = array.select { |x| !x.positive? }
# good
positives, negatives = array.partition { |x| x > 0 }
# good
positives, non_positives = array.partition { |x| x.positive? }
# good
positives, negatives = array.partition(&:positive?)
```
## Style/PercentLiteralDelimiters

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.19 | 0.48 |

Enforces the consistent usage of `%`-literal delimiters.
Using consistent delimiters across the codebase reduces
cognitive load when reading `%`-literals.

Specify the 'default' key to set all preferred delimiters at once. You can continue to specify individual preferred delimiters to override the default.

### Examples

```
# Style/PercentLiteralDelimiters:
# PreferredDelimiters:
# default: '[]'
# '%i': '()'
# good
%w[alpha beta] + %i(gamma delta)
# bad
%W(alpha #{beta})
# bad
%I(alpha beta)
```
## Style/PercentQLiterals

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.25 | - |

Checks for usage of the %Q() syntax when %q() would do.

## Style/PerlBackrefs

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.13 | - |

Looks for uses of Perl-style regexp match backreferences and their English versions like $1, $2, $&, $MATCH, $PREMATCH, etc.

## Style/PredicateWithKind

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.85 | - |

Looks for uses of `any?`, `all?`, `none?`, or `one?` with a block
containing only an `is_a?`, `kind_of?`, or `instance_of?` check, and
suggests using the predicate method with the class argument directly.

## Style/PreferredHashMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.41 | 0.70 |

Checks for uses of methods `Hash#has_key?` and
`Hash#has_value?`, and suggests using `Hash#key?` and `Hash#value?` instead.

It is configurable to enforce the verbose method names, by using the
`EnforcedStyle: verbose` configuration.

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
is a `Hash` or responds to the replacement methods.

## Style/Proc

## Style/QuotedSymbols

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.16 | - |

Checks if the quotes used for quoted symbols match the configured defaults.
By default uses the same configuration as `Style/StringLiterals`; if that
cop is not enabled, the default `EnforcedStyle` is `single_quotes`.

String interpolation is always kept in double quotes.

| `Lint/SymbolConversion`can be used in parallel to ensure that symbols
are not quoted that don’t need to be. This cop is for configuring the quoting
style to use for symbols that require quotes. |

## Style/RaiseArgs

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.14 | 1.61 |

Checks the args passed to `fail` and `raise`.

Exploded style (default) enforces passing the exception class and message arguments separately, rather than constructing an instance of the error.

Compact style enforces constructing an error instance.

Both styles allow passing just a message, or an error instance when there is more than one argument.

The exploded style has an `AllowedCompactTypes` configuration
option that takes an `Array` of exception name Strings.

### Examples

#### EnforcedStyle: exploded (default)

```
# bad
raise StandardError.new('message')
# good
raise StandardError, 'message'
fail 'message'
raise MyCustomError
raise MyCustomError.new(arg1, arg2, arg3)
raise MyKwArgError.new(key1: val1, key2: val2)
# With `AllowedCompactTypes` set to ['MyWrappedError']
raise MyWrappedError.new(obj)
raise MyWrappedError.new(obj), 'message'
```
## Style/RandomWithOffset

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for the use of randomly generated numbers, added/subtracted with integer literals, as well as those with Integer#succ and Integer#pred methods. Prefer using ranges instead, as it clearly states the intentions.

### Examples

```
# bad
rand(6) + 1
1 + rand(6)
rand(6) - 1
1 - rand(6)
rand(6).succ
rand(6).pred
Random.rand(6) + 1
Kernel.rand(6) + 1
rand(0..5) + 1
# good
rand(1..6)
rand(1...7)
```
## Style/ReduceToHash

| Requires Ruby version 2.6 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.85 | - |

Checks for `each_with_object`, `inject`, and `reduce` calls that build
a hash from an enumerable, where `to_h` with a block could be used instead.

This cop complements `Style/HashTransformKeys` and `Style/HashTransformValues`,
which handle hash-to-hash transformations with destructured key-value pairs.
This cop targets the case where a hash is built from individual elements
(non-destructured block parameter).

### Safety

This cop is unsafe because it cannot guarantee that the receiver
is an `Enumerable` by static analysis, so the correction may
not be actually equivalent. Additionally, `each_with_object` returns
the hash object while `to_h` returns a new hash, which could matter
if the hash object identity is important.

### Examples

```
# bad
array.each_with_object({}) { |elem, hash| hash[elem.id] = elem.name }
# bad
array.inject({}) { |hash, elem| hash[elem.id] = elem.name; hash }
# bad
array.reduce({}) { |hash, elem| hash[elem.id] = elem.name; hash }
# bad
array.each_with_object({}) { |elem, hash| hash[elem] = elem.to_s }
# good
array.to_h { |elem| [elem.id, elem.name] }
# good
array.to_h { |elem| [elem, elem.to_s] }
```
## Style/RedundantArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.4 | 1.55 |

Checks for a redundant argument passed to certain methods.

| This cop is limited to methods with single parameter. |

Method names and their redundant arguments can be configured like this:

```
Methods:
 join: ''
 sum: 0
 split: ' '
 chomp: "\n"
 chomp!: "\n"
 foo: 2
```
### Safety

This cop is unsafe because of the following limitations:

-
This cop matches by method names only and hence cannot tell apart methods with same name in different classes.
-
This cop may be unsafe if certain special global variables (e.g. `$;`,`$/`) are set. That depends on the nature of the target methods, of course. For example, the default argument to join is`$OUTPUT_FIELD_SEPARATOR`(or`$,`) rather than`''`, and if that global is changed,`''`is no longer a redundant argument.

### Examples

```
# bad
array.join('')
[1, 2, 3].join("")
array.sum(0)
exit(true)
exit!(false)
string.to_i(10)
string.split(" ")
"first\nsecond".split(" ")
string.chomp("\n")
string.chomp!("\n")
A.foo(2)
# good
array.join
[1, 2, 3].join
array.sum
exit
exit!
string.to_i
string.split
"first second".split
string.chomp
string.chomp!
A.foo
```
## Style/RedundantArrayConstructor

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.52 | - |

Checks for the instantiation of an array using a redundant `Array` constructor.
Autocorrect replaces it with an array literal which is the simplest and fastest.

## Style/RedundantArrayFlatten

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.76 | - |

Checks for redundant calls of `Array#flatten`.

`Array#join` joins nested arrays recursively, so flattening an array
beforehand is redundant.

## Style/RedundantAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.87 | - |

Checks for redundant assignment before returning.

When there are comments between the assignment and reference, the cop will report an offense but it will not autocorrect.

## Style/RedundantBegin

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.10 | 0.21 |

Checks for redundant `begin` blocks. A `begin` block is redundant
when the `rescue`/`ensure` can be handled by the enclosing method
or block definition directly, avoiding unnecessary indentation.

### Examples

```
# bad
def redundant
 begin
 ala
 bala
 rescue StandardError => e
 something
 end
end
# good
def preferred
 ala
 bala
rescue StandardError => e
 something
end
# bad
begin
 do_something
end
# good
do_something
# bad
# When using Ruby 2.5 or later.
do_something do
 begin
 something
 rescue => ex
 anything
 end
end
# good
# In Ruby 2.5 or later, you can omit `begin` in `do-end` block.
do_something do
 something
rescue => ex
 anything
end
# good
# Stabby lambdas don't support implicit `begin` in `do-end` blocks.
-> do
 begin
 foo
 rescue Bar
 baz
 end
end
```
## Style/RedundantCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.76 | 1.73 |

Checks for unnecessary conditional expressions.

| Since the intention of the comment cannot be automatically determined, autocorrection is not applied when a comment is used, as shown below: |

```
if b
 # Important note.
 b
else
 c
end
```
## Style/RedundantConstantBase

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.40 | 1.89 |

Avoid redundant `::` prefix on a constant.

How Ruby searches constants is a bit complicated, and it can often be difficult to
understand from the code whether the `::` is intended or not. Where `Module.nesting`
is empty, there is no need to prepend `::`, so it would be nice to consistently
avoid such meaningless `::` prefix to avoid confusion.

| This cop is disabled if `Lint/ConstantResolution`cop is enabled,
to prevent conflicting rules. This is because it respects user configurations
that want to enable`Lint/ConstantResolution`cop which is disabled by default. |

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, a leading `::` inside a namespace is also reported when the
constant provably resolves to the same declaration with and without
the base, i.e. nothing in the surrounding nesting shadows it.

## Style/RedundantCurrentDirectoryInPath

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.53 | - |

Checks for paths given to `require_relative` that start with
the current directory (`./`), which can be omitted.

## Style/RedundantDoubleSplatHashBraces

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.41 | - |

Checks for redundant uses of double splat hash braces.

## Style/RedundantEach

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.38 | - |

Checks for redundant `each`.

### Examples

```
# bad
array.each.each { |v| do_something(v) }
# good
array.each { |v| do_something(v) }
# bad
array.each.each_with_index { |v, i| do_something(v, i) }
# good
array.each.with_index { |v, i| do_something(v, i) }
array.each_with_index { |v, i| do_something(v, i) }
# bad
array.each.each_with_object { |v, o| do_something(v, o) }
# good
array.each.with_object { |v, o| do_something(v, o) }
array.each_with_object { |v, o| do_something(v, o) }
```
## Style/RedundantException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.14 | 0.29 |

Checks for `RuntimeError` as the argument of `raise`/`fail`.

## Style/RedundantFetchBlock

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.86 | - |

Identifies places where `fetch(key) { value }` can be replaced by `fetch(key, value)`.

In such cases `fetch(key, value)` method is faster than `fetch(key) { value }`.

| The block string `'value'`in`hash.fetch(:key) { 'value' }`is detected
when frozen string literal magic comment is enabled (i.e.`# frozen_string_literal: true`),
but not when disabled. |

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
does not have a different implementation of `fetch`.

## Style/RedundantFileExtensionInRequire

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.88 | - |

Checks for the presence of superfluous `.rb` extension in
the filename provided to `require` and `require_relative`.

| If the extension is omitted, Ruby tries adding '.rb', '.so',
 and so on to the name until found. If the file named cannot be found,
 a `LoadError`will be raised.
 There is an edge case where`foo.so`file is loaded instead of a`LoadError`if`foo.so`file exists when`require 'foo.rb'`will be changed to`require 'foo'`,
 but that seems harmless. |

## Style/RedundantFilterChain

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.52 | 1.57 |

Identifies usages of `any?`, `empty?` or `none?` predicate methods
chained to `select`/`filter`/`find_all` and change them to use predicate method instead.

### Safety

This cop’s autocorrection is unsafe because `array.select.any?` evaluates all elements
through the `select` method, while `array.any?` uses short-circuit evaluation.
In other words, `array.select.any?` guarantees the evaluation of every element,
but `array.any?` does not necessarily evaluate all of them.

### Examples

```
# bad
arr.select { |x| x > 1 }.any?
# good
arr.any? { |x| x > 1 }
# bad
arr.select { |x| x > 1 }.empty?
arr.select { |x| x > 1 }.none?
# good
arr.none? { |x| x > 1 }
# good
relation.select(:name).any?
arr.select { |x| x > 1 }.any?(&:odd?)
```
## Style/RedundantFormat

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.72 | 1.72 |

Checks for calls to `Kernel#format` or `Kernel#sprintf` that are redundant.

Calling `format` with only a single string or constant argument is redundant,
as it can be replaced by the string or constant itself.

Also looks for `format` calls where the arguments are literals that can be
inlined into a string easily. This applies to the `%s`, `%d`, `%i`, `%u`, and
`%f` format specifiers.

### Safety

This cop’s autocorrection is unsafe because string object returned by
`format` and `sprintf` are never frozen. If `format('string')` is autocorrected to
`'string'`, `FrozenError` may occur when calling a destructive method like `String#<<`.
Consider using `'string'.dup` instead of `format('string')`.
Additionally, since the necessity of `dup` cannot be determined automatically,
this autocorrection is inherently unsafe.

```
# frozen_string_literal: true
format('template').frozen? # => false
'template'.frozen? # => true
```
### Examples

```
# bad
format('the quick brown fox jumps over the lazy dog.')
sprintf('the quick brown fox jumps over the lazy dog.')
# good
'the quick brown fox jumps over the lazy dog.'
# bad
format(MESSAGE)
sprintf(MESSAGE)
# good
MESSAGE
# bad
format('%s %s', 'foo', 'bar')
sprintf('%s %s', 'foo', 'bar')
# good
'foo bar'
```
## Style/RedundantFreeze

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.34 | 0.66 |

Checks for uses of `Object#freeze` on immutable objects.

| `Regexp`and`Range`literals are frozen objects since Ruby 3.0. |

| From Ruby 3.0, this cop allows explicit freezing of interpolated
string literals when `# frozen-string-literal: true`is used. |

## Style/RedundantHeredocDelimiterQuotes

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.45 | - |

Checks for redundant heredoc delimiter quotes.

## Style/RedundantInitialize

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Command-line only (Unsafe) | 1.27 | 1.61 |

Checks for `initialize` methods that are redundant.

An initializer is redundant if it does not do anything, or if it only
calls `super` with the same arguments given to it. If the initializer takes
an argument that accepts multiple values (`restarg`, `kwrestarg`, etc.) it
will not register an offense, because it allows the initializer to take a different
number of arguments as its superclass potentially does.

| If an initializer takes any arguments and has an empty body, RuboCop
assumes it to notbe redundant. This is to prevent potential`ArgumentError`. |

| If an initializer argument has a default value, RuboCop assumes it
to notbe redundant. |

| Empty initializers are registered as offenses, but it is possible
to purposely create an empty `initialize`method to override a superclass’s
initializer. |

### Safety

This cop is unsafe because removing an empty initializer may alter the behavior of the code, particularly if the superclass initializer raises an exception. In such cases, the empty initializer may act as a safeguard to prevent unintended errors from propagating.

### Examples

```
# bad
def initialize
end
# bad
def initialize
 super
end
# bad
def initialize(a, b)
 super
end
# bad
def initialize(a, b)
 super(a, b)
end
# good
def initialize
 do_something
end
# good
def initialize
 do_something
 super
end
# good (different number of parameters)
def initialize(a, b)
 super(a)
end
# good (default value)
def initialize(a, b = 5)
 super
end
# good (default value)
def initialize(a, b: 5)
 super
end
# good (changes the parameter requirements)
def initialize(_)
end
# good (changes the parameter requirements)
def initialize(*)
end
# good (changes the parameter requirements)
def initialize(**)
end
# good (changes the parameter requirements)
def initialize(...)
end
```
## Style/RedundantInterpolation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.76 | 1.30 |

Checks for strings that are just an interpolated expression.

### Safety

Autocorrection is unsafe because when calling a destructive method to string,
the resulting string may have different behavior or raise `FrozenError`.

```
x = 'a'
y = "#{x}"
y << 'b' # return 'ab'
x # return 'a'
y = x.to_s
y << 'b' # return 'ab'
x # return 'ab'
x = 'a'.freeze
y = "#{x}"
y << 'b' # return 'ab'.
y = x.to_s
y << 'b' # raise `FrozenError`.
```
## Style/RedundantInterpolationUnfreeze

| Requires Ruby version 3.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.66 | - |

Before Ruby 3.0, interpolated strings followed the frozen string literal magic comment which sometimes made it necessary to explicitly unfreeze them. Ruby 3.0 changed interpolated strings to always be unfrozen which makes unfreezing them redundant.

## Style/RedundantLineContinuation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.49 | - |

Checks for redundant line continuation.

A line continuation is redundant when removing the backslash does not change how the program parses: the source is reparsed without the backslash and the resulting AST is compared to the original. Only backslashes that are pure noise are reported; backslashes that are significant — inside strings, for string concatenation, before an operator or argument that would otherwise start a new statement, and so on — are left alone, as are backslashes in comments.

### Examples

```
# bad
foo. \
 bar
foo \
 &.bar \
 .baz
# good
foo.
 bar
foo
 &.bar
 .baz
# bad
[foo, \
 bar]
{foo: \
 bar}
# good
[foo,
 bar]
{foo:
 bar}
# bad
foo(bar, \
 baz)
# good
foo(bar,
 baz)
# also good - backslash in string concatenation is not redundant
foo('bar' \
 'baz')
# also good - backslash at the end of a comment is not redundant
foo(bar, # \
 baz)
# also good - backslash at the line following the newline begins with a + or -,
# it is not redundant
1 \
 + 2 \
 - 3
# also good - backslash with newline between the method name and its arguments,
# it is not redundant.
some_method \
 (argument)
```
## Style/RedundantMinMaxBy

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.85 | - |

Identifies places where `max_by { … }`, `min_by { … }`, or
`minmax_by { … }` can be replaced by `max`, `min`, or `minmax`.

## Style/RedundantPercentQ

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.76 | - |

Checks for usage of the %q/%Q syntax when '' or "" would do.

### Examples

```
# bad
name = %q(Bruce Wayne)
time = %q(8 o'clock)
question = %q("What did you say?")
# good
name = 'Bruce Wayne'
time = "8 o'clock"
question = '"What did you say?"'
```
## Style/RedundantRegexpArgument

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.53 | - |

Identifies places where argument can be replaced from a deterministic regexp to a string.

### Examples

```
# bad
'foo'.byteindex(/f/)
'foo'.byterindex(/f/)
'foo'.gsub(/f/, 'x')
'foo'.gsub!(/f/, 'x')
'foo'.partition(/f/)
'foo'.rpartition(/f/)
'foo'.scan(/f/)
'foo'.split(/f/)
'foo'.start_with?(/f/)
'foo'.sub(/f/, 'x')
'foo'.sub!(/f/, 'x')
# good
'foo'.byteindex('f')
'foo'.byterindex('f')
'foo'.gsub('f', 'x')
'foo'.gsub!('f', 'x')
'foo'.partition('f')
'foo'.rpartition('f')
'foo'.scan('f')
'foo'.split('f')
'foo'.start_with?('f')
'foo'.sub('f', 'x')
'foo'.sub!('f', 'x')
```
## Style/RedundantRegexpConstructor

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.52 | - |

Checks for the instantiation of a regexp using a redundant `Regexp.new` or `Regexp.compile`.
Autocorrect replaces it with a regexp literal which is the simplest and fastest.

## Style/RedundantRegexpEscape

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.85 | - |

Checks for redundant escapes inside `Regexp` literals.

## Style/RedundantReturn

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.10 | 0.14 |

Checks for redundant `return` expressions. Ruby methods
implicitly return the value of the last evaluated expression,
so an explicit `return` at the end of a method body is unnecessary.

### Examples

```
# These bad cases should be extended to handle methods whose body is
# if/else or a case expression with a default branch.
# bad
def test
 return something
end
# bad
def test
 one
 two
 three
 return something
end
# bad
def test
 return something if something_else
end
# good
def test
 something if something_else
end
# good
def test
 if x
 elsif y
 else
 end
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowMultipleReturnValues |
 | Boolean |

## Style/RedundantSelf

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.10 | 0.13 |

Checks for redundant uses of `self`.

The usage of `self` is only needed when:

-
Sending a message to same object with zero arguments in presence of a method name clash with an argument or a local variable.
-
Calling an attribute writer to prevent a local variable assignment.

Note, with using explicit self you can only send messages with public or protected scope, you cannot send private messages this way.

Note we allow uses of `self` with operators because it would be awkward
otherwise. Also allows the use of `self.it` without arguments in blocks,
as in `0.times { self.it }`, following `Lint/ItWithoutArgumentsInBlock` cop.

## Style/RedundantSelfAssignment

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.90 | - |

Checks for places where redundant assignments are made for in place modification methods.

## Style/RedundantSelfAssignmentBranch

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.19 | - |

Checks for places where conditional branch makes redundant self-assignment.

It only detects local variable because it may replace state of instance variable,
class variable, and global variable that have state across methods with `nil`.

## Style/RedundantSort

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.76 | 1.22 |

Identifies instances of sorting and then
taking only the first or last element. The same behavior can
be accomplished without a relatively expensive sort by using
`Enumerable#min` instead of sorting and taking the first
element and `Enumerable#max` instead of sorting and taking the
last element. Similarly, `Enumerable#min_by` and
`Enumerable#max_by` can replace `Enumerable#sort_by` calls
after which only the first or last element is used.

### Safety

This cop is unsafe, because `sort…last` and `max` may not return the
same element in all cases.

In an enumerable where there are multiple elements where `a <⇒ b == 0`,
or where the transformation done by the `sort_by` block has the
same result, `sort.last` and `max` (or `sort_by.last` and `max_by`)
will return different elements. `sort.last` will return the last
element but `max` will return the first element.

For example:

```
class MyString < String; end
strings = [MyString.new('test'), 'test']
strings.sort.last.class #=> String
strings.max.class #=> MyString
```
```
words = %w(dog horse mouse)
words.sort_by { |word| word.length }.last #=> 'mouse'
words.max_by { |word| word.length } #=> 'horse'
```
### Examples

```
# bad
[2, 1, 3].sort.first
[2, 1, 3].sort[0]
[2, 1, 3].sort.at(0)
[2, 1, 3].sort.slice(0)
# good
[2, 1, 3].min
# bad
[2, 1, 3].sort.last
[2, 1, 3].sort[-1]
[2, 1, 3].sort.at(-1)
[2, 1, 3].sort.slice(-1)
# good
[2, 1, 3].max
# bad
arr.sort_by(&:foo).first
arr.sort_by(&:foo)[0]
arr.sort_by(&:foo).at(0)
arr.sort_by(&:foo).slice(0)
# good
arr.min_by(&:foo)
# bad
arr.sort_by(&:foo).last
arr.sort_by(&:foo)[-1]
arr.sort_by(&:foo).at(-1)
arr.sort_by(&:foo).slice(-1)
# good
arr.max_by(&:foo)
```
## Style/RedundantStringEscape

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.37 | - |

Checks for redundant escapes in string literals.

### Examples

```
# bad - no need to escape # without following {/$/@
"\#foo"
# bad - no need to escape single quotes inside double quoted string
"\'foo\'"
# bad - heredocs are also checked for unnecessary escapes
<<~STR
 \#foo \"foo\"
STR
# good
"#foo"
# good
"\#{no_interpolation}"
# good
"'foo'"
# good
"foo\
bar"
# good
<<~STR
 #foo "foo"
STR
```
## Style/RedundantStructKeywordInit

| Requires Ruby version 3.2 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always (Unsafe) | 1.85 | 1.86 |

Checks for redundant `keyword_init` option for `Struct.new`.

Since Ruby 3.2, `keyword_init` in `Struct.new` defaults to `nil` behavior.
Therefore, this cop detects and autocorrects redundant `keyword_init: nil`
and `keyword_init: true` in `Struct.new`.

This cop is disabled by default because `keyword_init: true` is not purely
redundant. It changes behavior in the following ways:

-
`Struct#keyword_init?`returns`true`instead of`nil`.
-
A `Struct`with`keyword_init: true`accepts a`Hash`argument and expands it as keyword arguments, whereas without it the`Hash`is treated as a positional argument.
-
`keyword_init: true`raises an`ArgumentError`for positional arguments, enforcing keyword-only initialization.

## Style/RegexpLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.30 |

Enforces using `//` or `%r` around regular expressions.

| The following `%r`cases using a regexp that starts with a blank or`=`as a method argument are allowed to prevent syntax errors. |

```
do_something %r{ regexp} # `do_something / regexp/` is an invalid syntax.
do_something %r{=regexp} # `do_something /=regexp/` is an invalid syntax.
```
### Examples

#### EnforcedStyle: slashes (default)

```
# bad
snake_case = %r{^[\dA-Z_]+$}
# bad
regex = %r{
 foo
 (bar)
 (baz)
}x
# good
snake_case = /^[\dA-Z_]+$/
# good
regex = /
 foo
 (bar)
 (baz)
/x
```
#### EnforcedStyle: percent_r

```
# bad
snake_case = /^[\dA-Z_]+$/
# bad
regex = /
 foo
 (bar)
 (baz)
/x
# good
snake_case = %r{^[\dA-Z_]+$}
# good
regex = %r{
 foo
 (bar)
 (baz)
}x
```
#### EnforcedStyle: mixed

```
# bad
snake_case = %r{^[\dA-Z_]+$}
# bad
regex = /
 foo
 (bar)
 (baz)
/x
# good
snake_case = /^[\dA-Z_]+$/
# good
regex = %r{
 foo
 (bar)
 (baz)
}x
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| AllowInnerSlashes |
 | Boolean |

## Style/RequireOrder

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always (Unsafe) | 1.40 | - |

Sort `require` and `require_relative` in alphabetical order.

### Examples

```
# bad
require 'b'
require 'a'
# good
require 'a'
require 'b'
# bad
require_relative 'b'
require_relative 'a'
# good
require_relative 'a'
require_relative 'b'
# good (sorted within each section separated by a blank line)
require 'a'
require 'd'
require 'b'
require 'c'
# good
require 'b'
require_relative 'c'
require 'a'
# bad
require 'a'
require 'c' if foo
require 'b'
# good
require 'a'
require 'b'
require 'c' if foo
# bad
require 'c'
if foo
 require 'd'
 require 'b'
end
require 'a'
# good
require 'c'
if foo
 require 'b'
 require 'd'
end
require 'a'
```
## Style/RescueModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.34 |

Checks for uses of `rescue` in its modifier form. It is added for the following
reasons:

-
The syntax of modifier form `rescue`can be misleading because it might lead us to believe that`rescue`handles the given exception but it actually rescues all exceptions to return the given rescue block. In this case, value returned by handle_error or SomeException.
-
Modifier form `rescue`would rescue all the exceptions. It would silently skip all exceptions or errors and handle the error. Example: If`NoMethodError`is raised, modifier form rescue would handle the exception.

### Examples

```
# bad
some_method rescue handle_error
# bad
some_method rescue SomeException
# good
begin
 some_method
rescue
 handle_error
end
# good
begin
 some_method
rescue SomeException
 handle_error
end
```
## Style/RescueStandardError

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for rescuing `StandardError`. There are two supported
styles `implicit` and `explicit`. This cop will not register an offense
if any error other than `StandardError` is specified.

### Examples

## Style/ReturnNil

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.50 | - |

Enforces consistency between `return nil` and `return`.

This cop is disabled by default. Because there seems to be a perceived semantic difference
between `return` and `return nil`. The former can be seen as just halting evaluation,
while the latter might be used when the return value is of specific concern.

Supported styles are `return` and `return_nil`.

## Style/ReturnNilInPredicateMethodDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.53 | 1.67 |

Checks for predicate method definitions that return `nil`.
A predicate method should only return a boolean value.

### Safety

Autocorrection is marked as unsafe because the change of the return value
from `nil` to `false` could potentially lead to incompatibility issues.

### Examples

```
# bad
def foo?
 return if condition
 do_something?
end
# bad
def foo?
 return nil if condition
 do_something?
end
# good
def foo?
 return false if condition
 do_something?
end
# bad
def foo?
 if condition
 nil
 else
 true
 end
end
# good
def foo?
 if condition
 false
 else
 true
 end
end
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowedMethods |
 | Array |
| AllowedPatterns |
 | Array |

## Style/ReverseFind

| Requires Ruby version 4.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.84 | - |

Identifies places where `array.reverse.find` can be replaced by `array.rfind`.

## Style/SafeNavigation

| Requires Ruby version 2.3 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.43 | 1.67 |

Transforms usages of a method call safeguarded by a non `nil`
check for the variable whose method is being called to
safe navigation (`&.`). If there is a method chain, all of the methods
in the chain need to be checked for safety, and all of the methods will
need to be changed to use safe navigation.

The default for `ConvertCodeThatCanStartToReturnNil` is `false`.
When configured to `true`, this will
check for code in the format `!foo.nil? && foo.bar`. As it is written,
the return of this code is limited to `false` and whatever the return
of the method is. If this is converted to safe navigation,
`foo&.bar` can start returning `nil` as well as what the method
returns.

The default for `MaxChainLength` is `2`.
We have limited the cop to not register an offense for method chains
that exceed this option’s value.

| This cop will recognize offenses but not autocorrect code when the
right hand side (RHS) of the `&&`statement is an`||`statement
(eg.`foo && (foo.bar? || foo.baz?)`). It can be corrected
manually by removing the`foo &&`and adding`&.`to each`foo`on the RHS. |

### Safety

Autocorrection is unsafe because if a value is `false`, the resulting
code will have different behavior or raise an error.

```
x = false
x && x.foo # return false
x&.foo # raises NoMethodError
```
Additionally, when a method chain is converted, a `NoMethodError` that
the original code raised for an intermediate `nil` value is suppressed:

```
x = Struct.new(:foo).new(nil)
x && x.foo.bar # raises NoMethodError
x&.foo&.bar # returns nil
```
### Examples

```
# bad
foo.bar if foo
foo.bar.baz if foo
foo.bar(param1, param2) if foo
foo.bar { |e| e.something } if foo
foo.bar(param) { |e| e.something } if foo
foo.bar if !foo.nil?
foo.bar unless !foo
foo.bar unless foo.nil?
foo && foo.bar
foo && foo.bar.baz
foo && foo.bar(param1, param2)
foo && foo.bar { |e| e.something }
foo && foo.bar(param) { |e| e.something }
foo ? foo.bar : nil
foo.nil? ? nil : foo.bar
!foo.nil? ? foo.bar : nil
!foo ? nil : foo.bar
# good
foo&.bar
foo&.bar&.baz
foo&.bar(param1, param2)
foo&.bar { |e| e.something }
foo&.bar(param) { |e| e.something }
foo && foo.bar.baz.qux # method chain with more than 2 methods
foo && foo.nil? # method that `nil` responds to
# Method calls that do not use `.`
foo && foo < bar
foo < bar if foo
# When checking `foo&.empty?` in a conditional, `foo` being `nil` will actually
# do the opposite of what the author intends.
foo && foo.empty?
# This could start returning `nil` as well as the return of the method
foo.nil? || foo.bar
!foo || foo.bar
# Methods that are used on assignment, arithmetic operation or
# comparison should not be converted to use safe navigation
foo.baz = bar if foo
foo.baz + bar if foo
foo.bar > 2 if foo
foo ? foo[index] : nil # Ignored `foo&.[](index)` due to unclear readability benefit.
foo ? foo[idx] = v : nil # Ignored `foo&.[]=(idx, v)` due to unclear readability benefit.
foo ? foo * 42 : nil # Ignored `foo&.*(42)` due to unclear readability benefit.
```
## Style/SafeNavigationChainLength

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | No | 1.68 | - |

Enforces safe navigation chains length to not exceed the configured maximum.
The longer the chain is, the harder it becomes to track what on it could be
returning `nil`.

There is a potential interplay with `Style/SafeNavigation` - if both are enabled
and their settings are "incompatible", one of the cops will complain about what
the other proposes.

E.g. if `Style/SafeNavigation` is configured with `MaxChainLength: 2` (default)
and this cop is configured with `Max: 1`, then for `foo.bar.baz if foo` the former
will suggest `foo&.bar&.baz`, which is an offense for the latter.

## Style/Sample

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.30 | - |

Identifies usages of `shuffle.first`,
`shuffle.last`, and `shuffle[]` and change them to use
`sample` instead.

### Examples

```
# bad
[1, 2, 3].shuffle.first
[1, 2, 3].shuffle.first(2)
[1, 2, 3].shuffle.last
[2, 1, 3].shuffle.at(0)
[2, 1, 3].shuffle.slice(0)
[1, 2, 3].shuffle[2]
[1, 2, 3].shuffle[0, 2] # sample(2) will do the same
[1, 2, 3].shuffle[0..2] # sample(3) will do the same
[1, 2, 3].shuffle(random: Random.new).first
# good
[1, 2, 3].shuffle
[1, 2, 3].sample
[1, 2, 3].sample(3)
[1, 2, 3].shuffle[1, 3] # sample(3) might return a longer Array
[1, 2, 3].shuffle[1..3] # sample(3) might return a longer Array
[1, 2, 3].shuffle[foo, bar]
[1, 2, 3].shuffle(random: Random.new)
```
## Style/SelectByKind

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.85 | - |

Looks for places where a subset of an Enumerable (array,
range, set, etc.; see note below) is calculated based on a class type
check, and suggests `grep` or `grep_v` instead.

| Hashes do not behave as you may expect with `grep`, which
means that`hash.grep`is not equivalent to`hash.select`. Although
RuboCop is limited by static analysis, this cop attempts to avoid
registering an offense when the receiver is a hash (hash literal,`Hash.new`,`Hash#[]`, or`to_h`/`to_hash`). |

## Style/SelectByRange

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.85 | - |

Looks for places where a subset of an Enumerable (array,
range, set, etc.; see note below) is calculated based on a range
check, and suggests `grep` or `grep_v` instead.

| Hashes do not behave as you may expect with `grep`, which
means that`hash.grep`is not equivalent to`hash.select`. Although
RuboCop is limited by static analysis, this cop attempts to avoid
registering an offense when the receiver is a hash (hash literal,`Hash.new`,`Hash#[]`, or`to_h`/`to_hash`). |

### Safety

Autocorrection is marked as unsafe because the cop cannot guarantee that the receiver is actually an array by static analysis, so the correction may not be actually equivalent.

### Examples

```
# bad (select or find_all)
array.select { |x| x.between?(1, 10) }
array.select { |x| (1..10).cover?(x) }
array.select { |x| (1..10).include?(x) }
# bad (reject)
array.reject { |x| x.between?(1, 10) }
# bad (find or detect)
array.find { |x| x.between?(1, 10) }
array.detect { |x| (1..10).cover?(x) }
# bad (negative form)
array.reject { |x| !x.between?(1, 10) }
array.find { |x| !(1..10).cover?(x) }
# good
array.grep(1..10)
array.grep_v(1..10)
array.grep(1..10).first
array.grep_v(1..10).first
```
## Style/SelectByRegexp

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.22 | - |

Looks for places where a subset of an Enumerable (array,
range, set, etc.; see note below) is calculated based on a `Regexp`
match, and suggests `grep` or `grep_v` instead.

| Hashes do not behave as you may expect with `grep`, which
means that`hash.grep`is not equivalent to`hash.select`. Although
RuboCop is limited by static analysis, this cop attempts to avoid
registering an offense when the receiver is a hash (hash literal,`Hash.new`,`Hash#[]`, or`to_h`/`to_hash`). |

| `grep`and`grep_v`were optimized when used without a block
in Ruby 3.0, but may be slower in previous versions.
See https://bugs.ruby-lang.org/issues/17030 |

### Safety

Autocorrection is marked as unsafe because `MatchData` will
not be created by `grep`, but may have previously been relied
upon after the `match?` or `=~` call.

Additionally, the cop cannot guarantee that the receiver of
`select` or `reject` is actually an array by static analysis,
so the correction may not be actually equivalent.

### Examples

```
# bad (select, filter, or find_all)
array.select { |x| x.match? /regexp/ }
array.select { |x| /regexp/.match?(x) }
array.select { |x| x =~ /regexp/ }
array.select { |x| /regexp/ =~ x }
# bad (reject)
array.reject { |x| x.match? /regexp/ }
array.reject { |x| /regexp/.match?(x) }
array.reject { |x| x =~ /regexp/ }
array.reject { |x| /regexp/ =~ x }
# bad (negative form)
array.reject { |x| !x.match? /regexp/ }
# good
array.grep(regexp)
array.grep_v(regexp)
```
## Style/SelfAssignment

## Style/Semicolon

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.19 |

Checks for multiple expressions placed on the same line. It also checks for lines terminated with a semicolon. In idiomatic Ruby, each expression should be on its own line for readability.

This cop has `AllowAsExpressionSeparator` configuration option.
It allows `;` to separate several expressions on the same line.

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| AllowAsExpressionSeparator |
 | Boolean |

## Style/Send

## Style/SendWithLiteralMethodName

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.64 | - |

Detects the use of the `public_send` method with a literal method name argument.
Since the `send` method can be used to call private methods, by default,
only the `public_send` method is detected.

| Writer methods with names ending in `=`are always permitted because their
behavior differs as follows: |

```
def foo=(foo)
 @foo = foo
 42
end
self.foo = 1 # => 1
send(:foo=, 1) # => 42
```
### Safety

This cop is not safe because it can incorrectly detect based on the receiver.
Additionally, when `AllowSend` is set to `true`, it cannot determine whether
the `send` method being detected is calling a private method.

## Style/SignalException

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.11 | 0.37 |

Checks for uses of `fail` and `raise`.

### Examples

#### EnforcedStyle: only_raise (default)

```
# The `only_raise` style enforces the sole use of `raise`.
# bad
begin
 fail
rescue Exception
 # handle it
end
def watch_out
 fail
rescue Exception
 # handle it
end
Kernel.fail
# good
begin
 raise
rescue Exception
 # handle it
end
def watch_out
 raise
rescue Exception
 # handle it
end
Kernel.raise
```
#### EnforcedStyle: only_fail

```
# The `only_fail` style enforces the sole use of `fail`.
# bad
begin
 raise
rescue Exception
 # handle it
end
def watch_out
 raise
rescue Exception
 # handle it
end
Kernel.raise
# good
begin
 fail
rescue Exception
 # handle it
end
def watch_out
 fail
rescue Exception
 # handle it
end
Kernel.fail
```
#### EnforcedStyle: semantic

```
# The `semantic` style enforces the use of `fail` to signal an
# exception, then will use `raise` to trigger an offense after
# it has been rescued.
# bad
begin
 raise
rescue Exception
 # handle it
end
def watch_out
 # Error thrown
rescue Exception
 fail
end
Kernel.fail
Kernel.raise
# good
begin
 fail
rescue Exception
 # handle it
end
def watch_out
 fail
rescue Exception
 raise 'Preferably with descriptive message'
end
explicit_receiver.fail
explicit_receiver.raise
```
## Style/SingleArgumentDig

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.89 | - |

Sometimes using `dig` method ends up with just a single
argument. In such cases, dig should be replaced with `[]`.

Since replacing `hash&.dig(:key)` with `hash[:key]` could potentially lead to error,
calls to the `dig` method using safe navigation will be ignored.

## Style/SingleLineBlockParams

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.16 | 1.6 |

Checks whether the block parameters of a single-line method accepting a block match the names specified via configuration.

For instance one can configure `reduce`(`inject`) to use |a, e| as
parameters.

Configuration option: Methods
Should be set to use this cop. `Array` of hashes, where each key is the
method name and value - array of argument names.

## Style/SingleLineDoEndBlock

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.57 | - |

Checks for single-line `do`…`end` block.

In practice a single line `do`…`end` is autocorrected when `EnforcedStyle: semantic`
is configured for `Style/BlockDelimiters`. The autocorrection maintains the
`do` … `end` syntax to preserve semantics and does not change it to `{`…`}` block.

| If `InspectBlocks`is set to`true`for`Layout/RedundantLineBreak`, blocks will
be autocorrected to be on a single line if possible. This cop respects that configuration
by not registering an offense if it would subsequently cause a`Layout/RedundantLineBreak`offense. |

## Style/SingleLineMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 1.8 |

Checks for single-line method definitions that contain a body. Single-line methods with a body are harder to read and debug than their multi-line equivalents. It will accept single-line methods with no body.

Endless methods added in Ruby 3.0 are also accepted by this cop.

If `Style/EndlessMethod` is enabled with `EnforcedStyle: allow_single_line`, `allow_always`,
`require_single_line`, or `require_always`, single-line methods will be autocorrected
to endless methods if there is only one statement in the body.

## Style/SlicingWithRange

| Requires Ruby version 2.6 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.83 | - |

Checks that arrays are not sliced with the redundant `ary[0..-1]`, replacing it with `ary`,
and ensures arrays are sliced with endless ranges instead of `ary[start..-1]` on Ruby 2.6+,
and with beginless ranges instead of `ary[nil..end]` on Ruby 2.7+.

### Safety

This cop is unsafe because `x..-1` and `x..` are only guaranteed to
be equivalent for `Array#[]`, `String#[]`, and the cop cannot determine what class
the receiver is.

For example:

```
sum = proc { |ary| ary.sum }
sum[-3..-1] # => -6
sum[-3..] # Hangs forever
```
### Examples

```
# bad
items[0..-1]
items[0..nil]
items[0...nil]
# good
items
# bad
items[1..-1] # Ruby 2.6+
items[1..nil] # Ruby 2.6+
# good
items[1..] # Ruby 2.6+
# bad
items[nil..42] # Ruby 2.7+
# good
items[..42] # Ruby 2.7+
items[0..42] # Ruby 2.7+
```
## Style/SoleNestedConditional

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.89 | 1.5 |

If the branch of a conditional consists solely of a conditional node, its conditions can be combined with the conditions of the outer branch. This helps to keep the nesting level from getting too deep.

### Examples

```
# bad
if condition_a
 if condition_b
 do_something
 end
end
# bad
if condition_b
 do_something
end if condition_a
# good
if condition_a && condition_b
 do_something
end
```
## Style/SpecialGlobalVars

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.13 | 0.36 |

Looks for uses of Perl-style global variables.
Perl-style global variables like `$;` or `$/` are cryptic
and hard to understand without consulting documentation.
The `English` library provides descriptive aliases like
`$FIELD_SEPARATOR` and `$INPUT_RECORD_SEPARATOR`.

Correcting to global variables in the `English` library
will add a require statement to the top of the file if
enabled by RequireEnglish config.

### Safety

Autocorrection is marked as unsafe because if `RequireEnglish` is not
true, replacing perl-style variables with english variables will break.

### Examples

#### EnforcedStyle: use_english_names (default)

```
# good
require 'English' # or this could be in another file.
puts $LOAD_PATH
puts $LOADED_FEATURES
puts $PROGRAM_NAME
puts $ERROR_INFO
puts $ERROR_POSITION
puts $FIELD_SEPARATOR # or $FS
puts $OUTPUT_FIELD_SEPARATOR # or $OFS
puts $INPUT_RECORD_SEPARATOR # or $RS
puts $OUTPUT_RECORD_SEPARATOR # or $ORS
puts $INPUT_LINE_NUMBER # or $NR
puts $LAST_READ_LINE
puts $DEFAULT_OUTPUT
puts $DEFAULT_INPUT
puts $PROCESS_ID # or $PID
puts $CHILD_STATUS
puts $LAST_MATCH_INFO
puts $IGNORECASE
puts $ARGV # or ARGV
```
#### EnforcedStyle: use_perl_names

```
# good
puts $:
puts $"
puts $0
puts $!
puts $@
puts $;
puts $,
puts $/
puts $\
puts $.
puts $_
puts $>
puts $<
puts $$
puts $?
puts $~
puts $=
puts $*
```
#### EnforcedStyle: use_builtin_english_names

```
# good
# Like `use_perl_names` but allows builtin global vars.
puts $LOAD_PATH
puts $LOADED_FEATURES
puts $PROGRAM_NAME
puts ARGV
puts $:
puts $"
puts $0
puts $!
puts $@
puts $;
puts $,
puts $/
puts $\
puts $.
puts $_
puts $>
puts $<
puts $$
puts $?
puts $~
puts $=
puts $*
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| RequireEnglish |
 | Boolean |
| EnforcedStyle |
 |
 |

## Style/StabbyLambdaParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.35 | - |

Checks for parentheses around stabby lambda arguments.
There are two different styles. Defaults to `require_parentheses`.

### Examples

## Style/StaticClass

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 1.3 | 1.89 |

Checks for places where classes with only class methods can be replaced with a module. Classes should be used only when it makes sense to create instances out of them.

When `AllCops/UseProjectIndex` is enabled and the `rubydex` gem is
installed, classes that are subclassed anywhere in the project are
not reported, since converting them to modules would break their
subclasses.

### Safety

This cop is unsafe, because it is possible that this class is a parent for some other subclass, monkey-patched with instance methods or a dummy instance is instantiated from it somewhere.

### Examples

```
# bad
class SomeClass
 def self.some_method
 # body omitted
 end
 def self.some_other_method
 # body omitted
 end
end
# good
module SomeModule
 module_function
 def some_method
 # body omitted
 end
 def some_other_method
 # body omitted
 end
end
# good - has instance method
class SomeClass
 def instance_method; end
 def self.class_method; end
end
```
## Style/StderrPuts

## Style/StringChars

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.12 | - |

Checks for uses of `String#split` with empty string or regexp literal argument.

### Safety

This cop is unsafe because it cannot be guaranteed that the receiver
is actually a string. If another class has a `split` method with
different behavior, it would be registered as a false positive.

## Style/StringConcatenation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.89 | 1.18 |

Checks for places where string concatenation can be replaced with string interpolation.

The cop can autocorrect simple cases but will skip autocorrecting more complex cases where the resulting code would be harder to read. In those cases, it might be useful to extract statements to local variables or methods which you can then interpolate in a string.

| When concatenation between two strings is broken over multiple
lines, this cop does not register an offense; instead, `Style/LineEndConcatenation`will pick up the offense if enabled. |

Two modes are supported:
1. `aggressive` style checks and corrects all occurrences of ```
` where
either the left or right side of `
```
 is a string literal.
2. `conservative` style on the other hand, checks and corrects only if
left side (receiver of `+` method call) is a string literal.
This is useful when the receiver is some expression that returns string like `Pathname`
instead of a string literal.

### Safety

This cop is unsafe in `aggressive` mode, as it cannot be guaranteed that
the receiver is actually a string, which can result in a false positive.

### Examples

#### Mode: aggressive (default)

```
# bad
email_with_name = user.name + ' <' + user.email + '>'
Pathname.new('/') + 'test'
# good
email_with_name = "#{user.name} <#{user.email}>"
email_with_name = format('%s <%s>', user.name, user.email)
"#{Pathname.new('/')}test"
# accepted, line-end concatenation
name = 'First' +
 'Last'
```
## Style/StringHashKeys

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 0.52 | 0.75 |

Checks for the use of strings as keys in hashes. The use of symbols is preferred instead.

### Safety

This cop is unsafe because while symbols are preferred for hash keys, there are instances when string keys are required.

## Style/StringLiterals

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.36 |

Checks if uses of quotes match the configured preference.

### Examples

## Style/StringLiteralsInInterpolation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.27 | - |

Checks that quotes inside string, symbol, and regexp interpolations match the configured preference.

### Examples

#### EnforcedStyle: single_quotes (default)

```
# bad
string = "Tests #{success ? "PASS" : "FAIL"}"
symbol = :"Tests #{success ? "PASS" : "FAIL"}"
heredoc = <<~TEXT
 Tests #{success ? "PASS" : "FAIL"}
TEXT
regexp = /Tests #{success ? "PASS" : "FAIL"}/
# good
string = "Tests #{success ? 'PASS' : 'FAIL'}"
symbol = :"Tests #{success ? 'PASS' : 'FAIL'}"
heredoc = <<~TEXT
 Tests #{success ? 'PASS' : 'FAIL'}
TEXT
regexp = /Tests #{success ? 'PASS' : 'FAIL'}/
```
#### EnforcedStyle: double_quotes

```
# bad
string = "Tests #{success ? 'PASS' : 'FAIL'}"
symbol = :"Tests #{success ? 'PASS' : 'FAIL'}"
heredoc = <<~TEXT
 Tests #{success ? 'PASS' : 'FAIL'}
TEXT
regexp = /Tests #{success ? 'PASS' : 'FAIL'}/
# good
string = "Tests #{success ? "PASS" : "FAIL"}"
symbol = :"Tests #{success ? "PASS" : "FAIL"}"
heredoc = <<~TEXT
 Tests #{success ? "PASS" : "FAIL"}
TEXT
regexp = /Tests #{success ? "PASS" : "FAIL"}/
```
## Style/StringMethods

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | Always | 0.34 | 0.34 |

Enforces the use of consistent method names
from the `String` class.

## Style/StructInheritance

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always (Unsafe) | 0.29 | 1.20 |

Checks for inheritance from `Struct.new`. Inheriting from `Struct.new`
adds a superfluous level in inheritance tree.

### Safety

Autocorrection is unsafe because it will change the inheritance
tree (e.g. return value of `Module#ancestors`) of the constant.

### Examples

```
# bad
class Person < Struct.new(:first_name, :last_name)
 def age
 42
 end
end
Person.ancestors
# => [Person, #<Class:0x000000010b4e14a0>, Struct, (...)]
# good
Person = Struct.new(:first_name, :last_name) do
 def age
 42
 end
end
Person.ancestors
# => [Person, Struct, (...)]
```
## Style/SuperArguments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.64 | - |

Checks for redundant argument forwarding when calling super with arguments identical to the method definition.

Using zero arity `super` within a `define_method` block results in `RuntimeError`:

```
def m
 define_method(:foo) { super() } # => OK
end
def m
 define_method(:foo) { super } # => RuntimeError
end
```
Furthermore, any arguments accompanied by a block may potentially be delegating to
`define_method`, therefore, `super` used within these blocks will be allowed.
This approach might result in false negatives, yet ensuring safe detection takes precedence.

| When forwarding the same arguments but replacing the block argument with a new inline block, it is not necessary to explicitly list the non-block arguments. As such, an offense will be registered in this case. |

### Examples

```
# bad
def method(*args, **kwargs)
 super(*args, **kwargs)
end
# good - implicitly passing all arguments
def method(*args, **kwargs)
 super
end
# good - forwarding a subset of the arguments
def method(*args, **kwargs)
 super(*args)
end
# good - forwarding no arguments
def method(*args, **kwargs)
 super()
end
# bad - forwarding with overridden block
def method(*args, **kwargs, &block)
 super(*args, **kwargs) { do_something }
end
# good - implicitly passing all non-block arguments
def method(*args, **kwargs, &block)
 super { do_something }
end
# good - assigning to the block variable before calling super
def method(&block)
 # Assigning to the block variable would pass the old value to super,
 # under this circumstance the block must be referenced explicitly.
 block ||= proc { 'fallback behavior' }
 super(&block)
end
```
## Style/SuperWithArgsParentheses

## Style/SwapValues

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always (Unsafe) | 1.1 | - |

Enforces the use of shorthand-style swapping of 2 variables.

### Safety

Autocorrection is unsafe, because the temporary variable used to swap variables will be removed, but may be referred to elsewhere.

## Style/SymbolArray

| Requires Ruby version 2.0 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.49 |

Checks for array literals made up of symbols that are not using the %i() syntax.

Alternatively, it checks for symbol arrays using the %i() syntax on projects which do not want to use that syntax, perhaps because they support a version of Ruby lower than 2.0.

Configuration option: MinSize
If set, arrays with fewer elements than this value will not trigger the
cop. For example, a `MinSize` of `3` will not enforce a style on an
array of 2 or fewer elements.

### Examples

### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| MinSize |
 | Integer |

## Style/SymbolProc

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.26 | 1.64 |

Use symbols as procs when possible.

If you prefer a style that allows block for method with arguments,
please set `true` to `AllowMethodsWithArguments`.
`define_method?` methods are allowed by default.
These are customizable with `AllowedMethods` option.

### Safety

This cop is unsafe because there is a difference that a `Proc`
generated from `Symbol#to_proc` behaves as a lambda, while
a `Proc` generated from a block does not.
For example, a lambda will raise an `ArgumentError` if the
number of arguments is wrong, but a non-lambda `Proc` will not.

For example:

```
class Foo
 def bar
 :bar
 end
end
def call(options = {}, &block)
 block.call(Foo.new, options)
end
call { |x| x.bar }
#=> :bar
call(&:bar)
# ArgumentError: wrong number of arguments (given 1, expected 0)
```
It is also unsafe because `Symbol#to_proc` does not work with
`protected` methods which would otherwise be accessible.

For example:

```
class Box
 def initialize
 @secret = rand
 end
 def normal_matches?(*others)
 others.map { |other| other.secret }.any?(secret)
 end
 def symbol_to_proc_matches?(*others)
 others.map(&:secret).any?(secret)
 end
 protected
 attr_reader :secret
end
boxes = [Box.new, Box.new]
Box.new.normal_matches?(*boxes)
# => false
boxes.first.normal_matches?(*boxes)
# => true
Box.new.symbol_to_proc_matches?(*boxes)
# => NoMethodError: protected method `secret' called for #<Box...>
boxes.first.symbol_to_proc_matches?(*boxes)
# => NoMethodError: protected method `secret' called for #<Box...>
```
### Examples

```
# bad
something.map { |s| s.upcase }
something.map { _1.upcase }
# good
something.map(&:upcase)
```
#### AllowMethodsWithArguments: false (default)

```
# bad
something.do_something(foo) { |o| o.bar }
# good
something.do_something(foo, &:bar)
```
#### AllowComments: false (default)

```
# bad
something.do_something do |s| # some comment
 # some comment
 s.upcase # some comment
 # some comment
end
```
#### AllowComments: true

```
# good - if there are comment in either position
something.do_something do |s| # some comment
 # some comment
 s.upcase # some comment
 # some comment
end
```
## Style/TallyMethod

| Requires Ruby version 2.7 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | No | Always (Unsafe) | 1.85 | - |

Checks for manual counting patterns that can be replaced by `Enumerable#tally`.

The cop detects the following patterns:

-
`each_with_object(Hash.new(0)) { |item, counts| counts[item] += 1 }`
-
`group_by(&:itself).transform_values(&:count)`
-
`group_by { |x| x }.transform_values(&:size)`
-
`group_by { |x| x }.transform_values { |v| v.length }`

## Style/TernaryParentheses

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.42 | 0.46 |

Checks for the presence of parentheses around ternary
conditions. It is configurable to enforce inclusion or omission of
parentheses using `EnforcedStyle`. Omission is only enforced when
removing the parentheses won’t cause a different behavior.

`AllowSafeAssignment` option for safe assignment.
By safe assignment we mean putting parentheses around
an assignment to indicate "I know I’m using an assignment
as a condition. It’s not a mistake."

### Examples

#### EnforcedStyle: require_no_parentheses (default)

```
# bad
foo = (bar?) ? a : b
foo = (bar.baz?) ? a : b
foo = (bar && baz) ? a : b
# good
foo = bar? ? a : b
foo = bar.baz? ? a : b
foo = bar && baz ? a : b
```
#### EnforcedStyle: require_parentheses

```
# bad
foo = bar? ? a : b
foo = bar.baz? ? a : b
foo = bar && baz ? a : b
# good
foo = (bar?) ? a : b
foo = (bar.baz?) ? a : b
foo = (bar && baz) ? a : b
```
## Style/TopLevelMethodDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 1.15 | - |

Newcomers to Ruby applications may write top-level methods, when ideally they should be organized in appropriate classes or modules. This cop looks for definitions of top-level methods and warns about them.

However, for Ruby scripts it is perfectly fine to use top-level methods. Hence this cop is disabled by default.

### Examples

```
# bad
def some_method
end
# bad
def self.some_method
end
# bad
define_method(:foo) { puts 1 }
# good
module Foo
 def some_method
 end
end
# good
class Foo
 def self.some_method
 end
end
# good
Struct.new do
 def some_method
 end
end
# good
class Foo
 define_method(:foo) { puts 1 }
end
```
## Style/TrailingBodyOnMethodDefinition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for trailing code after the method definition.

| It always accepts endless method definitions that are basically on the same line. |

## Style/TrailingCommaInArguments

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.36 | - |

Checks for trailing comma in argument lists. The supported styles are:

-
`consistent_comma`: Requires a comma after the last argument, for all parenthesized multi-line method calls with arguments.
-
`comma`: Requires a comma after the last argument, but only for parenthesized method calls where each argument is on its own line.
-
`diff_comma`: Requires a comma after the last argument, but only when that argument is followed by an immediate newline, even if there is an inline comment on the same line.
-
`no_comma`: Requires that there is no comma after the last argument.

Regardless of style, trailing commas are not allowed in single-line method calls.

### Examples

#### EnforcedStyleForMultiline: consistent_comma

```
# bad
method(1, 2,)
# good
method(1, 2)
# good
method(
 1, 2,
 3,
)
# good
method(
 1, 2, 3,
)
# good
method(
 1,
 2,
)
```
#### EnforcedStyleForMultiline: comma

```
# bad
method(1, 2,)
# good
method(1, 2)
# bad
method(
 1, 2,
 3,
)
# good
method(
 1, 2,
 3
)
# bad
method(
 1, 2, 3,
)
# good
method(
 1, 2, 3
)
# good
method(
 1,
 2,
)
```
## Style/TrailingCommaInArrayLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

Checks for trailing comma in array literals. The configuration options are:

-
`consistent_comma`: Requires a comma after the last item of all non-empty, multiline array literals.
-
`comma`: Requires a comma after the last item in an array, but only when each item is on its own line.
-
`diff_comma`: Requires a comma after the last item in an array, but only when that item is followed by an immediate newline, even if there is an inline comment on the same line.
-
`no_comma`: Does not require a comma after the last item in an array

### Examples

#### EnforcedStyleForMultiline: consistent_comma

```
# bad
a = [1, 2,]
# good
a = [1, 2]
# good
a = [
 1, 2,
 3,
]
# good
a = [
 1, 2, 3,
]
# good
a = [
 1,
 2,
]
# bad
a = [1, 2,
 3, 4]
# good
a = [1, 2,
 3, 4,]
```
#### EnforcedStyleForMultiline: comma

```
# bad
a = [1, 2,]
# good
a = [1, 2]
# bad
a = [
 1, 2,
 3,
]
# good
a = [
 1, 2,
 3
]
# bad
a = [
 1, 2, 3,
]
# good
a = [
 1, 2, 3
]
# good
a = [
 1,
 2,
]
```
## Style/TrailingCommaInBlockArgs

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 0.81 | - |

Checks whether trailing commas in block arguments are required. Blocks with only one argument and a trailing comma require that comma to be present. Blocks with more than one argument never require a trailing comma.

### Safety

This cop is unsafe because a trailing comma can indicate there are more parameters that are not used.

For example:

```
# with a trailing comma
{foo: 1, bar: 2, baz: 3}.map {|key,| key }
#=> [:foo, :bar, :baz]
# without a trailing comma
{foo: 1, bar: 2, baz: 3}.map {|key| key }
#=> [[:foo, 1], [:bar, 2], [:baz, 3]]
```
This can be fixed by replacing the trailing comma with a placeholder
argument (such as `|key, _value|`).

## Style/TrailingCommaInHashLiteral

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.53 | - |

Checks for trailing comma in hash literals. The configuration options are:

-
`consistent_comma`: Requires a comma after the last item of all non-empty, multiline hash literals.
-
`comma`: Requires a comma after the last item in a hash, but only when each item is on its own line.
-
`diff_comma`: Requires a comma after the last item in a hash, but only when that item is followed by an immediate newline, even if there is an inline comment on the same line.
-
`no_comma`: Does not require a comma after the last item in a hash

### Examples

#### EnforcedStyleForMultiline: consistent_comma

```
# bad
a = { foo: 1, bar: 2, }
# good
a = { foo: 1, bar: 2 }
# good
a = {
 foo: 1, bar: 2,
 qux: 3,
}
# good
a = {
 foo: 1, bar: 2, qux: 3,
}
# good
a = {
 foo: 1,
 bar: 2,
}
# bad
a = { foo: 1, bar: 2,
 baz: 3, qux: 4 }
# good
a = { foo: 1, bar: 2,
 baz: 3, qux: 4, }
```
#### EnforcedStyleForMultiline: comma

```
# bad
a = { foo: 1, bar: 2, }
# good
a = { foo: 1, bar: 2 }
# bad
a = {
 foo: 1, bar: 2,
 qux: 3,
}
# good
a = {
 foo: 1, bar: 2,
 qux: 3
}
# bad
a = {
 foo: 1, bar: 2, qux: 3,
}
# good
a = {
 foo: 1, bar: 2, qux: 3
}
# good
a = {
 foo: 1,
 bar: 2,
}
```
## Style/TrailingMethodEndStatement

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.52 | - |

Checks for trailing code after the method definition.

## Style/TrailingUnderscoreVariable

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.31 | 0.35 |

Checks for extra underscores in variable assignment.

## Style/TrivialAccessors

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 1.15 |

Looks for trivial reader/writer methods, that could
have been created with the attr_* family of functions automatically.
`to_ary`, `to_a`, `to_c`, `to_enum`, `to_h`, `to_hash`, `to_i`, `to_int`, `to_io`,
`to_open`, `to_path`, `to_proc`, `to_r`, `to_regexp`, `to_str`, `to_s`, and `to_sym` methods
are allowed by default. These are customizable with `AllowedMethods` option.

### Examples

```
# bad
def foo
 @foo
end
def bar=(val)
 @bar = val
end
def self.baz
 @baz
end
# good
attr_reader :foo
attr_writer :bar
class << self
 attr_reader :baz
end
```
#### AllowDSLWriters: false

```
# bad
def on_exception(action)
 @on_exception=action
end
# good
attr_writer :on_exception
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| ExactNameMatch |
 | Boolean |
| AllowPredicates |
 | Boolean |
| AllowDSLWriters |
 | Boolean |
| IgnoreClassMethods |
 | Boolean |
| AllowedMethods |
 | Array |

## Style/UnlessElse

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | - |

Looks for `unless` expressions with `else` clauses.

### Examples

```
# bad
unless foo_bar.nil?
 # do something...
else
 # do a different thing...
end
# good
if foo_bar.present?
 # do something...
else
 # do a different thing...
end
```
## Style/UnlessLogicalOperators

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | Yes | No | 1.11 | - |

Checks for the use of logical operators in an `unless` condition.
It discourages such code, as the condition becomes more difficult
to read and understand.

This cop supports two styles:

-
`forbid_mixed_logical_operators`(default)
-
`forbid_logical_operators`

`forbid_mixed_logical_operators` style forbids the use of more than one type
of logical operators. This makes the `unless` condition easier to read
because either all conditions need to be met or any condition needs to be met
in order for the expression to be truthy or falsey.

`forbid_logical_operators` style forbids any use of logical operators.
This makes it even easier to read the `unless` condition as
there is only one condition in the expression.

### Examples

#### EnforcedStyle: forbid_mixed_logical_operators (default)

```
# bad
return unless a || b && c
return unless a && b || c
return unless a && b and c
return unless a || b or c
return unless a && b or c
return unless a || b and c
# good
return unless a && b && c
return unless a || b || c
return unless a and b and c
return unless a or b or c
return unless a?
```
## Style/UnpackFirst

| Requires Ruby version 2.4 |

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.54 | - |

Checks for accessing the first element of `String#unpack`
which can be replaced with the shorter method `unpack1`.

## Style/VariableInterpolation

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.20 |

Checks for variable interpolation (like "#@ivar").

### Examples

```
# bad
"His name is #$name"
/check #$pattern/
"Let's go to the #@store"
# good
"His name is #{$name}"
/check #{$pattern}/
"Let's go to the #{@store}"
```
## Style/WhenThen

## Style/WhileUntilDo

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | - |

Checks for uses of `do` in multi-line `while/until` statements.

## Style/WhileUntilModifier

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 0.30 |

Checks for while and until statements that would fit on one line
if written as a modifier while/until. The maximum line length is
configured in the `Layout/LineLength` cop.

### Examples

```
# bad
while x < 10
 x += 1
end
# good
x += 1 while x < 10
# good
while x < 10
 y += 1 if x.odd?
end
# bad
until x > 10
 x += 1
end
# good
x += 1 until x > 10
# good
until x > 10
 y += 1 unless x.even?
end
# bad
x += 100 while x < 500 # a long comment that makes code too long if it were a single line
# good
while x < 500 # a long comment that makes code too long if it were a single line
 x += 100
end
```
## Style/WordArray

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | Yes | Always | 0.9 | 1.19 |

Checks for array literals made up of word-like strings, that are not using the %w() syntax.

Alternatively, it can check for uses of the %w() syntax, in projects which do not want to include that syntax.

| When using the `percent`style, %w() arrays containing a space
will be registered as offenses. |

Configuration option: MinSize
If set, arrays with fewer elements than this value will not trigger the
cop. For example, a `MinSize` of `3` will not enforce a style on an
array of 2 or fewer elements.

### Examples

#### EnforcedStyle: percent (default)

```
# good
%w[foo bar baz]
# bad
['foo', 'bar', 'baz']
# bad (contains spaces)
%w[foo\ bar baz\ quux]
# bad
[
 ['one', 'One'],
 ['two', 'Two']
]
# good
[
 %w[one One],
 %w[two Two]
]
# good (2d array containing spaces)
[
 ['one', 'One'],
 ['two', 'Two'],
 ['forty two', 'Forty Two']
]
```
### Configurable attributes

| Name | Default value | Configurable values |
|---|---|---|
| EnforcedStyle |
 |
 |
| MinSize |
 | Integer |
| WordRegex |
 |

## Style/YAMLFileRead

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Pending | Yes | Always | 1.53 | - |

Checks for the use of `YAML.load`, `YAML.safe_load`, and `YAML.parse` with
`File.read` argument.

| `YAML.safe_load_file`was introduced in Ruby 3.0. |

## Style/YodaCondition

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.49 | 0.75 |

Enforces or forbids Yoda conditions,
i.e. comparison operations where the order of expression is reversed.
eg. `5 == x`

### Safety

This cop is unsafe because comparison operators can be defined differently on different classes, and are not guaranteed to have the same result if reversed.

For example:

```
class MyKlass
 def ==(other)
 true
 end
end
obj = MyKlass.new
obj == 'string' #=> true
'string' == obj #=> false
```
### Examples

#### EnforcedStyle: forbid_for_all_comparison_operators (default)

```
# bad
99 == foo
"bar" != foo
42 >= foo
10 < bar
99 == CONST
# good
foo == 99
foo == "bar"
foo <= 42
bar > 10
CONST == 99
"#{interpolation}" == foo
/#{interpolation}/ == foo
```
#### EnforcedStyle: forbid_for_equality_operators_only

```
# bad
99 == foo
"bar" != foo
# good
99 >= foo
3 < a && a < 5
```
## Style/YodaExpression

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Disabled | No | Always (Unsafe) | 1.42 | 1.43 |

Forbids Yoda expressions, i.e. binary operations (using `*`, `+`, `&`, `|`,
and `^` operators) where the order of expression is reversed, eg. `1 + x`.
This cop complements `Style/YodaCondition` cop, which has a similar purpose.

This cop is disabled by default to respect user intentions such as:

`config.server_port = 9000 + ENV["TEST_ENV_NUMBER"].to_i`## Style/ZeroLengthPredicate

| Enabled by default | Safe | Supports autocorrection | Version Added | Version Changed |
|---|---|---|---|---|
| Enabled | No | Always (Unsafe) | 0.37 | 0.39 |

Checks for numeric comparisons that can be replaced
by a predicate method, such as `receiver.length == 0`,
`receiver.length > 0`, and `receiver.length != 0`,
`receiver.length < 1` and `receiver.size == 0` that can be
replaced by `receiver.empty?` and `!receiver.empty?`.

| `File`,`Tempfile`,`StringIO`, and`File::Stat`do not have`empty?`so allow`size == 0`and`size.zero?`. Note that when a`File::Stat`object
is stored in a variable (e.g.`stat = File.stat(path); stat.size.zero?`),
the cop cannot detect the type and may still register a false positive. |

# Plugins

RuboCop can be extended with additional cops and formatters. There are many official plugins maintained by RuboCop’s team, as well as third-party ones.

## Loading Plugins

Since RuboCop 1.72, the recommended way to load extensions is through the
**plugin** system, which is based on lint_roller. Older extensions that haven’t
been updated yet can still be loaded via `require`.

Refer to the Plugin Migration Guide when
migrating existing configurations from `require` to `plugins`.

### Loading Plugin Extensions

| The plugin system was introduced in RuboCop 1.72. |

Besides the `--plugin` command line option you can also specify ruby
files that should be loaded with the optional `plugins` directive in the
`.rubocop.yml` file:

```
plugins:
 - rubocop-performance
```
Most extension gems with plugin support should work with the example above.

The following is an example of a plugin that is not published as a gem:

```
plugins:
 - rubocop-extension:
 require_path: path/to/extension/plugin
 plugin_class_name: RuboCop::Extension::Plugin
```
There are other ways to specify this. For more details, please refer to the lint_roller documentation.

### Loading Legacy Extensions via `require`

Besides the `--require` command line option you can also specify ruby
files that should be loaded with the optional `require` directive in the
`.rubocop.yml` file:

```
require:
 - ../my/custom/file.rb
 - rubocop-extension
```
The extension loading via `require` remains compatible with the pre-plugin inject method.
Since `require` is used for internal extensions such as custom cops and formatters,
there are no plans to remove it in the future.
However, it is recommended that publicly available extension cops in gems migrate to the plugin system.

| The paths are directly passed to `Kernel.require`. If your
extension file is not in`$LOAD_PATH`, you need to specify the path as
relative path prefixed with`./`explicitly or absolute path. Paths
starting with a`.`are resolved relative to`.rubocop.yml`.
If a path containing`-`is given, it will be used as is, but if we
cannot find the file to load, we will replace`-`with`/`and try it
again as when Bundler loads gems. |

## Plugin Suggestions

Depending on what gems you have in your bundle, RuboCop might suggest plugins
that can be added to provide further functionality. For instance, if you are using
`rspec` without the corresponding `rubocop-rspec` plugin, RuboCop will suggest
enabling it.

This message can be disabled by adding the following to your configuration:

```
AllCops:
 SuggestExtensions: false
```
Default extensions are suggested when `SuggestExtensions: true`.

You can also opt-out of suggestions for a particular plugin as so (unspecified plugins will continue to be notified, as appropriate):

```
AllCops:
 SuggestExtensions:
 rubocop-rake: false
```
## Custom Cops

You can configure the custom cops in your `.rubocop.yml` just like any
other cop.

If you’d like to create a plugin gem, you can use rubocop-extension-generator.

For plugin specifications, please refer to lint_roller.

See Development to learn how to implement a cop.

## Available Plugins

The main RuboCop gem focuses on the core Ruby language and doesn’t include functionality related to any external Ruby libraries/frameworks. There are, however, many RuboCop plugins dedicated to those and a few of them are maintained by RuboCop’s Core Team.

### Official Plugins

-
rubocop-performance - Performance optimization analysis
-
rubocop-rails - Rails-specific analysis
-
rubocop-rspec - RSpec-specific analysis
-
rubocop-minitest - Minitest-specific analysis
-
rubocop-rake - Rake-specific analysis
-
rubocop-sequel - Code style checking for Sequel gem
-
rubocop-thread_safety - Thread-safety analysis
-
rubocop-capybara - Capybara-specific analysis
-
rubocop-factory_bot - factory_bot-specific analysis
-
rubocop-rspec_rails - RSpec Rails-specific analysis
-
rubocop-i18n - i18n wrapper function analysis ( `gettext`and`rails-i18n`)

### Third-party Plugins

-
rubocop-require_tools - Dynamic analysis for missing `require`statements
-
cookstyle - Custom cops and config defaults for Chef Infra Cookbooks
-
rubocop-packaging - Upstream best practices and coding conventions for downstream (e.g. Debian packages) compatibility.
-
rubocop-sorbet - Sorbet-specific analysis
-
rubocop-graphql - GraphQL-specific analysis
-
rubocop-changed - Reduced CI time by analyzing only changed files
-
rubocop-sketchup - SketchUp Ruby API specific analysis

Any plugins missing? Send us a Pull Request!

# Plugin Configuration

If you’re developing a plugin, you can tie its configuration into RuboCop.

## Config Obsoletions

When a cop that has been released is later renamed or removed, or one of its parameters is, RuboCop can output error messages letting users know to update their configuration to the newest values. If any obsolete configurations are encountered, RuboCop will output an error message and quit.

You can tie your plugin into this system by creating your own `obsoletions.yml` file and letting RuboCop know where to find it:

`RuboCop::ConfigObsoletion.files << File.expand_path(filename)`There are currently three types of obsoletions that can be defined for cops:

-
`renamed`: A cop was changed to have a new name, or moved to a different department.
-
`removed`: A cop was deleted (usually this is configured with`alternatives`or a`reason`why it was removed).
-
`split`: A cop was removed and replaced with multiple other cops.

Two additional types are available to be defined for parameter changes. These configurations can apply to multiple cops and multiple parameters at the same time (so they are expressed in YAML as an array of hashes):

-
`changed_parameters`: A parameter has been renamed.
-
`changed_enforced_styles`: A previously accepted`EnforcedStyle`value has been changed or removed.

| Parameter obsoletions can be set with `severity: warning`to deprecate an old parameter but still accept it. RuboCop will output a warning but continue to run. |

### Example Obsoletion Configuration

See `config/obsoletion.yml` for more examples.

| All plural keys (e.g. `cops`,`parameters`,`alternatives`, etc.) can either take a single value or an array. |

```
renamed:
 Layout/AlignArguments: Layout/ArgumentAlignment
 Lint/BlockAlignment: Layout/BlockAlignment
removed:
 Layout/SpaceAfterControlKeyword:
 alternatives: Layout/SpaceAroundKeyword
 Lint/InvalidCharacterLiteral:
 reason: it was never actually triggered
split:
 Style/MethodMissing:
 alternatives:
 - Style/MethodMissingSuper
 - Style/MissingRespondToMissing
changed_parameters: # must be an array of hashes
 - cops:
 - Metrics/BlockLength
 - Metrics/MethodLength
 parameters: ExcludedMethods
 alternative: IgnoredMethods
 severity: warning
changed_enforced_styles: # must be an array of hashes
 - cops: Layout/IndentationConsistency
 parameters: EnforcedStyle
 value: rails
 reason: >
 `EnforcedStyle: rails` has been renamed to
 `EnforcedStyle: indented_internal_methods`
```

# Plugin Migration Guide

| The plugin system was introduced in RuboCop 1.72. |

An improved API for RuboCop extensions, called plugins, is being introduced to provide officially documented ways to extend RuboCop. This API is provided by the lint_roller library and is intended to replace previous undocumented and unofficial methods that emerged out of necessity.

As a user, you will need to adjust your `.rubocop.yml` to load extensions that support this new API differently than before.
As an extension developer, you will need to modify your extension accordingly, as explained below.

This page provides guidance on migrating to RuboCop extensions using plugins. It is divided into two main sections: one for plugin users and the other for plugin developers.

## For Plugin Users

Please update the RuboCop extension cops specified in configuration files like `.rubocop.yml` from `require` to `plugins`.

```
plugins:
 - rubocop-performance
```
Specifically, those that trigger warnings like the following should be updated:

```
$ bundle exec rubocop
rubocop-performance extension supports plugin, specify `plugins: rubocop-performance`
instead of `require: rubocop-performance` in /Users/user/src/github.com/rubocop/rubocop/.rubocop.yml.
For more information, see https://docs.rubocop.org/rubocop/plugin_migration_guide.html.
```
| Non-plugin extension Ruby files can continue to use `require`for file requiring. |

## For Plugin Developers

Converting an existing RuboCop extension using `Inject.defaults!` into a plugin can generally be done in the following three steps.

-
Create `Plugin`class
-
Update gemspec file
-
Remove `Inject.defaults!`code

The `rubocop-performance` extension gem is used as an example here.
Replace "performance" with the name of your extension as appropriate.

### 1. Create `Plugin` class

Prepare the plugin class. It is typically placed under the department directory, such as `rubocop/performance/plugin.rb`.

```
# frozen_string_literal: true
require 'lint_roller'
module RuboCop
 module Performance
 # A plugin that integrates RuboCop Performance with RuboCop's plugin system.
 class Plugin < LintRoller::Plugin
 def about
 LintRoller::About.new(
 name: 'rubocop-performance',
 version: VERSION,
 homepage: 'https://github.com/rubocop/rubocop-performance',
 description: 'A collection of RuboCop cops to check for performance optimizations in Ruby code.'
 )
 end
 def supported?(context)
 context.engine == :rubocop
 end
 def rules(_context)
 LintRoller::Rules.new(
 type: :path,
 config_format: :rubocop,
 value: Pathname.new(__dir__).join('../../../config/default.yml')
 )
 end
 end
 end
end
```
For details on configurations, refer to the lint_roller documentation: https://github.com/standardrb/lint_roller?tab=readme-ov-file#how-to-make-a-plugin

### 2. Update gemspec file

Set the plugin class name in `spec.metadata['default_lint_roller_plugin']` and add a runtime dependency on `lint_roller`.

```
spec.metadata['default_lint_roller_plugin'] = 'RuboCop::Performance::Plugin'
spec.add_dependency('lint_roller')
```
Update `rubocop` to a version that supports plugins or higher.

`spec.add_dependency('rubocop', '>= 1.72.0')`### 3. Replace `Inject.defaults!` code

Replace `rubocop/performance/inject`, where the `RuboCop::Performance::Inject` class is defined.

`$ rm 'rubocop/performance/inject'`Remove the call to `RuboCop::Performance::Inject.defaults!` and the related `require_relative` code.

```
require_relative 'rubocop/performance/inject'
RuboCop::Performance::Inject.defaults!
```
And use `require_relative` for the plugin class file instead of the above.

`require_relative 'rubocop/performance/plugin'`

# Integration with Other Tools

## Speeding up integrations

RuboCop integrates with many other tools, including editors which may call
`rubocop` repeatedly. Since `rubocop` has to require its entire environment on
each invocation, this can result in noticeable slowness.

The recommended solution is RuboCop’s built-in Server mode, which keeps a daemonized process running so that subsequent runs avoid the boot time overhead.

| RuboCop’s built-in caching should also be used to ensure that source files that have not been changed are not being re-evaluated unnecessarily. |

## Editor integration

RuboCop includes a built-in Language Server Protocol (LSP) server (since 1.53), which is the recommended way to integrate with any editor. LSP provides real-time linting, formatting, and autocorrection as you type.

The following Ruby language servers use RuboCop internally:

-
Ruby LSP — the recommended option for most editors
-
Solargraph — also provides code completion and documentation

Any editor with an LSP client can use either of these. The per-editor sections below cover the recommended setup for each.

For editors that don’t support LSP, RuboCop can also be invoked by shelling out to the
`rubocop` binary. In that case, use the Server mode to avoid
the startup overhead on each invocation.

### Visual Studio Code

Use the Ruby LSP extension. There is also the dedicated vscode-rubocop extension that uses RuboCop’s built-in LSP server directly.

### Emacs

Alternative non-LSP options: rubocop.el and flycheck.

### Vim / Neovim

Neovim has a built-in LSP client. For Vim, ALE supports both LSP and direct RuboCop invocation.

### RubyMine / IntelliJ IDEA

RuboCop support is built-in.

### Helix

Helix has built-in LSP support. For formatting, see the Formatter Configurations for RuboCop.

### Sublime Text

Use the LSP package. See the LSP documentation for setup.

Alternative non-LSP options: SublimeLinter-rubocop.

### Other Editors

Any editor with an LSP client can use RuboCop — see the language servers listed above.

### MCP

RuboCop also provides a built-in Model Context Protocol server (since 1.85) for integration with AI-powered tools. See MCP for details.

## Git pre-commit hook integration with overcommit

overcommit is a fully configurable and
extendable Git commit hook manager. To use RuboCop with overcommit, add the
following to your `.overcommit.yml` file:

```
PreCommit:
 RuboCop:
 enabled: true
```
## Git pre-commit hook integration with pre-commit

pre-commit is a framework for managing and maintaining
multi-language pre-commit hooks. To use RuboCop with pre-commit, add the
following to your `.pre-commit-config.yaml` file:

```
- repo: https://github.com/rubocop/rubocop
 rev: v1.89.0
 hooks:
 - id: rubocop
```
If your RuboCop configuration uses extensions, be sure to include the gems as
entries in `additional_dependencies`:

```
- repo: https://github.com/rubocop/rubocop
 rev: v1.89.0
 hooks:
 - id: rubocop
 additional_dependencies:
 - rubocop-rails
 - rubocop-rspec
```
RuboCop supports pre-commit hook integration since version 1.9.

## Git hook integration with hk

hk is a fast, polyglot Git hooks manager with built-in RuboCop support.

## Guard integration

If you’re fond of Guard you might like guard-rubocop. It allows you to automatically check Ruby code style with RuboCop when files are modified.

## Mega-Linter integration

You can use Mega-Linter to run RuboCop automatically on every PR, and also lint all file types detected in your repository.

Please follow the installation instructions to activate RuboCop without any additional configuration.

Mega-Linter’s Ruby flavor is optimized for Ruby linting.

## Rake integration

To use RuboCop in your `Rakefile` add the following:

```
require 'rubocop/rake_task'
RuboCop::RakeTask.new
```
If you run `rake -T`, the following two RuboCop tasks should show up:

```
$ rake rubocop # Run RuboCop
$ rake rubocop:autocorrect # Autocorrect RuboCop offenses
```
The above will use default values.

```
require 'rubocop/rake_task'
desc 'Run RuboCop on the lib directory'
RuboCop::RakeTask.new(:rubocop) do |task|
 task.patterns = ['lib/**/*.rb']
 # only show the files with failures
 task.formatters = ['files']
 # don't abort rake on failure
 task.fail_on_error = false
end
```

# Automated Code Review

This section describes SaaS solutions that provide automated code reviews for Ruby based on RuboCop.

| The services are listed in alphabetical order. |

## Codacy

Codacy checks your code from style to security, duplication, complexity, and also integrates with coverage. Codacy is free for open source, and it provides RuboCop analysis out-of-the-box.

## Code Climate

Code Climate provides automated code review for test coverage, complexity, duplication, security, style, and more, and merge with confidence.

## CodeFactor

CodeFactor reports various code metrics like duplication, churn, and problems for code style, performance, complexity, and many others. CodeFactor is free for open source. It supports analysis and autocorrection for RuboCop.

## Codety

Codety detects code issues for 30+ programming languages and IaC frameworks. Codety Scanner is open source and it embeds 6,000+ code analysis rules (including RuboCop rules).

## Pronto

Pronto does quick automated code review of your changes. Created to be used on GitHub pull requests, but also works locally and integrates with GitLab and Bitbucket.

## ReviewDog

ReviewDog is similar to Pronto but with better support for GitHub Actions.

# Development

This guide covers developing new cops for RuboCop. It walks through:

-
Scaffolding a cop from a template
-
Understanding the AST that RuboCop operates on
-
Implementing detection with node patterns and callbacks
-
Autocorrect to fix offenses automatically
-
Constraining cops by Ruby or gem version
-
Common patterns and conventions
-
Testing your cop
-
Documenting the cop for users

## Create a new cop

| Clone the repository and run `bundle install`if not done yet.
The following rake task can only be run inside the rubocop project directory itself. |

Use the bundled rake task `new_cop` to generate a cop template:

```
$ bundle exec rake 'new_cop[Department/Name]'
[create] lib/rubocop/cop/department/name.rb
[create] spec/rubocop/cop/department/name_spec.rb
[modify] lib/rubocop/cop/department.rb - `register_cop :Name, "#{__dir__}/department/name"` was injected.
[modify] A configuration for the cop is added into config/default.yml.
Do 4 steps:
 1. Modify the description of Department/Name in config/default.yml
 2. Implement your new cop in the generated file!
 3. Commit your new cop with a message such as
 e.g. "Add new `Department/Name` cop"
 4. Run `bundle exec rake changelog:new` to generate a changelog entry
 for your new cop.
```
### Cop lazy loading

Cops are not required from `lib/rubocop.rb`. Instead, each department has a module file
(e.g. `lib/rubocop/cop/style.rb`) that registers its cops with the global registry through
`RuboCop::Cop::LazyLoader#register_cop`:

```
module RuboCop
 module Cop
 # Cops for the `Style` department. ...
 module Style
 extend LazyLoader
 register_cop :HashSyntax, "#{__dir__}/style/hash_syntax"
 # ...
 end
 end
end
```
A registered cop appears in the registry by name right away, but its file is loaded only
when the cop class is needed, e.g. when the cop is enabled for an inspection run.
This significantly speeds up loading RuboCop, since a typical run loads only the enabled cops.
The path passed to `register_cop` must be absolute so that loading does not depend on `$LOAD_PATH`;
interpolating `__dir__` as in the example above is the intended usage. Plugin authors can adopt
the same mechanism for their own department modules.

The `new_cop` rake task maintains the `register_cop` directives for you, so usually there is
no need to edit the department module file by hand.

## Understanding the AST

RuboCop uses the parser library to create the Abstract Syntax Tree (AST) representation of the code.

You can install `parser` gem and use `ruby-parse` command line utility to check
what the AST looks like in the output.

`$ gem install parser`And then try to parse a simple integer representation with `ruby-parse`:

```
$ ruby-parse -e '1'
(int 1)
```
Each expression surrounded by parentheses represents a node in the AST. The first element is the node type and the tail contains the children with all information needed to represent the code.

Here’s another example - a local variable `name` being assigned the
string value "John":

```
$ ruby-parse -e 'name = "John"'
(lvasgn :name
 (str "John"))
```
### Inspecting the AST

Let’s imagine we want to simplify statements from `!array.empty?` to
`array.any?`:

First, check what the bad code returns in the Abstract Syntax Tree representation.

```
$ ruby-parse -e '!array.empty?'
(send
 (send
 (send nil :array) :empty?) :!)
```
Now, it’s time to debug our expression using the REPL from RuboCop:

`$ bin/console`First we need to declare the code that we want to match, and use the ProcessedSource that is a simple wrap to make the parser interpret the code and build the AST:

```
code = '!something.empty?'
source = RuboCop::ProcessedSource.new(code, RUBY_VERSION.to_f)
node = source.ast
# => s(:send, s(:send, s(:send, nil, :something), :empty?), :!)
```
The node has a few attributes that can be useful in the journey:

```
node.type # => :send
node.children # => [s(:send, s(:send, nil, :something), :empty?), :!]
node.source # => "!something.empty?"
```
## Implementing the cop

### Writing node pattern rules

| You can write cops without using `NodePattern`(and many older cops don’t use it), but it
generally simplifies the code a lot, as manual node matching and destructuring can be
quite verbose. |

Now that you’re familiar with AST, you can learn a bit about the node pattern and use patterns to match with specific nodes that you want to match.

You can learn more about Node Pattern here.

Alias `NodePattern` to `RuboCop::AST::NodePattern` to make it easier to use:

`NodePattern = RuboCop::AST::NodePattern`Node pattern matches something very similar to the current output from AST representation, then let’s start with something very generic:

`NodePattern.new('send').match(node) # => true`It matches because the root is a `send` type. Now let’s match it deeply using
parentheses to define details for sub-nodes. If you don’t care about what an internal
node is, you can use `...` to skip it and just consider "a node".

```
NodePattern.new('(send ...)').match(node) # => true
NodePattern.new('(send (send ...) :!)').match(node) # => true
NodePattern.new('(send (send (send ...) :empty?) :!)').match(node) # => true
```
Sometimes it’s hard to comprehend complex expressions you’re building with the
pattern, then, if you got lost with the node pattern parens surrounding deeply,
try to use the `$` to capture the internal expression and check exactly each
piece of the expression:

`NodePattern.new('(send (send (send $...) :empty?) :!)').match(node) # => [nil, :something]`It’s not needed to strictly receive a send in the internal node because maybe it can also be a literal array like:

`![].empty?`The code above has the following representation:

`=> s(:send, s(:send, s(:array), :empty?), :!)`It’s possible to skip the internal node with `...` to make sure that it’s just
another internal node:

`NodePattern.new('(send (send (...) :empty?) :!)').match(node) # => true`In other words, it says: "Match code calling `!<expression>.empty?`".

Great! Now, let’s implement our cop to simplify such statements:

`$ rake 'new_cop[Style/SimplifyNotEmptyWithAny]'`After the cop scaffold is generated, change the node matcher to match with the expression achieved previously:

```
def_node_matcher :not_empty_call?, <<~PATTERN
 (send (send $(...) :empty?) :!)
PATTERN
```
Note that we added a `$` sign to capture the "expression" in `!<expression>.empty?`,
it will become useful later.

Get yourself familiar with the AST node hooks that
`parser`
and `rubocop-ast`
provide.

As it starts with a `send` type, it’s needed to implement the `on_send` method, as the
cop scaffold already suggested:

```
def on_send(node)
 return unless not_empty_call?(node)
 add_offense(node)
end
```
The `on_send` callback is the most used and can be optimized by restricting the acceptable
method names with a constant `RESTRICT_ON_SEND`.

The final cop code will look something like this:

```
module RuboCop
 module Cop
 module Style
 # `array.any?` is a simplified way to say `!array.empty?`
 #
 # @example
 # # bad
 # !array.empty?
 #
 # # good
 # array.any?
 #
 class SimplifyNotEmptyWithAny < Base
 MSG = 'Use `.any?` and remove the negation part.'
 RESTRICT_ON_SEND = [:!].freeze # optimization: don't call `on_send` unless
 # the method name is in this list
 def_node_matcher :not_empty_call?, <<~PATTERN
 (send (send $(...) :empty?) :!)
 PATTERN
 def on_send(node)
 return unless not_empty_call?(node)
 add_offense(node)
 end
 end
 end
 end
end
```
Callback ordering: `on_send` is called on a node *before* the `on_<type>`
callbacks for its children. There is also an `after_send` callback that runs
*after* children are processed. Every node type has a corresponding
`after_<type>` callback (except types that never have children).

Update the spec to cover the expected syntax:

```
describe RuboCop::Cop::Style::SimplifyNotEmptyWithAny, :config do
 it 'registers an offense when using `!a.empty?`' do
 expect_offense(<<~RUBY)
 !array.empty?
 ^^^^^^^^^^^^^ Use `.any?` and remove the negation part.
 RUBY
 end
 it 'does not register an offense when using `.any?` or `.empty?`' do
 expect_no_offenses(<<~RUBY)
 array.any?
 array.empty?
 RUBY
 end
end
```
If your code has variables of different lengths, you can use the following markers to format your template by passing the variables as keyword arguments:

-
`%{foo}`: Interpolates`foo`
-
`^{foo}`: Inserts`'^' * foo.size`for dynamic offense range length
-
`_{foo}`: Inserts`' ' * foo.size`for dynamic offense range indentation

You can also abbreviate offense messages with `[…]`.

```
%w[raise fail].each do |keyword|
 expect_offense(<<~RUBY, keyword: keyword)
 %{keyword}(RuntimeError, msg)
 ^{keyword}^^^^^^^^^^^^^^^^^^^ Redundant `RuntimeError` argument [...]
 RUBY
%w[has_one has_many].each do |type|
 expect_offense(<<~RUBY, type: type)
 class Book
 %{type} :chapter, foreign_key: 'book_id'
 _{type} ^^^^^^^^^^^^^^^^^^^^^^ Specifying the default [...]
 end
 RUBY
end
```
### Autocorrect

Autocorrect lets cops automatically fix the offenses they detect.
It’s necessary to `extend AutoCorrector`.
The method `add_offense` yields a corrector object that is a thin wrapper on
parser’s TreeRewriter
to which you can give instructions about what to do with the
offensive node.

Let’s start with a simple spec to cover it:

```
it 'corrects `!a.empty?`' do
 expect_offense(<<~RUBY)
 !array.empty?
 ^^^^^^^^^^^^^ Use `.any?` and remove the negation part.
 RUBY
 expect_correction(<<~RUBY)
 array.any?
 RUBY
end
```
And then add the autocorrecting block on the cop side:

```
extend AutoCorrector
def on_send(node)
 expression = not_empty_call?(node)
 return unless expression
 add_offense(node) do |corrector|
 corrector.replace(node, "#{expression.source}.any?")
 end
end
```
The corrector allows you to `insert_after`, `insert_before`, `wrap` or
`replace` a specific node or in any specific range of the code.

Range can be determined on `node.location` where it brings specific
ranges for expression or other internal information that the node holds.

#### Preventing clobbering

The corrector detects and prevents correcting overlapping nodes, to prevent one correction from clobbering another.
Supporting nested corrections is done by taking multiple passes, and skipping corrections for nested nodes.
This can be implemented using the `IgnoredNode` module:

```
 extend AutoCorrector
+include IgnoredNode
 def on_send(node)
 return unless some_condition?(node)
 add_offense(node) do |corrector|
+ next if part_of_ignored_node?(node)
+
 corrector.replace(node, "...")
 end
+
+ ignore_node(node)
 end
```
This works because file correction is implemented by repeating investigation and correction until the file no longer requires correction, meaning all nested nodes will eventually be processed.

Note that `expect_correction` in `Cop` specs asserts the result after all passes.

### Limit by Ruby or gem versions

Some cops apply changes that only apply in particular contexts, such as if the user has a minimum Ruby version. There are helpers that let you constrain your cops automatically, to only run where applicable.

#### Requiring a minimum Ruby version

If your cop uses new Ruby syntax or standard library APIs, it should only register offenses if the user has the proper target Ruby version, which you can require with `TargetRubyVersion#minimum_target_ruby_version`.

For example, the `Performance/SelectMap` cop requires Ruby 2.7, which introduced `Enumerable#filter_map`:

```
class RuboCop::Cop::Performance::SelectMap < Base
 extend TargetRubyVersion
 minimum_target_ruby_version 2.7
 # ...
end
```
This cop won’t register offenses on Ruby 2.6 or older.

#### Requiring a maximum Ruby version

Mirroring `minimum_target_ruby_version`, you can also specify a maximum Ruby version your cop should analyze.

For example, the `Lint/CircularArgumentReference` cop only runs when analyzing code for Ruby before 2.7. The code it looks for can never be written in more recent Rubies — it would be a syntax error:

```
class RuboCop::Cop::Lint::CircularArgumentReference < Base
 extend TargetRubyVersion
 maximum_target_ruby_version 2.6
 # ...
end
```
#### Requiring a gem

If your cop depends on the presence of a gem, you can declare that with `RuboCop::Cop::Base.requires_gem`.

For example, to declare that `MyCop` should only apply if the bundle is using `my-gem` with a version between `1.2.3` and `4.5.6`:

```
class MyCop < Base
 requires_gem "my-gem", ">= 1.2.3", "< 4.5.6"
 # ...
end
```
You can specify any gem requirement using the same syntax as your `Gemfile`.

You can also handle multiple versions of a gem with `target_gem_version`. It behaves similarly to `target_ruby_version`, allowing you to inspect a gem version at runtime:

```
class MyCop < Base
 requires_gem "my-gem"
 def on_send(node)
 if target_gem_version("my-gem") < "2.0"
 # ...
 else
 # ...
 end
 end
end
```
When writing tests, you can specify the gem version to run your example against through the `gem_versions` RSpec helper:

```
describe RuboCop::Cop::Style::MyCop, :config do
 context 'when `my-gem` is at version `1.X`' do
 let(:gem_versions) { { 'my-gem' => '1.0.0' } }
 it 'registers no offense' do
 expect_no_offenses(<<~RUBY)
 MyGem.foo
 RUBY
 end
 end
 context 'when `my-gem` is at version `2.X`' do
 let(:gem_versions) { { 'my-gem' => '2.0.0' } }
 it 'registers an offense' do
 expect_offense(<<~RUBY)
 MyGem.foo
 ^^^^^^^^^ Instead of `foo`, use the newer `bar` method.
 RUBY
 end
 end
end
```
#### Special case: Rails

Historically, many cops in `rubocop-rails` aren’t actually specific to Rails itself, but some of its components (e.g., Active Support). These dependencies are declared with `TargetRailsVersion.minimum_target_rails_version`.

For example, the `Rails/Pluck` cop requires Active Support 6.0, which introduces `Enumerable#pluck`:

```
class RuboCop::Cop::Rails::Pluck < Base
 extend TargetRailsVersion
 minimum_target_rails_version 6.0
 #...
end
```
### Configuration

Each cop can hold a configuration and you can refer to `cop_config` in the
instance and it will bring a hash with options declared in the `.rubocop.yml`
file.

For example, let’s imagine we want to make the replacement method configurable,
so it works with a method other than `.any?`:

```
Style/SimplifyNotEmptyWithAny:
 Enabled: true
 ReplaceAnyWith: "size > 0"
```
And then in the autocorrect method, you just need to use `cop_config`:

```
def on_send(node)
 expression = not_empty_call?(node)
 return unless expression
 add_offense(node) do |corrector|
 replacement = cop_config['ReplaceAnyWith'] || 'any?'
 corrector.replace(node, "#{expression.source}.#{replacement}")
 end
end
```
## Common patterns

This section covers conventions and patterns you should follow when writing cops.

### Handling safe navigation (`&.`)

If your cop defines `on_send`, you should almost always also handle the
safe navigation operator by aliasing `on_csend`:

```
def on_send(node)
 # ...
end
alias on_csend on_send
```
Without this, your cop will silently ignore code like `foo&.bar`. Only skip
this if the cop explicitly does not apply to safe navigation.

### Handling numbered-parameter and `it`-parameter blocks

Similarly, if your cop defines `on_block`, alias the numbered-parameter and
`it`-parameter variants:

```
def on_block(node)
 # ...
end
alias on_numblock on_block
alias on_itblock on_block
```
### Using `RESTRICT_ON_SEND`

When your cop uses `on_send`, define a `RESTRICT_ON_SEND` constant listing
the method names the cop cares about. This is a performance optimization —
RuboCop will skip calling `on_send` entirely for method names not in the list:

`RESTRICT_ON_SEND = %i[bad_method other_bad_method].freeze`### Documenting node matchers with `@!method`

Every `def_node_matcher` and `def_node_search` should have a `@!method`
YARD tag above it so that documentation tools can find the generated method:

```
# @!method not_empty_call?(node)
def_node_matcher :not_empty_call?, <<~PATTERN
 (send (send $(...) :empty?) :!)
PATTERN
```
## Running tests

RuboCop supports two parser engines: the Parser gem and Prism. By default, tests are executed with the Parser:

`$ bundle exec rake spec`To run all tests with Prism, use `bundle exec rake prism_spec`. To run a single spec file with Prism,
set the `PARSER_ENGINE` environment variable:

`$ PARSER_ENGINE=parser_prism bundle exec rspec spec/rubocop/cop/style/hash_syntax_spec.rb``bundle exec rake` runs the full CI suite: specs for both parser engines, self-linting, and documentation checks. Always run this before submitting a PR.

## Documentation

Every cop needs YARD documentation with examples directly in the source file.
The CI `documentation_syntax_check` task parses every `@example` block, so all
examples must contain valid Ruby syntax.

Key rules:

-
The first line of the YARD comment must be a complete sentence starting with a verb and ending with a period (e.g., "Checks for …", "Enforces …").
-
Every cop must have at least one `# bad`/`# good`example pair.
-
For each `SupportedStyle`or unique configuration key, add a separate`@example`block.
-
Mark the default style with `(default)`.
-
List config keys in alphabetical order.

```
module RuboCop
 module Cop
 module Style
 # Simplifies `!array.empty?` to `array.any?`.
 #
 # @example EnforcedStyle: any? (default)
 # # bad
 # !array.empty?
 #
 # # good
 # array.any?
 #
 # @example EnforcedStyle: size
 # # bad
 # !array.empty?
 #
 # # good
 # array.size > 0
 #
 class SimplifyNotEmptyWithAny < Base
 # ...
```
Add additional `@example` blocks following the same pattern for each style or
config value.

## Testing your cop in a real codebase

It’s generally good practice to check if your cop is working properly over a significant codebase (e.g. Rails or some big project you’re working on) to guarantee it’s working in a range of different syntaxes.

There are several ways to do this. Two common approaches:

-
From within your local `rubocop`repo, run`exe/rubocop ~/your/other/codebase`.
-
From within the other codebase’s `Gemfile`, set a path to your local repo like this:`gem 'rubocop', path: '/full/path/to/rubocop'`. Then run`rubocop`within your codebase.

With approach #2, you can use local versions of RuboCop extension repos such as `rubocop-rspec` as well.

| Use `--only`to run just your cop and avoid noise from other cops: |

`$ rubocop --only Style/SimplifyNotEmptyWithAny`## Custom formatters

Beyond cops, RuboCop can also be extended with custom output formatters.

### Creating a custom formatter

To implement a custom formatter, you need to subclass
`RuboCop::Formatter::BaseFormatter` and override some methods,
or implement all formatter API methods by duck typing.

Please see the documents below for more formatter API details.

### Using a custom formatter from the command line

You can tell RuboCop to use your custom formatter with a combination of
`--format` and `--require` option.
For example, when you have defined `MyCustomFormatter` in
`./path/to/my_custom_formatter.rb`, you would type this command:

`$ rubocop --require ./path/to/my_custom_formatter --format MyCustomFormatter`## Template support

RuboCop can also analyze Ruby embedded in templates (ERB, Haml, Slim, etc.) through a Ruby extractor API.

A template file contains multiple embedded Ruby snippets, unlike a regular
Ruby file. RuboCop solves this with `RuboCop::Runner.ruby_extractors` — a
list of callable extractors that plugins can prepend to.

A Ruby extractor takes a `RuboCop::ProcessedSource` and returns either an
`Array` of `Hash`-es containing extracted Ruby source and offsets, or `nil`
if the file is not relevant.

`ruby_extractor.call(processed_source)`An example returned value from a Ruby extractor would be as follows:

```
[
 {
 offset: 2,
 processed_source: #<RuboCop::ProcessedSource>
 },
 {
 offset: 10,
 processed_source: #<RuboCop::ProcessedSource>
 }
]
```
On the extension side, the code would be something like this:

`RuboCop::Runner.ruby_extractors.unshift(ruby_extractor)``RuboCop::Runner.ruby_extractors` is processed from the beginning and ends when one of them returns a non-nil value. By default, there is a Ruby extractor that returns the given Ruby source code with offset 0, so you can unshift any Ruby extractor before it.

| This is still an experimental feature and may change in the future. |

# Contributing

## Issues

Report issues and suggest features and improvements on the GitHub issue tracker.

If you want to file an issue (bug report or feature request), please provide all the necessary info listed in our issue reporting template (it’s loaded automatically when you create a new GitHub issue).

| Please, don’t ask support-like questions on the issue tracker (e.g. "How can I configure RuboCop to do X?") - use the support channels instead. |

## Code contributions

Patches in any form are always welcome! GitHub pull requests are even better! :-)

Before submitting a pull request make sure all tests are passing and that your patch is in line with the contribution guidelines.

### Quick start

```
$ git clone https://github.com/rubocop/rubocop.git
$ cd rubocop
$ bundle install
$ bundle exec rake # run the full CI suite
```
A typical contribution workflow:

-
Pick an issue from the issue tracker (look for `good first issue`or`help wanted`labels).
-
Create a feature branch: `git checkout -b fix-1234-short-description`.
-
Make your changes and add tests (see Development).
-
Run the full test suite: `bundle exec rake`.
-
Add a changelog entry: `bundle exec rake changelog:fix`(or`changelog:new`/`changelog:change`).
-
Open a pull request.

## Documentation

Good documentation is just as important as good code. Please, help us improve RuboCop’s documentation.

You should also check out the cop documentation section of the docs and consider adding or improving Cop descriptions.

### Working on the Docs

The manual is generated from the AsciiDoc files in the docs folder of RuboCop’s GitHub repo and is published to https://docs.rubocop.org. Antora is used to convert the manual into HTML. The filesystem layout is described at https://docs.antora.org/antora/3.1/standard-directories/.

To make changes to the manual you simply have to change the files under `docs`.
The manual will be regenerated manually periodically.

The site configuration and build infrastructure live in a separate repo: docs.rubocop.org.

#### Installing Antora

| The instructions here assume you already have (the right version of) node.js installed. |

Installing Antora is super simple — just clone the docs site repo and run:

```
$ cd docs.rubocop.org
$ make install
```
Check out the detailed installation instructions if you run into any problems.

#### Building the Site

You can build the documentation locally from the docs.rubocop.org repo.

```
$ cd docs.rubocop.org
$ make build
```
| You can preview your changes by opening `build/site/index.html`in your favourite browser. |

The site is automatically deployed to GitHub Pages via a GitHub Action on push to `master`.

If you want to make changes to the manual’s page structure you’ll have to edit nav.adoc.

### Working on rubocop.org

The rubocop.org landing page is a separate project, built with Bridgetown and Tailwind CSS. Its source lives at github.com/rubocop/rubocop.org.

To run the site locally you’ll need Ruby (>= 3.2) and Node.js (>= 20):

```
$ cd rubocop.org
$ bundle install && npm install
$ bin/bridgetown start
```
Then visit localhost:4000 to preview your changes.

The site is automatically deployed to GitHub Pages on push to `main`.

## Funding

While RuboCop is free software and will always be, the project would benefit immensely from some funding. Raising a monthly budget of a couple of thousand dollars would make it possible to pay people to work on certain complex features, fund other development related stuff (e.g. hardware, conference trips) and so on. Raising a monthly budget of over $5000 would open the possibility of someone working full-time on the project which would speed up the pace of development significantly.

We welcome both individual and corporate sponsors! We also offer a wide array of funding channels to account for your preferences (although currently Open Collective is our preferred funding platform).

If you’re working in a company that’s making significant use of RuboCop we’d appreciate it if you suggest to your company to become a RuboCop sponsor.

You can support the development of RuboCop via GitHub Sponsors, Patreon, PayPal, and Open Collective.

# Support

RuboCop currently has several official & unofficial support channels. For questions, suggestions, and support refer to one of them.

| Please, don’t use the support channels to report issues, as this makes them harder to track. |

## GitHub Discussions

These days our primary discussion board is here.

It’s a great place to share ideas, help other RuboCop users and just be up-to-date with interesting developments.

## Discord

If you’re into chat you can drop by RuboCop’s Discord server. You can often find some RuboCop maintainers there and get some interesting news from the project’s kitchen.

## Mailing List

The official mailing list is hosted at Google Groups. It’s a low-traffic list, so don’t be too hesitant to subscribe.

## StackOverflow

We’re also encouraging users to ask RuboCop-related questions on StackOverflow.

When doing so you should use the
RuboCop tag (ideally combined
with the tag `ruby`).

# Versioning

RuboCop is stable between major versions, both in terms of API and cop configuration. We aim to ease the maintenance of RuboCop extensions (by keeping the API stable) and the upgrades between RuboCop releases (by not enabling new cops and changing the configuration of existing cops). All big (breaking) changes are reserved for major releases.

## Release Policy

We’re following Semantic Versioning. API compatibility between major releases is a big concern, as there are many RuboCop extensions that can be affected by breaking API changes.

The development cycle for the next minor (feature) release starts immediately after the previous one has been shipped. Bug-fix (point) releases (if any) address only serious bugs and never contain new features.

Here are a few examples:

-
1.1.0 - Feature release
-
1.1.1 - Bug-fix release
-
1.1.2 - Bug-fix release
-
1.2.0 - Feature release

### What goes in each release type

| Change type | Release |
|---|---|
| Bug fixes to existing cops | Patch |
| New cops (added as | Minor |
| New configuration options | Minor |
| New CLI flags or formatters | Minor |
| Renaming or moving cops (with backwards-compatible obsoletions) | Minor |
| Dropping runtime Ruby version support | Minor |
| Enabling pending cops by default |
 |
| Changing cop defaults (e.g. |
 |
| Removing cops without a replacement |
 |
| Breaking changes to the Ruby API |
 |
| Dropping analysis support for a Ruby version |
 |

| Dropping runtimesupport for a particular Ruby version is not considered a breaking change,
as it doesn’t affect clients in any way. They are simply restricted to the last version of
RuboCop supporting their Ruby runtime. |

| Prior to RuboCop 1.0 bumps of the minor (second) version number were considered major releases and always included new features and/or changes to existing features. |

## Pending Cops

In the early versions of RuboCop a common source of frustration was that new cops were added to pretty much every release, and as they were enabled by default, every upgrade resulted in broken CI builds and trying to figure out what exactly was changed. After considering many options to address this eventually we opted for an approach that limits these types of changes to major RuboCop releases.

Now new cops introduced between major versions are set to a special pending status and are not enabled by default. A warning is emitted if such cops are not explicitly enabled or disabled in the user configuration. Here’s one such message:

The following cops were added to RuboCop, but are not configured. Please set Enabled to either `true` or `false` in your `.rubocop.yml` file: - Style/HashEachMethods (0.80) - Style/HashTransformKeys (0.80) - Style/HashTransformValues (0.80) For more information: https://docs.rubocop.org/rubocop/versioning.html

You can see that 3 new cops were added in RuboCop 0.80 and it’s up to you to decide if you want to enable or disable them.

| Occasionally, some new cops will be introduced as disabledby
default. Usually, this means that we believe that the cop is useful,
but not for everyone. Typical cases might be the enforcement of
programming styles that are not very common in the wild, or cops that
yield too many false positives (so you’d run them manually from time
to time, instead of running them all the time). |

### Enabling/Disabling Pending Cops in Bulk

To suppress this message set `NewCops` to either `enable` or `disable` in your `.rubocop.yml` file.
You can use the following configuration or the `--enable-pending-cops` command-line option to enable all pending cops in bulk:

```
AllCops:
 NewCops: enable
```
Alternatively, you can use the following configuration or the `--disable-pending-cops` command-line option to disable all pending cops in bulk:

```
AllCops:
 NewCops: disable
```
| The command-line option takes precedence over the `.rubocop.yml`file. |

### Enabling/Disabling Pending Cops per Department

A single version in `AllCops` cannot express which pending cops to enable, because extension gems follow their own version schemes.
Instead, `NewCops` can be set for a department, in which case it takes precedence over the `AllCops` setting for the cops of that department.
In addition to `enable`, `disable`, and `pending`, a department accepts a version, which enables all pending cops of the department
that were added in the specified version or earlier:

```
AllCops:
 NewCops: disable # keywords only, as before
Lint:
 NewCops: enable # a department overrides `AllCops`
Style:
 NewCops: '1.19' # enables pending cops with `VersionAdded` <= 1.19
Minitest:
 NewCops: '0.10' # an extension is pinned to its own version
```
Pending cops added in a later version keep emitting the warning, so you can review them and bump the version at your own pace. Since each department is pinned to its own version, this also works for extensions whose version schemes differ from RuboCop’s.

| Quote the version to avoid YAML interpreting it as a number (for example, an unquoted `1.20`is read as`1.2`). |

### Enabling/Disabling Individual Pending Cops

Finally, you can enable/disable individual pending cops by setting their `Enabled` configuration to either `true` or `false` in your `.rubocop.yml` file:

`Style/ANewCop` is an example of a newly added pending cop:

```
Style/ANewCop:
 Enabled: true
```
or

```
Style/ANewCop:
 Enabled: false
```
| On major RuboCop version updates (e.g. 1.0 → 2.0), allpending cops are enabled in bulk. |

# v1 Upgrade Notes

## Cop upgrade guide

Your custom cops should continue to work in v1.

Nevertheless, it is suggested that you tweak them to use the v1 API by following these steps:

1) Your class should inherit from `RuboCop::Cop::Base` instead of `RuboCop::Cop::Cop`.

2) Locate your calls to `add_offense` and make sure that you pass as the first argument either an `AST::Node`, a `::Parser::Source::Comment` or a `::Parser::Source::Range`, and no `location:` named parameter.

#### Example:

```
# Before
class MySillyCop < Cop
 def on_send(node)
 if node.method_name == :-
 add_offense(node, location: :selector, message: "Be positive")
 end
 end
end
# After
class MySillyCop < Base
 def on_send(node)
 if node.method_name == :-
 add_offense(node.loc.selector, message: "Be positive")
 end
 end
end
```
### If your class supports autocorrection

Your class must `extend AutoCorrector`.

The `corrector` is now yielded from `add_offense`. Move the code of your method `autocorrect` into that block and do not wrap your correction in a lambda. `Corrector`s are more powerful and can now be `merge`d.

#### Example:

```
# Before
class MySillyCorrectingCop < Cop
 def on_send(node)
 if node.method_name == :-
 add_offense(node, location: :selector, message: 'Be positive')
 end
 end
 def autocorrect(node)
 lambda do |corrector|
 corrector.replace(node.loc.selector, '+')
 end
 end
end
# After
class MySillyCorrectingCop < Base
 extend AutoCorrector
 def on_send(node)
 if node.method_name == :-
 add_offense(node.loc.selector, message: 'Be positive') do |corrector|
 corrector.replace(node.loc.selector, '+')
 end
 end
 end
end
```
### Instance variables

Do not use RuboCop’s internal instance variables. If you used `@processed_source`, use `processed_source`. If you have a need to access an instance variable, open an issue with your use case.

By default, a Cop instance will be called only once for a given `processed_source`, so instance variables will be uninitialized when the investigation starts. Using `@cache ||= …` is fine. If you want to initialize some instance variable, the callback `on_new_investigation` is the best place to do so.

```
class MyCachingCop < Base
 def on_send(node)
 if my_cached_data[node]
 @counts(node.method_name) += 1
 #...
 end
 end
 # One way:
 def my_cached_data
 @data ||= processed_source.comments.map { # ... }
 end
 # Another way:
 def on_new_investigation
 @counts = Hash.new(0)
 super # Be nice and call super for callback
 end
end
```
### Other API changes

If your cop uses `investigate`, `investigate_post_walk`, `join_force?`, or internal classes like `Corrector`, `Commissioner`, `Team`, these have changed. See the Detailed API Changes.

### Upgrading specs

It is highly recommended you use `expect_offense` / `expect_correction` / `expect_no_offenses` in your specs, e.g.:

```
require 'rubocop/rspec/support'
RSpec.describe RuboCop::Cop::Custom::MySillyCorrectingCop, :config do
 # No need for `let(:cop)`
 it 'is positive' do
 expect_offense(<<~RUBY)
 42 + 2 - 2
 ^ Be positive
 RUBY
 expect_correction(<<~RUBY)
 42 + 2 + 2
 RUBY
 end
 it 'does not register an offense for calls to `despair`' do
 expect_no_offenses(<<~RUBY)
 "don't".despair
 RUBY
 end
end
```
In the unlikely case where you use the class `RuboCop::Cop::Corrector` directly, it has changed a bit but you can ease your transition with `RuboCop::Cop::Legacy::Corrector` that is meant to be somewhat backwards compatible. You will need to `require 'rubocop/cop/legacy/corrector'`.

## Detailed API Changes

This section lists all changes (big or small) to the API. It is meant for maintainers of the nuts & bolts of RuboCop; most cop writers will not be impacted by these and are thus not the target audience.

### Base class

*Legacy*: Cops inherit from `Cop::Cop`.

*Current*: Cops inherit from `Cop::Base`. Having a different base class makes the implementation much cleaner and makes it easy to signal which API is being used. `Cop::Cop` inherits from `Cop::Base` and refines some methods for backward compatibility.

`add_offense` API

#### arguments

*Legacy:* interface allowed for a `node`, with an optional `location` (symbol or range) or a range with a mandatory range as the location. Some cops were abusing the `node` argument and passing very different things.

*Current:* pass a range (or node as a shortcut for `node.loc.expression`), no `location:`. No abuse tolerated.

#### deduping changes

Both dedupe on `range` and won’t process the duplicated offenses at all.

*Legacy:* if offenses on same `node` but different `range`: considered as multiple offenses but a single autocorrect call.

*Current:* not applicable and not needed with autocorrection’s API.

#### yield

Both yield under the same conditions (unless cop is disabled for that line), but:

*Legacy:* yields after offense added to `#offenses`

*Current:* yields before offense is added to `#offenses`.

Even the legacy mode yields a corrector, but if a developer uses it an error will be raised asking them to inherit from `Cop::Base` instead.

### Autocorrection

`#autocorrect`

*Legacy:* calls `autocorrect` unless it is disabled / autocorrect is off.

*Current:* yields a corrector unless it is disabled. The corrector will be ignored if autocorrecting is off, etc. No support for `autocorrect` method, but a warning is issued if that method is still defined.

#### Empty corrections

*Legacy:* `autocorrect` could return `nil` / `false` in cases where it couldn’t actually make a correction.

*Current:* No special API. Cases where no corrections are made are automatically detected.

#### Correction timing

*Legacy:* the lambda was called only later in the process, and only under specific conditions (if the autocorrect setting is turned on, etc.)

*Current:* correction is built immediately (assuming the cop isn’t disabled for the line) and applied later in the process.

#### Exception handling

Both: `Commissioner` will rescue all `StandardError`s during analysis (unless `option[:raise_error]`) and store a corresponding `ErrorWithAnalyzedFileLocation` in its error list. This is done when calling the cop’s `on_send` & al., or when calling `investigate` / `investigate_post_walk` callback.

*Legacy:* autocorrecting cops were treating errors differently depending on when they occurred. Some errors were silently ignored. Others were rescued as above. Others crashed. Some code in `Team` would rescue errors and add them to the list of errors but I don’t think the code worked.

*Current:* `Team` no longer has any special error handling to do as potential exceptions happen when `Commissioner` is running.

#### Other error handling

*Legacy:* Clobbering errors are silently ignored. Calling `insert_before` with ranges that extend beyond the source code was silently fixed.

*Current:* Such errors are not ignored. It is still ok that a given Cop’s corrections clobber another Cop’s, but any given Cop should not issue corrections that clobber each other, or with invalid ranges, otherwise these will be listed in the processing errors.

### Cop persistence

Cops can now be persisted between files. By default new cop instances are created for each source. See `support_multiple_source?` documentation.

### Internal classes

### Misc API changes

-
internal API clarified for Commissioner. It calls `begin_investigation`and receives the results in`complete_investigation`.
-
New method `add_global_offense`for offenses that are not attached to a location in particular.
-
`#offenses`: No longer accessible.
-
Callbacks `investigate(processed_source)`and`investigate_post_walk(processed_source)`are renamed`on_new_investigation`and`on_investigation_end`and don’t accept an argument; all`on_`callbacks should rely on`processed_source`.
-
`#find_location`is deprecated.
-
`Correction`is deprecated.
-
A few registry access methods were moved from `Cop`to`Registry`both for correctness (e.g.`MyCop.qualified_cop_name`did not work nor made sense) and so that`Cop::Cop`no longer holds any necessary code anymore. Backwards compatibility is maintained.-
`Cop.registry`=>`Registry.global`
-
`Cop.all`=>`Registry.all`
-
`Cop.qualified_cop_name`=>`Registry.qualified_cop_name`

-
-
The `ConfigurableMax`mixin for tracking exclude limits of configuration options is deprecated. Use`exclude_limit ParameterName`instead.

# History

You don’t know where you’re going until you know where you’ve been.

RuboCop was created by Bozhidar Batsov in the spring of 2012, as an attempt
to make it easy to apply consistently the guidelines from the community Ruby Style Guide.[1]

## Notable Milestones

-
12 Sep 2011 - Creation of the community Ruby Style Guide.
-
21 Apr 2012 - Initial commit by Bozhidar Batsov.
-
03 May 2012 - RuboCop 0.0.0 is released.
-
28 May 2013 - RuboCop 0.8 is released. It’s the first version fully powered by the powerful `parser`gem.
-
01 Jul 2013 - RuboCop 0.9 introduces autocorrection and output formatters.
-
31 May 2018 - At RubyKaigi 2018 Bozhidar announces the transition of the project to a RuboCop GitHub organization. He also outlines the plans for a 1.0 release that guide the next couple of years of development.
-
21 Oct 2020 - RuboCop 1.0 is released (exactly 7.5 years after the first commit).
-
06 Mar 2024 - RuboCop 1.62 adds support for Prism.
-
14 Feb 2025 - RuboCop 1.72 introduces the plugin system for extensions.
-
26 Mar 2025 - RuboCop 1.75 makes Prism the default parser for Ruby 3.4+.

# Team

## The Core Team

The direction of the project and its official extensions is being stewarded by the RuboCop core team. This group of long-term contributors manages releases, evaluates pull-requests, and does a lot of the groundwork on major new features. Here are the current members of the RuboCop core team, listed in the order of joining it:

-
Bozhidar Batsov (author & head maintainer)
-
Koichi Ito (also head maintainer of RuboCop Rails, RuboCop Performance and RuboCop Minitest)
-
Benjamin Quorning (also head maintainer of RuboCop RSpec)
-
Marc-André Lafortune (also head maintainer of RuboCop AST)

## RuboCop Alumni

In addition, we’d like to extend a special thanks to the following retired RuboCop core team members. Lovingly known as The Alumni:

-
Yuji Nakayama (also author of `guard-rubocop`)

## Joining the Core Team

We’re always looking for more people to join our Core Team, as there’s plenty of work to go around. There’s no formal procedure for applying to the team, but it’s fairly straightforward to get an invitation to join it:

-
contribute consistently over an extended period of time (e.g. 6-12 months)
-
demonstrate a passion for what we’re doing
-
eventually an invitation will be extended to you to become part of RuboCop’s team

# Changelog

An extensive changelog is available here.

| Only user-visible changes are documented there. |

## Changelog Generation

The changelog is automatically generated from the files in the changelog folder at release time. Direct editing of the changelog file is discouraged as it often creates merge conflicts (pretty much every PR would modify this file and this used to be quite the annoyance in the past). One can create new changelog entries like this:

```
$ bundle exec rake changelog:new
$ bundle exec rake changelog:fix
$ bundle exec rake changelog:change
```
Those commands correspond to "new feature", "bug-fix" and "changed" entries in the changelog. To update the changelog file you can run:

`$ bundle exec rake changelog:merge`| Typically only the RuboCop maintainers would need to do this. |

### Background

Our `CHANGELOG.md` file was previously updated manually by each
contributor that felt their change warranted an entry. When two merge
requests added their own entries at the same spot in the list, it
created a merge conflict in one as soon as the other was merged. When
we had dozens of merge requests fighting for the same changelog entry
location, this quickly became a major source of merge conflicts and
delays in development.

Eventually we adopted GitLab’s solution to this common problem, that they discussed here.

## Release Notes

You can also peruse the release notes for individual releases over at GitHub.

# Logo

RuboCop’s logo was created by Dimiter Petrov. You can find the logo in various formats here.

The logo is licensed under a Creative Commons Attribution-NonCommercial 4.0 International License.

RuboCop’s logo was created by Dimiter Petrov. You can find the logo in various formats here.

The logo is licensed under a Creative Commons Attribution-NonCommercial 4.0 International License.

# License & Copyright

## Code License

Use of RuboCop is granted under the terms of the MIT license. Check
out the `LICENSE.txt` file in RuboCop’s code repository for more details.

Use of RuboCop is granted under the terms of the MIT license. Check
out the `LICENSE.txt` file in RuboCop’s code repository for more details.
