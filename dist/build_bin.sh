#!/bin/bash
#
# Build binaries for all target architrectures specified in the HTTPR_ARCH_LIST variable
#
HTTPR_ARCH_LIST="darwin linux windows"

for HTTPR_ARCH in $HTTPR_ARCH_LIST
do
  HTTPR_BIN_PATH="dist/bin/$HTTPR_ARCH"
  if [[ ! -d $HTTPR_BIN_PATH ]];
  then
    mkdir -p $HTTPR_BIN_PATH
  fi
  HTTPR_BIN_DIST=$HTTPR_BIN_PATH/httpr
  echo -n "Building for $HTTPR_ARCH ... "
  GOOS=$HTTPR_ARCH go build -o $HTTPR_BIN_DIST
  echo "Done: $(file $HTTPR_BIN_DIST)"
done;

