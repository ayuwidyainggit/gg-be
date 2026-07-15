#!/bin/bash

###############################################################################
# Start All Scylla Backend Services
# Quick script untuk start semua service di Docker
#
# Usage:
#   ./docker/start-all.sh
#
###############################################################################

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Starting All Scylla Services${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

cd "$(dirname "$0")"

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}Warning: .env file not found${NC}"
    echo "Creating .env from .env.example..."
    cp .env.example .env 2>/dev/null || true
    echo ""
fi

# Start all services
echo -e "${BLUE}Building and starting all services...${NC}"
docker-compose up -d --build

echo ""
echo -e "${GREEN}All services started!${NC}"
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
echo -e "Useful commands:"
echo -e "  ${BLUE}View logs:${NC}     docker-compose logs -f"
echo -e "  ${BLUE}Stop all:${NC}      docker-compose down"
echo -e "  ${BLUE}Restart:${NC}       docker-compose restart"
echo ""

