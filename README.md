# KEEPTRAK

Track recon/OSINT commands in an organized, grepable, fasion.

**No more saving and organizing terminal output manually.**

Keeptrak is a case management system for ethical hackers, bug bounty hunters, OSINT investigators, and cyber security professionals. It runs a nested shell where you can run your usual scans and commands but in the background it saves and timestamps the output, keeping an organized. grepable, history and list of records. You can also store timestamped notes in text format and evidence records in a csv file. Keeptrak is designed to work not only as a nested shell but also integrates with other programs by accepting and echoing data via pipes. This means you can pipe information into keeptrak without being in the nested shell and even pass information through it to another program. 

Example Output:
- casedirectory (folder of case files)
  - **db.csv** (csv file of records: LABEL, VALUE, DATA_TYPE, CONFIRMED, DATE_ADDED)
  - **dump** (dump of all terminal output)
  - **history** (list of commands that were used with timestamps)
  - **notes** (notes with timestamps)
  - **nmap** (all output by said program)
  - **dnsrecon** (all output by said program)
  - etc...


## Features

- [Nested Shell](#nested-shell)
- [Pipes](#pipes)
- [Scripting](#scripting)
- [Notes](#quick-notes)
- [Findings Records](#store-findings-records)
- [Trust Chains](#trust-chains)

## Installation

Go v19 or higher is the only prerequisite.

### Install using go

```
go install github.com/xenophonsec/keeptrak@latest
```

Make sure the go bin path has been added to your PATH global variable.

### Build from source

```
git clone https://github.com/xenophonsec/keeptrak.git
cd keeptrak
go build
# or install globally
go install
```

If installing globally is not working, you can do so mannually.

```
export PATH="$PATH:/path/to/keeptrak"
```
Replace "/path/to/keeptrak" with the actual path to the keeptrak directory.
Running this command will only work temporarily. If you want the install to be permenant, you should add this command to the bottom of your .bashrc file.
```
echo "export PATH=\"\$PATH:/path/to/keeptrak\"" >> ~/.bashrc
```

## Nested Shell

Running keeptrak with no commands will launch the nested shell, which stores your command hitory and records all output.

```
keeptrak
Enter Case Name: mywebsite
KEEPTRAK> nmap -sV mywebsite
```

To exit the nested shell, type and run `exit`.

The output of every command that is run will be stored in a file named after the command within a directory named after the case name provided.
Every command is also added to the history file and timestamped.

## Pipes

You can pipe data into keeptrak and assign a case and label to track the output of whatever you run.
```
cat file | keeptrak CASENAME LABEL
```

Keeptrak will relay whatever is piped into it so you don't lose visibililty and can continue to pipe that into anything else.
```
cat test | keeptrak testcase test
this is my test file
it has multiple lines
I am tracking it in keeptrak
```
```
cat file | keeptrak mycase mylabel | sort | uniq
apple
box
cat
```

## Scripting

You can use pipes to feed data into keeptrak in a script:
```
#!/bin/bash

TARGET=$1
nmap -sV $TARGET | keeptrak $TARGET nmap
```

You can also store data you found in keeptrak by passing it as arguments:
```
keeptrak $CASE Server $IPADDRESS "IP Address" Y
```

## Quick Notes

You can save timestamped notes quickly using the "note: " prefix in the nested shell.

```
KEEPTRAK> note: The netmask is 255.255.255.0
```

If you want to create timestamped notes in a script you can use the following syntax:
```
keeptrak mycase note "this is an important note"
```

If you don't want to create a case folder you can simply use a period as the case name and it will crate the note in the current directory.

```
keeptrak . note "This is my note"
```

### Tags

Just as a friendly recommendation, you can use tags to create more searchable notes.

- You can use "#" syntax to create a tag system that grep can use
- `#important this is a super important note`
- `cat notes | grep "#important"`

## Store Findings Records

You can store individual records in a csv case database by providing the values as arguments.
This can be used for any data found in an investigation: usernames, passwords, hashes, etc...

Usage:
```
keeptrak CASE LABEL VALUE DATA_TYPE CONFIRMED
```
Example:
```
keeptrak myosint username happyfeet credential Y
```

You can even do this within the nested shell
```
KEEPTRAK> keeptrak myosint username happyfeet credential Y
```

## Trust Chains

Keeptrak signs the history of commands with hashes that include all previous commands, the current command, and a timestamp. This creates a chain of trust for verifying that the evidence has not been tampered with.

Signatures consist of the letter D followed by the date (YYYYMMDD), the letter T followed by the time (HHMMSS) and the letter H followed by the hash (SHA256). Due to the format, these signatures are referred to as DTH signatures.

