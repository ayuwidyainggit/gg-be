#!/bin/bash

###############################################################################
# Start Development Mode
# Start semua service dengan volume mounting untuk hot reload
# Code changes di local langsung terupdate di container
#
# Usage:
#   ./docker/start-dev.sh
#
###############################################################################

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Starting Development Mode${NC}"
echo -e "${BLUE}  (Hot Reload Enabled)${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

cd "$(dirname "$0")"

# Check if Air config exists (optional check)
if [ ! -f "../system/.air.toml" ]; then
    echo -e "${YELLOW}Warning: Air config not found${NC}"
    echo "Run ./docker/setup-air.sh to setup Air hot reload"
    echo ""
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}Warning: .env file not found${NC}"
    echo "Creating .env from .env.example..."
    cp .env.example .env 2>/dev/null || true
    echo ""
fi

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_info "Starting services in development mode..."
print_info "Code changes will be automatically reflected in containers"
echo ""

# Start all services
docker-compose -f docker-compose.dev.yml up -d

echo ""
print_success "All services started in development mode!"
echo ""
echo -e "Service URLs:"
echo -e "  ${BLUE}System:${NC}    http://localhost:9001/ping"
echo -e "  ${BLUE}Master:${NC}     http://localhost:9002/ping"
echo -e "  ${BLUE}Inventory:${NC}  http://localhost:9003/ping"
echo -e "  ${BLUE}Sales:${NC}      http://localhost:9004/ping"
echo -e "  ${BLUE}Finance:${NC}    http://localhost:9005/ping"
echo -e "  ${BLUE}Mobile:${NC}     http://localhost:9008/ping"
echo -e "  ${BLUE}Cronjob:${NC}    http://localhost:9100/ping"
echo ""
echo -e "Database:"
echo -e "  ${BLUE}PostgreSQL:${NC} localhost:54321"
echo -e "  ${BLUE}Redis:${NC}      localhost:6379"
echo ""
echo -e "${YELLOW}Note:${NC} Services use Air for auto-reload when code changes"
echo -e "${YELLOW}      ${NC} Make sure to run ./docker/setup-air.sh first (if not done)"
echo ""
echo -e "Useful commands:"
echo -e "  ${BLUE}View logs:${NC}     docker-compose -f docker-compose.dev.yml logs -f [service]"
echo -e "  ${BLUE}Stop all:${NC}      docker-compose -f docker-compose.dev.yml down"
echo -e "  ${BLUE}Restart service:${NC} docker-compose -f docker-compose.dev.yml restart [service]"
echo ""

