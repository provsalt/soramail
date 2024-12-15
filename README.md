# Soramail
**Soramail** is a Terminal User Interface (TUI) application that allows you to efficiently generate and manage forwarded email addresses using [Cloudflare Email Routing](https://www.cloudflare.com/products/email-routing/).

With Soramail, you can quickly create forwarded email rules and reduce spam emails without navigating through the Cloudflare web interface. It's fast, intuitive, and designed for developers who prefer working within the terminal.

---

## Features

- **Multiple Zones**: Select from different domains configured in your Cloudflare account.
- **High Performance**: Built with Go and the Bubbletea framework for a smooth, fast TUI experience.

---

## Installation

### Prerequisites

- Go (latest version)
- Cloudflare API Token with appropriate permissions for Email Routing

### Install Soramail

#### With Go Install (recommended)

```bash
go install https://github.com/provsalt/soramail
```


#### Manually

Clone the repository and build the binary:

```bash
git clone https://github.com/yourusername/soramail.git
cd soramail
go build -o soramail
```

Move the binary to a location in your PATH:

```bash
sudo mv soramail /usr/local/bin/
```

---

## Configuration

1. **Set your Cloudflare API token**:

   Create a new folder at ~/.config/soramail or %LOCALAPPDATA%\soramail (for windows)
   ```bash
   mkdir ~/.config/soramail
   touch ~/.config/soramail/config.toml 
   ```

2. **Configure**
    Edit the config file with your preferred editor
    ```toml
    APIKey = "your-api-key"
    ```

3. **Run Soramail**:

   ```bash
   soramail
   ```

   The TUI will guide you through selecting a zone and creating a new email alias.

---

## Usage

- **Navigate** through the TUI using arrow keys or using `j` and `k` keys.
- **Select** options with `Enter` or `l`.
- **Backtrack** by pressing `h` or `esc`.
- **Quit** at any time with `Ctrl+C` or `q`.

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with improvements or bug fixes.

### Planned features
- Setup wizard
- Managing of email forwarding rules.
- Settings for configuring API key and randomize function.
- Colours using lipgloss

---

## License

This project is licensed under the GPLv3 License. See the [LICENSE](./LICENSE) file for details.

---

## Acknowledgments

- [Bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI framework (thanks for sending over stickers awhile back!)
- Cloudflare for their powerful Email Routing service

---

Happy mailing! ✉️