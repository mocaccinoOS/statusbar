#!/bin/sh
set -e
APP=mocaccino-statusbar
APPDIR=${DESTDIR:-${APP}_1.0.0}

mkdir -p $APPDIR/usr/bin
mkdir -p $APPDIR/usr/share/applications
mkdir -p $APPDIR/usr/share/icons/hicolor/1024x1024/apps
mkdir -p $APPDIR/usr/share/icons/hicolor/256x256/apps

go generate
go build -o $APPDIR/usr/bin/$APP

cp icon/icon.png $APPDIR/usr/share/icons/hicolor/1024x1024/apps/${APP}.png
cp icon/icon.png $APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png

cat > $APPDIR/usr/share/applications/${APP}.desktop << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=$APP
Exec=$APP
Icon=$APP
Terminal=false
StartupWMClass=Lorca
EOF