# anonfileCLI
 AnonFiles CLI is a command-line tool written in Go that allows you to upload files to the AnonFiles service and retrieve information about uploaded files.

## Usage

To use the AnonFiles CLI, follow these steps:

1. Make sure you have Go installed on your machine.
2. Clone this repository or download the `main.go` file.
3. Open a terminal or command prompt and navigate to the directory containing the `main.go` file.
4. Run the following command to build the executable:

    go build -o anonfiles main.go

The above command will create an executable file named anonfiles in the current directory.
Run the CLI with the desired options and arguments. Here are the available options:
-f, --file: Specifies the file to upload.
-i, --info: Retrieves information about a specific file using the file's ID.
-v, --version: Displays the version of the AnonFiles CLI.

Examples:

To upload a file:
    ./anonfiles -f path/to/file.txt

To retrieve information about a file using its ID:
    ./anonfiles -i fileID

To display the version of the AnonFiles CLI:
    ./anonfiles -v

License

This project is licensed under the MIT License.
