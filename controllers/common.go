package controllers

import (
	"flight-details/utils"
)

var Client *utils.ESClient

func Init() {
	Client = utils.GetElasticClient()
}