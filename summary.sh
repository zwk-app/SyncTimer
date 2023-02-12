#!/bin/bash

MD_OUT=SyncTimer.summary.md

function MdHeader {
  echo -ne "## Release Notes\n\n" > "${MD_OUT}"
  # shellcheck disable=SC2129
  echo -ne "### 🛠 Breaking Changes\n\n- LoremIpsum\n" >> "${MD_OUT}"
  echo -ne "### 🎯 Features\n\n- LoremIpsum\n" >> "${MD_OUT}"
  echo -ne "### 🩹 Fix:\n\n- LoremIpsum\n" >> "${MD_OUT}"
  echo -ne "### 🧹 Other:\n\n- LoremIpsum\n" >> "${MD_OUT}"
}

function MdSumHead {
  echo -ne "\n## Checksums\n\n" >> "${MD_OUT}"
  echo -ne "|OS|File|SHA256|\n|---|---|---|\n" >> "${MD_OUT}"
}

function MdSumFile {
  [ ! -f "${2}" ] && return
  SUM="$(sha256sum "${2}" | awk '{print $1}')"
  echo -ne "|${1}|${2}|${SUM}|\n" >> "${MD_OUT}"
}

MdHeader
MdSumHead
MdSumFile "Windows" "SyncTimer.exe"
MdSumFile "Linux" "SyncTimer.tar.xz"
MdSumFile "Android" "SyncTimer.apk"

exit 0
