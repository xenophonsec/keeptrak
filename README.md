# KEEPTRAK

Track recon/OSINT commands in an organized, grepable, fasion.

Output:
- casedirectory (folder of case files)
  - db.csv (csv file of records: LABEL, VALUE, DATA_TYPE, CONFIRMED, DATE_ADDED)
  - dump (dump of all terminal output)
  - history (list of commands run)
  - nmap (all output by said program)
  - dnsrecon (all output by said program)

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

## Store specific data

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
