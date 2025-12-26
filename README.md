# ATemplateB

A lightweight, Helm-inspired Go template engine for generating configuration files from templates and values files. Built on Go's powerful `text/template` package with [Sprig](https://masterminds.github.io/sprig/) template functions.

## Features

- **Helm-like workflow** - Use templates with values files (YAML or JSON)
- **Go templates + Sprig** - Full access to Go template syntax and 100+ Sprig functions
- **Custom functions** - Built-in `getFile` and `getYaml` for dynamic file inclusion
- **Multi-template support** - Compose templates using `{{ template }}` directive
- **Auto-format detection** - Automatically detects JSON or YAML values files
- **Cross-platform** - Pre-built binaries for Windows, Linux, and macOS
- **Zero dependencies** - Single binary, no installation required

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/at-blacknight/A-Template-B/releases):

```bash
# Linux (amd64)
wget https://github.com/at-blacknight/A-Template-B/releases/latest/download/ATemplateB-linux-amd64
chmod +x ATemplateB-linux-amd64
sudo mv ATemplateB-linux-amd64 /usr/local/bin/ATemplateB

# macOS (arm64)
wget https://github.com/at-blacknight/A-Template-B/releases/latest/download/ATemplateB-darwin-arm64
chmod +x ATemplateB-darwin-arm64
sudo mv ATemplateB-darwin-arm64 /usr/local/bin/ATemplateB

# Windows
# Download ATemplateB-windows-amd64.exe and add to your PATH
```

### Build from Source

```bash
git clone https://github.com/at-blacknight/A-Template-B.git
cd A-Template-B
go build -o ATemplateB ATemplateB.go
```

## Quick Start

### 1. Create a template file

Create `app.tmpl.yaml`:

```yaml
# Application Configuration
app:
  name: {{ .Values.app_name }}
  port: {{ .Values.port }}
  environment: {{ .Values.env }}

{{ if .Values.debug }}
logging:
  level: debug
{{ else }}
logging:
  level: info
{{ end }}

database:
  host: {{ .Values.db_host }}
  port: {{ .Values.db_port | default "5432" }}
```

### 2. Create a values file

Create `values.yaml`:

```yaml
app_name: MyApp
port: 8080
env: production
debug: false
db_host: postgres.example.com
```

### 3. Generate the output

```bash
ATemplateB --template=app.tmpl.yaml --values=values.yaml --output=app.yaml
```

**Output** (`app.yaml`):

```yaml
# Application Configuration
app:
  name: MyApp
  port: 8080
  environment: production

logging:
  level: info

database:
  host: postgres.example.com
  port: 5432
```

## Usage

### Command-Line Syntax

```bash
ATemplateB --template=<file>.tmpl[.ext] [--values=<file>.yaml/json] --output=<file>.ext
```

### Flags

| Flag | Required | Description | Default |
|------|----------|-------------|---------|
| `--template` | Yes | Path to template file | - |
| `--values` | No | Path to YAML or JSON values file | Empty map |
| `--output` | No | Path to output file | `nginx.conf` |
| `--version` | No | Display version and exit | - |

### Template File Naming

Templates can use two naming patterns:

- **Basic**: `template.tmpl` → generates any output filename
- **With extension**: `template.tmpl.yaml` → helps IDEs with syntax highlighting

The extension after `.tmpl` is optional and only affects editor behavior.

## Template Syntax

### Accessing Values

Values are accessed via the `.Values` object:

```yaml
# values.yaml
server:
  name: web-01
  port: 8080
```

```go
# template.tmpl
Server: {{ .Values.server.name }}
Port: {{ .Values.server.port }}
```

### Conditionals

```go
{{ if .Values.ssl_enabled }}
ssl_certificate {{ .Values.ssl_cert }};
ssl_certificate_key {{ .Values.ssl_key }};
{{ end }}
```

### Loops

```yaml
# values.yaml
upstreams:
  - 192.168.1.10:5000
  - 192.168.1.11:5000
```

```go
# template.tmpl
upstream backend {
{{- range .Values.upstreams }}
    server {{ . }};
{{- end }}
}
```

### Default Values

Use the Sprig `default` function:

```go
port: {{ .Values.port | default "8080" }}
timeout: {{ .Values.timeout | default 30 }}
```

### Whitespace Control

- `{{-` trims whitespace before
- `-}}` trims whitespace after

```go
{{- range .Values.items }}
item: {{ . }}
{{- end }}
```

## Custom Functions

ATemplateB includes custom functions beyond standard Sprig:

### `getFile`

Reads raw file contents into the template:

```go
{{ getFile "path/to/file.txt" }}
```

**Example**:
```go
# Include an SSH public key
authorized_keys: {{ getFile "/home/user/.ssh/id_rsa.pub" }}
```

### `getYaml`

Reads and pretty-formats a YAML file:

```yaml
# main.tmpl.yaml
sitecode: syd
nodes:
- {{ getYaml "config.yaml" | indent 2 | trim }}
```

**Useful for**:
- Including external YAML snippets
- Composing multiple configuration files
- Dynamic file inclusion based on environment

## Sprig Functions

ATemplateB includes all [Sprig v3 functions](https://masterminds.github.io/sprig/). Here are commonly used ones:

### String Functions
```go
{{ .Values.name | upper }}              # UPPERCASE
{{ .Values.name | lower }}              # lowercase
{{ .Values.name | title }}              # Title Case
{{ .Values.text | quote }}              # "quoted"
{{ .Values.text | indent 4 }}           # indent 4 spaces
{{ .Values.text | trim }}               # remove whitespace
```

### Default & Conditionals
```go
{{ .Values.port | default "8080" }}     # default value
{{ .Values.config | required "config is required" }}  # required value
{{ .Values.name | empty }}              # check if empty
```

### Lists
```go
{{ .Values.items | join "," }}          # join with separator
{{ list "a" "b" "c" }}                  # create list
```

### Dictionaries
```go
{{ dict "key" "value" "foo" "bar" }}    # create dict
```

### Encryption & Encoding
```go
{{ .Values.password | b64enc }}         # base64 encode
{{ .Values.data | b64dec }}             # base64 decode
{{ .Values.text | sha256sum }}          # SHA256 hash
```

[→ View all Sprig functions](https://masterminds.github.io/sprig/)

## Multi-Template Composition

ATemplateB automatically loads all `.tmpl*` files in the same directory, allowing template composition:

**File structure**:
```
templates/
├── main.tmpl.conf
├── server.tmpl
└── upstream.tmpl
```

**main.tmpl.conf**:
```go
http {
    {{ template "upstream" . }}

    {{ template "server" . }}
}
```

**upstream.tmpl**:
```go
{{ define "upstream" }}
upstream backend {
    {{- range .Values.backends }}
    server {{ . }};
    {{- end }}
}
{{ end }}
```

**server.tmpl**:
```go
{{ define "server" }}
server {
    listen {{ .Values.port }};
    server_name {{ .Values.server_name }};
}
{{ end }}
```

Pass data to sub-templates:
```go
{{ template "server" . }}                          # pass entire context
{{ template "server" .Values.myserver }}           # pass specific value
{{ template "server" (dict "port" 8080) }}         # pass custom dict
```

## Examples

### Example 1: Simple Nginx Reverse Proxy

See [examples/nginx.tmpl](examples/nginx.tmpl) and [examples/config.yaml](examples/config.yaml)

**Generate**:
```bash
ATemplateB --template=examples/nginx.tmpl --values=examples/config.yaml --output=nginx.conf
```

### Example 2: Multi-Environment Configuration

**values-dev.yaml**:
```yaml
env: development
debug: true
db_host: localhost
replicas: 1
```

**values-prod.yaml**:
```yaml
env: production
debug: false
db_host: postgres.prod.example.com
replicas: 3
```

**Generate for each environment**:
```bash
ATemplateB --template=app.tmpl.yaml --values=values-dev.yaml --output=app-dev.yaml
ATemplateB --template=app.tmpl.yaml --values=values-prod.yaml --output=app-prod.yaml
```

### Example 3: Including External Files

**template.tmpl.yaml**:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  config.yaml: |
{{ getYaml "external-config.yaml" | indent 4 }}

  script.sh: |
{{ getFile "scripts/deploy.sh" | indent 4 }}
```

### Example 4: Complex Nginx with Template Composition

See [examples/sca-nginx.tmpl.conf](examples/sca-nginx.tmpl.conf) for a production-grade example showing:
- Multiple server blocks from arrays
- Template composition with `{{ template }}`
- Conditional SSL configuration
- Dynamic upstream generation

## Comparison to Helm

| Feature | ATemplateB | Helm |
|---------|------------|------|
| **Purpose** | General-purpose config generator | Kubernetes package manager |
| **Template Engine** | Go templates + Sprig | Go templates + Sprig |
| **Values Files** | YAML or JSON | YAML |
| **Output** | Any text file | Kubernetes manifests |
| **Dependencies** | None | Kubernetes cluster |
| **Use Case** | Configs, manifests, any text files | Kubernetes deployments |
| **Learning Curve** | Low (single binary, simple CLI) | Medium (K8s knowledge needed) |

**When to use ATemplateB**:
- Generating application configs (nginx, apache, etc.)
- Creating environment-specific configurations
- Templating infrastructure-as-code files
- Non-Kubernetes manifest generation
- Lightweight templating without K8s overhead

**When to use Helm**:
- Deploying to Kubernetes
- Need package management features
- Require dependency management
- Using Helm charts ecosystem

## Common Patterns

### Environment-Specific Values

```bash
# Organize values by environment
values/
├── common.yaml
├── dev.yaml
├── staging.yaml
└── prod.yaml

# Generate per environment (note: currently single values file supported)
ATemplateB --template=app.tmpl --values=values/prod.yaml --output=app-prod.conf
```

### Using JSON Values

```json
{
  "server": {
    "port": 8080,
    "host": "example.com"
  },
  "features": {
    "cache": true,
    "debug": false
  }
}
```

```bash
ATemplateB --template=app.tmpl --values=values.json --output=app.conf
```

### Template Without Values

Templates can work without values files:

```bash
# Uses empty map for .Values
ATemplateB --template=static.tmpl --output=output.conf
```

## Error Handling

### Template Errors

Errors in template execution show the issue clearly:

```
ERROR: could not execute template nginx.tmpl: template: nginx.tmpl:10:5:
executing "nginx.tmpl" at <.Values.missing>: map has no entry for key "missing"
```

### Custom Function Errors

Errors from `getFile` or `getYaml` are collected and reported:

```
Template rendering completed with errors:
# ERROR[getFile]: could not read /missing/file.txt: no such file or directory
# ERROR[getYaml]: could not parse YAML invalid.yaml: yaml: line 5: did not find expected key
```

### Validation

Check for required values:

```go
{{ .Values.critical_setting | required "critical_setting must be set in values file" }}
```

## Tips & Best Practices

1. **Use descriptive value names** - `db_connection_string` vs `dcs`
2. **Provide defaults** - `{{ .Values.port | default "8080" }}`
3. **Comment your templates** - Explain complex logic
4. **Test incrementally** - Build templates step by step
5. **Use `.tmpl.ext` naming** - Helps IDEs with syntax highlighting
6. **Validate output** - Check generated configs before deployment
7. **Keep values flat when possible** - Easier to override
8. **Use template composition** - Break complex templates into reusable parts

## Troubleshooting

**Problem**: `map has no entry for key`
**Solution**: Value doesn't exist in values file. Add it or use `| default`

**Problem**: Template renders but output is wrong
**Solution**: Check whitespace control with `{{-` and `-}}`

**Problem**: Values file not loading
**Solution**: Ensure valid YAML/JSON syntax. Check file path is correct.

**Problem**: Template inheritance not working
**Solution**: Ensure all `.tmpl*` files are in same directory. Use `{{ define "name" }}` blocks.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- **GitHub**: https://github.com/at-blacknight/A-Template-B
- **Sprig Functions**: https://masterminds.github.io/sprig/
- **Go Templates**: https://pkg.go.dev/text/template
- **Issues**: https://github.com/at-blacknight/A-Template-B/issues
