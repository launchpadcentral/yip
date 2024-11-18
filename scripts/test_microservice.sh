#!/usr/bin/env bash
set -e -o pipefail
truncate -s 0 coverage.txt

trap "go-junit-report <${TEST_RESULTS}/go-test.out > ${TEST_RESULTS}/go-test-report.xml" EXIT

for d in $(find . -maxdepth 10 -type d -not -path "*vendor*" -not -path "*extras*"); do
  if ls $d/*.go &> /dev/null; then
    go test -v -coverprofile=profile.out -covermode=atomic $d | tee -a ${TEST_RESULTS}/go-test.out

    if [ -f profile.out ]; then
      cat profile.out >> coverage.txt
      rm profile.out
    fi
  fi

done

