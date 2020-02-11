cat ~/Downloads/ril_export.html | pup 'ul:first-of-type a json{}' | jq -r '.[] | "\(.href) \(.time_added)"' | xargs -P4 -L1 go run .
find articles -type f -name "*.html" -exec ls -1t "{}" +;
