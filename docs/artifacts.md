# A2A Artifacts Implementation

This document describes the implementation of A2A (Agent-to-Agent) Artifacts support in the browser-agent.

## Overview

The A2A Artifacts feature allows the browser-agent to create, store, and serve artifacts generated during skill execution. Artifacts can be screenshots, data exports, or any other files created during browser automation tasks.

## Architecture

### Components

1. **Artifact Server** (`internal/artifacts/server.go`)
   - HTTP server running on port 8081 (configurable)
   - Serves artifacts via REST endpoints
   - Provides artifact metadata and file downloads

2. **Artifact Registry** (`internal/artifacts/server.go`)
   - In-memory registry mapping artifact IDs to file information
   - Thread-safe operations with mutex protection
   - Stores metadata, file paths, and artifact properties

3. **Enhanced Artifact Helper** (`internal/artifacts/helper.go`)
   - Extends ADK's ArtifactHelper with registry integration
   - Creates different artifact types (FilePart, TextPart, DataPart)
   - Handles file storage and registry updates

4. **Global Manager** (`internal/artifacts/manager.go`)
   - Provides global access to artifact services
   - Singleton pattern for server-wide artifact management
   - Thread-safe initialization and access

## Configuration

### Environment Variables

```bash
# Enable/disable artifacts server
ARTIFACTS_ENABLED=true

# Port for artifacts server (default: 8081)
ARTIFACTS_PORT=8081

# Base URL for artifact access
ARTIFACTS_BASE_URL=http://localhost:8081
```

### agent.yaml Configuration

```yaml
spec:
  config:
    artifacts:
      enabled: true
      port: 8081
      base_url: "http://localhost:8081"
```

## API Endpoints

### 1. Download Artifact
```http
GET /artifacts/{artifactId}
```

**Response**: Binary file content with appropriate MIME type headers

**Example**:
```bash
curl http://localhost:8081/artifacts/screenshot_12345
```

### 2. Get Artifact Metadata
```http
GET /artifacts/{artifactId}/metadata
```

**Response**:
```json
{
  "id": "screenshot_12345",
  "file_path": "/tmp/artifacts/screenshot.png",
  "file_name": "screenshot.png",
  "mime_type": "image/png",
  "size": 45678,
  "created_at": "2025-09-24T12:00:00Z",
  "title": "Screenshot: screenshot.png",
  "description": "Screenshot captured from browser session",
  "metadata": {
    "full_page": false,
    "image_type": "png",
    "quality": 80
  }
}
```

### 3. List All Artifacts
```http
GET /artifacts/
```

**Response**:
```json
{
  "artifacts": [
    {
      "id": "screenshot_12345",
      "file_name": "screenshot.png",
      "mime_type": "image/png",
      "size": 45678,
      "created_at": "2025-09-24T12:00:00Z",
      "title": "Screenshot: screenshot.png"
    }
  ],
  "count": 1
}
```

### 4. Health Check
```http
GET /health
```

**Response**:
```json
{
  "status": "healthy",
  "server": "artifacts",
  "port": 8081
}
```

## Artifact Types

### FilePart
Binary files such as screenshots, PDFs, or other documents.

**Properties**:
- `type`: "FilePart"
- `fileUri`: URL to download the file
- `mimeType`: MIME type of the file
- `filename`: Original filename

### TextPart
Plain text content.

**Properties**:
- `type`: "TextPart"
- `content`: Text content
- `mimeType`: "text/plain"

### DataPart
Structured data in JSON format.

**Properties**:
- `type`: "DataPart"
- `data`: JSON object
- `mimeType`: "application/json"

## Usage in Skills

Skills automatically register artifacts when they create files. Here's how the `take_screenshot` skill integrates:

```go
// Create artifact using ADK helper
screenshotArtifact := s.artifactHelper.CreateFileArtifactFromBytes(
    fmt.Sprintf("Screenshot: %s", filename),
    fmt.Sprintf("Screenshot captured from browser session %s", session.ID),
    filename,
    screenshotData,
    &mimeType,
)

// Register with global artifact manager
s.registerWithGlobalManager(screenshotArtifact, generatedPath, filename, mimeType, metadata)
```

## File Storage

Artifacts are stored in the filesystem at the configured data directory:

**Default Location**: `/tmp/playwright/artifacts/`

**Directory Structure**:
```
/tmp/artifacts/
├── playwright/           # Original files from skills
│   ├── screenshot_1.png
│   └── data_export.csv
└── runtime/             # Generated artifacts
    └── {artifact_id}/
        └── {filename}
```

## Integration with A2A Protocol

### Current Implementation

The artifacts infrastructure is fully implemented and integrated with:

- **Skill Execution**: Artifacts are automatically created and registered
- **File Serving**: HTTP server provides access to all artifacts
- **Metadata Management**: Rich metadata support for all artifact types

### Future Enhancement: `includeArtifacts` Parameter

The `includeArtifacts` parameter for task polling is prepared for implementation but requires deeper integration with the ADK framework's task handling system. The current infrastructure supports this feature once the appropriate ADK hooks are available.

**Planned Usage**:
```json
{
  "method": "task.get",
  "params": {
    "taskId": "task_12345",
    "includeArtifacts": true
  }
}
```

**Expected Response**:
```json
{
  "result": {
    "taskId": "task_12345",
    "status": "completed",
    "artifacts": [
      {
        "type": "FilePart",
        "artifactId": "screenshot_12345",
        "title": "Screenshot",
        "description": "Page screenshot",
        "fileUri": "http://localhost:8081/artifacts/screenshot_12345",
        "mimeType": "image/png",
        "size": 45678
      }
    ]
  }
}
```

## Testing

Run the artifact tests:

```bash
# Run all artifact tests
go test -v ./internal/artifacts/...

# Run specific test
go test -v ./internal/artifacts/ -run TestArtifactServer
```

## Security Considerations

1. **File Access**: Only registered artifacts are accessible via the API
2. **Path Validation**: All file paths are validated before serving
3. **MIME Type Detection**: Proper MIME types are set for security
4. **No Directory Traversal**: Artifact IDs cannot be used to access arbitrary files

## Troubleshooting

### Common Issues

1. **Server Not Starting**
   - Check if port 8081 is available
   - Verify ARTIFACTS_ENABLED=true
   - Check logs for startup errors

2. **Artifacts Not Found**
   - Verify artifact ID is correct
   - Check if artifact was properly registered
   - Ensure file exists on disk

3. **Permission Issues**
   - Verify write permissions on data directory
   - Check file permissions for created artifacts

### Debug Endpoints

```bash
# Check server health
curl http://localhost:8081/health

# List all artifacts
curl http://localhost:8081/artifacts/

# Get artifact metadata
curl http://localhost:8081/artifacts/{id}/metadata
```

## Configuration Examples

### Docker Compose

```yaml
version: '3.8'
services:
  browser-agent:
    image: browser-agent:latest
    ports:
      - "8080:8080"    # A2A server
      - "8081:8081"    # Artifacts server
    environment:
      - ARTIFACTS_ENABLED=true
      - ARTIFACTS_PORT=8081
      - ARTIFACTS_BASE_URL=http://localhost:8081
    volumes:
      - ./artifacts:/tmp/playwright/artifacts
```

### Kubernetes

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: browser-agent-config
data:
  ARTIFACTS_ENABLED: "true"
  ARTIFACTS_PORT: "8081"
  ARTIFACTS_BASE_URL: "http://browser-agent-artifacts:8081"
---
apiVersion: v1
kind: Service
metadata:
  name: browser-agent-artifacts
spec:
  ports:
  - port: 8081
    targetPort: 8081
  selector:
    app: browser-agent
```

## Performance Considerations

1. **Memory Usage**: Artifact registry is in-memory; consider Redis for high-volume scenarios
2. **File Storage**: Monitor disk usage in the artifacts directory
3. **Concurrent Access**: Server handles concurrent requests efficiently
4. **Cleanup**: Implement artifact cleanup policies for long-running deployments

## Future Enhancements

1. **Persistent Registry**: Redis or database backing for the artifact registry
2. **Artifact Expiration**: TTL-based cleanup for old artifacts
3. **Authentication**: Access control for artifact downloads
4. **Compression**: Automatic compression for large artifacts
5. **Cloud Storage**: S3/GCS integration for scalable storage
6. **Artifact Versioning**: Support for multiple versions of the same artifact