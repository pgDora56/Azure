curl -H "Content-Type: application/json" -X POST -d '{"username": "AzureWatcher", "content": "Restart azure"}' https://discord.com/api/webhooks/908894969073385512/TeKXr1Yv-BBQ-SxHgb7iJHTRMzHba83jLBOXCin8maQLX3aMob8nmLAbyWHPd4gN8qtR
cd /server/azure
./azure > ./azlog.log
