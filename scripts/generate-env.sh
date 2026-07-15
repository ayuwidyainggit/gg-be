#!/bin/bash

###############################################################################
# Generate .env Files for All Services
# Generate .env file untuk semua service dengan konfigurasi yang sesuai
#
# Usage:
#   ./scripts/generate-env.sh [docker|local]
#   ./scripts/generate-env.sh docker   # Connect ke Docker services
#   ./scripts/generate-env.sh local    # Connect ke local PostgreSQL
#
###############################################################################

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Get mode
MODE=${1:-docker}

if [ "$MODE" != "docker" ] && [ "$MODE" != "local" ]; then
    echo "Usage: $0 [docker|local]"
    echo "  docker - Generate .env untuk connect ke Docker services"
    echo "  local  - Generate .env untuk connect ke local PostgreSQL"
    exit 1
fi

print_info "Generating .env files for all services (mode: $MODE)"
echo ""

# Service ports mapping
declare -A SERVICE_PORTS=(
    ["system"]="9001"
    ["master"]="9002"
    ["inventory"]="9003"
    ["sales"]="9004"
    ["finance"]="9005"
    ["mobile"]="9008"
    ["cronjob"]="9100"
)

# Database configuration based on mode
if [ "$MODE" == "docker" ]; then
    DB_HOST="localhost"
    DB_PORT="54321"
    DB_USER="postgres"
    DB_PASS="postgres"
    print_info "Mode: Docker (DB_PORT=54321)"
else
    DB_HOST="localhost"
    DB_PORT="5432"
    read -p "Local PostgreSQL User [default: postgres]: " DB_USER
    DB_USER=${DB_USER:-postgres}
    read -sp "Local PostgreSQL Password: " DB_PASS
    echo
    print_info "Mode: Local (DB_PORT=5432)"
fi

# Generate .env for each service
for service in system master inventory sales finance mobile cronjob; do
    SERVICE_PORT=${SERVICE_PORTS[$service]}
    
    print_info "Generating .env for $service service..."
    
    ENV_FILE="../$service/.env"
    
    # Create .env file
    cat > "$ENV_FILE" <<EOF
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=$SERVICE_PORT

# Database Configuration
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASS=$DB_PASS
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Huawei OBS Configuration (optional)
OBS_HUAWEI_AK=
OBS_HUAWEI_SK=
OBS_HUAWEI_ENDPOINT=
OBS_HUAWEI_BUCKET=
EOF

    # Add service-specific config
    if [ "$service" == "cronjob" ]; then
        echo "" >> "$ENV_FILE"
        echo "# Cronjob Specific" >> "$ENV_FILE"
        echo "INTERVAL_RELOAD_JOBS_IN_SECOND=30" >> "$ENV_FILE"
    fi
    
    print_success "Generated $ENV_FILE"
done

echo ""
print_success "All .env files generated successfully!"
echo ""
print_info "Next steps:"
echo "  1. Review and update .env files if needed"
echo "  2. Add OBS credentials if required"
echo "  3. Start services: go run main.go"
echo ""

