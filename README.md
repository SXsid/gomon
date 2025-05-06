# Gomon

Gomon is a nodemon alternative CLI tool specifically designed for Go applications. It monitors your Go files for changes and automatically rebuilds and restarts your application, making development faster and more efficient.


## Features

- Monitors your Go application for file changes
- Automatically rebuilds and restarts on changes
- Supports web server applications with custom port configuration
- Works on multiple platforms (Windows, macOS, Linux)

## Installation

You have two options to install Gomon:

### Option 1: Install directly using Go

```bash
go install github.com/SXsid/gomon/cmd/gomon@latest
```

This will download and install Gomon to your Go bin directory. Make sure your Go bin directory is in your PATH.

### Option 2: Clone the repository

```bash
git clone https://github.com/SXsid/gomon.git
cd gomon
go mod tidy
```

After cloning, you can:

- On Unix/Linux/macOS:
  ```bash
  make install
  ```

- On Windows:
  ```bash
  go build ./cmd/gomon
  ```
  This will create a `gomon.exe` file in your current directory.

## Usage

### For standard Go applications

If you installed using `go install` or `make install`:

```bash
gomon --main <path_to_your_main.go>
```

If you built the executable on Windows:

```bash
 --watch <directory_path>
```

### For web server applications

For web applications, you can specify a custom port and watching directory:
```bash
gomon --main <path_to_your_main.go> --port <port_number> --watch <directory_path>
```

Example:

```bash
gomon --main ./cmd/server/main.go --port 8080 --watch .
```

## Configuration Options

| Option | Description |
|--------|-------------|
| `--main` | Path to your main Go file (required) |
| `--port` | Port number for web applications (optional) |
| `--watch` | Directory to watch for changes (default: current directory) |

## Example

Running a simple Go web server:

```bash
gomon --main ./main.go --port 3000
```

## Troubleshooting

- Ensure your Go bin directory is in your PATH
- Make sure you have write permissions to the directory where you're running Gomon
- If changes are not being detected, check if you're monitoring the correct directory

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.