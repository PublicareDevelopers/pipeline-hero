#!/bin/bash

# Get the root directory
root_dir="$1"

# Initialize an error status variable
error_status=0

# Find all .mk files in the root directory
find "$root_dir" -name '*.mk' | while read -r makefile
do
  # Get the directory of the Makefile
  makefile_dir=$(dirname "$makefile")

  # Change to the directory of the Makefile
  cd "$makefile_dir" || exit

  # Check the Makefile for errors
  if ! make -n -f "$makefile" >/dev/null 2>&1; then
    echo "Error in $makefile"
    error_status=1
  else
    # Execute the Makefile with 'build' and 'zip' targets
    make -f "$makefile" build >/dev/null # Redirect stdout to /dev/null
    output=$(make -f "$makefile" zip 2>&1) # Redirect stderr to stdout to capture warnings and errors

    # Check the output for the zip warning
    if echo "$output" | grep -q "zip warning: name not matched"; then
      echo "Warning in $makefile: name not matched"
      error_status=1
    fi
  fi

  # Change back to the root directory
  cd "$root_dir" || exit
done

# Print a newline at the end for better readability
echo

# Exit with the error status
exit $error_status