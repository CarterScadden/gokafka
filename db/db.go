package db

import "gokafka/topic"

var (
	// AllTopics Holds all the topics
	AllTopics = make(chan []topic.Topic)
)
