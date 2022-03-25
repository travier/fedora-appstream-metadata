// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type PkgDbCollections struct {
	Collections []Collection `json:"collections"`
	Output      string       `json:"output"`
}

type Collection struct {
	AllowRetire bool   `json:"allow_retire"`
	Branchname  string `json:"branchname"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	DistTag     string `json:"dist_tag"`
	KojiName    string `json:"koji_name"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Version     string `json:"version"`
}

type Component struct {
	XMLName     xml.Name `xml:"component"`
	Text        string   `xml:",chardata"`
	Type        string   `xml:"type,attr"`
	ID          string   `xml:"id"`
	Name        string   `xml:"name"`
	Summary     string   `xml:"summary"`
	Description struct {
		Text string `xml:",chardata"`
		P    string `xml:"p"`
	} `xml:"description"`
	URL struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"url"`
	MetadataLicense string `xml:"metadata_license"`
	DeveloperName   string `xml:"developer_name"`
	Releases        struct {
		Text    string    `xml:",chardata"`
		Release []Release `xml:"release"`
	} `xml:"releases"`
}

type Release struct {
	Text        string `xml:",chardata"`
	Version     string `xml:"version,attr"`
	Type        string `xml:"type,attr"`
	Date        string `xml:"date,attr"`
	DateEol     string `xml:"date_eol,attr"`
	Description struct {
		Text string `xml:",chardata"`
		P    string `xml:"p"`
	} `xml:"description"`
}

type Schedule struct {
	Changelog struct {
	} `json:"changelog"`
	End     string `json:"end"`
	ExtAttr struct {
		Text1 int `json:"Text1"`
		Text2 int `json:"Text2"`
		Text3 int `json:"Text3"`
	} `json:"ext_attr"`
	FlagsAttrID interface{} `json:"flags_attr_id"`
	Name        string      `json:"name"`
	Resources   struct {
	} `json:"resources"`
	Slug  string `json:"slug"`
	Start string `json:"start"`
	Tasks []struct {
		Level      int           `json:"_level"`
		Complete   float64       `json:"complete"`
		End        string        `json:"end"`
		Flags      []interface{} `json:"flags"`
		Index      int           `json:"index"`
		Name       string        `json:"name"`
		ParentTask string        `json:"parentTask"`
		Priority   int           `json:"priority"`
		Slug       string        `json:"slug"`
		Start      string        `json:"start"`
		Type       string        `json:"type"`
		Link       string        `json:"link,omitempty"`
	} `json:"tasks"`
	UsedFlags []string `json:"used_flags"`
}

var pkgdb string = "https://admin.fedoraproject.org/pkgdb/api/collections/"
var schedule string = "https://fedorapeople.org/groups/schedule/f-%s/f-%s-all-milestones.json"

func getJSON(url string) []byte {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func main() {
	var body []byte
	if len(os.Args[1:]) > 0 {
		jsonFile, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		body, _ = ioutil.ReadAll(jsonFile)
	} else {
		body = getJSON(pkgdb)
	}
	collections := PkgDbCollections{}
	err := json.Unmarshal([]byte(body), &collections)
	if err != nil {
		log.Fatal(err)
	}
	if collections.Output != "ok" {
		log.Fatalln("Invalid value for output in JSON document")
	}

	m := Component{}
	m.Type = "operating-system"
	m.ID = "org.fedoraproject.fedora"
	m.Name = "Fedora Linux"
	m.Summary = "Fedora Linux distribution from the Fedora Project"
	m.Description.P = "Fedora creates an innovative, free, and open source platform for hardware, clouds, and containers that enables software developers and community members to build tailored solutions for their users."
	m.URL.Text = "https://fedoraproject.org/"
	m.URL.Type = "homepage"
	m.MetadataLicense = "MIT"
	m.DeveloperName = "The Fedora Project"

	for _, c := range collections.Collections {
		// Filter only Fedora & Fedora Linux
		if !(c.Name == "Fedora" || c.Name == "Fedora Linux") {
			continue
		}
		// Ignore too old versions
		if c.Name == "Fedora" && c.Version < "30" {
			continue
		}
		// Ignore duplicates following the rename to Fedora Linux
		if c.Version >= "35" && c.Name == "Fedora" {
			continue
		}
		r := Release{}
		if c.Version == "devel" {
			r.Version = "Rawhide"
			r.Type = "development"
			r.Description.P = "This is the current version for Rawhide, which is a continuous rolling development branch. No releases are ever made directly from Rawhide, and it never freezes. There is no guarantee of stability. Rawhide is intended for initial testing of the very latest code under active development."
		} else if c.Status == "Under Development" {
			r.Version = c.Version
			r.Type = "development"
			r.Description.P = "The upcoming release of Fedora Linux."
		} else if c.Status == "Active" || c.Status == "EOL" {
			r.Version = c.Version
			r.Type = "stable"
			r.Description.P = fmt.Sprintf("https://fedoramagazine.org/announcing-fedora-%s/", c.Version)

			// Get the release & EOL dates from the schedule
			schedRaw := getJSON(fmt.Sprintf(schedule, c.Version, c.Version))
			schedule := Schedule{}
			err = json.Unmarshal(schedRaw, &schedule)
			if err != nil {
				log.Fatal(err)
			}

			for _, t := range schedule.Tasks {
				if t.Name == "Current Final Target date" && t.Type == "Milestone" {
					res, err := strconv.Atoi(t.Start)
					if err != nil {
						log.Fatal(err)
					}
					tm := time.Unix(int64(res), 0)
					// log.Printf("Current Final Target date: %d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
					r.Date = fmt.Sprintf("%d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
				}
			}
			if c.Status == "EOL" {
				for _, t := range schedule.Tasks {
					if t.Name == "EOL" && t.Type == "Milestone" {
						res, err := strconv.Atoi(t.Start)
						if err != nil {
							log.Fatal(err)
						}
						tm := time.Unix(int64(res), 0)
						// log.Printf("EOL: %d-%02d-%02d", tm.UTC().Year(), tm.UTC().Month(), tm.UTC().Day())
						r.DateEol = fmt.Sprintf("%d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
					}
				}
			}
		} else {
			log.Fatal("Should never be reached")
		}
		m.Releases.Release = append(m.Releases.Release, r)
	}

	out, _ := xml.MarshalIndent(m, "", "  ")

	f, err := os.Create("org.fedoraproject.fedora.metainfo.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(out)
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString("\n")
	f.Sync()
}
