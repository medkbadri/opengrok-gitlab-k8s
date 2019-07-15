#!/bin/bash

LOCKFILE=/var/run/opengrok-indexer
LOCKSYNCFILE="$SOURCE_DIR/lock-opengrok-sync"
URI="http://localhost:8080"

if [ -f "$LOCKFILE" ]; then
	date +"%F %T Indexer still locked, skipping indexing"
	exit 1
fi

# If a new pod is created, init container will skip Synchronization and initial indexing could happen in parallel with an ongoing Synchronization.
# That's why we should prevent this use case.
if [ -f "$LOCKSYNCFILE" ]; then
	date +"%F %T Synchronization still ongoing, skipping for now"
	exit 1
fi

touch $LOCKFILE

date +"%F %T Indexing starting"
opengrok-indexer \
    -a /opengrok/lib/opengrok.jar -- \
    -s $SRC_ROOT\
    -d /opengrok/data \
    -H -P -S -G \
    --leadingWildCards on \
    -W /var/opengrok/etc/configuration.xml \
    -U "$URI" \
    $INDEXER_OPT "$@"
date +"%F %T Indexing finished"
echo "Pausing until new sync finishes"

rm -f $LOCKFILE