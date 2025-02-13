package utils

var Client *ESClient

var isClientExists bool = false

func Init() {
	if !isClientExists {
		Client = getElasticClient()
		isClientExists = true
	}
}
