#!/bin/sh

CRON_TYPE=$1

if [ -z "$CRON_TYPE" ]; then
    echo "Usage: $0 <cron_type>"
    exit 1
fi

# Use env vars passed to container for auth
# silent fail (-f), show error (-S), max time 10s
curl -f -s -S -u "${CRON_USERNAME}:${CRON_PASSWORD}" \
     "http://core-api:8080/v1/cron/${CRON_TYPE}"

echo " [$(date)] Triggered cron: ${CRON_TYPE} - Exit Code: $?"
