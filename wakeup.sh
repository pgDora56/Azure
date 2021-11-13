curl -H "Content-Type: application/json" -X POST -d '{"username": "AzureWatcher", "content": "Restart azure"}' $AZURE_NOTIFY_URL
echo $AZURE_NOTIFY_URL
cd /server/azure
./azure > ./azlog.log
