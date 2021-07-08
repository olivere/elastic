module github.com/olivere/elastic/recipes/aws-es-client

go 1.16

require (
	github.com/aws/aws-sdk-go v1.39.2
	github.com/olivere/elastic/v7 v7.0.26
	github.com/olivere/env v1.1.0
)

replace github.com/olivere/elastic/v7 => ../..
