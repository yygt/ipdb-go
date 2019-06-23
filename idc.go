package ipdb

import (
	"reflect"
	"time"
	"net"
	"os"
)

type IDCInfo struct {
	CountryName	string	`json:"country_name"`
	RegionName string 	`json:"region_name"`
	CityName string 	`json:"city_name"`
	OwnerDomain string 	`json:"owner_domain"`
	IspDomain string 	`json:"isp_domain"`
	IDC string 			`json:"idc"`
}

type IDC struct {
	reader *reader
}

func NewIDC(name string) (*IDC, error) {

	r, e := newReader(name, &IDCInfo{})
	if e != nil {
		return nil, e
	}

	return &IDC{
		reader: r,
	}, nil
}

func (db *IDC) Reload(name string) error {

	_, err := os.Stat(name)
	if err != nil {
		return err
	}

	reader, err := newReader(name, &IDCInfo{})
	if err != nil {
		return err
	}

	db.reader = reader

	return nil
}

func (db *IDC) Find(ip net.IP, language string) ([]string, error) {
	return db.reader.find1(ip, language)
}

func (db *IDC) FindMap(ip net.IP, language string) (map[string]string, error) {

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

func (db *IDC) FindInfo(ip net.IP, language string) (*IDCInfo, error) {

	data, err := db.reader.FindMap(ip, language)
	if err != nil {
		return nil, err
	}

	info := &IDCInfo{}

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

func (db *IDC) IsIPv4() bool {
	return db.reader.IsIPv4Support()
}

func (db *IDC) IsIPv6() bool {
	return db.reader.IsIPv6Support()
}

func (db *IDC) Languages() []string {
	return db.reader.Languages()
}

func (db *IDC) Fields() []string {
	return db.reader.meta.Fields
}

func (db *IDC) BuildTime() time.Time {
	return db.reader.Build()
}
