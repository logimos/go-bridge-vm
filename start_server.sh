#!/bin/bash

# Server management script for myllm

SERVER_PORT=8080
PID_FILE="/tmp/myllm.pid"

# Function to check if server is running
is_server_running() {
    if [ -f "$PID_FILE" ]; then
        local pid=$(cat "$PID_FILE")
        if kill -0 "$pid" 2>/dev/null; then
            return 0
        else
            rm -f "$PID_FILE"
        fi
    fi
    return 1
}

# Function to stop server
stop_server() {
    if is_server_running; then
        local pid=$(cat "$PID_FILE")
        echo "Stopping server (PID: $pid)..."
        kill -TERM "$pid"
        
        # Wait for graceful shutdown
        local count=0
        while kill -0 "$pid" 2>/dev/null && [ $count -lt 10 ]; do
            sleep 1
            count=$((count + 1))
        done
        
        # Force kill if still running
        if kill -0 "$pid" 2>/dev/null; then
            echo "Force killing server..."
            kill -KILL "$pid"
        fi
        
        rm -f "$PID_FILE"
        echo "Server stopped."
    else
        echo "Server is not running."
    fi
}

# Function to start server
start_server() {
    if is_server_running; then
        echo "Server is already running (PID: $(cat $PID_FILE))"
        return 1
    fi
    
    echo "Starting server..."
    
    # Set environment variables
    export AI_PROVIDER=enhanced_local
    export INTENT_CONFIG_PATH=configs/personal_assistant.json
    
    # Start server and capture PID
    ./bin/myllm &
    local pid=$!
    echo $pid > "$PID_FILE"
    
    # Wait a moment to see if it starts successfully
    sleep 2
    if kill -0 "$pid" 2>/dev/null; then
        echo "Server started successfully (PID: $pid)"
        echo "Server is running on http://localhost:$SERVER_PORT"
        echo "Use './start_server.sh stop' to stop the server"
    else
        echo "Failed to start server"
        rm -f "$PID_FILE"
        return 1
    fi
}

# Function to restart server
restart_server() {
    echo "Restarting server..."
    stop_server
    sleep 1
    start_server
}

# Function to show status
show_status() {
    if is_server_running; then
        local pid=$(cat "$PID_FILE")
        echo "Server is running (PID: $pid)"
        echo "Server URL: http://localhost:$SERVER_PORT"
    else
        echo "Server is not running"
    fi
}

# Main script logic
case "${1:-start}" in
    start)
        start_server
        ;;
    stop)
        stop_server
        ;;
    restart)
        restart_server
        ;;
    status)
        show_status
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        echo "  start   - Start the server"
        echo "  stop    - Stop the server"
        echo "  restart - Restart the server"
        echo "  status  - Show server status"
        exit 1
        ;;
esac 