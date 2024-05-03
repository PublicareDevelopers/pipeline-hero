#!/bin/bash

# Root directory
root_dir="$1"

# Directories to check
directories=("pipelines/teststage" "pipelines/productivestage")

# Initialize a variable to keep track of whether any files were not found
file_not_found=0

# Function to trim leading and trailing spaces
trim() {
  local var="$*"
  # remove leading whitespace characters
  var="${var#"${var%%[![:space:]]*}"}"
  # remove trailing whitespace characters
  var="${var%"${var##*[![:space:]]}"}"
  echo -n "$var"
}

# Loop over directories
for dir in "${directories[@]}"; do
  # Find all .yml files in the directory
  find "$dir" -name '*.yml' | while read -r yml_file; do
    # Extract all lines that contain 'artifact: bin/' and include line numbers
    awk '/artifact: bin\// {print NR, $0}' "$yml_file" | while read -r line_number line; do
      # Extract the path
      bin_path=$(echo "$line" | awk -F'artifact: ' '{print $2}' | tr -d '[:space:]')

      # Add the root directory to the path
      full_path="$root_dir/$bin_path"

      # Check if the file exists
      if test -f "$full_path"; then
        # Do nothing
        :
      else
        echo "File does not exist: $full_path at line $line_number in $yml_file"
        echo "goland://open?file=$root_dir/$yml_file&line=$line_number" # Output the file and line number as a link
        file_not_found=1
      fi
    done
  done
done

# Exit with the value of file_not_found as the exit status
exit $file_not_found