package lib

import "mqueues"

var queue *mqueues.Queue

func SetQueue (mqueue *mqueues.Queue)  {
	queue = mqueue
}

func GetQueue () *mqueues.Queue {
	return queue
}