#!/bin/bash

set -e

CLIENT="emporous"
COLLECTION_TAG="latest"
CLIENT_ARGS=""
PUSH="false"

for i in "$@"
do
  case $i in
    -d=* | --directory=* )
      COLLECTIONS_DIRECTORY="${i#*=}"
      shift
      ;;
    -c=* | --client=* )
      CLIENT="${i#*=}"
      shift
      ;;
    -g=* | --gitops-namespace=* )
      GITOPS_NAMESPACE="${i#*=}"
      shift
      ;;
    -i | --insecure )
      CLIENT_ARGS+=" --insecure"
      shift
      ;;
    -p | --plain-http )
      CLIENT_ARGS+=" --plain-http"
      shift
      ;;
    --push )
      PUSH="true"
      shift
      ;;
    -r=* | --repository=* )
      REPOSITORY="${i#*=}"
      shift
      ;;
    -t=* | --tag=* )
      COLLECTION_TAG="${i#*=}"
      shift
      ;;
  esac
done

if [ -z "${COLLECTIONS_DIRECTORY}" ] || [ ! -d ${COLLECTIONS_DIRECTORY} ]; then
  echo "ERROR: Directory '${COLLECTIONS_DIRECTORY}' Not Found!"
  exit 1
fi

if [ -z "${REPOSITORY}" ]; then
  echo "ERROR: Repository Value is required!"
  exit 1
fi

if ! command -v ${CLIENT} &> /dev/null
then
    echo "ERROR: ${CLIENT} could not be found"
    exit
fi

function build_destination_artifact() {
    repository=$1
    collection=$2
    tag=$3

    echo -n ${repository}/${collection}:${tag}
}

for directory in `ls -d -- ${COLLECTIONS_DIRECTORY}/*`; do

  client_build_cmd="${CLIENT} build collection${CLIENT_ARGS}"

  dirname="$(basename "${directory}")"

  destination=$(build_destination_artifact ${REPOSITORY} ${dirname} ${COLLECTION_TAG})
  
  client_build_cmd+=" . ${destination}"

  # Find dsconfig file
  dsconfig_file=$(find ${directory} -name "*.yaml" | head -n 1)

  if [ "${dsconfig_file}" != "" ]; then
    client_build_cmd+=" --dsconfig=$(basename ${dsconfig_file})"
  fi

  pushd "${directory}" >/dev/null 2>&1
  
  echo
  echo "== Building Collection '${dirname}' =="
  echo
  eval "${client_build_cmd}"
  
  popd >/dev/null 2>&1

done

if [ "${PUSH}" == "true" ]; then

    for directory in `ls -d -- ${COLLECTIONS_DIRECTORY}/*`; do

        client_push_cmd="${CLIENT} push${CLIENT_ARGS}"

        dirname="$(basename "${directory}")"

        destination=$(build_destination_artifact ${REPOSITORY} ${dirname} ${COLLECTION_TAG})
        
        client_push_cmd+=" ${destination}"

        pushd "${directory}" >/dev/null 2>&1
        
        echo
        echo "== Pushing Collection '${dirname}' =="
        echo
        eval "${client_push_cmd}"
        
        popd >/dev/null 2>&1

    done

fi
