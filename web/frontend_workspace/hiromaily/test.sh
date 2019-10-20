#!/bin/sh

WATCH="watchify -t babelify ./src/app.js -o 'exorcist ./build.js.map > ./dist/build.js' -d"
LINT="eslint js/**/*.jsx"

cat package.json.bk |
  jq --arg WATCH "$WATCH" --arg LINT "$LINT" 'to_entries |
       map(if .key == "scripts"
          then . + {"value":
                      {
                  "watch": $WATCH,
                  "lint": $LINT
                      }
                    }
          else .
          end
      ) | from_entries' |

  jq 'del(.main) | del(.keywords) | del(.author) | del(.license)'
