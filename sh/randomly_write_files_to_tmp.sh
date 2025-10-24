#! /bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir="${scriptdir%/*}"
tmpdir="${basedir}/tmp"

max="${1}"

[[ -z "${max}" ]] && max=16

for n in $(seq 1 ${max}); do
  dd \
    if=/dev/urandom \
    of=${tmpdir}/file$(printf %03d "$n").bin \
    bs=1 count=$((RANDOM + 1024))
done

rm "${tmpdir}/*"
