#!/bin/bash

###############################################################################
# Start PostgreSQL with Citus Extension (Docker)
# Quick script untuk start PostgreSQL dengan Citus di Docker
#
# Usage:
#   ./docker/start-citus.sh
#
###############################################################################

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}Starting PostgreSQL with Citus extension...${NC}"

cd "$(dirname "$0")"

# Start container
docker-compose -f docker-compose.postgres-citus.yml up -d

echo -e "${GREEN}Container started!${NC}"
echo ""
echo -e "Connection info:"
echo -e "  Host: localhost"
echo -e "  Port: ${YELLOW}54321${NC}"
echo -e "  User: postgres"
echo -e "  Password: postgres"
echo -e "  Database: scylla_db"
echo ""
echo -e "Connect using:"
echo -e "  ${BLUE}psql -h localhost -p 54321 -U postgres -d scylla_db${NC}"
echo ""
echo -e "Or using Docker:"
echo -e "  ${BLUE}docker exec -it scylla-postgres-citus psql -U postgres -d scylla_db${NC}"
echo ""
echo -e "View logs:"
echo -e "  ${BLUE}docker-compose -f docker-compose.postgres-citus.yml logs -f${NC}"

