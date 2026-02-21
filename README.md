<div align="center">
  <h1>Tiny Server</h1>

  <p><i>✨ A minimal HTTP server for local development and file sharing.</i></p>

  <p>
    <a href="https://github.com/cqroot/tinyserver/actions">
      <img src="https://github.com/cqroot/tinyserver/workflows/test/badge.svg" alt="Action Status" />
    </a>
    <a href="https://github.com/cqroot/tinyserver/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/cqroot/tinyserver" alt="LICENSE"/>
    </a>
    <a href="https://github.com/cqroot/tinyserver/issues">
      <img src="https://img.shields.io/github/issues/cqroot/tinyserver" alt="Issues"/>
    </a>
  </p>
</div>

## Usage

Start the server in the current directory:

```bash
tinyserver
```

By default, the server listens on all interfaces (`""`), port `9000`, and serves files from the current working directory.

### Command-line Options

| Option            | Flag              | Default | Description                                 |
| :---------------- | :---------------- | :------ | :------------------------------------------ |
| Working directory | `-d, --work_dir`  | `"."`   | Directory to serve files from               |
| Bind IP           | `-i, --bind_ip`   | `""`    | IP address to bind (empty = all interfaces) |
| Bind port         | `-p, --bind_port` | `9000`  | Port to listen on                           |
| Whitelist         | `-w, --whitelist` | `nil`   | Comma-separated list of allowed source IPs  |

### Configuration File

You can also provide settings via a YAML configuration file. Place a file named `tinyserver.yaml` in the server's directory (or in the directory you intend to serve) with the following structure:

```yaml
bind_ip: 0.0.0.0  # Bind IP
bind_port: 9000   # Bind port
whitelist:        # List of allowed source IPs
  - 192.168.1.10
  - 192.168.1.11
```

To start the server with a specific working directory, use:

```bash
tinyserver -d /path/to/serve
```

## License

Tiny Server is licensed under the [GNU General Public License v3.0](LICENSE). See the [LICENSE](LICENSE) file for details.
