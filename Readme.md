# GO BUILDER CLI

This cli copies directory content to the destination directory. Also compiles
your go application into binary.

## Getting started

### Installation

    go get github.com/sudiptob2/go-builder-cli

### Usage

    Usage:

        go-builder-cli [command]

    Available Commands:
        buildexecute build and execute the project of specified
        directory help Help about any command

    Flags:
            --config string config file (default is $HOME/.go-builder-cli.yaml)
        -h, --help help for go-builder-cli -t, --toggle Help message for toggle


        go-builder-cli buildexecute [flags]

    Flags:
    -b, --builddir string   Specify the destination directory
    -c, --copydir string    Copies the content of the specified dirtectory
    -x, --exclude-tests     excludes the golang test files
    -e, --exe string        compile the code of the directory as a binary
    -h, --help              help for buildexecute

### Example

Do not specify builddir. In this case files will be copied to the current working diretory.

        go run main.go buildexecute --copydir D:\test\sudiptob2\source\.

Specify builddir

        go run main.go buildexecute --copydir D:\test\sudiptob2\source\. --builddir D:\test\sudiptob2\dest

Copy and compile

        go run main.go buildexecute --copydir D:\test\sudiptob2\source\. --exe <binary name>

Exclude go test files

        go run main.go buildexecute --copydir D:\test\sudiptob2\source\. --exe <binary name>  --exclude-tests
        


