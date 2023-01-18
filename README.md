# Stylist

Stylist aggregates all your code quality tools into a single, fast,
and easy to use executable. You define the tools your project uses
in `.stylist.yml`:

```yml
---
# A "processor" is a code quality tool (linter, scanner, formatter, etc...)
processors:

  ############################################################
  # Stylist comes with presets for many common tools...
  ############################################################

  - name: cspell
    preset: cspell
    includes:
      - "**/*"
    excludes:
      - "problem-file1.txt"
      - "problem-file2.txt"

  - name: gofmt
    preset: gofmt
    includes:
      - "**/*.go"

  - name: golangci-lint
    preset: golangci-lint
    includes:
      - "**/*.go"
    excludes:
      - "gen/*"

  - name: markdownlint
    preset: markdownlint
    includes:
      - "**/*.md"

  ############################################################
  # ... or you can provide custom configuration.
  ############################################################

  - name: custom-tool
    includes:
      - "**/*.{ext1,ext2,ext3}"
    # Processors define check and/or fix commands
    # (check commands report issues, fix commands correct them).
    # See the documentation for more detail.
    check:
      command: "custom-tool lint --output=sarif"
      input: variadic
      output: stdout
      format: sarif
    fix:
      command: "custom-tool lint --auto-fix --output=sarif"
      input: variadic
      output: stdout
      format: sarif
```

Then run stylist's `check` and `fix` commands:

```shell
$ stylist check
[cspell] main.go:8:16 Unknown word 'someFunc' (spelling)
[golangci-lint] main.go:8:16 Error return value of `someFunc` is not checked (errcheck)
[markdownlint] README.md:1:9 Trailing punctuation in heading [Punctuation: ':'] (MD026)
[markdownlint] README.md:3:81 Line length [Expected: 80; Actual: 153] (MD013)

$ stylist fix
[markdownlint] README.md
```
