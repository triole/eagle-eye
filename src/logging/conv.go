package logging

import "github.com/sirupsen/logrus"

type F map[string]interface{}

func (lg Logging) conv(fields F) logrus.Fields {
	if fields != nil {
		return logrus.Fields(fields)
	}
	return logrus.Fields{}
}
