package express

import (
	"fmt"
	"net/http"

	"log"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/skitta/goquery/db"
)

// Package represents the structure of our resource
type Package struct {
	ID        string
	IsChecked string
	Value     string
}

// GetAllPackages list packages
func GetAllPackages() []Package {
	result := []Package{}
	Db := db.Redis{}
	Db.Init()
	defer Db.Close()

	index, err := Db.Index()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(index); i++ {
		value, err := Db.Client.Get(index[i]).Result()
		if err != nil {
			log.Fatal(err)
		}
		p := Package{index[i], "0", value}
		result = append(result, p)
	}

	return result
}

// Update a Package in to sqlite
func (p *Package) Update() error {
	Db := db.Redis{}
	Db.Init()
	defer Db.Close()

	err := Db.Update(p.ID, p.Value)
	if err != nil {
		return err
	}
	return nil
}

// Del a Package from Redis
func (p *Package) Del() error {
	Db := db.Redis{}
	Db.Init()
	defer Db.Close()

	err := Db.Del(p.ID)
	if err != nil {
		return err
	}
	return nil
}

// Checked a package
func (p *Package) Checked() bool {
	if p.IsChecked == "1" {
		return true
	}
	return false
}

// Track return a full Package by fetching data from kuaidi100
func Track(nu string) (p Package) {
	comCode := getComCode(nu)
	status, ischeck := getStatus(comCode, nu)
	p.ID = nu
	p.Value = status
	p.IsChecked = ischeck
	return
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
			status = fmt.Sprintf("%s | %s", stime, context)
			return
		}
	}
	return
}
