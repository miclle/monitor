#!/bin/sh
# Generate test coverage statistics for Go packages.
# Usage: sh coverage.sh --xml
#

set -e

workdir=.cover
profile="$workdir/coverage.out"
mode=count

generate_cover_data() {
    rm -rf "$workdir"
    mkdir "$workdir"

    for pkg in "$@"; do
        f="$workdir/$(echo $pkg | tr / -).cover"
        CGO_ENABLED=0 go test -v -covermode="$mode" -coverprofile="$f" "$pkg"
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

generate_xml_report(){
    echo "convert stout to json|convert json to xml"
    gocov convert "$profile"|gocov-xml >coverage.xml
    echo "done"
}

generate_html_report(){
    echo "convert stout to json|convert json to html"
    gocov convert "$profile"|gocov-html >coverage.html
    echo "done"
}



generate_cover_data $(go list ./...)
    case "$1" in
        "")
           ;;
        --html )
          generate_html_report  ;;
        --xml )
          generate_xml_report  ;;
          *)
        echo >&2 "error:invalid option:$1";exit 1;;
    esac