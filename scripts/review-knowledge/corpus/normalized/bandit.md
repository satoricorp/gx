# Configuration

## Bandit Settings

Projects may include an INI file named `.bandit`, which specifies
command line arguments that should be supplied for that project.
In addition or alternatively, you can use a YAML or TOML file, which
however needs to be explicitly specified using the `-c` option.
The currently supported arguments are:

`targets`
 comma separated list of target dirs/files to run bandit on
`exclude`
 comma separated list of excluded paths -- *INI only*
`exclude_dirs`
 comma separated list of excluded paths (directories or files) -- *YAML and TOML only*
`skips`
 comma separated list of tests to skip
`tests`
 comma separated list of tests to run

To use this, put an INI file named `.bandit` in your project's directory.
Command line arguments must be in `[bandit]` section.
For example:

```ini
# FILE: .bandit
[bandit]
exclude = tests,path/to/file
tests = B201,B301
skips = B101,B601
```

Alternatively, put a YAML or TOML file anywhere, and use the `-c` option.
For example:

```yaml
# FILE: bandit.yaml
exclude_dirs: ['tests', 'path/to/file']
tests: ['B201', 'B301']
skips: ['B101', 'B601']
```

```toml
# FILE: pyproject.toml
[tool.bandit]
exclude_dirs = ["tests", "path/to/file"]
tests = ["B201", "B301"]
skips = ["B101", "B601"]
```

Then run bandit like this:

```console
bandit -c bandit.yaml -r .
```

```console
bandit -c pyproject.toml -r .
```

Note that Bandit will look for `.bandit` file only if it is invoked with `-r` option.
If you do not use `-r` or the INI file's name is not `.bandit`, you can specify
the file's path explicitly with `--ini` option, e.g.

```console
bandit --ini tox.ini
```

If Bandit is used via pre-commit and a config file, you have to specify the config file
and optional additional dependencies in the pre-commit configuration:

```yaml
repos:
- repo: https://github.com/PyCQA/bandit
 rev: '' # Update me!
 hooks:
 - id: bandit
 args: ["-c", "pyproject.toml"]
 additional_dependencies: ["bandit[toml]"]
```

## Exclusions

In the event that a line of code triggers a Bandit issue, but that the line
has been reviewed and the issue is a false positive or acceptable for some
other reason, the line can be marked with a `# nosec` and any results
associated with it will not be reported.

For example, although this line may cause Bandit to report a potential
security issue, it will not be reported:

```python
self.process = subprocess.Popen('/bin/echo', shell=True) # nosec
```

Because multiple issues can be reported for the same line, specific tests may
be provided to suppress those reports. This will cause other issues not
included to be reported. This can be useful in preventing situations where a
nosec comment is used, but a separate vulnerability may be added to the line
later causing the new vulnerability to be ignored.

For example, this will suppress the report of B602 and B607:

```python
self.process = subprocess.Popen('/bin/ls *', shell=True) # nosec B602, B607
```

Full test names rather than the test ID may also be used.

For example, this will suppress the report of B101 and continue to report B506
as an issue.

```python
assert yaml.load("{}") == [] # nosec assert_used
```

## Scanning Behavior

Bandit is designed to be configurable and cover a wide range of needs, it may
be used as either a local developer utility or as part of a full CI/CD
pipeline. To provide for these various usage scenarios bandit can be configured
via a YAML file. This file is completely optional and in many cases not
needed, it may be specified on the command line by using `-c`.

A bandit configuration file may choose the specific test plugins to run and
override the default configurations of those tests. An example config might
look like the following:

```yaml
### profile may optionally select or skip tests

exclude_dirs: ['tests', 'path/to/file']

# (optional) list included tests here:
tests: ['B201', 'B301']

# (optional) list skipped tests here:
skips: ['B101', 'B601']

### override settings - used to set settings for plugins to non-default values

any_other_function_with_shell_equals_true:
 no_shell: [os.execl, os.execle, os.execlp, os.execlpe, os.execv, os.execve,
 os.execvp, os.execvpe, os.spawnl, os.spawnle, os.spawnlp, os.spawnlpe,
 os.spawnv, os.spawnve, os.spawnvp, os.spawnvpe, os.startfile]
 shell: [os.system, os.popen, os.popen2, os.popen3, os.popen4,
 popen2.popen2, popen2.popen3, popen2.popen4, popen2.Popen3,
 popen2.Popen4, commands.getoutput, commands.getstatusoutput]
 subprocess: [subprocess.Popen, subprocess.call, subprocess.check_call,
 subprocess.check_output]
```

Run with:

```console
bandit -c bandit.yaml -r .
```

If you require several sets of tests for specific tasks, then you should create
several config files and pick from them using `-c`. If you only wish to control
the specific tests that are to be run (and not their parameters) then using
`-s` or `-t` on the command line may be more appropriate.

Also, you can configure bandit via a pyproject.toml file. In this case you
would explicitly specify the path to configuration via `-c`, too. For example:

```toml
[tool.bandit]
exclude_dirs = ["tests", "path/to/file"]
tests = ["B201", "B301"]
skips = ["B101", "B601"]

[tool.bandit.any_other_function_with_shell_equals_true]
no_shell = [
 "os.execl",
 "os.execle",
 "os.execlp",
 "os.execlpe",
 "os.execv",
 "os.execve",
 "os.execvp",
 "os.execvpe",
 "os.spawnl",
 "os.spawnle",
 "os.spawnlp",
 "os.spawnlpe",
 "os.spawnv",
 "os.spawnve",
 "os.spawnvp",
 "os.spawnvpe",
 "os.startfile"
]
shell = [
 "os.system",
 "os.popen",
 "os.popen2",
 "os.popen3",
 "os.popen4",
 "popen2.popen2",
 "popen2.popen3",
 "popen2.popen4",
 "popen2.Popen3",
 "popen2.Popen4",
 "commands.getoutput",
 "commands.getstatusoutput"
]
subprocess = [
 "subprocess.Popen",
 "subprocess.call",
 "subprocess.check_call",
 "subprocess.check_output"
]
```

Run with:

```console
bandit -c pyproject.toml -r .
```

## Skipping Tests

The bandit config may contain optional lists of test IDs to either include
(`tests`) or exclude (`skips`). These lists are equivalent to using `-t` and
`-s` on the command line. If only `tests` is given then bandit will include
only those tests, effectively excluding all other tests. If only `skips`
is given then bandit will include all tests not in the skips list. If both are
given then bandit will include only tests in `tests` and then remove `skips`
from that set. It is an error to include the same test ID in both `tests` and
`skips`.

Note that command line options `-t`/`-s` can still be used in conjunction with
`tests` and `skips` given in a config. The result is to concatenate `-t` with
`tests` and likewise for `-s` and `skips` before working out the tests to run.

## Suppressing Individual Lines

If you have lines in your code triggering vulnerability errors and you are
certain that this is acceptable, they can be individually silenced by appending
`# nosec` to the line:

```python
# The following hash is not used in any security context. It is only used
# to generate unique values, collisions are acceptable and "data" is not
# coming from user-generated input
the_hash = md5(data).hexdigest() # nosec
```

In such cases, it is good practice to add a comment explaining *why* a given
line was excluded from security checks.

## Generating a Config

Bandit ships the tool `bandit-config-generator` designed to take the leg work
out of configuration. This tool can generate a configuration file
automatically. The generated configuration will include default config blocks
for all detected test and blacklist plugins. This data can then be deleted or
edited as needed to produce a minimal config as desired. The config generator
supports `-t` and `-s` command line options to specify a list of test IDs that
should be included or excluded respectively. If no options are given then the
generated config will not include `tests` or `skips` sections (but will provide
a complete list of all test IDs for reference when editing).

## Configuring Test Plugins

Bandit's configuration file is written in YAML and options
for each plugin test are provided under a section named to match the test
method. For example, given a test plugin called 'try_except_pass' its
configuration section might look like the following:

```yaml
try_except_pass:
 check_typed_exception: True
```

The specific content of the configuration block is determined by the plugin
test itself. See the plugin test list for complete information on
configuring each one.

# Frequently Asked Questions

## Under Which Version of Python Should I Install Bandit?

The answer to this question depends on the project(s) you will be running
Bandit against. If your project is only compatible with Python 3.9, you
should install Bandit to run under Python 3.9. If your project is only
compatible with Python 3.10, then use 3.10 respectively. If your project
supports both, you *could* run Bandit with both versions but you don't have to.

Bandit uses the `ast` module from Python's standard library in order to
analyze your Python code. The `ast` module is only able to parse Python code
that is valid in the version of the interpreter from which it is imported. In
other words, if you try to use Python 2.7's `ast` module to parse code written
for 3.5 that uses, for example, `yield from` with asyncio, then you'll have
syntax errors that will prevent Bandit from working properly. Alternatively,
if you are relying on 2.7's octal notation of `0777` then you'll have a syntax
error if you run Bandit on 3.x.

# Welcome to Bandit

Bandit is a tool designed to find common security issues in Python code. To do
this, Bandit processes each file, builds an AST from it, and runs appropriate
plugins against the AST nodes. Once Bandit has finished scanning all the files,
it generates a report.

# Using and Extending Bandit

# Contributing

* Source code
* Issue tracker
* Join us on Discord

# Indices and tables

* `genindex`
* `modindex`
* `search`

# License

The `bandit` library is provided under the terms and conditions of the
[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.txt)

# Integrations

Bandit can be integrated into a wide variety of developer tools, editors,
CI/CD systems, and code quality pipelines. This page outlines popular
integrations to help you seamlessly incorporate Bandit into your development
workflow.

## IDE Integrations

**List table:**
* - Visual Studio Code
 - `Bandit by PyCQA <https://marketplace.visualstudio.com/items?itemName=pycqa.bandit-pycqa>`_
* - Sublime Text
 - `SublimeLinter-bandit <https://github.com/SublimeLinter/SublimeLinter-bandit>`_
* - Vim/Neovim
 - `Asynchronous Lint Engine <https://github.com/dense-analysis/ale>`_
* - Emacs
 - `flycheck-pycheckers <https://github.com/msherry/flycheck-pycheckers>`_

## CI/CD Integrations

**List table:**
* - GitHub Action
 - `Bandit by PyCQA <https://github.com/marketplace/actions/bandit-by-pycqa>`_
* - Hudson/Jenkins
 - `Bandit Plugin <https://github.com/mewz/bandit-plugin->`_

## Linters

**List table:**
* - Ruff
 - `flake8-bandit (S) <https://docs.astral.sh/ruff/rules/#flake8-bandit-s>`_
* - Flake8
 - `flake8-bandit <https://github.com/tylerwince/flake8-bandit>`_

## Packages

**List table:**
* - Ubuntu
 - `bandit <https://packages.ubuntu.com/search?keywords=bandit&searchon=names&section=all>`_
* - Homebrew
 - `bandit <https://formulae.brew.sh/formula/bandit>`_
* - FreeBSD
 - `py-bandit <https://www.freshports.org/devel/py-bandit/>`_
🙌 Contributions Welcome

If you’ve integrated Bandit into another platform or tool, feel free to open
a PR and update this page!

# Getting Started

## Installation

Bandit is distributed on PyPI. The best way to install it is with pip.

Create a virtual environment and activate it using `virtualenv` (optional):

```console
virtualenv bandit-env
source bandit-env/bin/activate
```

Alternatively, use `venv` instead of `virtualenv` (optional):

```console
python3 -m venv bandit-env
source bandit-env/bin/activate
```

Install Bandit:

```console
pip install bandit
```

If you want to include TOML support, install it with the `toml` extras:

```console
pip install bandit[toml]
```

If you want to use the bandit-baseline CLI, install it with the `baseline`
extras:

```console
pip install bandit[baseline]
```

If you want to include SARIF output formatter support, install it with the
`sarif` extras:

```console
pip install bandit[sarif]
```

Run Bandit:

```console
bandit -r path/to/your/code
```

Bandit can also be installed from source. To do so, either clone the
repository or download the source tarball from PyPI, then install it:

```console
python setup.py install
```

Alternatively, let pip do the downloading for you, like this:

```console
pip install git+https://github.com/PyCQA/bandit#egg=bandit
```

## Usage

Example usage across a code tree:

```console
bandit -r ~/your_repos/project
```

Two examples of usage across the `examples/` directory, showing three lines of
context and only reporting on the high-severity issues:

```console
bandit examples/*.py -n 3 --severity-level=high
```

```console
bandit examples/*.py -n 3 -lll
```

Bandit can be run with profiles. To run Bandit against the examples directory
using only the plugins listed in the `ShellInjection` profile:

```console
bandit examples/*.py -p ShellInjection
```

Bandit also supports passing lines of code to scan using standard input. To
run Bandit with standard input:

```console
cat examples/imports.py | bandit -
```

For more usage information:

```console
bandit -h
```

## Baseline

Bandit allows specifying the path of a baseline report to compare against using the base line argument (i.e. `-b BASELINE` or `--baseline BASELINE`).

```console
bandit -b BASELINE
```

This is useful for ignoring known vulnerabilities that you believe are non-issues (e.g. a cleartext password in a unit test). To generate a baseline report simply run Bandit with the output format set to `json` (only JSON-formatted files are accepted as a baseline) and output file path specified:

```console
bandit -f json -o PATH_TO_OUTPUT_FILE
```

## Version control integration

Use pre-commit. Once you have it installed, add this to the
`.pre-commit-config.yaml` in your repository
(be sure to update `rev` to point to a `real git tag/revision`_!):

```yaml
repos:
- repo: https://github.com/PyCQA/bandit
 rev: '' # Update me!
 hooks:
 - id: bandit
```

Then run `pre-commit install` and you're ready to go.
