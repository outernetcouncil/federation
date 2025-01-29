#set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
path_before_bazel_out="../"

process_virtual_imports() {
  result=""
  local directory="$1"
  local absolute_directory="$(readlink -f "$directory")"

  local seen=()

  while IFS= read -r -d '' file; do
    path="$file"
    if [[ "$path" == *"_virtual_imports/"* ]]; then
      virtual_imports_index=${path##*_virtual_imports/}
      first_subdir=${virtual_imports_index%%/*}

	  intermediate_result="${path%/_virtual_imports/*}/_virtual_imports/$first_subdir/"
	  if [[ ! " ${seen[@]} " =~ " $intermediate_result " ]]; then
        seen+=("$intermediate_result")
		result="$result -I $intermediate_result"
      fi
    fi
  done < <(find "$absolute_directory" -type l -print0)
  echo "$result"
}

all_paths=$(process_virtual_imports "../org_outernetcouncil_nmts+")

../gazelle++go_deps+com_github_googleapis_api_linter/cmd/api-linter/api-linter_/api-linter outernet/federation/interconnect/v1alpha/interconnect.proto \
    -I "$path_before_bazel_out" \
	-I "../googleapis+" \
	--descriptor-set-in "$path_before_bazel_out"protobuf+/src/google/protobuf/duration_proto-descriptor-set.proto.bin \
	--descriptor-set-in "$path_before_bazel_out"protobuf+/src/google/protobuf/empty_proto-descriptor-set.proto.bin \
	$all_paths \
	--set-exit-status
