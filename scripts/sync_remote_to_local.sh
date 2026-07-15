#!/usr/bin/env bash

set -Eeuo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[OK]${NC} $1"; }
print_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1" >&2; }

die() { print_error "$1"; exit "${2:-1}"; }

usage() {
  cat <<'EOF'
Usage: ./scripts/sync_remote_to_local.sh [options]

Sync remote PostgreSQL DB to local DB using pg_dump -Fc + pg_restore.

Remote defaults:
  SOURCE_DB_HOST=103.28.219.73
  SOURCE_DB_PORT=25431
  SOURCE_DB_USER=postgres
  SOURCE_DB_NAME=scylla_citus_dev
  SOURCE_DB_PASSWORD=<required>

Local defaults:
  LOCAL_DB_HOST=localhost
  LOCAL_DB_PORT=5432
  LOCAL_DB_USER=postgres
  LOCAL_DB_NAME=ggn_scyllax
  LOCAL_DB_PASSWORD=<required>

Options:
  --install-citus        Install/prepare local citus if source uses citus
  --drop                 Drop and recreate local target DB (requires --yes)
  --yes                  Confirm destructive action for --drop
  --preflight-only       Run checks only (no dump/restore)
  --dry-run              Print planned actions and run preflight checks only
  --keep-dump            Keep dump file after restore
  --dump-file PATH       Custom dump file path (default: /tmp/scylla_citus_dev_YYYYMMDD_HHMMSS.dump)
  --allow-nonlocal-target Allow destructive restore to non-local LOCAL_DB_HOST
  --allow-custom-target  Allow destructive restore to custom LOCAL_DB_NAME besides ggn_scyllax
  --help, -h             Show this help
EOF
}

require_commands() {
  local cmds=(pg_dump pg_restore psql createdb dropdb)
  local missing=0
  for c in "${cmds[@]}"; do
    if ! command -v "$c" >/dev/null 2>&1; then
      print_error "Missing required command: $c"
      missing=1
    fi
  done
  [[ "$missing" -eq 0 ]] || die "Install PostgreSQL client tools first"
}

is_dangerous_db_name() {
  local n="$1"
  [[ -z "$n" || "$n" == "postgres" || "$n" == "template0" || "$n" == "template1" ]]
}

is_local_host() {
  local h="$1"
  [[ "$h" == "localhost" || "$h" == "127.0.0.1" ]]
}

validate_local_target_guard() {
  if [[ "$DROP_DB" == "1" ]]; then
    is_local_host "$LOCAL_DB_HOST" || [[ "$ALLOW_NONLOCAL_TARGET" == "1" ]] || die "Refuse destructive non-local target host '$LOCAL_DB_HOST'. Use --allow-nonlocal-target to override"

    if [[ "$LOCAL_DB_NAME" != "ggn_scyllax" && "$ALLOW_CUSTOM_TARGET" != "1" ]]; then
      die "Refuse destructive custom LOCAL_DB_NAME '$LOCAL_DB_NAME'. Use --allow-custom-target to override"
    fi
  fi
}

db_exists_local() {
  PGPASSWORD="$LOCAL_DB_PASSWORD" psql \
    -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" -d postgres \
    -tAc "SELECT 1 FROM pg_database WHERE datname='${LOCAL_DB_NAME}';" | grep -q '^1$'
}

remote_query() {
  local sql="$1"
  PGPASSWORD="$SOURCE_DB_PASSWORD" psql \
    -h "$SOURCE_DB_HOST" -p "$SOURCE_DB_PORT" -U "$SOURCE_DB_USER" -d "$SOURCE_DB_NAME" \
    -v ON_ERROR_STOP=1 -tAc "$sql"
}

local_query() {
  local db="$1"
  local sql="$2"
  PGPASSWORD="$LOCAL_DB_PASSWORD" psql \
    -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" -d "$db" \
    -v ON_ERROR_STOP=1 -tAc "$sql"
}

check_remote_preflight() {
  print_info "Remote preflight SELECT check"
  remote_query "SELECT current_database(), current_user;" >/dev/null
}

check_local_preflight() {
  print_info "Local preflight connection check"
  local_query "postgres" "SELECT current_database(), current_user;" >/dev/null
}

check_remote_duplicate_primary_keys() {
  print_info "Remote duplicate PK preflight check"

  local m_week_dups
  local m_work_day_dups

  m_week_dups="$(remote_query "SELECT string_agg(format('%s|%s|%s|%s x%s', cust_id, per_year, per_id, week_id, dup_count), E'\n') FROM (SELECT cust_id, per_year, per_id, week_id, COUNT(*) AS dup_count FROM mst.m_week GROUP BY 1,2,3,4 HAVING COUNT(*) > 1 ORDER BY dup_count DESC, cust_id, per_year, per_id, week_id LIMIT 20) d;")"

  m_work_day_dups="$(remote_query "SELECT string_agg(format('%s|%s|%s|%s|%s x%s', cust_id, per_year, per_id, week_id, work_date, dup_count), E'\n') FROM (SELECT cust_id, per_year, per_id, week_id, work_date, COUNT(*) AS dup_count FROM mst.m_work_day GROUP BY 1,2,3,4,5 HAVING COUNT(*) > 1 ORDER BY dup_count DESC, cust_id, per_year, per_id, week_id, work_date LIMIT 20) d;")"

  if [[ -n "${m_week_dups//[[:space:]]/}" || -n "${m_work_day_dups//[[:space:]]/}" ]]; then
    print_error "Remote DB contains duplicate rows that violate declared primary keys. Restore will fail until source data is cleaned."
    if [[ -n "${m_week_dups//[[:space:]]/}" ]]; then
      print_error "Duplicates in mst.m_week (showing up to 20):"
      printf '%b\n' "$m_week_dups" >&2
    fi
    if [[ -n "${m_work_day_dups//[[:space:]]/}" ]]; then
      print_error "Duplicates in mst.m_work_day (showing up to 20):"
      printf '%b\n' "$m_work_day_dups" >&2
    fi
    die "Source DB is not restorable with current PK constraints"
  fi
}

source_uses_citus() {
  local count
  count="$(remote_query "SELECT COUNT(*) FROM pg_extension WHERE extname IN ('citus','citus_columnar');")"
  [[ "${count//[[:space:]]/}" != "0" ]]
}

local_has_citus_available() {
  local_query "postgres" "SELECT 1 FROM pg_available_extensions WHERE name='citus' LIMIT 1;" | grep -q '^1$'
}

local_has_citus_enabled_in_target() {
  local_query "$LOCAL_DB_NAME" "SELECT 1 FROM pg_extension WHERE extname='citus' LIMIT 1;" | grep -q '^1$'
}

install_citus_package() {
  print_info "Installing/preparing local Citus package"
  PGPASSWORD="$LOCAL_DB_PASSWORD" ./scripts/install_citus.sh --non-interactive
}

create_citus_extension_in_target() {
  print_info "Creating Citus extension in target DB"
  PGPASSWORD="$LOCAL_DB_PASSWORD" ./scripts/install_citus.sh --non-interactive --create-extension --db "$LOCAL_DB_NAME" --host "$LOCAL_DB_HOST" --port "$LOCAL_DB_PORT" --user "$LOCAL_DB_USER"
}

prepare_local_db() {
  if is_dangerous_db_name "$LOCAL_DB_NAME"; then
    die "Refuse dangerous LOCAL_DB_NAME: '$LOCAL_DB_NAME'"
  fi

  if db_exists_local; then
    if [[ "$DROP_DB" == "1" && "$ASSUME_YES" == "1" ]]; then
      print_warn "Dropping local DB: $LOCAL_DB_NAME"
      PGPASSWORD="$LOCAL_DB_PASSWORD" psql -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" -d postgres \
        -v ON_ERROR_STOP=1 -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='${LOCAL_DB_NAME}' AND pid <> pg_backend_pid();"
      PGPASSWORD="$LOCAL_DB_PASSWORD" dropdb -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" "$LOCAL_DB_NAME"
      PGPASSWORD="$LOCAL_DB_PASSWORD" createdb -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" "$LOCAL_DB_NAME"
      print_success "Recreated local DB: $LOCAL_DB_NAME"
    else
      die "Local DB exists. Use --drop --yes to recreate: $LOCAL_DB_NAME"
    fi
  else
    print_info "Creating local DB: $LOCAL_DB_NAME"
    PGPASSWORD="$LOCAL_DB_PASSWORD" createdb -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" "$LOCAL_DB_NAME"
  fi
}

dump_remote() {
  print_info "Dumping remote DB to: $DUMP_FILE"
  PGPASSWORD="$SOURCE_DB_PASSWORD" pg_dump \
    -h "$SOURCE_DB_HOST" -p "$SOURCE_DB_PORT" -U "$SOURCE_DB_USER" -d "$SOURCE_DB_NAME" \
    -Fc --no-owner --no-privileges -f "$DUMP_FILE"
}

restore_local() {
  print_info "Restoring dump into local DB: $LOCAL_DB_NAME"
  PGPASSWORD="$LOCAL_DB_PASSWORD" pg_restore \
    -h "$LOCAL_DB_HOST" -p "$LOCAL_DB_PORT" -U "$LOCAL_DB_USER" -d "$LOCAL_DB_NAME" \
    --clean --if-exists --no-owner --no-privileges "$DUMP_FILE"
}

post_validate() {
  print_info "Post-restore validation"
  local_query "$LOCAL_DB_NAME" "SELECT current_database(), current_user;"
  local_query "$LOCAL_DB_NAME" "SELECT COUNT(*) AS user_schema_table_count FROM information_schema.tables WHERE table_schema NOT IN ('pg_catalog','information_schema');"
  local_query "$LOCAL_DB_NAME" "SELECT extname FROM pg_extension ORDER BY extname;"
}

cleanup_dump() {
  if [[ "$KEEP_DUMP" == "1" ]]; then
    print_info "Keeping dump file: $DUMP_FILE"
  else
    rm -f "$DUMP_FILE"
    print_info "Removed dump file: $DUMP_FILE"
  fi
}

# defaults
SOURCE_DB_HOST="${SOURCE_DB_HOST:-103.28.219.73}"
SOURCE_DB_PORT="${SOURCE_DB_PORT:-25431}"
SOURCE_DB_USER="${SOURCE_DB_USER:-postgres}"
SOURCE_DB_NAME="${SOURCE_DB_NAME:-scylla_citus_dev}"
SOURCE_DB_PASSWORD="${SOURCE_DB_PASSWORD:-Ar3m4n1a}"

LOCAL_DB_HOST="${LOCAL_DB_HOST:-localhost}"
LOCAL_DB_PORT="${LOCAL_DB_PORT:-5432}"
LOCAL_DB_USER="${LOCAL_DB_USER:-postgres}"
LOCAL_DB_NAME="${LOCAL_DB_NAME:-ggn_scyllax}"
LOCAL_DB_PASSWORD="${LOCAL_DB_PASSWORD:-postgres}"

INSTALL_CITUS=0
DROP_DB=0
ASSUME_YES=0
PREFLIGHT_ONLY=0
DRY_RUN=0
KEEP_DUMP=0
ALLOW_NONLOCAL_TARGET=0
ALLOW_CUSTOM_TARGET=0
DUMP_FILE="/tmp/scylla_citus_dev_$(date +%Y%m%d_%H%M%S).dump"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --install-citus) INSTALL_CITUS=1 ;;
    --drop) DROP_DB=1 ;;
    --yes) ASSUME_YES=1 ;;
    --preflight-only) PREFLIGHT_ONLY=1 ;;
    --dry-run) DRY_RUN=1 ;;
    --keep-dump) KEEP_DUMP=1 ;;
    --allow-nonlocal-target) ALLOW_NONLOCAL_TARGET=1 ;;
    --allow-custom-target) ALLOW_CUSTOM_TARGET=1 ;;
    --dump-file)
      shift
      [[ $# -gt 0 ]] || die "--dump-file requires path"
      DUMP_FILE="$1"
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *) die "Unknown option: $1" ;;
  esac
  shift
done

require_commands

[[ -n "$SOURCE_DB_PASSWORD" ]] || die "SOURCE_DB_PASSWORD is required"
[[ -n "$LOCAL_DB_PASSWORD" ]] || die "LOCAL_DB_PASSWORD is required"

if [[ "$DROP_DB" == "1" && "$ASSUME_YES" != "1" ]]; then
  die "--drop requires --yes in automation mode"
fi

validate_local_target_guard

print_info "Source: ${SOURCE_DB_USER}@${SOURCE_DB_HOST}:${SOURCE_DB_PORT}/${SOURCE_DB_NAME}"
print_info "Local : ${LOCAL_DB_USER}@${LOCAL_DB_HOST}:${LOCAL_DB_PORT}/${LOCAL_DB_NAME}"

check_remote_preflight
check_local_preflight
check_remote_duplicate_primary_keys

SOURCE_HAS_CITUS=0
if source_uses_citus; then
  SOURCE_HAS_CITUS=1
  print_info "Source uses Citus extension"
  if ! local_has_citus_available; then
    if [[ "$INSTALL_CITUS" == "1" ]]; then
      install_citus_package
      local_has_citus_available || die "Local Citus still not available after install attempt"
    else
      die "Local Citus not available. Re-run with --install-citus"
    fi
  fi
fi

if [[ "$PREFLIGHT_ONLY" == "1" || "$DRY_RUN" == "1" ]]; then
  preflight_modes=()
  [[ "$PREFLIGHT_ONLY" == "1" ]] && preflight_modes+=("--preflight-only")
  [[ "$DRY_RUN" == "1" ]] && preflight_modes+=("--dry-run")
  print_success "Preflight complete (${preflight_modes[*]})"
  exit 0
fi

prepare_local_db

if [[ "$SOURCE_HAS_CITUS" == "1" ]] && ! local_has_citus_enabled_in_target; then
  if [[ "$INSTALL_CITUS" == "1" ]]; then
    create_citus_extension_in_target
    local_has_citus_enabled_in_target || die "Target DB missing citus extension after create attempt"
  else
    die "Target DB missing citus extension. Re-run with --install-citus"
  fi
fi

dump_remote
restore_local
post_validate
cleanup_dump

print_success "Sync complete"
