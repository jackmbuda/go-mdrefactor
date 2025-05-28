# Go Markdown Refactor CLI

**Go Markdown Refactor CLI** is a command-line tool developed in Go that improves Markdown files using the OpenAI API, such as GPT-3.5-turbo or GPT-4. It enhances the clarity, structure, and formatting of Markdown documents.

## Description

The tool accepts a Markdown file as input, sends its content to the specified OpenAI model through the API, and then provides the refactored Markdown output.

## Features

- Reads Markdown content from a local file.
- Sends content to the OpenAI API for refactoring.
- Supports various OpenAI models, configurable through a flag.
- Allows customization of system prompts to guide the AI's refactoring style.
- Outputs refactored Markdown content to a specified file or standard output (stdout).
- Allows providing the API key via a command-line flag or an environment variable.

## Prerequisites

- **Go (for building from source)**: Version 1.18 or higher is recommended.
- **OpenAI API Key**: An active key from OpenAI is necessary to use this tool. Obtain one from platform.openai.com.
- **Internet Connection**: Required for communicating with the OpenAI API.

## Installation

You can use the tool by downloading a pre-compiled binary (if available in releases) or by building it from the source.

### Option 1: Using Pre-compiled Binaries (Recommended for most users)

For Linux/macOS:
```bash
chmod +x ./markdown-refactor-linux-amd64
sudo mv ./markdown-refactor-linux-amd64 /usr/local/bin/markdown-refactor
```

For Windows:
- Place the downloaded .exe file in a chosen directory.
- Optionally, add that directory to your system's PATH environment variable.

### Option 2: Building from Source

```bash
git clone https://github.com/jackmbuda/go-mdrefactor.git
cd go-mdrefactor
go build -o go-mdrefactor mdrefactor.go
go install .
```

## Configuration

### OpenAI API Key

To enable the tool to function, you must provide the OpenAI API key in either of the following ways:

- **Environment Variable (Recommended)**:
```bash
export OPENAI_API_KEY="your_openai_api_key_here"
```

- **Command-Line Flag**:
```bash
./markdown-refactor -apikey "your_openai_api_key_here" -input ...
```

## Usage

```bash
./markdown-refactor [flags]
```

Flags:
- `-input <filepath>`: Path to the input Markdown file.
- `-output <filepath>`: Path to the output Markdown file. If absent, the refactored content is printed to stdout.
- `-apikey <key>`: Your OpenAI API key, overriding the environment variable.
- `-model <model_name>`: The OpenAI model for refactoring.
- `-prompt "<system_prompt_text>"`: System prompt to guide the AI's refactoring style.

## Examples

```bash
./mdrefactor -input mydoc.md
./mdrefactor -input mydoc.md -output refactored_doc.md -apikey "sk-yourkey"
./mdrefactor -input draft.md -output final.md -model "gpt-4" -prompt "Refactor this Markdown to be more concise and suitable for a technical audience."
```

## Building for Distribution (Cross-Compilation)

If you wish to create binaries for various operating systems and architectures, use the provided build script or `go build` with appropriate environment variables.

## Contributing

Contributions are encouraged! Fork the repository, create your feature branch, commit your changes, and open a pull request.