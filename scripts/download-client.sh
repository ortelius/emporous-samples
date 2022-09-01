#!/bin/bash

CLIENT_NAME=uor-client-go
GITHUB_ORGANIZATION=uor-framework
GITHUB_REPOSITORY=uor-client-go
PLATFORM=$(uname | tr '[:upper:]' '[:lower:]')
ARCH="amd64"
DESTINATION=""

if ! command -v jq &> /dev/null
then
    echo "ERROR: jq could not be found"
    exit
fi

for i in "$@"
do
  case $i in
    -a=* | --architecture=* )
      ARCH="${i#*=}"
      shift
      ;;
    -d=* | --destination=* )
      DESTINATION="${i#*=}"
      shift
      ;;
    -o=* | --github-organization=* )
      GITHUB_ORGANIZATION="${i#*=}"
      shift
      ;;
    -r=* | --github-repository=* )
      GITHUB_REPOSITORY="${i#*=}"
      shift
      ;;
    -g=* | --gitops-namespace=* )
      GITOPS_NAMESPACE="${i#*=}"
      shift
      ;;
    -p=* | --platform=* )
      PLATFORM="${i#*=}"
      shift
      ;;

  esac
done

#TODO Add failure logic
RELEASE_JSON=$(curl -s https://api.github.com/repos/${GITHUB_ORGANIZATION}/${GITHUB_REPOSITORY}/releases/latest)

if [ $? -ne 0 ]; then
    echo "Error: Failed to retrieve release information"
    exit 1
fi

DOWNLOAD_URL=$(echo $RELEASE_JSON | jq -r '.assets[] | select(.name | endswith("linux-amd64")) | {browser_download_url} | .browser_download_url')

if [[ "$DOWNLOAD_URL" == "null" ]] || [[ -z "$DOWNLOAD_URL" ]]; then
  echo "Error: Unable to obtain release URL"
  exit 1
fi

if [ "$DESTINATION" == "" ]; then
  DESTINATION=$(echo $(pwd)/"${DOWNLOAD_URL##*/}")
fi

curl -sL -o "${DESTINATION}" "$DOWNLOAD_URL"
chmod +x "${DESTINATION}"

echo "File Downloaded to ${DESTINATION}"
