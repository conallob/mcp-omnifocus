# Release Process

This document describes how to create a new release of mcp-omnifocus.

## Prerequisites

- Push access to the GitHub repository
- All changes committed and pushed to `main` branch
- CHANGELOG or release notes prepared

## Creating a Release

The release process is automated using GitHub Actions and GoReleaser v2.

### 1. Create and Push a Version Tag

```bash
# Ensure you're on the main branch and up to date
git checkout main
git pull origin main

# Create a new tag (following semantic versioning)
git tag -a v0.1.0 -m "Release v0.1.0"

# Push the tag to GitHub
git push origin v0.1.0
```

### 2. Automatic Build

Once you push the tag, GitHub Actions will automatically:

1. Checkout the code
2. Set up Go 1.21
3. Run GoReleaser v2
4. Build binaries for:
   - macOS AMD64 (Intel Macs)
   - macOS ARM64 (Apple Silicon Macs)
5. Create archives with:
   - Binary
   - JXA scripts
   - README, LICENSE, QUICKSTART
   - Config example
6. Generate checksums
7. Create a GitHub release with:
   - Release notes
   - Downloadable binaries
   - Changelog

### 3. Monitor the Build

1. Go to: https://github.com/conall/mcp-omnifocus/actions
2. Find the "Release" workflow run for your tag
3. Monitor the build progress
4. If successful, the release will appear at: https://github.com/conall/mcp-omnifocus/releases

### 4. Verify the Release

After the workflow completes:

1. Check the release page
2. Verify both binaries are present:
   - `mcp-omnifocus_VERSION_darwin_x86_64.tar.gz` (Intel)
   - `mcp-omnifocus_VERSION_darwin_arm64.tar.gz` (Apple Silicon)
3. Download and test one of the archives
4. Verify checksums.txt is present

## Installation and Usage

### Installing via Homebrew

Users can install the server via Homebrew:

```bash
brew tap conall/mcp-omnifocus
brew install mcp-omnifocus
```

This installs:
- Binary: `/opt/homebrew/bin/mcp-omnifocus` (Apple Silicon) or `/usr/local/bin/mcp-omnifocus` (Intel)
- Scripts: `/opt/homebrew/share/mcp-omnifocus/scripts/` (Apple Silicon) or `/usr/local/share/mcp-omnifocus/scripts/` (Intel)

### Installing from Release Archive

Users can also download the release archive and extract it:

```bash
# Download the appropriate archive for your architecture
curl -LO https://github.com/conall/mcp-omnifocus/releases/download/v0.1.0/mcp-omnifocus_0.1.0_darwin_arm64.tar.gz

# Extract
tar -xzf mcp-omnifocus_0.1.0_darwin_arm64.tar.gz

# Move to desired location
mv mcp-omnifocus /usr/local/bin/
mv scripts /usr/local/share/mcp-omnifocus/
```

### Configuring Claude Desktop

Users must configure Claude Desktop to use the server. The configuration depends on the installation method:

**For Homebrew Installation (Apple Silicon):**
```json
{
  "mcpServers": {
    "omnifocus": {
      "command": "/opt/homebrew/bin/mcp-omnifocus",
      "args": ["-scripts", "/opt/homebrew/share/mcp-omnifocus/scripts"]
    }
  }
}
```

**For Homebrew Installation (Intel Mac):**
```json
{
  "mcpServers": {
    "omnifocus": {
      "command": "/usr/local/bin/mcp-omnifocus",
      "args": ["-scripts", "/usr/local/share/mcp-omnifocus/scripts"]
    }
  }
}
```

**For Manual Download/Extraction:**
```json
{
  "mcpServers": {
    "omnifocus": {
      "command": "/path/to/mcp-omnifocus",
      "args": ["-scripts", "/path/to/scripts"]
    }
  }
}
```

**Important**: The `-scripts` flag explicitly specifies where the JXA automation scripts are located. While the server includes auto-detection logic that works for most installation scenarios, using the `-scripts` flag ensures reliability across all installation methods and is the recommended approach.

## Semantic Versioning

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v1.0.0 → v2.0.0): Breaking changes
- **MINOR** (v0.1.0 → v0.2.0): New features, backward compatible
- **PATCH** (v0.1.0 → v0.1.1): Bug fixes, backward compatible

Examples:
- `v0.1.0` - Initial release
- `v0.2.0` - Added new MCP tools
- `v0.1.1` - Fixed bug in task creation
- `v1.0.0` - First stable release

## Pre-releases

For alpha/beta releases, use suffixes:

```bash
git tag -a v0.2.0-alpha.1 -m "Alpha release for testing"
git push origin v0.2.0-alpha.1
```

GoReleaser will automatically mark these as pre-releases.

## Release Notes Template

When creating the tag message, use this template:

```
Release vX.Y.Z

## New Features
- Feature 1
- Feature 2

## Bug Fixes
- Fix 1
- Fix 2

## Improvements
- Improvement 1

## Breaking Changes (if any)
- Breaking change description
```

## Troubleshooting

### Build Fails

1. Check the Actions log for errors
2. Common issues:
   - Go compilation errors (check your code compiles locally first)
   - Missing dependencies (run `make deps` locally)
   - GoReleaser configuration errors

### Tag Already Exists

If you need to re-release the same version:

```bash
# Delete local tag
git tag -d v0.1.0

# Delete remote tag
git push origin :refs/tags/v0.1.0

# Create new tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

**Warning**: Only do this if the release hasn't been published yet or you have a very good reason.

### Release is a Draft

If you want releases to be drafts for manual review before publishing, edit `.goreleaser.yaml`:

```yaml
release:
  draft: true  # Change from false to true
```

## Manual Release (Not Recommended)

If you need to release manually:

```bash
# Install GoReleaser
brew install goreleaser

# Create a GitHub token with repo access
export GITHUB_TOKEN=your_token_here

# Run GoReleaser locally
goreleaser release --clean
```

## Testing GoReleaser Config

To test the GoReleaser configuration without creating a release:

```bash
# Build snapshot (no tag required)
goreleaser build --snapshot --clean

# Check the dist/ directory for output
ls -la dist/
```

This builds the binaries without creating a release, useful for testing changes to `.goreleaser.yaml`.
