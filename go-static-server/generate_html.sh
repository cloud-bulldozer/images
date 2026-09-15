#!/usr/bin/env bash
# generate_html.sh — Generate static HTML files of specific byte sizes
# for benchmarking. Each file is padded with 'A' characters inside a
# minimal HTML wrapper so the total file size matches the target exactly.
#
# Usage: ./generate_html.sh [output_dir]
#   output_dir defaults to ./html

set -euo pipefail

OUTDIR="${1:-html}"
mkdir -p "${OUTDIR}"

# 15 target sizes in bytes.
SIZES=(128 256 512 1024 2048 4096 8192 16384 32768 65536 131072 262144 524288 1048576 2097152)

# Minimal HTML skeleton — the padding goes between the <body> tags.
HEAD='<!DOCTYPE html><html><head><title>Static</title></head><body>'
TAIL='</body></html>'

OVERHEAD=$(( ${#HEAD} + ${#TAIL} ))

for size in "${SIZES[@]}"; do
    file="${OUTDIR}/${size}.html"
    pad=$(( size - OVERHEAD ))
    if (( pad < 0 )); then
        pad=0
    fi
    # Write head, then padding, then tail.
    printf '%s' "${HEAD}" > "${file}"
    # Use dd for efficient padding generation.
    dd if=/dev/zero bs=1 count="${pad}" 2>/dev/null | tr '\0' 'A' >> "${file}"
    printf '%s' "${TAIL}" >> "${file}"
done

# index.html — simple directory listing linking to every size.
{
    echo '<!DOCTYPE html><html><head><title>Static File Index</title></head><body>'
    echo '<h1>Static Files</h1><ul>'
    for size in "${SIZES[@]}"; do
        echo "  <li><a href=\"/${size}.html\">${size} bytes</a></li>"
    done
    echo '</ul></body></html>'
} > "${OUTDIR}/index.html"

echo "Generated ${#SIZES[@]} static files + index.html in ${OUTDIR}/"
