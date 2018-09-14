Cl is a tool for filtering data by columns on the command line.

# What it does
Usage: `cl [options...] <column_indexes...>`

Examples:
```bash
# Filter a simple table of data for the second column:
$ echo "1 2 3
> 4 5 6
> 7 8 9" | cl 2
2
5
8

# Grab a list of process IDs from ps, ignoring the header row (-i):
$ ps | cl 1 -i
7958
29855

# Grab the third column of output when there may be spaces, in values, and tabs are the separator (-t):
$ netstat | cl 3 -t

# Or the first 2 columns:
$ netstat | cl 1 2 -t

# Or the first 4 columns (in bash):
$ netstat | cl {1..4} -t
```

Options:
```
  -i    ignore the header row (first row)
  -s string
        a character or regex to split lines (default: whitespace)
  -t    use tabs as separator (alias of -s \t)
```

# Install

```bash
go get github.com/dankinder/cl
```

# Why
Existing commonly-available bash commands like `awk` and `tail` can accomplish
what `cl` does. But they become quite verbose.

To select a column with `awk`, you use:
```sh
my_command | awk '{ print $1 }'
```

But then, you want to delete the header row too (another very common thing). For that you need:
```sh
my_command | awk '{ print $1 }' | tail -n +2
```

I used to have the following bash function in my .bashrc, and I'm sure many
others have equivalents:
```bash
# Grab column N from stdin, ex. `ps aux | cl 2` => Just Process IDs
function cl() {
    awk "{print \$$1}"
}

# Ignore the header row of the input
function ih() {
	tail -n +2
}
```

Further, some tools worked around this specific problem by building another
tool (`pgrep`) giving us yet another command to learn.

It felt to me like we need a small, simple, generalized tool for this task.

So I build `cl`.
