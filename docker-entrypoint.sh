#!/bin/bash
set -e

# Configuration from environment variables
XVFB_ENABLED="${BROWSER_XVFB_ENABLED:-false}"
XVFB_DISPLAY="${BROWSER_XVFB_DISPLAY:-:99}"
XVFB_SCREEN="${BROWSER_XVFB_SCREEN_RESOLUTION:-1920x1080x24}"

# Function to check if Xvfb is ready
wait_for_xvfb() {
    local max_attempts=10
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if xdpyinfo -display "$XVFB_DISPLAY" >/dev/null 2>&1; then
            echo "Xvfb is ready on display $XVFB_DISPLAY"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 0.5
    done

    echo "Warning: Xvfb failed to start within timeout"
    return 1
}

# Start Xvfb if enabled
if [ "$XVFB_ENABLED" = "true" ]; then
    echo "Starting Xvfb on display $XVFB_DISPLAY with screen resolution $XVFB_SCREEN"

    # Clean up stale lock files from previous runs
    DISPLAY_NUM=$(echo "$XVFB_DISPLAY" | sed 's/://g')
    LOCK_FILE="/tmp/.X${DISPLAY_NUM}-lock"
    if [ -f "$LOCK_FILE" ]; then
        echo "Removing stale lock file: $LOCK_FILE"
        rm -f "$LOCK_FILE"
    fi

    # Start Xvfb without -ac flag for security
    # Use -nolisten tcp to prevent network access
    Xvfb "$XVFB_DISPLAY" -screen 0 "$XVFB_SCREEN" -nolisten tcp &
    XVFB_PID=$!

    # Wait for Xvfb to be ready
    if wait_for_xvfb; then
        export DISPLAY="$XVFB_DISPLAY"
        echo "Xvfb started successfully (PID: $XVFB_PID)"
    else
        echo "Error: Xvfb failed to start properly"
        kill "$XVFB_PID" 2>/dev/null || true
        exit 1
    fi

    # Trap to cleanup Xvfb on exit
    trap "echo 'Stopping Xvfb...'; kill $XVFB_PID 2>/dev/null || true" EXIT
else
    echo "Xvfb disabled, using native headless mode"
fi

# Log configuration
echo "Browser configuration:"
echo "  Engine: ${BROWSER_ENGINE:-chromium}"
echo "  Headless: ${BROWSER_HEADLESS:-true}"
echo "  Stealth Mode: ${BROWSER_STEALTH_MODE:-false}"
echo "  Xvfb Enabled: $XVFB_ENABLED"

# Start the main application
exec ./main
