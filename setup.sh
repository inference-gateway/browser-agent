#!/bin/bash

# Setup script for playwright-agent development environment
# This script installs the necessary system dependencies for Playwright

set -e

echo "🎭 Setting up Playwright Agent development environment..."

# Check if running on supported OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "✅ Detected Linux environment"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "✅ Detected macOS environment"
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
    echo "✅ Detected Windows environment"
else
    echo "⚠️  Unsupported OS: $OSTYPE"
    echo "Please install Playwright dependencies manually:"
    echo "  npx playwright install-deps"
    exit 1
fi

# Check if npx is available
if ! command -v npx &> /dev/null; then
    echo "❌ npx is not available. Please install Node.js first."
    echo "Visit: https://nodejs.org/"
    exit 1
fi

echo "📦 Installing Playwright system dependencies..."

# Install Playwright dependencies
if npx playwright install-deps; then
    echo "✅ Playwright system dependencies installed successfully!"
else
    echo "⚠️  Failed to install with npx playwright install-deps"
    echo "Try installing manually:"
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "  sudo apt-get update && sudo apt-get install -y \\"
        echo "    libxcursor1 libgtk-3-0t64 libpangocairo-1.0-0 \\"
        echo "    libcairo-gobject2 libgdk-pixbuf-2.0-0"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "  brew install playwright"
        echo "  npx playwright install"
    fi
    exit 1
fi

echo ""
echo "🎉 Setup complete! You can now run the agent with:"
echo "  go run ."
echo ""
echo "Or use the task runner:"
echo "  task run"