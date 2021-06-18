package weather_test

import (
	"errors"
	"testing"
	"time"

	"github.com/keep94/marvin2/weather"
	asserts "github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := weather.NewCache()
	defer cache.Close()
	observation, stale := cache.Get()
	if observation != nil {
		t.Error("Expected nil observation")
	}
	go func() {
		cache.Set(&weather.Observation{Temperature: 30.0})
	}()
	<-stale
	observation, stale = cache.Get()
	if observation.Temperature != 30.0 {
		t.Error("Expected 30.0 temperature")
	}
	go func() {
		cache.Set(&weather.Observation{Temperature: 35.0})
	}()
	<-stale
	observation, _ = cache.Get()
	if observation.Temperature != 35.0 {
		t.Error("Expected 35.0 temperature")
	}
}

func TestReportCache(t *testing.T) {
	assert := asserts.New(t)
	cache := weather.NewReportCache()
	defer cache.Close()
	var report weather.Report
	stale := cache.Get(&report)
	assert.Zero(report)
	go func() {
		report := weather.Report{Temperature: 25.0}
		cache.Set(&report)
		report.Temperature = 95.0
	}()
	<-stale
	stale = cache.Get(&report)
	assert.Equal(25.0, report.Temperature)
	go func() {
		cache.Set(&weather.Report{Temperature: 35.0})
	}()
	<-stale
	cache.Get(&report)
	assert.Equal(35.0, report.Temperature)
}

func TestAvgAQI(t *testing.T) {
	assert := asserts.New(t)
	conn := fakeConn{1001: 35, 1002: 100, 1003: 45}
	aqi, err := weather.AvgAQI(conn, time.Millisecond, 1001, 1002, 1003)
	assert.Equal(60, aqi)
	assert.NoError(err)
	aqi, err = weather.AvgAQI(conn, time.Millisecond, 1001, 1002, 9999)
	assert.Equal(68, aqi)
	assert.NoError(err)
	aqi, err = weather.AvgAQI(conn, time.Millisecond, 9999, 9998, 9997)
	assert.Error(err)
}

func TestAvgAQIPanics(t *testing.T) {
	assert := asserts.New(t)
	conn := fakeConn{1001: 35, 1002: 100, 1003: 45}
	assert.Panics(func() { weather.AvgAQI(conn, time.Millisecond) })
}

type fakeConn map[int64]int

func (f fakeConn) GetAQI(stationId int64) (int, error) {
	aqi, ok := f[stationId]
	if !ok {
		return 0, errors.New("No such stationId")
	}
	return aqi, nil
}
