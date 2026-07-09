#!/bin/sh
## Example: a typical script with several problems
for f in $(ls *.m3u)
do
grep -qi hq.*mp3 $f \
&& echo -e 'Playlist $f contains a HQ file in mp3 format'
done

#!/bin/sh
## Example: The shebang says 'sh' so shellcheck warns about portability
## Change it to '#!/bin/bash' to allow bashisms
for n in {1..$RANDOM}
do
str=""
if (( n % 3 == 0 ))
then
str="fizz"
fi
if [ $[n%5] == 0 ]
then
str="$strbuzz"
fi
if [[ ! $str ]]
then
str="$n"
fi
echo "$str"
done

#!/bin/bash
## Example: ShellCheck can detect some higher level semantic problems
while getopts "nf:" param
do
case "$param" in
f) file="$OPTARG" ;;
v) set -x ;;
esac
done
case "$file" in
*.gz) gzip -d "$file" ;;
*.zip) unzip "$file" ;;
*.tar.gz) tar xzf "$file" ;;
*) echo "Unknown filetype" ;;
esac
if [[ "$$(uname)" == "Linux" ]]
then
echo "Using Linux"
fi

#!/bin/bash
## Example: ShellCheck can detect many different kinds of quoting issues
if ! grep -q backup=true.* "~/.myconfig"
then
echo 'Backup not enabled in $HOME/.myconfig, exiting'
exit 1
fi
if [[ $1 =~ "-v(erbose)?" ]]
then
verbose='-printf "Copying %f\n"'
fi
find backups/ \
-iname *.tar.gz \
$verbose \
-exec scp {} “myhost:backups” +
