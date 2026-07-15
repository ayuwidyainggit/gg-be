#!/bin/bash

###############################################################################
# Database Restore Script
# Restore database dari backup file (.dump atau .sql) ke local PostgreSQL
#
# Usage:
#   ./scripts/restore_db.sh [backup_file] [database_name]
#   ./scripts/restore_db.sh backup.dump scylla_db
#   ./scripts/restore_db.sh backup.sql scylla_db
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

# Check required commands
check_requirements() {
    if ! command -v pg_restore >/dev/null 2>&1 && ! command -v psql >/dev/null 2>&1; then
        print_error "PostgreSQL client tools not found. Please install PostgreSQL."
        exit 1
    fi
}

# Main function
main() {
    print_info "========================================="
    print_info "  Database Restore Script"
    print_info "========================================="
    echo
    
    check_requirements
    
    # Get backup file
    if [ -z "$1" ]; then
        # Try to find backup file in current directory
        if [ -f "backup.dump" ]; then
            BACKUP_FILE="backup.dump"
        elif [ -f "backup.sql" ]; then
            BACKUP_FILE="backup.sql"
        else
            print_error "Backup file not specified and no default backup file found."
            echo "Usage: $0 [backup_file] [database_name]"
            exit 1
        fi
    else
        BACKUP_FILE=$1
    fi
    
    # Check if backup file exists
    if [ ! -f "$BACKUP_FILE" ]; then
        print_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    
    # Get database name
    DB_NAME=${2:-scylla_db}
    
    # Local database config
    LOCAL_HOST=${LOCAL_DB_HOST:-localhost}
    LOCAL_PORT=${LOCAL_DB_PORT:-5432}
    LOCAL_USER=${LOCAL_DB_USER:-postgres}
    
    print_info "Backup file: $BACKUP_FILE"
    print_info "Target database: $DB_NAME"
    print_info "Host: $LOCAL_HOST:$LOCAL_PORT"
    echo
    
    # Check if database exists
    if psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
        print_warning "Database $DB_NAME already exists"
        read -p "Do you want to drop and recreate it? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_info "Dropping database $DB_NAME..."
            dropdb -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" "$DB_NAME" || true
            print_info "Creating database $DB_NAME..."
            createdb -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" "$DB_NAME"
        fi
    else
        print_info "Creating database $DB_NAME..."
        createdb -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" "$DB_NAME"
    fi
    
    # Restore based on file extension
    if [[ "$BACKUP_FILE" == *.dump ]] || [[ "$BACKUP_FILE" == *.backup ]]; then
        print_info "Restoring from custom format dump..."
        print_warning "Note: Extension errors (citus, citus_columnar) will be ignored if not installed"
        
        # Run restore and filter out citus extension errors
        pg_restore -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$DB_NAME" \
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
        
        # Check if restore was successful by checking table count
        TABLE_COUNT=$(psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | xargs)
        if [ -n "$TABLE_COUNT" ] && [ "$TABLE_COUNT" -gt 0 ]; then
            print_success "Restore completed successfully ($TABLE_COUNT tables restored)"
        else
            print_error "Restore may have failed - no tables found in database"
            exit 1
        fi
    elif [[ "$BACKUP_FILE" == *.sql ]] || [[ "$BACKUP_FILE" == *.sql.gz ]]; then
        print_info "Restoring from SQL dump..."
        if [[ "$BACKUP_FILE" == *.gz ]]; then
            gunzip -c "$BACKUP_FILE" | psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$DB_NAME"
        else
            psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$DB_NAME" < "$BACKUP_FILE"
        fi
    else
        print_error "Unknown backup file format. Supported: .dump, .backup, .sql, .sql.gz"
        exit 1
    fi
    
    if [ $? -eq 0 ]; then
        print_success "Database restored successfully!"
        
        # Show database info
        print_info "Database statistics:"
        psql -h "$LOCAL_HOST" -p "$LOCAL_PORT" -U "$LOCAL_USER" -d "$DB_NAME" -c "
            SELECT 
                schemaname,
                COUNT(*) as table_count
            FROM pg_tables 
            WHERE schemaname = 'public'
            GROUP BY schemaname;
        " 2>/dev/null || true
    else
        print_error "Database restore failed"
        exit 1
    fi
    
    echo
    print_success "Done!"
}

# Run main function
main "$@"

