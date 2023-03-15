# KEEPTRAK

Track recon/OSINT commands in an organized, grepable, fasion.

**No more saving and organizing terminal output manually.**

Keeptrak is a case management system for hackers, bug bounty hunters, and general cyber security professionals. It runs a nested shell where you can run your usual scans and commands but in the background it saves the output and commands run in various ways. You can also store notes and evidence records with timestamps for backtracking later on.

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

## Installation

Go v19 is the only prerequisite.

### Install using go

```
go install github.com/xenophonsec/keeptrak@latest
```

### Build from source

```
git clone https://github.com/xenophonsec/keeptrak.git
cd keeptrak
go build
# or install globally
go install
```


## Nested Shell

Running keeptrak with no commands will launch the nested shell, which stores your command hitory and records all output.

```
keeptrak
Enter Case Name: mywebsite
KEEPTRAK> nmap -sV mywebsite
```

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

You can store individual records in a case database by providing the values as arguments.
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
