#!/bin/bash

###############################################################################
# Setup Air Configuration untuk semua service
# Copy .air.toml template ke setiap service directory
#
# Usage:
#   ./docker/setup-air.sh
#
###############################################################################

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Setting up Air Hot Reload${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

cd "$(dirname "$0")/.."

SERVICES=("system" "master" "inventory" "sales" "finance" "mobile" "cronjob")
TEMPLATE_FILE="docker/.air.toml.template"

if [ ! -f "$TEMPLATE_FILE" ]; then
    echo -e "${YELLOW}Error: Template file not found: $TEMPLATE_FILE${NC}"
    exit 1
fi

for service in "${SERVICES[@]}"; do
    if [ -d "$service" ]; then
        target_file="$service/.air.toml"
        if [ ! -f "$target_file" ]; then
            cp "$TEMPLATE_FILE" "$target_file"
            echo -e "${GREEN}✓${NC} Created $target_file"
        else
            echo -e "${YELLOW}⚠${NC} $target_file already exists, skipping"
        fi
    else
        echo -e "${YELLOW}⚠${NC} Service directory not found: $service"
    fi
done

echo ""
echo -e "${GREEN}Setup complete!${NC}"
echo ""
echo -e "Air configuration files created for all services."
echo -e "Now you can use hot reload with: ./docker/start-dev.sh"
echo ""

