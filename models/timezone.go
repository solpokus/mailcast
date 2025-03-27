package models

import "gorm.io/gorm"

type Timezonelist struct {
	gorm.Model
	AirportCode string `gorm:"column:airportcode;unique"`
	AirportName string `gorm:"column:airport_name"`
	CityCode    string `gorm:"column:citycode"`
	CityName    string `gorm:"column:cityname"`
	GmtTz       string `gorm:"column:gmttz"`
	TzName      string `gorm:"column:tzname"`
}

// TableName overrides the default table name
func (Timezonelist) TableName() string {
	return "new.m_timezone_list"
}
