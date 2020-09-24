package logx

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type LogX struct {
	logger    *zap.Logger
	atomLevel *zap.AtomicLevel
}

type Field = zap.Field
type FieldMap map[string]interface{}

// 公共的列
func getField(args ...interface{}) []Field {
	_, res := getFields("", args...)
	return res
}

// 判断其他类型--start
func getFields(format string, args ...interface{}) (string, []Field) {
	l := len(args)
	var str []interface{}
	var fields []Field
	if l > 0 {
		for _, v := range args {
			if f, ok := v.(Field); ok {
				fields = append(fields, f)
			} else {
				str = append(str, v)
			}
		}
		return fmt.Sprintf(format, str...), fields
	}
	return format, []Field{}
}

func (l *LogX) Debug(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.Debug(es, getField(args...)...)
	}
	return e
}
func (l *LogX) Info(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.Info(es, getField(args...)...)
	}
	return e
}
func (l *LogX) Warn(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.Warn(es, getField(args...)...)
	}
	return e
}
func (l *LogX) Error(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.Error(es, getField(args...)...)
	}
	return e
}
func (l *LogX) DPanic(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.DPanic(es, getField(args...)...)
	}
	return e
}
func (l *LogX) Panic(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.Panic(es, getField(args...)...)

	}
	return e
}
func (l *LogX) Fatal(s interface{}, args ...interface{}) error {
	es, e := checkErr(s)
	if es != "" {
		l.logger.Fatal(es, getField(args...)...)
	}
	return e
}

func checkErr(s interface{}) (string, error) {
	switch e := s.(type) {
	case error:
		return e.Error(), e
	case string:
		return e, errors.New(e)
	case []byte:
		return string(e), nil
	default:
		return "", nil
	}
}

func (l *LogX) LogError(err error) error {
	return l.Error(err.Error())
}

func (l *LogX) Debugf(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.Debug(s, f...)
	return errors.New(s)
}

func (l *LogX) Infof(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.Info(s, f...)
	return errors.New(s)
}

func (l *LogX) Warnf(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.Warn(s, f...)
	return errors.New(s)
}

func (l *LogX) Errorf(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.Error(s, f...)
	return errors.New(s)
}

func (l *LogX) DPanicf(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.DPanic(s, f...)
	return errors.New(s)
}

func (l *LogX) Panicf(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.Panic(s, f...)
	return errors.New(s)
}

func (l *LogX) Fatalf(format string, args ...interface{}) error {
	s, f := getFields(format, args...)
	l.logger.Fatal(s, f...)
	return errors.New(s)
}
