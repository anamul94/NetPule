#!/bin/bash

# Unset all snap-related variables
unset SNAP_LIBRARY_PATH GTK_PATH GDK_PIXBUF_MODULEDIR GDK_PIXBUF_MODULE_FILE
unset GTK_EXE_PREFIX GTK_IM_MODULE_FILE GIO_MODULE_DIR GSETTINGS_SCHEMA_DIR
unset XDG_DATA_DIRS LOCPATH

# Clean environment
export PATH="/usr/local/go/bin:/usr/local/bin:/usr/bin:/bin"
export LD_LIBRARY_PATH="/lib/x86_64-linux-gnu:/usr/lib/x86_64-linux-gnu"
export XDG_DATA_DIRS="/usr/local/share:/usr/share"

# Preserve display environment
export DISPLAY="${DISPLAY:-:0}"
export XAUTHORITY="${XAUTHORITY:-$HOME/.Xauthority}"

# Run the application
./netpulse