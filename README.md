cat ril_export.html | pup 'ul:first-of-type a attr{href}' | tail -r | xargs -L1 go run .
find articles -type f -name "*.html" -exec ls -1t "{}" +;
