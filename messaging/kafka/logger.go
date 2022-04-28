package kafka

import "github.com/simiancreative/simiango/logger"

type Logger struct{}

func (l Logger) Printf(pattern string, replacements ...interface{}) {
	logger.Debugf("Kafka: "+pattern, replacements...)
}
