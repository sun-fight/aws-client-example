

all:
	protoc --go_out=paths=source_relative:. ./*.proto
	sed -i '' 's/"`/" dynamodbav:",omitempty"`/g' *.go