#! /bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir="${scriptdir%/*}"
tmpdir="${basedir}/tmp"

for n in {1..16}; do
  dd \
    if=/dev/urandom \
    of=${tmpdir}/file$(printf %03d "$n").bin \
    bs=1 count=$((RANDOM + 1024))
done

rm "${tmpdir}/*"
