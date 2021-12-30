curl -H "Content-Type: application/json" -X POST -d '{"username": "AzureWatcher", "content": "Restart azure"}' $AZURE_NOTIFY_URL
echo $AZURE_NOTIFY_URL
cd /server/azure
outfile=`date '+logs/azoutlog%Y%m%d%H%M%S.log'`
errfile=`date '+logs/azerrlog%Y%m%d%H%M%S.log'`
mv azoutlog.log $outfile
mv azerrlog.log $errfile
./azure 1> ./azoutlog.log 2> ./azerrlog.log
