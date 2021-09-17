#!/bin/sh
set -e
APP=mocaccino-statusbar
VERSION=${VERSION:-0.4}
APPDIR=${DESTDIR:-${APP}_$VERSION}

mkdir -p $APPDIR/usr/bin
mkdir -p $APPDIR/usr/share/applications
mkdir -p $APPDIR/usr/share/icons/hicolor/1024x1024/apps
mkdir -p $APPDIR/usr/share/icons/hicolor/256x256/apps
mkdir -p $APPDIR/etc/xdg/autostart

go generate
go build -ldflags "-X \"main.Version=$VERSION\"" -o $APPDIR/usr/bin/$APP

cp icon/menu_icon.png $APPDIR/usr/share/icons/hicolor/1024x1024/apps/${APP}.png
cp icon/menu_icon.png $APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png

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

cp -rf $APPDIR/usr/share/applications/${APP}.desktop $APPDIR/etc/xdg/autostart/${APP}.desktop
