#!/usr/bin/env bash

###############################################################################
# Install Citus Extension Script
# Install Citus extension untuk PostgreSQL local
#
# Usage:
#   ./scripts/install_citus.sh
#   ./scripts/install_citus.sh --non-interactive --create-extension --db ggn_scyllax
#
###############################################################################

set -Eeuo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

die() { print_error "$1"; exit "${2:-1}"; }

NON_INTERACTIVE=0
CREATE_EXTENSION_FLAG=0
DB_NAME="postgres"
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="postgres"

usage() {
  cat <<'EOF'
Usage: ./scripts/install_citus.sh [options]

Options:
  --non-interactive      No prompt. Fail if manual choice needed.
  --create-extension     Create extension after install/check.
  --db NAME              Target DB for extension create (default: postgres)
  --host HOST            DB host (default: localhost)
  --port PORT            DB port (default: 5432)
  --user USER            DB user (default: postgres)
  --help, -h             Show help
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --non-interactive) NON_INTERACTIVE=1 ;;
    --create-extension) CREATE_EXTENSION_FLAG=1 ;;
    --db)
      shift; [[ $# -gt 0 ]] || die "--db needs value"; DB_NAME="$1"
      ;;
    --host)
      shift; [[ $# -gt 0 ]] || die "--host needs value"; DB_HOST="$1"
      ;;
    --port)
      shift; [[ $# -gt 0 ]] || die "--port needs value"; DB_PORT="$1"
      ;;
    --user)
      shift; [[ $# -gt 0 ]] || die "--user needs value"; DB_USER="$1"
      ;;
    --help|-h)
      usage; exit 0
      ;;
    *) die "Unknown option: $1" ;;
  esac
  shift
done

detect_os() {
  if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "macos"
  elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if [[ -f /etc/debian_version ]]; then
      echo "debian"
    elif [[ -f /etc/redhat-release ]]; then
      echo "rhel"
    else
      echo "linux"
    fi
  else
    echo "unknown"
  fi
}

get_pg_major() {
  psql --version | sed -E 's/.* ([0-9]+)(\.[0-9]+)?.*/\1/'
}

install_citus_macos() {
  print_info "Installing Citus on macOS using Homebrew..."
  command -v brew >/dev/null 2>&1 || die "Homebrew not found"
  brew install citus
}

install_citus_debian() {
  print_info "Installing Citus on Debian/Ubuntu..."
  command -v sudo >/dev/null 2>&1 || die "sudo not found"
  curl -fsSL https://install.citusdata.com/community/deb.sh | sudo bash
  local pg_major
  pg_major="$(get_pg_major)"
  sudo apt-get update
  sudo apt-get install -y "postgresql-${pg_major}-citus"
}

install_citus_rhel() {
  print_info "Installing Citus on RHEL/CentOS..."
  command -v sudo >/dev/null 2>&1 || die "sudo not found"
  curl -fsSL https://install.citusdata.com/community/rpm.sh | sudo bash
  local pg_major
  pg_major="$(get_pg_major)"
  sudo yum install -y "citus_${pg_major}"
}

create_extension() {
  print_info "Creating Citus extension in database: $DB_NAME"
  if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -c "CREATE EXTENSION IF NOT EXISTS citus;"; then
    print_success "Citus extension created in database $DB_NAME"
  else
    die "Failed to create Citus extension"
  fi

  if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -c "CREATE EXTENSION IF NOT EXISTS citus_columnar;"; then
    print_success "Citus Columnar extension created in database $DB_NAME"
  else
    print_warning "Citus Columnar extension not available (optional)"
  fi
}

main() {
  print_info "========================================="
  print_info "  Install Citus Extension Script"
  print_info "========================================="

  local os
  os="$(detect_os)"
  print_info "Detected OS: $os"

  if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d postgres -tAc "SELECT 1 FROM pg_available_extensions WHERE name='citus' LIMIT 1;" | grep -q '^1$'; then
    print_success "Citus extension already available"
    if [[ "$CREATE_EXTENSION_FLAG" == "1" ]]; then
      create_extension
    elif [[ "$NON_INTERACTIVE" == "0" ]]; then
      read -r -p "Do you want to create extension in database '$DB_NAME'? (y/N): " reply
      if [[ "$reply" =~ ^[Yy]$ ]]; then
        create_extension
      fi
    fi
    exit 0
  fi

  case "$os" in
    macos) install_citus_macos ;;
    debian) install_citus_debian ;;
    rhel) install_citus_rhel ;;
    *) die "Unsupported OS: $os. Install manually from https://github.com/citusdata/citus" ;;
  esac

  print_success "Citus installation completed"

  if [[ "$CREATE_EXTENSION_FLAG" == "1" ]]; then
    create_extension
  elif [[ "$NON_INTERACTIVE" == "0" ]]; then
    read -r -p "Do you want to create extension in database '$DB_NAME'? (y/N): " reply
    if [[ "$reply" =~ ^[Yy]$ ]]; then
      create_extension
    fi
  fi

  print_success "Done"
}

main "$@"
