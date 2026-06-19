#!/bin/bash

# Security and Linting Pre-commit Hook
# Runs gitleaks, semgrep, gosec, and golangci-lint before commit.

EXIT_CODE=0

echo "Running pre-commit security and linting checks..."

# 1. Gitleaks Check
if command -v gitleaks &> /dev/null; then
    echo "Running gitleaks..."
    gitleaks detect --source=. --verbose --redact
    if [ $? -ne 0 ]; then
        echo "Error: Gitleaks detected secrets!"
        EXIT_CODE=1
    fi
else
    echo "Warning: gitleaks is not installed. Skipping secrets check."
fi

# 2. Semgrep Check
if command -v semgrep &> /dev/null; then
    echo "Running semgrep..."
    semgrep --config=auto --error --quiet
    if [ $? -ne 0 ]; then
        echo "Error: Semgrep found security issues!"
        EXIT_CODE=1
    fi
else
    echo "Warning: semgrep is not installed. Skipping semgrep check."
fi

# 3. Gosec Check
if command -v gosec &> /dev/null; then
    echo "Running gosec..."
    gosec -quiet ./...
    if [ $? -ne 0 ]; then
        echo "Error: Gosec found vulnerabilities!"
        EXIT_CODE=1
    fi
else
    echo "Warning: gosec is not installed. Skipping gosec check."
fi

# 4. Golangci-lint Check
if command -v golangci-lint &> /dev/null; then
    echo "Running golangci-lint..."
    golangci-lint run
    if [ $? -ne 0 ]; then
        echo "Error: Golangci-lint checks failed!"
        EXIT_CODE=1
    fi
else
    echo "Warning: golangci-lint is not installed. Skipping lint check."
fi

if [ $EXIT_CODE -ne 0 ]; then
    echo "Pre-commit checks failed! Please fix the errors before committing."
    exit 1
fi

echo "All pre-commit checks passed successfully!"
exit 0
