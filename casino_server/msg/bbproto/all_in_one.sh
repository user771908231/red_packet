OUT_FILE=./client/bbproto.proto
rm -f $OUT_FILE
echo 'package bbproto;' > $OUT_FILE
cat *.proto|grep -v "^import.*"|sed 's/^package.*/\/\/=====================================================/g' >> $OUT_FILE
echo 'all in one done.'
ls -lrt ./client
cp -vf $OUT_FILE /Users/kory/Desktop/Share/client_proto/