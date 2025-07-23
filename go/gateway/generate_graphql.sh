#!/bin/bash

# Script to generate Go code from GraphQL schema using gqlgen

set -e  # Exit on any error

echo "🚀 Generating Go code from GraphQL schema..."

# We're already in the gateway directory, no need to change directories
# The gqlgen.yml file is in the current directory

# Check if gqlgen is installed
if ! command -v gqlgen &> /dev/null; then
    echo "📦 Installing gqlgen..."
    go install github.com/99designs/gqlgen@latest
    
    # Add Go bin to PATH if not already there
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Generate Go code from GraphQL schema
echo "⚡ Running gqlgen generate..."
gqlgen generate

echo "✅ GraphQL code generation completed!"
echo "📁 Generated files:"
echo "   - graph/generated.go (GraphQL server code)"
echo "   - graph/model/models_gen.go (GraphQL models)"
echo "   - graph/schema.resolvers.go (Resolver implementations)"

# Optional: Run go mod tidy to clean up dependencies
echo "🧹 Running go mod tidy..."
go mod tidy

echo "🎉 All done! Your GraphQL Go code has been generated."
