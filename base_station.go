package ipdb

import (
	"reflect"
	"time"
	"net"
	"os"
)

type BaseStationInfo struct {
	CountryName	string	`json:"country_name"`
	RegionName string 	`json:"region_name"`
	CityName string 	`json:"city_name"`
	OwnerDomain string 	`json:"owner_domain"`
	IspDomain string 	`json:"isp_domain"`
	BaseStation string 	`json:"base_station"`
}

type BaseStation struct {
	reader *reader
}

func NewBaseStation(name string) (*BaseStation, error) {

	r, e := newReader(name, &BaseStationInfo{})
	if e != nil {
		return nil, e
	}

	return &BaseStation{
		reader: r,
	}, nil
}

func (db *BaseStation) Reload(name string) error {

	_, err := os.Stat(name)
	if err != nil {
		return err
	}

	reader, err := newReader(name, &BaseStationInfo{})
	if err != nil {
		return err
	}

	db.reader = reader

	return nil
}

func (db *BaseStation) Find(ip net.IP, language string) ([]string, error) {
	return db.reader.find1(ip, language)
}

func (db *BaseStation) FindMap(ip net.IP, language string) (map[string]string, error) {

	data, err := db.reader.find1(ip, language)
	if err != nil {
		return nil, err
	}
	info := make(map[string]string, len(db.reader.meta.Fields))
	for k, v := range data {
		info[db.reader.meta.Fields[k]] = v
	}

	return info, nil
}

func (db *BaseStation) FindInfo(ip net.IP, language string) (*BaseStationInfo, error) {

	data, err := db.reader.FindMap(ip, language)
	if err != nil {
		return nil, err
	}

	info := &BaseStationInfo{}

	for k, v := range data {
		sv := reflect.ValueOf(info).Elem()
		sfv := sv.FieldByName(db.reader.refType[k])

		if !sfv.IsValid() {
			continue
		}
		if !sfv.CanSet() {
			continue
		}

		sft := sfv.Type()
		fv := reflect.ValueOf(v)
		if sft == fv.Type() {
			sfv.Set(fv)
		}
	}

	return info, nil
}

func (db *BaseStation) IsIPv4() bool {
	return db.reader.IsIPv4Support()
}

func (db *BaseStation) IsIPv6() bool {
	return db.reader.IsIPv6Support()
}

func (db *BaseStation) Languages() []string {
	return db.reader.Languages()
}

func (db *BaseStation) Fields() []string {
	return db.reader.meta.Fields
}

func (db *BaseStation) BuildTime() time.Time {
	return db.reader.Build()
}
