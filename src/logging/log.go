package logging

func (lg Logging) Debug(msg string, fields F) {
	lg.Logrus.WithFields(lg.conv(fields)).Debug(msg)
}

func (lg Logging) Info(msg string, fields F) {
	lg.Logrus.WithFields(lg.conv(fields)).Info(msg)
}

func (lg Logging) Error(msg interface{}, fields F) {
	lg.Logrus.WithFields(lg.conv(fields)).Error(msg)
}

func (lg Logging) Fatal(msg string, fields F) {
	lg.Logrus.WithFields(lg.conv(fields)).Fatal(msg)
}
