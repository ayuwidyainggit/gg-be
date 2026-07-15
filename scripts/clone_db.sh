#!/bin/bash

###############################################################################
# Database Clone Script
# Clone database dari remote (berdasarkan .env) ke local PostgreSQL
#
# Usage:
#   ./scripts/clone_db.sh [service-name]
#   ./scripts/clone_db.sh              # Akan mencari .env di root atau service pertama
#   ./scripts/clone_db.sh finance      # Clone dari finance/.env
#
###############################################################################

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check required commands
check_requirements() {
    if ! command_exists pg_dump; then
        print_error "pg_dump not found. Please install PostgreSQL client tools."
        exit 1
    fi
    
    if ! command_exists psql; then
        print_error "psql not found. Please install PostgreSQL client tools."
        exit 1
    fi
    
    if ! command_exists createdb; then
        print_error "createdb not found. Please install PostgreSQL client tools."
        exit 1
    fi
}

# Find .env file
find_env_file() {
    local service_name=$1
    local env_file=""
    
    if [ -n "$service_name" ]; then
        # Check in specific service directory
        if [ -f "$service_name/.env" ]; then
            env_file="$service_name/.env"
        else
            print_error ".env file not found in $service_name/"
            exit 1
        fi
    else
        # Check in root first
        if [ -f ".env" ]; then
            env_file=".env"
        else
            # Check in first available service
            for service in cronjob finance inventory master mobile sales system; do
                if [ -f "$service/.env" ]; then
                    env_file="$service/.env"
                    print_info "Using .env from $service/"
                    break
                fi
            done
        fi
        
        if [ -z "$env_file" ]; then
            print_error ".env file not found. Please create .env file or specify service name."
            exit 1
        fi
    fi
    
    echo "$env_file"
}

# Load environment variables from .env file
load_env() {
    local env_file=$1
    
    print_info "Loading environment from $env_file"
    
    # Export variables from .env (skip comments and empty lines)
    set -a
    source "$env_file"
    set +a
    
    # Set defaults if not set
    REMOTE_HOST=${DB_HOST:-localhost}
    REMOTE_PORT=${DB_PORT:-5432}
    REMOTE_USER=${DB_USER:-postgres}
    REMOTE_PASS=${DB_PASS:-}
    REMOTE_DB=${DB_NAME:-scylla_db}
    
    # Local database config (can be overridden by environment)
    LOCAL_HOST=${LOCAL_DB_HOST:-localhost}
    LOCAL_PORT=${LOCAL_DB_PORT:-5432}
    LOCAL_USER=${LOCAL_DB_USER:-postgres}
    LOCAL_DB=${LOCAL_DB_NAME:-scylla_db}
}

# Validate database connection
test_connection() {
    local host=$1
    local port=$2
    local user=$3
    local pass=$4
    local db=$5
    
    print_info "Testing connection to $host:$port/$db..."
    
    if PGPASSWORD="$pass" psql -h "$host" -p "$port" -U "$user" -d "$db" -c "SELECT 1;" >/dev/null 2>&1; then
        print_success "Connection successful"
        return 0
    else
        print_error "Connection failed. Please check your credentials."
        return 1
    fi
}

# Create local database if not exists
create_local_db() {
    print_info "Checking local database $LOCAL_DB..."
    
    if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -lqt | cut -d \| -f 1 | grep -qw "$LOCAL_DB"; then
        print_warning "Database $LOCAL_DB already exists"
        read -p "Do you want to drop and recreate it? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_info "Dropping database $LOCAL_DB..."
            dropdb -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" "$LOCAL_DB" || true
            print_info "Creating database $LOCAL_DB..."
            createdb -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" "$LOCAL_DB"
            print_success "Database $LOCAL_DB created"
        else
            print_info "Using existing database $LOCAL_DB"
        fi
    else
        print_info "Creating database $LOCAL_DB..."
        createdb -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" "$LOCAL_DB"
        print_success "Database $LOCAL_DB created"
    fi
}

# Clone database
clone_database() {
    print_info "Starting database clone..."
    print_info "Source: $REMOTE_HOST:$REMOTE_PORT/$REMOTE_DB"
    print_info "Target: $LOCAL_HOST:$LOCAL_PORT/$LOCAL_DB"
    print_warning "Note: Extension errors (citus, citus_columnar) will be ignored if not installed"
    
    # Ask for backup option
    read -p "Do you want to create a backup file? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).dump"
        print_info "Creating backup file: $BACKUP_FILE"
        
        # Create backup
        PGPASSWORD="$REMOTE_PASS" pg_dump -h "$REMOTE_HOST" -p "$REMOTE_PORT" -U "$REMOTE_USER" -d "$REMOTE_DB" \
            -F c -f "$BACKUP_FILE" 2>&1 | while IFS= read -r line; do
            echo "  $line"
        done
        
        if [ $? -eq 0 ]; then
            print_success "Backup created: $BACKUP_FILE"
            
            # Restore from backup
            print_info "Restoring from backup to local database..."
            pg_restore -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" \
                --clean --if-exists --no-owner --no-acl "$BACKUP_FILE" 2>&1 | \
            while IFS= read -r line; do
                # Filter out citus extension errors but show other errors
                if [[ ! "$line" =~ "extension.*citus" ]] && [[ ! "$line" =~ "extension.*does not exist" ]]; then
                    if [[ "$line" == *"ERROR"* ]]; then
                        print_error "$line"
                    elif [[ "$line" == *"WARNING"* ]]; then
                        print_warning "$line"
                    else
                        echo "  $line"
                    fi
                fi
            done
        else
            print_error "Backup failed"
            exit 1
        fi
    else
        # Direct pipe method
        print_info "Cloning directly (no backup file)..."
        
        PGPASSWORD="$REMOTE_PASS" pg_dump -h "$REMOTE_HOST" -p "$REMOTE_PORT" -U "$REMOTE_USER" -d "$REMOTE_DB" \
            --clean --if-exists --no-owner --no-acl 2>&1 | \
        while IFS= read -r line; do
            # Filter out citus extension errors but show other errors
            if [[ ! "$line" =~ "extension.*citus" ]] && [[ ! "$line" =~ "extension.*does not exist" ]]; then
                if [[ "$line" == *"ERROR"* ]]; then
                    print_error "$line"
                elif [[ "$line" == *"WARNING"* ]]; then
                    print_warning "$line"
                else
                    echo "  $line"
                fi
            fi
        done | psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" 2>&1 | \
        while IFS= read -r line; do
            if [[ ! "$line" =~ "extension.*citus" ]] && [[ ! "$line" =~ "extension.*does not exist" ]]; then
                if [[ "$line" == *"ERROR"* ]]; then
                    print_error "$line"
                elif [[ "$line" == *"WARNING"* ]]; then
                    print_warning "$line"
                else
                    echo "  $line"
                fi
            fi
        done
    fi
    
    # Check if clone was successful by checking table count
    TABLE_COUNT=$(psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | xargs)
    if [ -n "$TABLE_COUNT" ] && [ "$TABLE_COUNT" -gt 0 ]; then
        print_success "Database cloned successfully! ($TABLE_COUNT tables cloned)"
        
        # Show database info
        print_info "Database statistics:"
        psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" -c "
            SELECT 
                schemaname,
                COUNT(*) as table_count
            FROM pg_tables 
            WHERE schemaname = 'public'
            GROUP BY schemaname;
        " 2>/dev/null || true
    else
        print_error "Database clone may have failed - no tables found in database"
        exit 1
    fi
}

# Main function
main() {
    print_info "========================================="
    print_info "  Database Clone Script"
    print_info "========================================="
    echo
    
    # Check requirements
    check_requirements
    
    # Find .env file
    SERVICE_NAME=$1
    ENV_FILE=$(find_env_file "$SERVICE_NAME")
    
    # Load environment
    load_env "$ENV_FILE"
    
    # Display configuration
    echo
    print_info "Configuration:"
    echo "  Remote: $REMOTE_HOST:$REMOTE_PORT/$REMOTE_DB (user: $REMOTE_USER)"
    echo "  Local:  $LOCAL_HOST:$LOCAL_PORT/$LOCAL_DB (user: $LOCAL_USER)"
    echo
    
    # Test remote connection
    if ! test_connection "$REMOTE_HOST" "$REMOTE_PORT" "$REMOTE_USER" "$REMOTE_PASS" "$REMOTE_DB"; then
        exit 1
    fi
    
    # Test local connection
    if ! test_connection "$LOCAL_HOST" "$LOCAL_PORT" "$LOCAL_USER" "" "postgres"; then
        print_error "Cannot connect to local PostgreSQL. Please check your local PostgreSQL is running."
        exit 1
    fi
    
    # Create local database
    create_local_db
    
    # Clone database
    echo
    clone_database
    
    echo
    print_success "Done!"
}

# Run main function
main "$@"

