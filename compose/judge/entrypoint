#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

QUIET=true
ISOLATE_CHECK_EXECUTE=false
STRICT=false

ARGS=$(getopt -n "$0" -o vesh: -l verbose,execute-patches,strict,help -- "$@")
eval set -- "$ARGS"

while true; do
    case "$1" in
        -v|--verbose)
            QUIET=false
            shift ;;
        -e|--execute-patches)
            ISOLATE_CHECK_EXECUTE=true
            shift ;;
        -s|--strict)
            STRICT=true
            shift ;;
        --)
            shift
            break ;;
        -h|--help)
            echo "$(basename "$0")"
            echo "Usage: [-v|--verbose] [-e|--execute-patches] [-s|--strict] [-h|--help] [--] <command>"
            echo "  --verbose           Print every thing"
            echo "  --execute-patches   Run isolate-check-environment --execute --quiet which increases reproducibility"
            echo "  --strict            Fail if isolate-check-environment fails"
            echo "  --help              Show this help message and exit"
            echo "  -- <command>        Optional command to be executed after isolate-cg-keeper is started"
            exit 0 ;;
    esac
done

print() {
    if [ $QUIET = false ]; then
        echo "$@"
    fi
}

if ! mount -t cgroup2 | grep -E "\(rw\)|\(rw,|,rw\)|,rw,"; then
    print "/sys/fs/cgroup read-only. Remounting as read-write."
    mount -o remount,rw /sys/fs/cgroup/
fi

# Run isolate daemon
print "Running isolate daemon. This will move all processes to the /daemon control group."
isolate-cg-keeper --move-cg-neighbors & DAEMON_PID=$!

if [ $ISOLATE_CHECK_EXECUTE = true ]; then
    print "Running isolate-check-environment --execute --quiet"
    isolate-check-environment --execute --quiet > /dev/null 2> /dev/null || true
fi

if [ $STRICT = true ]; then
    print "Running isolate-check-environment"
    if [ $QUIET = true ]; then
        isolate-check-environment --quiet > /dev/null 2> /dev/null
    else
        isolate-check-environment
    fi
else
    print "Skipping isolate-check-environment"
fi

if [ $# -eq 0 ]; then
    print "No command to execute. Waiting for isolate-cg-keeper to finish."
    wait $DAEMON_PID
    exit 0
fi

print "Executing $@"
exec "$@"
