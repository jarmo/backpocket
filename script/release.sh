#!/usr/bin/env bash

set -e
set -x

if [[ $(git status --porcelain | wc -c) -ne 0 ]]; then
  echo "Cannot release - uncommitted changes found!"
  exit 1
fi

git push

VERSION=v`grep "VERSION =" backpocket.go | awk '{print $4}' | tr -d '"'`
read -rp "Enter changelog to release version $VERSION: " CHANGELOG
RESPONSE=`http -b "https://api.github.com/repos/jarmo/backpocket/releases" Authorization:"Bearer $GITHUB_RELEASE_TOKEN" tag_name="$VERSION" draft:=true name="backpocket $VERSION" body="$CHANGELOG"`

rm -rf dist
mkdir -p dist

for file in `find bin -type f`; do
  DIST_FILE_BASE=`echo $file | awk -F "/" '{name=$3 "-" $2; print name}'`
  DIST_FILE_PATH=dist/$DIST_FILE_BASE-$VERSION.zip
  zip -j $DIST_FILE_PATH $file
  shasum -a 512 $DIST_FILE_PATH > $DIST_FILE_PATH.sha512
done

RELEASE_ID=`echo $RESPONSE | jq -r .id`
for file in `ls -d dist/*`; do
  http -b POST "https://uploads.github.com/repos/jarmo/backpocket/releases/$RELEASE_ID/assets?name=`basename $file`" Authorization:"Bearer $GITHUB_RELEASE_TOKEN" @$file > /dev/null
done

RESPONSE=`http -b PATCH "https://api.github.com/repos/jarmo/backpocket/releases/$RELEASE_ID" Authorization:"Bearer $GITHUB_RELEASE_TOKEN" draft:=false`

echo "Release done:"
echo $RESPONSE | jq -r .html_url
