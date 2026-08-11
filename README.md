# OpenAPI Enhancer

A CLI tool that automatically adds validator tags to OpenAPI 3.0 schemas. Perfect for integration with [`oapi-codegen`](https://github.com/deepmap/oapi-codegen) and [`go-playground/validator`](https://github.com/go-playground/validator).

---

## Features

- **Automatic Validator Tags** – Adds `validate` tags based on OpenAPI specifications
- **Go-Validator Support** – Compatible with `go-playground/validator`
- **OpenAPI 3.0** – Full support for schemas, properties, and constraints
- **Easy Integration** – Seamlessly integrates into existing Makefile workflows
- **Zero Configuration** – Works out-of-the-box

---

## Installation

### Via `go install` (recommended)

```bash
go install github.com/deinusername/openapi-enhancer/cmd/enhancer@latest
```

### From Source

```bash
git clone https://github.com/deinusername/openapi-enhancer.git
cd openapi-enhancer
go build -o enhancer ./cmd/enhancer
sudo mv enhancer /usr/local/bin/
```

## Usage

### Basic Usage

```bash
openapi-enhancer -file api/openapi.bundled.yaml
```

### Options

| Flag    | Short | Description               | Default      |
| ------- | ----- | ------------------------- | ------------ |
| `-file` | `-f`  | Path to OpenAPI YAML file | **required** |

### Examples

```bash
# Standard
openapi-enhancer -file api/openapi.bundled.yaml

# With relative path
openapi-enhancer -f ./specs/api.yaml

# With absolute path
openapi-enhancer -file /home/user/project/api/openapi.yaml
```

## Supported Validator Rules

The tool automatically generates validator tags based on OpenAPI fields:

### String Fields

| OpenAPI Field          | Validator Tag | Example                                         |
| ---------------------- | ------------- | ----------------------------------------------- |
| `required` (in schema) | `required`    | `validate:"required"`                           |
| `minLength`            | `min`         | `validate:"min=3"`                              |
| `maxLength`            | `max`         | `validate:"max=50"`                             |
| `format: email`        | `email`       | `validate:"email"`                              |
| `format: uuid`         | `uuid`        | `validate:"uuid"`                               |
| `format: uri` / `url`  | `uri`         | `validate:"uri"`                                |
| `format: date`         | `datetime`    | `validate:"datetime=2006-01-02"`                |
| `format: date-time`    | `datetime`    | `validate:"datetime=2006-01-02T15:04:05Z07:00"` |
| `enum`                 | `oneof`       | `validate:"oneof=active inactive"`              |

### Numeric Fields (integer, number)

| OpenAPI Field | Validator Tag | Example              |
| ------------- | ------------- | -------------------- |
| `minimum`     | `min`         | `validate:"min=0"`   |
| `maximum`     | `max`         | `validate:"max=100"` |

### Arrays

| OpenAPI Field | Validator Tag | Example             |
| ------------- | ------------- | ------------------- |
| `minItems`    | `min`         | `validate:"min=1"`  |
| `maxItems`    | `max`         | `validate:"max=10"` |

## Example

### Input (`openapi.bundled.yaml`)

```yaml
components:
  schemas:
    User:
      type: object
      required:
        - email
        - name
      properties:
        email:
          type: string
          format: email
        name:
          type: string
          minLength: 3
          maxLength: 50
        age:
          type: integer
          minimum: 0
          maximum: 150
        status:
          type: string
          enum:
            - active
            - inactive
```

### Output (nach `openapi-enhancer`)

```yaml
components:
  schemas:
    User:
      type: object
      required:
        - email
        - name
      properties:
        email:
          type: string
          format: email
          x-oapi-codegen-extra-tags:
            validate: "required,email"
        name:
          type: string
          minLength: 3
          maxLength: 50
          x-oapi-codegen-extra-tags:
            validate: "required,min=3,max=50"
        age:
          type: integer
          minimum: 0
          maximum: 150
          x-oapi-codegen-extra-tags:
            validate: "min=0,max=150"
        status:
          type: string
          enum:
            - active
            - inactive
          x-oapi-codegen-extra-tags:
            validate: "oneof=active inactive"
```

### Generated Go Code (with `oapi-codegen`)

```go
type User struct {
    Email  string `json:"email" validate:"required,email"`
    Name   string `json:"name" validate:"required,min=3,max=50"`
    Age    *int   `json:"age,omitempty" validate:"min=0,max=150"`
    Status string `json:"status" validate:"oneof=active inactive"`
}
```

## Makefile Integration

### Complete Workflow

```makefile
.PHONY: bundle enhance gen-code all install-tools

# Install tools
install-tools:
	@go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	@go install github.com/deinusername/openapi-enhancer/cmd/enhancer@latest
	@npm install -g @redocly/cli

# 1. Bundle OpenAPI spec
bundle:
	@redocly bundle api/openapi.yaml -o api/openapi.bundled.yaml

# 2. Add validator tags
enhance:
	@openapi-enhancer -file api/openapi.bundled.yaml

# 3. Generate Go code
gen-code: enhance
	@oapi-codegen --config=api/codegen.yaml api/openapi.bundled.yaml

# Complete workflow
all: bundle gen-code
```

### Usage

```bash
# Einmalig
make install-tools

# Bei jeder Änderung
make all
```

## Notes

- **Backup**: The tool overwrites the input file by default. Create a backup if needed.
- **Formatting**: YAML structure is preserved, but formatting may change slightly.
- **Idempotency**: The tool can be run multiple times without issues.

## Troubleshooting

### Error: `-file flag is required`

```bash
# Solution: Set flag explicitly
openapi-enhancer -file api/openapi.yaml
```

### Error: `file not found`

```bash
# Solution: Check path
ls -la api/openapi.bundled.yaml
```

### Error: `no components section found`

```bash
# Solution: Ensure your OpenAPI spec is valid
# Must contain components.schemas
```

## License

MIT License – see [LICENSE](LICENSE) file.

<!-- ## Contributing

Pull requests are welcome! For major changes, please open an issue first.

---

## Support

- **Issues**: [GitHub Issues](https://github.com/deinusername/openapi-enhancer/issues)
- **Documentation**: This README

--- -->

## 🔗 Related Tools

| Tool                                                                  | Description                          |
| --------------------------------------------------------------------- | ------------------------------------ |
| [oapi-codegen](https://github.com/deepmap/oapi-codegen)               | Generates Go code from OpenAPI specs |
| [redocly](https://github.com/Redocly/redocly-cli)                     | Bundles and validates OpenAPI specs  |
| [go-playground/validator](https://github.com/go-playground/validator) | Go validation library                |
