#!/bin/bash

frontend_pid=""
backend_pid=""

# Function to handle the termination of the script
cleanup() {
    echo "Stopping servers..."
    if [ ! -z "$frontend_pid" ]; then
        kill $frontend_pid
    fi
    if [ ! -z "$backend_pid" ]; then
        kill $backend_pid
    fi
    exit
}

# Trap the SIGINT signal (CTRL+C) and call the cleanup function
trap cleanup SIGINT

# Allow non-root users to bind to privileged ports (0-1023)
sudo sysctl net.ipv4.ip_unprivileged_port_start=80

# Create frontend.log and backend.log if they don't exist
touch ./frontend/frontend.log
touch ./backend/backend.log

# Navigate to frontend directory and run npm commands
cd frontend
npm install
npm run dev &> frontend.log &
frontend_pid=$!

# Navigate to backend directory and run go command
cd ../backend
go run . &> backend.log &
backend_pid=$!

# Tail the logs in the foreground to see the outputs
tail -f ../frontend/frontend.log -f ../backend/backend.log

