#!/bin/bash
set -e

XVFB_ENABLED="${BROWSER_XVFB_ENABLED:-false}"
XVFB_DISPLAY="${BROWSER_XVFB_DISPLAY:-:99}"
XVFB_SCREEN="${BROWSER_XVFB_SCREEN_RESOLUTION:-1920x1080x24}"

BG_PIDS=""
cleanup() {
    [ -n "$BG_PIDS" ] || return 0
    echo "Stopping background processes:$BG_PIDS"
    kill $BG_PIDS 2>/dev/null || true
}
trap cleanup EXIT

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

if [ "$XVFB_ENABLED" = "true" ]; then
    echo "Starting Xvfb on display $XVFB_DISPLAY with screen resolution $XVFB_SCREEN"

    DISPLAY_NUM=$(echo "$XVFB_DISPLAY" | sed 's/://g')
    LOCK_FILE="/tmp/.X${DISPLAY_NUM}-lock"
    if [ -f "$LOCK_FILE" ]; then
        echo "Removing stale lock file: $LOCK_FILE"
        rm -f "$LOCK_FILE"
    fi

    if [ "${BROWSER_XVFB_ALLOW_TCP:-false}" = "true" ]; then
        echo "Starting Xvfb with TCP listening enabled (for VNC)"
        Xvfb "$XVFB_DISPLAY" -screen 0 "$XVFB_SCREEN" -listen tcp -ac &
    else
        echo "Starting Xvfb with TCP disabled (secure mode)"
        Xvfb "$XVFB_DISPLAY" -screen 0 "$XVFB_SCREEN" -nolisten tcp &
    fi
    XVFB_PID=$!

    if wait_for_xvfb; then
        export DISPLAY="$XVFB_DISPLAY"
        BG_PIDS="$BG_PIDS $XVFB_PID"
        echo "Xvfb started successfully (PID: $XVFB_PID)"
    else
        echo "Error: Xvfb failed to start properly"
        kill "$XVFB_PID" 2>/dev/null || true
        exit 1
    fi
else
    echo "Xvfb disabled, using native headless mode"
fi

wait_for_port() {
    local attempt=0

    while [ $attempt -lt 20 ]; do
        if (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null; then
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 0.5
    done

    return 1
}

if [ "${BROWSER_ENGINE:-chromium}" = "lightpanda" ]; then
    case "${BROWSER_CDP_URL:-}" in
        "" | ws://127.0.0.1:* | ws://localhost:* | http://127.0.0.1:* | http://localhost:*)
            export BROWSER_CDP_URL="${BROWSER_CDP_URL:-ws://127.0.0.1:9222}"
            LP_PORT="${BROWSER_CDP_URL##*:}"

            echo "Starting bundled Lightpanda on port $LP_PORT"
            lightpanda serve --host 127.0.0.1 --port "$LP_PORT" &
            BG_PIDS="$BG_PIDS $!"

            if ! wait_for_port "$LP_PORT"; then
                echo "Error: Lightpanda failed to accept connections on port $LP_PORT"
                exit 1
            fi
            echo "Lightpanda is ready"
            ;;
        *)
            echo "Using remote CDP endpoint, not starting the bundled Lightpanda"
            ;;
    esac
fi

echo "Browser configuration:"
echo "  Engine: ${BROWSER_ENGINE:-chromium}"
echo "  Headless: ${BROWSER_HEADLESS:-true}"
echo "  Stealth Mode: ${BROWSER_STEALTH_MODE:-false}"
echo "  Xvfb Enabled: $XVFB_ENABLED"
if [ -n "${BROWSER_CDP_URL:-}" ]; then
    echo "  CDP URL: $BROWSER_CDP_URL"
fi

exec ./main start
