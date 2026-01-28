# Project GNS
###### INSA Lyon Telecoms Department 3rd year Group 2, GNS Group 26

This project goal is to develop a program that automate the configuration of networks in the GNS3 software using a "Network Intent" JSON file. It parses the intent file, create .cfg router config files and applies them to the specified GNS3 project.

## Installation

### Prerequisites

- Go 1.24.*
- GNS3
- A GNS3 project with the routers and their links matching the intent file.

### Installation

1. Clone the repository:
```sh
git clone https://github.com/Shaa42/GNS-26.git
cd GNS-26
```

2. Build the project (optionnal):
```sh
go build -o gns-config cmd/main.go
```

### How to use it

1. Open GNS3 with the target project open

2. Choose the right Network Intent file

3. Run the program :
```sh
./gns-config
```

Or directly:
```sh
go run cmd/main.go
```

## Project tree
```
.
├── cmd
│   └── main.go                 # Entry point
├── go.mod
├── internal
│   ├── applycfg                # Package for applying configs file to GNS3
│   ├── parseintent             # Package to parse the JSON Network Intent file
│   └── writeconf               # Package for writing config files
├── networks                    # Network Intent files with their txt file to resume the 
│   ├── network_intent_1.json   # network structure
│   ├── network_intent_2.json
│   ├── nw1_scheme.txt
│   └── nw2_scheme.txt
├── README.md
└── tests
    └── writeconf_test.go
```

## Commit convention :
- FEAT: Added README.md
- FIX: Corrected bug 01214 in README.md
- OTHER: Folder organization

## Dev rules :
For each new feature that need to be added please create a new branch and develop on it. Then create a pull request and don't merge it until we gather in TD.

We also need to find who is going to be responsible of merging branches.
