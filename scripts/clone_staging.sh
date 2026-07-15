#!/bin/bash

###############################################################################
# Clone Staging Database Script
# Clone database dari staging ke local PostgreSQL
# Support multiple database selection
#
# Usage:
#   ./scripts/clone_staging.sh
#
###############################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check required commands
check_requirements() {
    if ! command -v pg_dump >/dev/null 2>&1; then
        print_error "pg_dump not found. Please install PostgreSQL client tools."
        exit 1
    fi
    
    if ! command -v psql >/dev/null 2>&1; then
        print_error "psql not found. Please install PostgreSQL client tools."
        exit 1
    fi
}

# Check and install Citus extension
check_citus_extension() {
    # Set local config defaults if not set
    local local_host=${LOCAL_HOST:-localhost}
    local local_port=${LOCAL_PORT:-5432}
    local local_user=${LOCAL_USER:-postgres}
    
    print_info "Checking Citus extension availability..."
    
    # Check if citus extension is available
    if psql -h "$local_host" -p "$local_port" -U "$local_user" -d "postgres" -c "SELECT * FROM pg_available_extensions WHERE name = 'citus';" 2>/dev/null | grep -q citus; then
        print_success "Citus extension is available"
        return 0
    else
        print_warning "Citus extension is not available in local PostgreSQL"
        echo
        print_info "To install Citus extension, run:"
        echo "  ./scripts/install_citus.sh"
        echo
        print_info "Or install manually:"
        echo "  macOS: brew install citus"
        echo "  Ubuntu/Debian: See https://docs.citusdata.com/en/stable/installation/"
        echo
        read -p "Do you want to continue without Citus extension? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Please install Citus extension first, then run this script again."
            exit 0
        fi
        return 1
    fi
}

# Install Citus extension in database
install_citus_extension() {
    local db_name=$1
    
    print_info "Installing Citus extension in database $db_name..."
    
    # Check if extension already exists
    if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$db_name" -c "SELECT * FROM pg_extension WHERE extname = 'citus';" 2>/dev/null | grep -q citus; then
        print_info "Citus extension already installed"
        return 0
    fi
    
    # Try to create extension
    if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$db_name" -c "CREATE EXTENSION IF NOT EXISTS citus;" 2>/dev/null; then
        print_success "Citus extension installed successfully"
        
        # Try to install citus_columnar if available
        if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$db_name" -c "CREATE EXTENSION IF NOT EXISTS citus_columnar;" 2>/dev/null; then
            print_success "Citus Columnar extension installed successfully"
        else
            print_warning "Citus Columnar extension not available (this is optional)"
        fi
        return 0
    else
        print_warning "Failed to install Citus extension (will continue without it)"
        return 1
    fi
}

# Get staging database credentials
get_staging_credentials() {
    print_info "Please enter staging database credentials:"
    echo
        
    read -p "Staging DB Host [default: localhost]: " STAGING_HOST
    STAGING_HOST=${STAGING_HOST:-localhost}
    
    read -p "Staging DB Port [default: 5432]: " STAGING_PORT
    STAGING_PORT=${STAGING_PORT:-5432}
    
    read -p "Staging DB User [default: postgres]: " STAGING_USER
    STAGING_USER=${STAGING_USER:-postgres}
    
    read -sp "Staging DB Password: " STAGING_PASS
    echo
    
    # List available databases
    print_info "Fetching available databases from staging..."
    DATABASES=$(PGPASSWORD="$STAGING_PASS" psql -h "$STAGING_HOST" -p "$STAGING_PORT" -U "$STAGING_USER" -lqt 2>/dev/null | cut -d \| -f 1 | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' | grep -v '^$' | grep -v template | grep -v postgres)
    
    if [ -z "$DATABASES" ]; then
        print_error "Cannot connect to staging database or no databases found"
        exit 1
    fi
    
    # Show available databases
    echo
    print_info "Available databases:"
    DB_ARRAY=()
    INDEX=1
    while IFS= read -r db; do
        if [ -n "$db" ]; then
            echo "  $INDEX. $db"
            DB_ARRAY+=("$db")
            ((INDEX++))
        fi
    done <<< "$DATABASES"
    
    # Select database
    echo
    read -p "Select database number [1-${#DB_ARRAY[@]}]: " SELECTED_INDEX
    if [ "$SELECTED_INDEX" -ge 1 ] && [ "$SELECTED_INDEX" -le "${#DB_ARRAY[@]}" ]; then
        REMOTE_DB="${DB_ARRAY[$((SELECTED_INDEX-1))]}"
        print_info "Selected database: $REMOTE_DB"
    else
        print_error "Invalid selection"
        exit 1
    fi
    
    # Local database name
    echo
    read -p "Local database name [default: ${REMOTE_DB}_local]: " LOCAL_DB
    LOCAL_DB=${LOCAL_DB:-${REMOTE_DB}_local}
    
    # Local database config
    LOCAL_HOST=${LOCAL_DB_HOST:-localhost}
    LOCAL_PORT=${LOCAL_DB_PORT:-5432}
    LOCAL_USER=${LOCAL_DB_USER:-postgres}
}

# Test connections
test_connections() {
    print_info "Testing staging connection..."
    if PGPASSWORD="$STAGING_PASS" psql -h "$STAGING_HOST" -p "$STAGING_PORT" -U "$STAGING_USER" -d "$REMOTE_DB" -c "SELECT 1;" >/dev/null 2>&1; then
        print_success "Staging connection successful"
    else
        print_error "Cannot connect to staging database"
        exit 1
    fi
    
    print_info "Testing local PostgreSQL..."
    if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "postgres" -c "SELECT 1;" >/dev/null 2>&1; then
        print_success "Local PostgreSQL connection successful"
    else
        print_error "Cannot connect to local PostgreSQL"
        exit 1
    fi
}

# Create local database
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
    print_info "Source: $STAGING_HOST:$STAGING_PORT/$REMOTE_DB"
    print_info "Target: $LOCAL_HOST:$LOCAL_PORT/$LOCAL_DB"
    print_warning "Note: Extension errors (citus, citus_columnar) will be ignored if not installed"
    echo
    
    # Ask for backup option
    read -p "Do you want to create a backup file? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        BACKUP_FILE="backup_${REMOTE_DB}_$(date +%Y%m%d_%H%M%S).dump"
        print_info "Creating backup file: $BACKUP_FILE"
        
        PGPASSWORD="$STAGING_PASS" pg_dump -h "$STAGING_HOST" -p "$STAGING_PORT" -U "$STAGING_USER" -d "$REMOTE_DB" \
            -F c -f "$BACKUP_FILE" 2>&1 | while IFS= read -r line; do
            echo "  $line"
        done
        
        if [ $? -eq 0 ]; then
            print_success "Backup created: $BACKUP_FILE"
            
            print_info "Restoring from backup to local database..."
            
            # Check if Citus is installed, if not filter out citus errors
            if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" -c "SELECT * FROM pg_extension WHERE extname = 'citus';" 2>/dev/null | grep -q citus; then
                # Citus is installed, show all errors
                pg_restore -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" \
                    --clean --if-exists --no-owner --no-acl "$BACKUP_FILE" 2>&1 | \
                while IFS= read -r line; do
                    if [[ "$line" == *"ERROR"* ]]; then
                        print_error "$line"
                    elif [[ "$line" == *"WARNING"* ]]; then
                        print_warning "$line"
                    else
                        echo "  $line"
                    fi
                done
            else
                # Citus not installed, filter out citus errors
                pg_restore -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" \
                    --clean --if-exists --no-owner --no-acl "$BACKUP_FILE" 2>&1 | \
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
        else
            print_error "Backup failed"
            exit 1
        fi
    else
        print_info "Cloning directly (no backup file)..."
        
        PGPASSWORD="$STAGING_PASS" pg_dump -h "$STAGING_HOST" -p "$STAGING_PORT" -U "$STAGING_USER" -d "$REMOTE_DB" \
            --clean --if-exists --no-owner --no-acl 2>&1 | \
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
    
    # Check if clone was successful
    TABLE_COUNT=$(psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$LOCAL_DB" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | xargs)
    if [ -n "$TABLE_COUNT" ] && [ "$TABLE_COUNT" -gt 0 ]; then
        print_success "Database cloned successfully! ($TABLE_COUNT tables cloned)"
        
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
        print_error "Database clone may have failed - no tables found"
        exit 1
    fi
}

# Main function
main() {
    print_info "========================================="
    print_info "  Clone Staging Database Script"
    print_info "========================================="
    echo
    
    check_requirements
    get_staging_credentials
    echo
    test_connections
    echo
    create_local_db
    echo
    
    # Check Citus extension (after local config is set)
    check_citus_extension
    CITUS_AVAILABLE=$?
    echo
    
    # Install Citus extension if available
    if [ $CITUS_AVAILABLE -eq 0 ]; then
        install_citus_extension "$LOCAL_DB"
        echo
    fi
    
    clone_database
    
    echo
    print_success "Done!"
}

main "$@"

