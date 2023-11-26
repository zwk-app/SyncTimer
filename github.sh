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

APP_ID=""
APP_NAME=""
APP_VERS=""
APP_BUILD=""
APP_VERSION=""
MOD_NAME=""
BUILD_SUMMARY_MD=BuildSummary.md

NDK_VERSION="r25c"
NDK_URL="https://dl.google.com/android/repository/android-ndk-${NDK_VERSION}-linux.zip"

function Usage() {
  ERROR_MSG="${1}"
  cat << EOF_USAGE_TXT >&2
Usage : ${SCRIPT_FILE} <action> [param]
Having as <action> one of the following:
    github-environment
    build-increment
    build-requirements
    package linux|android|windows
    release linux|android|windows
    build-summary
EOF_USAGE_TXT
  [ -n "${ERROR_MSG}" ] && echo -e "\033[0;31m${ERROR_MSG}\033[0m"
  exit 1
}

function AppId() {
  APP_ID="$(grep "ApplicationId" version.go | awk '{split($0,a,"="); print a[2]}' | xargs)"
  APP_ID="${APP_ID//\"/}"
  APP_ID="${APP_ID//\'/}"
  echo -ne "${APP_ID}"
}

function AppName() {
  APP_NAME="$(grep "ApplicationName" version.go | awk '{split($0,a,"="); print a[2]}' | xargs)"
  APP_NAME="${APP_NAME//\"/}"
  APP_NAME="${APP_NAME//\'/}"
  echo -ne "${APP_NAME}"
}

function AppVers() {
  MAJOR="$(grep "MajorVersion" version.go | awk '{split($0,a,"="); print a[2]}' | xargs)"
  MINOR="$(grep "MinorVersion" version.go | awk '{split($0,a,"="); print a[2]}' | xargs)"
  APP_VERS=$(echo -ne "${MAJOR}.${MINOR}")
  echo -ne "${APP_VERS}"
}

function AppBuild() {
  APP_BUILD="$(grep "BuildNumber" version.go | awk '{split($0,a,"="); print a[2]}' | xargs)"
  echo -ne "${APP_BUILD}"
}

function AppVersion() {
  APP_VERSION="$(AppVers).$(AppBuild)"
  echo -ne "${APP_VERSION}"
}

function ModName() {
  MOD_NAME=$(grep -m 1 "module" go.mod | awk '{print $2}')
  echo -ne "${MOD_NAME}"
}

function GithubEnvironment() {
  APP_NAME="$(AppName)"
  APP_VERSION="$(AppVersion)"
  echo "RELEASE_NAME=${APP_NAME} v${APP_VERSION}" >> "${GITHUB_ENV}"
  echo "RELEASE_TAG=v${APP_VERSION}" >> "${GITHUB_ENV}"
}

function BuildIncrement() {
  CURR_BUILD=$(grep "BuildNumber" version.go | awk '{print $4}')
  # shellcheck disable=SC2004
  NEXT_BUILD=$((${CURR_BUILD}+1))
  sed -i -e "s/const\ BuildNumber\ =\ ${CURR_BUILD}/const\ BuildNumber\ =\ ${NEXT_BUILD}/" version.go
}

function AndroidNDK() {
  echo -e "ANDROID_NDK_HOME=${ANDROID_NDK_HOME}"
  export ANDROID_NDK_HOME="${SCRIPT_PATH}/android-ndk-${NDK_VERSION}/"
  echo -e "ANDROID_NDK_HOME=${ANDROID_NDK_HOME}"
  if [ ! -d "${ANDROID_NDK_HOME}" ]; then
    echo -e "Downloading Android NDK version ${NDK_VERSION}"
    NDK_ZIP="android.ndk.zip"
    curl -s "${NDK_URL}" -o "${NDK_ZIP}"
    unzip "${NDK_ZIP}" > /dev/null 2>&1
    rm "${NDK_ZIP}" > /dev/null 2>&1
  fi
  ls -lha "${ANDROID_NDK_HOME}"
}

function BuildRequirements() {
  # Update & Upgrade
  sudo apt-get -y update && apt-get -y upgrade
  # Go & Co
  sudo apt-get -y install golang gcc libgl1-mesa-dev xorg-dev
  sudo apt-get -y install gcc-mingw-w64-x86-64
  # Fyne
  go get fyne.io/fyne/v2
  go get github.com/fyne-io/fyne-cross@latest
  go install fyne.io/fyne/v2/cmd/fyne@latest
  # TextToSpeech
  sudo apt-get -y install libasound2-dev
  go get github.com/hajimehoshi/go-mp3
  go get github.com/hajimehoshi/oto/v2
  # Tidy
  go mod tidy
}

function FyneExec() {
  COMMAND="${1}"
  TARGET="${2}"
  [[ ! ${COMMAND} =~ package|release ]] && Usage "Wrong command '${COMMAND}'"
  [[ ! ${TARGET} =~ linux|android|windows ]] && Usage "Wrong target '${TARGET}'"
  unset CC CGO_ENABLED CGO_LDFLAGS GOOS
  unset
  case "${TARGET}" in
  	windows)
  	  export CC=/usr/bin/x86_64-w64-mingw32-gcc-win32
  	  export CGO_LDFLAGS="-static-libgcc -static"
  		;;
  esac
  export CGO_ENABLED=1
  export GOOS="${TARGET}"
  fyne "${COMMAND}" -os "${TARGET}" --name "$(AppName)" --appID "$(AppId)" --appVersion "$(AppVers)" --appBuild "$(AppBuild)" -icon "./res/icon.png"
}

function FileCheckSum() {
  [ ! -f "${2}" ] && return
  SUM="$(sha256sum "${2}" | awk '{print $1}')"
  echo -ne "|${1}|${2}|${SUM}|\n" >> "${BUILD_SUMMARY_MD}"
}

function BuildSummary() {
  BUILD_PATH="${1}"
  BUILD_SUMMARY_MD="${BUILD_PATH}$(ModName).md"
  cat << EOF_BUILD_SUMMARY_MD > "${BUILD_SUMMARY_MD}"
# $(AppName) v$(AppVersion)
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
  FileCheckSum "Windows" "${BUILD_PATH}SyncTimer.exe"
  FileCheckSum "Linux" "${BUILD_PATH}SyncTimer.tar.xz"
  FileCheckSum "Android" "${BUILD_PATH}SyncTimer.apk"
}

[ ${#} -lt 1 ] && Usage "Missing parameter"
SCRIPT_ACTION="${1}"
SCRIPT_SUBACTION=""
[ ${#} -gt 1 ] && SCRIPT_SUBACTION="${2}"
case "${SCRIPT_ACTION}" in
	github-environment)
	  [ ${#} -gt 1 ] && Usage "Too many parameters for ${SCRIPT_ACTION}"
		GithubEnvironment
		;;
	build-increment)
	  [ ${#} -gt 1 ] && Usage "Too many parameters for ${SCRIPT_ACTION}"
		BuildIncrement
		;;
	build-requirements)
	  [ ${#} -gt 1 ] && Usage "Too many parameters for ${SCRIPT_ACTION}"
		BuildRequirements
		;;
	package|release)
	  [ ${#} -gt 2 ] && Usage "Too many parameters for ${SCRIPT_ACTION}"
	  [ ${#} -lt 2 ] && Usage "Missing parameters for ${SCRIPT_ACTION}"
	  FyneExec "${SCRIPT_ACTION}" "${SCRIPT_SUBACTION}"
		;;
	build-summary)
	  [ ${#} -gt 2 ] && Usage "Too many parameters for ${SCRIPT_ACTION}"
	  [ ${#} -lt 2 ] && Usage "Missing parameters for ${SCRIPT_ACTION}"
		BuildSummary "${SCRIPT_SUBACTION}"
		;;
  test)
    AndroidNDK
    ;;
	*)
		Usage "Unknown action '${SCRIPT_ACTION}'"
		;;
esac

exit 0