# rena

A simple utility to rename multiple files

## 🚀 Usage

```
$ rena
A simple utility to rename multiple files

Usage:
  rena [flags] [file]...

Flags:
  -c, --command strings          specify the rename command(s) to execute
  -s, --command-script strings   load commands from a script file
  -f, --force                    execute without confirmation prompts
  -h, --help                     help for rena
  -v, --version                  version for rena
```

## 🛠️ Commands

### Search and Replace

**Format:**

```
s,search,replace,[flags]
```

**Flags (optional):**

- `r` – enable regex
- `m` – match case
- `w` - keep duplicate whitespace

**Variables:**

- `%index` – regex capture group index
- `%name` – regex capture group name

**Example:**

```bash
# guthib.txt -> github.txt
rena 'guthib.txt' -c 's/gut/git' -c 's/hib/hub'
```

### Delete Pattern

**Format:**

```
d,pattern,[flags]
```

**Flags (optional):**

- `r` – enable regex
- `m` – match case
- `w` - keep duplicate whitespace

**Example:**

```bash
# github.txt -> git.txt
rena 'github.txt' -c 'd/hub'
```

### Template Rename

**Format:**

```
t,template
```

**Variables:**

- `%f` – full filename with extension
- `%n` – filename without extension
- `%x` – file extension

**Example:**

```bash
# github.txt -> 01.github.txt
rena 'github.txt' -c 't/01.%f'
```

### Move Files

**Format:**

```
m,pattern,destinationDirectory,[flags]
```

**Flags (optional):**

- `r` – enable regex
- `m` – match case

**Example:**

```bash
# Move github.txt & gitlab.txt to /home/projects
rena 'github.txt' 'gitlab.txt' -c 'm|git|/home/projects'
```
