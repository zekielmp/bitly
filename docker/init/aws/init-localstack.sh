#1/bin/bash

#Create bucket
awslocal s3 nb s3://bitly-uploads

# Create SQS queue
awsLocal sqs create-queue --queue-name bitly-events

echo "Localstack initialization complete"
