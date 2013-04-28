#!/bin/bash

root=$(dirname $0)/..
function norm {
  gunzip -c $1 | xsltproc --nonet --novalid $root/xslt/identity.xslt -
}
diff -C 3 <(norm $1) <(norm $2)
