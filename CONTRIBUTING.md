# Contributing to Text2SQL Skill Engine

Thank you for your interest in contributing to the Text2SQL Skill Engine! This document provides guidelines and instructions for contributing to the project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)

## 📜 Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to ensure a welcoming environment for everyone.

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- MySQL or PostgreSQL (for testing)

### Setting Up Development Environment

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/text2sql-skill.git
   cd text2sql-skill
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/original-owner/text2sql-skill.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Create a development branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 🔧 Development Workflow

### Branch Naming Convention
- `feature/` - New features or enhancements
- `bugfix/` - Bug fixes
- `hotfix/` - Critical production fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test improvements

### Commit Messages
Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat(security): add input validation for SQL injection prevention

- Add maximum length validation
- Add entropy analysis
- Add forbidden keyword detection

Closes #123
```

## 💻 Code Style

### Go Code
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use meaningful variable and function names
- Add comments for exported functions and types

### Configuration Files
- Use YAML format for configuration
- Include both English and Chinese comments
- Follow the existing structure in `config.yaml.example`

### Error Handling
- Use proper error wrapping with `fmt.Errorf`
- Log errors at appropriate levels
- Provide meaningful error messages

## 🧪 Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test suite
go test ./tests/...

# Run integration tests (requires database)
go test -tags=integration ./tests/...
```

### Writing Tests
- Write unit tests for new functionality
- Use table-driven tests when appropriate
- Mock external dependencies
- Test both success and failure cases

### Test Structure
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "result",
            wantErr:  false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.expected {
                t.Errorf("FunctionName() = %v, expected %v", got, tt.expected)
            }
        })
    }
}
```

## 📚 Documentation

### Code Documentation
- Document all exported functions, types, and methods
- Use clear, concise comments
- Include examples for complex functions

### User Documentation
- Update README.md for user-facing changes
- Add configuration examples
- Document new features

### API Documentation
- Document REST endpoints (if applicable)
- Include request/response examples
- Document error codes and messages

## 🔄 Pull Request Process

1. **Ensure tests pass**:
   ```bash
   go test ./...
   ```

2. **Update documentation** if needed

3. **Update CHANGELOG.md** with your changes

4. **Create Pull Request**:
   - Use a clear, descriptive title
   - Reference related issues
   - Describe changes in detail
   - Include test results

5. **Address review comments** promptly

6. **Wait for CI checks to pass**

7. **Merge after approval**

### PR Template
```markdown
## Description
Brief description of the changes

## Related Issues
Fixes #123

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Test update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] All tests pass

## Documentation
- [ ] README updated
- [ ] Code comments added
- [ ] API documentation updated

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] All tests pass
- [ ] Documentation updated
```

## 🐛 Issue Reporting

### Bug Reports
When reporting bugs, include:
1. **Description**: Clear description of the issue
2. **Steps to Reproduce**: Detailed reproduction steps
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**: OS, Go version, database version
6. **Logs**: Relevant error logs
7. **Screenshots**: If applicable

### Feature Requests
When requesting features, include:
1. **Use Case**: Why this feature is needed
2. **Proposed Solution**: How it should work
3. **Alternatives Considered**: Other approaches
4. **Additional Context**: Any other relevant information

## 🏗️ Project Structure

```
text2sql-skill/
├── config/          # Configuration structures and validation
├── core/           # Core engine components
├── drivers/        # Database drivers
├── interfaces/     # Public interfaces
├── tests/         # Test suites
├── utils/         # Utility functions
├── main.go        # Application entry point
└── README.md      # Project documentation
```

## 🛠️ Development Tools

### Recommended Tools
- **Go**: Official Go tools
- **gofmt**: Code formatting
- **go vet**: Static analysis
- **golangci-lint**: Linting
- **Docker**: Containerization

### Editor Configuration
- Use Go modules support
- Enable auto-formatting on save
- Install Go language server

## 🤝 Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **Discussions**: For questions and discussions
- **Documentation**: Check README.md and code comments

## 📄 License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).

---

Thank you for contributing to Text2SQL Skill Engine! 🎉
