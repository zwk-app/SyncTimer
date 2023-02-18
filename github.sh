#!/bin/bash
#
# github.sh
#

SCRIPT_FILE="$( basename "${0}" )"
# shellcheck disable=SC2034
SCRIPT_NAME="${SCRIPT_FILE:0:${#SCRIPT_FILE}-3}"
# shellcheck disable=SC2034
SCRIPT_PATH="$( dirname "${0}" )"
SCRIPT_PATH="$( cd "${SCRIPT_PATH}" && pwd)"
if [ -z "${SCRIPT_PATH}" ] ; then
	# For some reason, the path is not accessible to the script (e.g. permissions re-evaled after suid)
	echo "[ERROR] path is not accessible to the script"
	exit 1
fi

APP_NAME=""
APP_VERS=""
BUILD_SUMMARY_MD=BuildSummary.md

function AppName {
  APP_NAME=$(grep "ApplicationName" version.go | awk '{print $4}')
  APP_NAME="${APP_NAME//\"/}"
  APP_NAME="${APP_NAME//\'/}"
  echo -ne "${APP_NAME}"
}

function AppVersion {
  MAJOR=$(grep "MajorVersion" version.go | awk '{print $4}')
  MINOR=$(grep "MinorVersion" version.go | awk '{print $4}')
  BUILD=$(grep "BuildNumber" version.go | awk '{print $4}')
  APP_VERS=$(echo -ne "${MAJOR}.${MINOR}.${BUILD}")
  echo -ne "${APP_VERS}"
}

function SetReleaseEnv {
  APP_NAME="$(AppName)"
  APP_VERS="$(AppVersion)"
  echo "RELEASE_NAME=${APP_NAME} v${APP_VERS}" >> "${GITHUB_ENV}"
  echo "RELEASE_TAG=v${APP_VERS}" >> "${GITHUB_ENV}"
}

function SetNextBuildNumber() {
  CURR_BUILD=$(grep "BuildNumber" version.go | awk '{print $4}')
  NEXT_BUILD=$((${CURR_BUILD}+1))
  sed -i -e "s/const\ BuildNumber\ =\ ${CURR_BUILD}/const\ BuildNumber\ =\ ${NEXT_BUILD}/" version.go
}

function FileCheckSum {
  [ ! -f "${2}" ] && return
  SUM="$(sha256sum "${2}" | awk '{print $1}')"
  echo -ne "|${1}|${2}|${SUM}|\n" >> "${BUILD_SUMMARY_MD}"
}

function BuildSummary {
  BUILD_SUMMARY_MD="$(AppName).md"
  cat << EOF_BUILD_SUMMARY_MD > "${BUILD_SUMMARY_MD}"
## What's Changed
### 🛠 Breaking Changes
* LoremIpsum

### 🎯 Features
* LoremIpsum

### 🩹 Fix:
* LoremIpsum

### 🧹 Other:
* LoremIpsum

## CheckSums
|OS|File|SHA256|
|:---:|:----|:----|
EOF_BUILD_SUMMARY_MD
  FileCheckSum "Windows" "SyncTimer.exe"
  FileCheckSum "Linux" "SyncTimer.tar.xz"
  FileCheckSum "Android" "SyncTimer.apk"
}

function Usage() {
  ERROR_MSG="${1}"
  cat << EOF_USAGE_TXT >&2
Usage : ${SCRIPT_FILE} <action>
Having as <action> one of the following:
    SetReleaseEnv
    SetNextBuildNumber
    BuildSummary
EOF_USAGE_TXT
  [ -n "${ERROR_MSG}" ] && echo -e "\033[0;31m${ERROR_MSG}\033[0m"
  exit 1
}

[ ! ${#} -eq 1 ] && Usage "Missing parameter"
SCRIPT_ACTION="${1}"

case "${SCRIPT_ACTION}" in
	SetReleaseEnv)
		SetReleaseEnv
		;;
	SetNextBuildNumber)
		SetNextBuildNumber
		;;
	BuildSummary)
		BuildSummary
		;;
	*)
		Usage "Unknown action '${SCRIPT_ACTION}'"
		;;
esac

exit 0