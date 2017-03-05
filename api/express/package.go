package express

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/skitta/goquery/db"
)

// Package represents the structure of our resource
type Package struct {
	Key       string `json:"key"`
	IsChecked string `json:"ischeck"`
	Value     string `json:"value"`
}

// Track return a full Package by fetching data from kuaidi100
func Track(nu string) (p Package) {
	comCode := getComCode(nu)
	status, ischeck := getStatus(comCode, nu)
	p.Key = nu
	p.Value = status
	p.IsChecked = ischeck
	return
}

// GetAllPackages list packages
func GetAllPackages() []Package {
	Db := db.MgoDb{}
	Db.Init()
	result := []Package{}

	if err := Db.C("brain").Find(nil).All(&result); err != nil {
		return nil
	}
	Db.Close()
	return result
}

// Update a Package in to MongoDB
func (p *Package) Update() error {
	Db := db.MgoDb{}
	Db.Init()

	if err := Db.C("brain").Update(bson.M{"key": p.Key}, &p); err != nil {
		return err
	}

	Db.Close()
	return nil
}

// Del a Package from Redis
func (p *Package) Del() error {
	Db := db.MgoDb{}
	Db.Init()

	if err := Db.C("brain").Remove(bson.M{"key": p.Key}); err != nil {
		return err
	}

	Db.Close()
	return nil
}

// Checked a package
func (p *Package) Checked() bool {
	if p.IsChecked == "1" {
		return true
	}
	return false
}

// Private functions for getting package data from kuaidi100.com
// Test passed
func getComCode(nu string) string {
	qurl := "https://www.kuaidi100.com/autonumber/autoComNum?text=" + nu
	res, err := http.Get(qurl)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		data, _ := simplejson.NewFromReader(res.Body)
		comCode, _ := data.Get("auto").GetIndex(0).Get("comCode").String()
		return comCode
	}
	return ""
}

func getStatus(com, nu string) (status, ischeck string) {
	const packQueryURL = "http://m.kuaidi100.com/query?id=1&temp=0.3416398304980248"
	qurl := packQueryURL + "&type=" + com + "&postid=" + nu
	ischeck = "0"

	res, err := http.Get(qurl)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		data, _ := simplejson.NewFromReader(res.Body)
		ischeck, _ = data.Get("ischeck").String()
		msg, _ := data.Get("message").String()
		if msg == "ok" {
			statusJSON := data.Get("data").GetIndex(0)
			stime, _ := statusJSON.Get("time").String()
			context, _ := statusJSON.Get("context").String()
			status = fmt.Sprintf("快递 %s | %s | %s", nu, stime, context)
			return
		}
	}
	return
}
