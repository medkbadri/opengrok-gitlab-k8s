#!/bin/bash

LOCKFILE="$SOURCE_DIR/lock-opengrok-sync"
LOCKINITFILE="$SOURCE_DIR/lock-init-sync"
LOCKINITFILEDONE="$SOURCE_DIR/lock-init-sync-done"

if test "$SYNC_TYPE" = "INITIAL"; then
	if [ -f "$LOCKINITFILE" ]; then
		date +"%F %T Initial synchronization still locked, skipping for now"
		inotifywait -e delete $SOURCE_DIR| while read dir op file
		# Waiting for the first init container to finish synchronization so that all app containers work on valid data
		do
			[[ "${dir}" == "$SOURCE_DIR/" && "${file}" == "lock-init-sync" ]] && date +"%F %T Waiting for first init sync to finish"
		done
		exit
	fi

	if [ -f "$LOCKINITFILEDONE" ]; then
		date +"%F %T Initial synchronization of data shared volume was already performed, skipping for now"
		exit
	fi
        
	touch $LOCKINITFILE
	date +"%F %T Synchronization starting"
	opengrok-gitlab $GROUP_ID $GITLAB_TOKEN $SOURCE_DIR
	date +"%F %T Synchronization finished"
	touch $LOCKINITFILEDONE
	rm -f $LOCKINITFILE
else
	if [ -f "$LOCKINITFILE" ]; then
		date +"%F %T Initial synchronization still locked, skipping for now"
		exit
	fi
		
	date +"%F %T Synchronization starting"
	opengrok-gitlab $GROUP_ID $GITLAB_TOKEN $SOURCE_DIR
	date +"%F %T Synchronization finished"
	rm -f $LOCKFILE
fi