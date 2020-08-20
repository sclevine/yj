#!/bin/sh

# find the man-dir
list="/usr/share/man /usr/man /usr/local/share/man /usr/local/man /tmp"
for mandir in $list; do
	[ -d $mandir ] && break
done	

# install... 
file=yj.1
cp yj.man $mandir/man1/$file
# for gzipped version
#gzip $mandir/man1/$file
#file=$file.gz
chmod 644 $mandir/man1/$file
chown root:root $mandir/man1/$file
