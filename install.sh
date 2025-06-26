#!/bin/bash

set -e  # Stop jika ada error

TARGET_DIR="/Volumes/Haliim-SSD/project"

echo "📥 Meng-clone repo Spago..."
git clone https://github.com/nlpodyssey/spago.git
cd spago

echo "🔄 Checkout ke commit CLI terakhir (50d8191)..."
git checkout v0.3.0



echo "🔨 Build CLI Spago..."
go build -o spago ./cmd/spago

echo "📁 Pindahkan binary ke $TARGET_DIR..."
mv spago "$TARGET_DIR"

echo "✅ Instalasi selesai!"
echo "🔧 Sekarang jalankan Spago dari folder project: "
echo "    $TARGET_DIR/spago --help"
