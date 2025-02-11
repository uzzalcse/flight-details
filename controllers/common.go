package controllers

import (
	"flight_api/utils"
)

var Client *utils.ESClient

func Init() {
	Client = utils.GetElasticClient()
}