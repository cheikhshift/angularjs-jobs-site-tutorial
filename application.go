package main

import (
	//iogos-replace
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cheikhshift/db"
	netform "github.com/cheikhshift/form"
	"github.com/cheikhshift/gos/core"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/fatih/color"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/opentracing/opentracing-go"
	"html"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
	"strings"
	"time"
	"unsafe"
)

var store = sessions.NewCookieStore([]byte("a very very very very secret key"))

type NoStruct struct {
	/* emptystruct */
}

func NetsessionGet(key string, s *sessions.Session) string {
	return s.Values[key].(string)
}

func UrlAtZ(url, base string) (isURL bool) {
	isURL = strings.Index(url, base) == 0
	return
}

func NetsessionDelete(s *sessions.Session) string {
	//keys := make([]string, len(s.Values))

	//i := 0
	for k := range s.Values {
		// keys[i] = k.(string)
		NetsessionRemove(k.(string), s)
		//i++
	}

	return ""
}

func NetsessionRemove(key string, s *sessions.Session) string {
	delete(s.Values, key)
	return ""
}
func NetsessionKey(key string, s *sessions.Session) bool {
	if _, ok := s.Values[key]; ok {
		//do something here
		return true
	}

	return false
}

func Netadd(x, v float64) float64 {
	return v + x
}

func Netsubs(x, v float64) float64 {
	return v - x
}

func Netmultiply(x, v float64) float64 {
	return v * x
}

func Netdivided(x, v float64) float64 {
	return v / x
}

func NetsessionGetInt(key string, s *sessions.Session) interface{} {
	return s.Values[key]
}

func NetsessionSet(key string, value string, s *sessions.Session) string {
	s.Values[key] = value
	return ""
}
func NetsessionSetInt(key string, value interface{}, s *sessions.Session) string {
	s.Values[key] = value
	return ""
}

func dbDummy() {
	smap := db.O{}
	smap["key"] = "set"
	log.Println(smap)
}

func Netimportcss(s string) string {
	return fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\" /> ", s)
}

func Netimportjs(s string) string {
	return fmt.Sprintf("<script type=\"text/javascript\" src=\"%s\" ></script> ", s)
}

func formval(s string, r *http.Request) string {
	return r.FormValue(s)
}

func renderTemplate(w http.ResponseWriter, p *Page, span opentracing.Span) bool {
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path : web%s.tmpl reason : %s", p.R.URL.Path, n))

			DebugTemplate(w, p.R, fmt.Sprintf("web%s", p.R.URL.Path))
			w.WriteHeader(http.StatusInternalServerError)

			pag, err := loadPage("/your-500-page")

			if err != nil {
				log.Println(err.Error())
				return
			}

			if pag.isResource {
				w.Write(pag.Body)
			} else {
				pag.R = p.R
				pag.Session = p.Session
				renderTemplate(w, pag, span) ///your-500-page"

			}
		}
	}()

	var sp opentracing.Span
	opName := fmt.Sprintf("Building template %s%s", p.R.URL.Path, ".tmpl")

	if true {
		carrier := opentracing.HTTPHeadersCarrier(p.R.Header)
		wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			sp = opentracing.StartSpan(opName)
		} else {
			sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
		}
	}
	defer sp.Finish()

	t := template.New("PageWrapper")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(p.Body))
	outp := new(bytes.Buffer)
	err := t.Execute(outp, p)
	if err != nil {
		log.Println(err.Error())
		DebugTemplate(w, p.R, fmt.Sprintf("web%s", p.R.URL.Path))
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		pag, err := loadPage("/your-500-page")

		if err != nil {
			log.Println(err.Error())
			return false
		}
		pag.R = p.R
		pag.Session = p.Session
		p = nil
		if pag.isResource {
			w.Write(pag.Body)
		} else {
			renderTemplate(w, pag, span) // "/your-500-page"

		}
		return false
	}

	p.Session.Save(p.R, w)

	fmt.Fprintf(w, html.UnescapeString(outp.String()))
	p.Session = nil
	p.Body = nil
	p.R = nil
	p = nil
	return true

}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string, *sessions.Session, opentracing.Span)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		span := opentracing.StartSpan(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer span.Finish()
		carrier := opentracing.HTTPHeadersCarrier(r.Header)
		if err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier); err != nil {
			log.Fatalf("Could not inject span context into header: %v", err)
		}

		var session *sessions.Session
		var er error
		if session, er = store.Get(r, "session-"); er != nil {
			session, _ = store.New(r, "session-")
		}
		if attmpt := apiAttempt(w, r, session, span); !attmpt {
			fn(w, r, "", session, span)
		} else {
			context.Clear(r)
		}

	}
}

func mResponse(v interface{}) string {
	data, _ := json.Marshal(&v)
	return string(data)
}
func apiAttempt(w http.ResponseWriter, r *http.Request, session *sessions.Session, span opentracing.Span) (callmet bool) {
	var response string
	response = ""

	if strings.Contains(r.URL.Path, "") {

		lastLine := ""
		var sp opentracing.Span
		opName := fmt.Sprintf(" [ValidateForm]%s %s", r.Method, r.URL.Path)

		if true {
			carrier := opentracing.HTTPHeadersCarrier(r.Header)
			wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
			if err != nil {
				sp = opentracing.StartSpan(opName)
			} else {
				sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
			}
		}
		defer sp.Finish()
		defer func() {
			if n := recover(); n != nil {
				log.Println("Web request () failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml :", strings.TrimSpace(lastLine))
				log.Println("Reason : ", n)
				//wheredefault
				span.SetTag("error", true)
				span.LogEvent(fmt.Sprintf("%s request at %s, reason : %s ", r.Method, r.URL.Path, n))

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/html")
				pag, err := loadPage("/your-500-page")

				if err != nil {
					log.Println(err.Error())
					callmet = true
					return
				}
				pag.R = r
				pag.Session = session
				if pag.isResource {
					w.Write(pag.Body)
				} else {
					// renderTemplate(w, pag, span)

				}

				callmet = true
			}
		}()
		lastLine = `if _, ok := session.Values["formtoken"]; !ok {`
		if _, ok := session.Values["formtoken"]; !ok {
			lastLine = `session.Values["formtoken"] = core.NewLen(10)`
			session.Values["formtoken"] = core.NewLen(10)
			lastLine = `session.Save(r, w)`
			session.Save(r, w)
			lastLine = `}`
		}
		lastLine = `if r.ContentLength > 0	{`
		if r.ContentLength > 0 {
			lastLine = `if !netform.ValidateRequest(r, session.Values["formtoken"].(string) ) || r.ContentLength > int64(netform.MaxSize * netform.MB) {`
			if !netform.ValidateRequest(r, session.Values["formtoken"].(string)) || r.ContentLength > int64(netform.MaxSize*netform.MB) {
				lastLine = `w.WriteHeader(http.StatusBadRequest)`
				w.WriteHeader(http.StatusBadRequest)
				lastLine = `w.Header().Set("Content-Type",  "text/xml")`
				w.Header().Set("Content-Type", "text/xml")
				lastLine = `w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Error>Invalid request sent</Error>"))`
				w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Error>Invalid request sent</Error>"))
				lastLine = `return true`
				return true
				lastLine = `}`
			}
			lastLine = `}`
		}

	}

	if isURL := (r.URL.Path == "/test/form-build" && r.Method == strings.ToUpper("GET")); !callmet && isURL {

		lastLine := ""
		var sp opentracing.Span
		opName := fmt.Sprintf(" []%s %s", r.Method, r.URL.Path)

		if true {
			carrier := opentracing.HTTPHeadersCarrier(r.Header)
			wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
			if err != nil {
				sp = opentracing.StartSpan(opName)
			} else {
				sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
			}
		}
		defer sp.Finish()
		defer func() {
			if n := recover(); n != nil {
				log.Println("Web request (/test/form-build) failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml :", strings.TrimSpace(lastLine))
				log.Println("Reason : ", n)
				//wheredefault
				span.SetTag("error", true)
				span.LogEvent(fmt.Sprintf("%s request at %s, reason : %s ", r.Method, r.URL.Path, n))
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/html")
				pag, err := loadPage("/your-500-page")

				if err != nil {
					log.Println(err.Error())
					callmet = true
					return
				}
				pag.R = r
				pag.Session = session
				if pag.isResource {
					w.Write(pag.Body)
				} else {
					renderTemplate(w, pag, span) //"s"

				}
				callmet = true

			}
		}()

		lastLine = `w.Header().Set("Content-Type",  "text/html")`
		w.Header().Set("Content-Type", "text/html")
		lastLine = `SampleForm := SampleForm{Text: "Sample",Created:"2017-05-05" ,Emal:"sample",FieldF:"orange",Count:500,Terms : true}`
		SampleForm := SampleForm{Text: "Sample", Created: "2017-05-05", Emal: "sample", FieldF: "orange", Count: 500, Terms: true}
		lastLine = `w.Write([]byte(NetBuild(&SampleForm, "/target/url", "POST", "Update",session) ))`
		w.Write([]byte(NetBuild(&SampleForm, "/target/url", "POST", "Update", session)))
		callmet = true
	}

	if isURL := (r.URL.Path == "/target/url" && r.Method == strings.ToUpper("POST")); !callmet && isURL {

		lastLine := ""
		var sp opentracing.Span
		opName := fmt.Sprintf(" []%s %s", r.Method, r.URL.Path)

		if true {
			carrier := opentracing.HTTPHeadersCarrier(r.Header)
			wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
			if err != nil {
				sp = opentracing.StartSpan(opName)
			} else {
				sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
			}
		}
		defer sp.Finish()
		defer func() {
			if n := recover(); n != nil {
				log.Println("Web request (/target/url) failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml :", strings.TrimSpace(lastLine))
				log.Println("Reason : ", n)
				//wheredefault
				span.SetTag("error", true)
				span.LogEvent(fmt.Sprintf("%s request at %s, reason : %s ", r.Method, r.URL.Path, n))
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/html")
				pag, err := loadPage("/your-500-page")

				if err != nil {
					log.Println(err.Error())
					callmet = true
					return
				}
				pag.R = r
				pag.Session = session
				if pag.isResource {
					w.Write(pag.Body)
				} else {
					renderTemplate(w, pag, span) //"s"

				}
				callmet = true

			}
		}()

		lastLine = `var sampleform SampleForm`
		var sampleform SampleForm
		lastLine = `if err := netform.Form(r, &sampleform); err != nil {`
		if err := netform.Form(r, &sampleform); err != nil {
			lastLine = `w.WriteHeader(http.StatusBadRequest)`
			w.WriteHeader(http.StatusBadRequest)
			lastLine = `w.Header().Set("Content-Type",  "text/xml")`
			w.Header().Set("Content-Type", "text/xml")
			lastLine = `w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Error>%s Value :%v</Error>", err.Error(),sampleform )))`
			w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" encoding=\"UTF-8\"?><Error>%s Value :%v</Error>", err.Error(), sampleform)))
			lastLine = `return true`
			return true
			lastLine = `}`
		}
		lastLine = `response = mResponse(sampleform)`
		response = mResponse(sampleform)
		lastLine = `//http.ServeFile(w, r, netform.Path(sampleform.Photo))`
		//http.ServeFile(w, r, netform.Path(sampleform.Photo))
		callmet = true
	}

	if callmet {
		session.Save(r, w)
		if response != "" {
			//Unmarshal json
			//w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(response))
		}
		return
	}
	return
}
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}

	structFieldValue.Set(val)
	return nil
}
func DebugTemplate(w http.ResponseWriter, r *http.Request, tmpl string) {
	lastline := 0
	linestring := ""
	defer func() {
		if n := recover(); n != nil {
			log.Println()
			// log.Println(n)
			log.Println("Error on line :", lastline+1, ":"+strings.TrimSpace(linestring))
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()

	p, err := loadPage(r.URL.Path)
	filename := tmpl + ".tmpl"
	body, err := Asset(filename)
	session, er := store.Get(r, "session-")

	if er != nil {
		session, er = store.New(r, "session-")
	}
	p.Session = session
	p.R = r
	if err != nil {
		log.Print(err)

	} else {

		lines := strings.Split(string(body), "\n")
		// log.Println( lines )
		linebuffer := ""
		waitend := false
		open := 0
		for i, line := range lines {

			processd := false

			if strings.Contains(line, "{{with") || strings.Contains(line, "{{ with") || strings.Contains(line, "with}}") || strings.Contains(line, "with }}") || strings.Contains(line, "{{range") || strings.Contains(line, "{{ range") || strings.Contains(line, "range }}") || strings.Contains(line, "range}}") || strings.Contains(line, "{{if") || strings.Contains(line, "{{ if") || strings.Contains(line, "if }}") || strings.Contains(line, "if}}") || strings.Contains(line, "{{block") || strings.Contains(line, "{{ block") || strings.Contains(line, "block }}") || strings.Contains(line, "block}}") {
				linebuffer += line
				waitend = true

				endstr := ""
				processd = true
				if !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {

					open++

				}
				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
				t, _ = t.Parse(ReadyTemplate(body))
				lastline = i
				linestring = line
				erro := t.Execute(outp, p)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if waitend && !processd && !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end")) {
				linebuffer += line

				endstr := ""
				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
				t, _ = t.Parse(ReadyTemplate(body))
				lastline = i
				linestring = line
				erro := t.Execute(outp, p)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}

			}

			if !waitend && !processd {
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
				t, _ = t.Parse(ReadyTemplate(body))
				lastline = i
				linestring = line
				erro := t.Execute(outp, p)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if !processd && (strings.Contains(line, "{{end") || strings.Contains(line, "{{ end")) {
				open--

				if open == 0 {
					waitend = false

				}
			}
		}

	}

}

func DebugTemplatePath(tmpl string, intrf interface{}) {
	lastline := 0
	linestring := ""
	defer func() {
		if n := recover(); n != nil {

			log.Println("Error on line :", lastline+1, ":"+strings.TrimSpace(linestring))
			log.Println(n)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()

	filename := tmpl
	body, err := Asset(filename)

	if err != nil {
		log.Print(err)

	} else {

		lines := strings.Split(string(body), "\n")
		// log.Println( lines )
		linebuffer := ""
		waitend := false
		open := 0
		for i, line := range lines {

			processd := false

			if strings.Contains(line, "{{with") || strings.Contains(line, "{{ with") || strings.Contains(line, "with}}") || strings.Contains(line, "with }}") || strings.Contains(line, "{{range") || strings.Contains(line, "{{ range") || strings.Contains(line, "range }}") || strings.Contains(line, "range}}") || strings.Contains(line, "{{if") || strings.Contains(line, "{{ if") || strings.Contains(line, "if }}") || strings.Contains(line, "if}}") || strings.Contains(line, "{{block") || strings.Contains(line, "{{ block") || strings.Contains(line, "block }}") || strings.Contains(line, "block}}") {
				linebuffer += line
				waitend = true

				endstr := ""
				if !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {

					open++

				}

				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec

				processd = true
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
				t, _ = t.Parse(ReadyTemplate([]byte(fmt.Sprintf("%s%s", linebuffer, endstr))))
				lastline = i
				linestring = line
				erro := t.Execute(outp, intrf)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if waitend && !processd && !(strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {
				linebuffer += line

				endstr := ""
				for i := 0; i < open; i++ {
					endstr += "\n{{end}}"
				}
				//exec
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
				t, _ = t.Parse(ReadyTemplate([]byte(fmt.Sprintf("%s%s", linebuffer, endstr))))
				lastline = i
				linestring = line
				erro := t.Execute(outp, intrf)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}

			}

			if !waitend && !processd {
				outp := new(bytes.Buffer)
				t := template.New("PageWrapper")
				t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
				t, _ = t.Parse(ReadyTemplate([]byte(fmt.Sprintf("%s%s", linebuffer))))
				lastline = i
				linestring = line
				erro := t.Execute(outp, intrf)
				if erro != nil {
					log.Println("Error on line :", i+1, line, erro.Error())
				}
			}

			if !processd && (strings.Contains(line, "{{end") || strings.Contains(line, "{{ end") || strings.Contains(line, "end}}") || strings.Contains(line, "end }}")) {
				open--

				if open == 0 {
					waitend = false

				}
			}
		}

	}

}
func Handler(w http.ResponseWriter, r *http.Request, contxt string, session *sessions.Session, span opentracing.Span) {
	var p *Page
	p, err := loadPage(r.URL.Path)

	var sp opentracing.Span
	opName := fmt.Sprintf(fmt.Sprintf("Web:/%s", r.URL.Path))

	if true {
		carrier := opentracing.HTTPHeadersCarrier(r.Header)
		wireContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err != nil {
			sp = opentracing.StartSpan(opName)
		} else {
			sp = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
		}
	}
	defer sp.Finish()

	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)
		span.SetTag("error", true)
		span.LogEvent(fmt.Sprintf("%s request at %s, reason : %s ", r.Method, r.URL.Path, err))
		pag, err := loadPage("/your-404-page")

		if err != nil {
			log.Println(err.Error())
			//context.Clear(r)
			return
		}
		pag.R = r
		pag.Session = session
		if p != nil {
			p.Session = nil
			p.Body = nil
			p.R = nil
			p = nil
		}
		if pag.isResource {
			w.Write(pag.Body)
		} else {
			renderTemplate(w, pag, span) //"/your-500-page"
		}
		context.Clear(r)
		return
	}

	if !p.isResource {
		w.Header().Set("Content-Type", "text/html")
		p.Session = session
		p.R = r
		renderTemplate(w, p, span) //fmt.Sprintf("web%s", r.URL.Path)

		// log.Println(w)
	} else {
		w.Header().Set("Cache-Control", "public")
		if strings.Contains(r.URL.Path, ".css") {
			w.Header().Add("Content-Type", "text/css")
		} else if strings.Contains(r.URL.Path, ".js") {
			w.Header().Add("Content-Type", "application/javascript")
		} else {
			w.Header().Add("Content-Type", http.DetectContentType(p.Body))
		}

		w.Write(p.Body)
	}

	p.Session = nil
	p.Body = nil
	p.R = nil
	p = nil
	context.Clear(r)
	return
}

func loadPage(title string) (*Page, error) {

	if roottitle := (title == "/"); roottitle {
		webbase := "web/"
		fname := fmt.Sprintf("%s%s", webbase, "index.html")
		body, err := Asset(fname)
		if err != nil {
			fname = fmt.Sprintf("%s%s", webbase, "index.tmpl")
			body, err = Asset(fname)
			if err != nil {
				return nil, err
			}
			return &Page{Body: body, isResource: false}, nil
		}

		return &Page{Body: body, isResource: true}, nil

	}

	filename := fmt.Sprintf("web%s.tmpl", title)

	if body, err := Asset(filename); err != nil {
		filename = fmt.Sprintf("web%s.html", title)

		if body, err = Asset(filename); err != nil {
			filename = fmt.Sprintf("web%s", title)

			if body, err = Asset(filename); err != nil {
				return nil, err
			} else {
				if strings.Contains(title, ".tmpl") {
					return nil, nil
				}
				return &Page{Body: body, isResource: true}, nil
			}
		} else {
			return &Page{Body: body, isResource: true}, nil
		}
	} else {
		return &Page{Body: body, isResource: false}, nil
	}

	//wheredefault

}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
func equalz(args ...interface{}) bool {
	if args[0] == args[1] {
		return true
	}
	return false
}
func nequalz(args ...interface{}) bool {
	if args[0] != args[1] {
		return true
	}
	return false
}

func netlt(x, v float64) bool {
	if x < v {
		return true
	}
	return false
}
func netgt(x, v float64) bool {
	if x > v {
		return true
	}
	return false
}
func netlte(x, v float64) bool {
	if x <= v {
		return true
	}
	return false
}

func GetLine(fname string, match string) int {
	intx := 0
	file, err := os.Open(fname)
	if err != nil {
		color.Red("Could not find a source file")
		return -1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		intx = intx + 1
		if strings.Contains(scanner.Text(), match) {

			return intx
		}

	}

	return -1
}
func netgte(x, v float64) bool {
	if x >= v {
		return true
	}
	return false
}

type Page struct {
	Title      string
	Body       []byte
	isResource bool
	R          *http.Request
	Session    *sessions.Session
}

func ReadyTemplate(body []byte) string {
	return strings.Replace(strings.Replace(strings.Replace(string(body), "/{", "\"{", -1), "}/", "}\"", -1), "`", "\"", -1)
}

var jobs = []Job{}

func init() {

}

type SampleForm struct {
	TestField string           `title:"Hi world!",valid:"unique",placeholder:"Testfield prompt"`
	Count     int              `placeholder:"Count"`
	Name      netform.Password `valid:"required",title:"Input title"`
	FieldTwo  netform.Radio    `title:"Enter Email",valid:"email,unique,required",select:"blue,orange,red,green"`
	FieldF    netform.Select   `placeholder:"Prompt?",valid:"email,unique,required",select:"blue,orange,red,green"`
	Created   netform.Date
	Text      netform.Paragraph `title:"Enter a description."`
	Photo     netform.File      `file:"image/*"`
	Emal      netform.Email
	Terms     bool `title:"Accept terms of use."`
}

func NetcastSampleForm(args ...interface{}) *SampleForm {

	s := SampleForm{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructSampleForm() *SampleForm { return &SampleForm{} }

type fInput struct {
	Name, Type, Placeholder, Title, Class, Value string
	Choices                                      []string
	Required                                     bool
}

func NetcastfInput(args ...interface{}) *fInput {

	s := fInput{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructfInput() *fInput { return &fInput{} }

type afInput struct {
	Name, Type, Placeholder, Title, Class, Value string
	Choices                                      []string
	Required                                     bool
	ModelName                                    string
}

func NetcastafInput(args ...interface{}) *afInput {

	s := afInput{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructafInput() *afInput { return &afInput{} }

type fForm struct {
	Target, Method, Token, ButtonClass, CTA string
	Input                                   []fInput
}

func NetcastfForm(args ...interface{}) *fForm {

	s := fForm{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructfForm() *fForm { return &fForm{} }

type aForm struct {
	Target, Method, Token, ButtonClass, CTA string
	Input                                   []afInput
	ModelName                               string
}

func NetcastaForm(args ...interface{}) *aForm {

	s := aForm{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructaForm() *aForm { return &aForm{} }

type Job struct {
	Title       string            `title:"Title of post",valid:"unique",placeholder:"Testfield prompt"`
	Location    string            `title:"Address of job",placeholder:"123 Smith St."`
	Author      netform.Email     `title:"Contact e-mail address",valid:"email,required",placeholder:"@",`
	Description netform.Paragraph `title:"Job description :",valid:"unique"`
	Time        time.Time
}

func NetcastJob(args ...interface{}) *Job {

	s := Job{}
	mapp := args[0].(db.O)
	if _, ok := mapp["_id"]; ok {
		mapp["Id"] = mapp["_id"]
	}
	data, _ := json.Marshal(&mapp)

	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err.Error())
	}

	return &s
}
func NetstructJob() *Job { return &Job{} }
func NetLoadWebAsset(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `data,err := Asset( fmt.Sprintf("web%s", args[0].(string) ) )`
	data, err := Asset(fmt.Sprintf("web%s", args[0].(string)))
	lastLine = `if err != nil {`
	if err != nil {
		lastLine = `return err.Error()`
		return err.Error()
		lastLine = `}`
	}
	lastLine = `return string(data)`
	return string(data)
}
func NetaC(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `return "}}"`
	return "}}"
}
func NetaO(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `return "{{"`
	return "{{"
}
func NetIsIn(args ...interface{}) bool {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `return strings.Contains(args[0].(string), args[1].(string))`
	return strings.Contains(args[0].(string), args[1].(string))
}
func NetHasBody(args ...interface{}) bool {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `return args[0].(*http.Request).ContentLength > 0`
	return args[0].(*http.Request).ContentLength > 0
}
func NetForm(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `err := netform.Form(args[0].(*http.Request) , args[1])`
	err := netform.Form(args[0].(*http.Request), args[1])
	lastLine = `return err.Error()`
	return err.Error()
}
func NetTokenizeForm(args ...interface{}) (form fForm) {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `v := reflect.ValueOf(args[0]).Elem()`
	v := reflect.ValueOf(args[0]).Elem()
	lastLine = `//t := reflect.TypeOf(item)`
	//t := reflect.TypeOf(item)
	lastLine = `bso :=  netform.ToBson( mResponse(args[0]) )`
	bso := netform.ToBson(mResponse(args[0]))
	lastLine = `for i := 0; i < v.NumField(); i++ {`
	for i := 0; i < v.NumField(); i++ {
		lastLine = `field := v.Type().Field(i)`
		field := v.Type().Field(i)
		lastLine = `fieldtype := strings.ToLower(field.Type.String())`
		fieldtype := strings.ToLower(field.Type.String())
		lastLine = `requird := strings.Contains(string(field.Tag), "required")`
		requird := strings.Contains(string(field.Tag), "required")
		lastLine = `opts := netform.GetSel(string(field.Tag))`
		opts := netform.GetSel(string(field.Tag))
		lastLine = `title := field.Tag.Get("title")`
		title := field.Tag.Get("title")
		lastLine = `placehlder := netform.GetPl(string(field.Tag))`
		placehlder := netform.GetPl(string(field.Tag))
		lastLine = `if strings.Contains(fieldtype, "bool"){`
		if strings.Contains(fieldtype, "bool") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "checkbox",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "checkbox", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `} else if strings.Contains(fieldtype, "int") || strings.Contains(fieldtype, "float") {`
		} else if strings.Contains(fieldtype, "int") || strings.Contains(fieldtype, "float") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "number",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "number", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `}  else if strings.Contains(fieldtype, "string"){`
		} else if strings.Contains(fieldtype, "string") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "text",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "text", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `} else if strings.Contains(fieldtype, "email"){`
		} else if strings.Contains(fieldtype, "email") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "email",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "email", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `} else if strings.Contains(fieldtype, "password"){`
		} else if strings.Contains(fieldtype, "password") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "password",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "password", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `}  else if strings.Contains(fieldtype, "select"){`
		} else if strings.Contains(fieldtype, "select") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "select",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), strings.Split(opts,",") , requird })`
			form.Input = append(form.Input, fInput{field.Name, "select", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird})
			lastLine = `}  else if strings.Contains(fieldtype, "radio"){`
		} else if strings.Contains(fieldtype, "radio") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "radio",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), strings.Split(opts,",") , requird })`
			form.Input = append(form.Input, fInput{field.Name, "radio", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird})
			lastLine = `} else if strings.Contains(fieldtype, "selectmult"){`
		} else if strings.Contains(fieldtype, "selectmult") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "selectmult",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), strings.Split(opts,",") , requird })`
			form.Input = append(form.Input, fInput{field.Name, "selectmult", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird})
			lastLine = `}  else if strings.Contains(fieldtype, "date"){`
		} else if strings.Contains(fieldtype, "date") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "date",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name] ), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "date", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `}  else if strings.Contains(fieldtype, "file"){`
		} else if strings.Contains(fieldtype, "file") {
			lastLine = `placehlder = field.Tag.Get("file")`
			placehlder = field.Tag.Get("file")
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "file",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name] ), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "file", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `}  else if strings.Contains(fieldtype, "paragraph"){`
		} else if strings.Contains(fieldtype, "paragraph") {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "textarea",placehlder,title,netform.InputClass,fmt.Sprintf("%v", bso[field.Name] ), nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "textarea", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird})
			lastLine = `} else {`
		} else {
			lastLine = `form.Input = append(form.Input, fInput{field.Name, "invalid",placehlder,title,netform.InputClass,"", nil , requird })`
			form.Input = append(form.Input, fInput{field.Name, "invalid", placehlder, title, netform.InputClass, "", nil, requird})
			lastLine = `}`
		}
		lastLine = `}`
	}
	lastLine = `return`
	return
}
func NetTokenizeFormAng(args ...interface{}) (form aForm) {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `v := reflect.ValueOf(args[0]).Elem()`
	v := reflect.ValueOf(args[0]).Elem()
	lastLine = `//t := reflect.TypeOf(item)`
	//t := reflect.TypeOf(item)
	lastLine = `modelClass := args[1].(string)`
	modelClass := args[1].(string)
	lastLine = `bso :=  netform.ToBson( mResponse(args[0]) )`
	bso := netform.ToBson(mResponse(args[0]))
	lastLine = `for i := 0; i < v.NumField(); i++ {`
	for i := 0; i < v.NumField(); i++ {
		lastLine = `field := v.Type().Field(i)`
		field := v.Type().Field(i)
		lastLine = `fieldtype := strings.ToLower(field.Type.String())`
		fieldtype := strings.ToLower(field.Type.String())
		lastLine = `requird := strings.Contains(string(field.Tag), "required")`
		requird := strings.Contains(string(field.Tag), "required")
		lastLine = `opts := netform.GetSel(string(field.Tag))`
		opts := netform.GetSel(string(field.Tag))
		lastLine = `title := field.Tag.Get("title")`
		title := field.Tag.Get("title")
		lastLine = `placehlder := netform.GetPl(string(field.Tag))`
		placehlder := netform.GetPl(string(field.Tag))
		lastLine = `if strings.Contains(fieldtype, "bool"){`
		if strings.Contains(fieldtype, "bool") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "checkbox",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird , modelClass})`
			form.Input = append(form.Input, afInput{field.Name, "checkbox", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `} else if strings.Contains(fieldtype, "int") || strings.Contains(fieldtype, "float") {`
		} else if strings.Contains(fieldtype, "int") || strings.Contains(fieldtype, "float") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "number",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "number", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `}  else if strings.Contains(fieldtype, "string"){`
		} else if strings.Contains(fieldtype, "string") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "text",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird , modelClass})`
			form.Input = append(form.Input, afInput{field.Name, "text", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `} else if strings.Contains(fieldtype, "email"){`
		} else if strings.Contains(fieldtype, "email") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "email",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "email", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `} else if strings.Contains(fieldtype, "password"){`
		} else if strings.Contains(fieldtype, "password") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "password",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "password", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `}  else if strings.Contains(fieldtype, "select"){`
		} else if strings.Contains(fieldtype, "select") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "select",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), strings.Split(opts,",") , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "select", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird, modelClass})
			lastLine = `}  else if strings.Contains(fieldtype, "radio"){`
		} else if strings.Contains(fieldtype, "radio") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "radio",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), strings.Split(opts,",") , requird , modelClass})`
			form.Input = append(form.Input, afInput{field.Name, "radio", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird, modelClass})
			lastLine = `} else if strings.Contains(fieldtype, "selectmult"){`
		} else if strings.Contains(fieldtype, "selectmult") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "selectmult",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name]), strings.Split(opts,",") , requird , modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "selectmult", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), strings.Split(opts, ","), requird, modelClass})
			lastLine = `}  else if strings.Contains(fieldtype, "date"){`
		} else if strings.Contains(fieldtype, "date") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "date",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name] ), nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "date", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `}  else if strings.Contains(fieldtype, "file"){`
		} else if strings.Contains(fieldtype, "file") {
			lastLine = `placehlder = field.Tag.Get("file")`
			placehlder = field.Tag.Get("file")
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "file",placehlder,title,netform.InputClass,fmt.Sprintf("%v",bso[field.Name] ), nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "file", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `}  else if strings.Contains(fieldtype, "paragraph"){`
		} else if strings.Contains(fieldtype, "paragraph") {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "textarea",placehlder,title,netform.InputClass,fmt.Sprintf("%v", bso[field.Name] ), nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "textarea", placehlder, title, netform.InputClass, fmt.Sprintf("%v", bso[field.Name]), nil, requird, modelClass})
			lastLine = `} else {`
		} else {
			lastLine = `form.Input = append(form.Input, afInput{field.Name, "invalid",placehlder,title,netform.InputClass,"", nil , requird, modelClass })`
			form.Input = append(form.Input, afInput{field.Name, "invalid", placehlder, title, netform.InputClass, "", nil, requird, modelClass})
			lastLine = `}`
		}
		lastLine = `}`
	}
	lastLine = `return`
	return
}
func NetBuild(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `if len(args) < 5 {`
	if len(args) < 5 {
		lastLine = `return "<h1>No enough arguments to build form.</h1>"`
		return "<h1>No enough arguments to build form.</h1>"
		lastLine = `}`
	}
	lastLine = `form := NetTokenizeForm(args[0])`
	form := NetTokenizeForm(args[0])
	lastLine = `target := args[1].(string)`
	target := args[1].(string)
	lastLine = `session := args[4].(*sessions.Session)`
	session := args[4].(*sessions.Session)
	lastLine = `form.Token = netform.GenerateToken(target, session.Values["formtoken"].(string) )`
	form.Token = netform.GenerateToken(target, session.Values["formtoken"].(string))
	lastLine = `form.Method = args[2].(string)`
	form.Method = args[2].(string)
	lastLine = `form.Target = target`
	form.Target = target
	lastLine = `form.ButtonClass = netform.ButtonClass`
	form.ButtonClass = netform.ButtonClass
	lastLine = `form.CTA = args[3].(string)`
	form.CTA = args[3].(string)
	lastLine = `return btForm(form)`
	return btForm(form)
}
func NetAngularForm(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `if len(args) < 6 {`
	if len(args) < 6 {
		lastLine = `return "<h1>No enough arguments to build form.</h1>"`
		return "<h1>No enough arguments to build form.</h1>"
		lastLine = `}`
	}
	lastLine = `modelclass := args[4].(string)`
	modelclass := args[4].(string)
	lastLine = `form := NetTokenizeFormAng(args[0], modelclass)`
	form := NetTokenizeFormAng(args[0], modelclass)
	lastLine = `target := args[1].(string)`
	target := args[1].(string)
	lastLine = `session := args[5].(*sessions.Session)`
	session := args[5].(*sessions.Session)
	lastLine = `form.ModelName = modelclass`
	form.ModelName = modelclass
	lastLine = `form.Token = netform.GenerateToken(target, session.Values["formtoken"].(string) )`
	form.Token = netform.GenerateToken(target, session.Values["formtoken"].(string))
	lastLine = `form.Method = args[2].(string)`
	form.Method = args[2].(string)
	lastLine = `form.Target = target`
	form.Target = target
	lastLine = `form.ButtonClass = netform.ButtonClass`
	form.ButtonClass = netform.ButtonClass
	lastLine = `form.CTA = args[3].(string)`
	form.CTA = args[3].(string)
	lastLine = `return batForm(form)`
	return batForm(form)
}
func NetGenerateToken(args ...interface{}) string {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `session := args[1].(*sessions.Session)`
	session := args[1].(*sessions.Session)
	lastLine = `return netform.GenerateToken(args[0].(string), session.Values["formtoken"].(string) )`
	return netform.GenerateToken(args[0].(string), session.Values["formtoken"].(string))
}
func NetJobs() (returnjobs []Job) {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `returnjobs = jobs`
	returnjobs = jobs
	lastLine = `return`
	return
}
func NetAddJob(job Job) (returnjobs []Job) {

	lastLine := ""

	defer func() {
		if n := recover(); n != nil {
			log.Println("Pipeline failed at line :", GetLine(".//gos.gxml", lastLine), "Of file:.//gos.gxml:", strings.TrimSpace(lastLine))
			log.Println("Reason : ", n)

		}
	}()
	lastLine = `job.Time = time.Now()`
	job.Time = time.Now()
	lastLine = `jobs = append([]Job{job}, jobs...)`
	jobs = append([]Job{job}, jobs...)
	lastLine = `//New list of jobs`
	//New list of jobs
	lastLine = `returnjobs = jobs`
	returnjobs = jobs
	lastLine = `return`
	return
}

func NettInput(args ...interface{}) string {

	var d fInput
	filename := "tmpl/form/input.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tInput) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = fInput{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func btInput(d fInput) string {
	return NetbtInput(d)
}

func NetbtInput(d fInput) string {

	filename := "tmpl/form/input.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tInput) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetctInput(args ...interface{}) (d fInput) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = fInput{}
	}
	return
}

func ctInput(args ...interface{}) (d fInput) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = fInput{}
	}
	return
}

func BtInput(intstr interface{}) string {
	return NettInput(intstr)
}

func NetatInput(args ...interface{}) string {

	var d afInput
	filename := "tmpl/form/ainput.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atInput) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = afInput{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func batInput(d afInput) string {
	return NetbatInput(d)
}

func NetbatInput(d afInput) string {

	filename := "tmpl/form/ainput.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atInput")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atInput) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetcatInput(args ...interface{}) (d afInput) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = afInput{}
	}
	return
}

func catInput(args ...interface{}) (d afInput) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = afInput{}
	}
	return
}

func BatInput(intstr interface{}) string {
	return NetatInput(intstr)
}

func NettForm(args ...interface{}) string {

	var d fForm
	filename := "tmpl/form/form.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tForm) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = fForm{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func btForm(d fForm) string {
	return NetbtForm(d)
}

func NetbtForm(d fForm) string {

	filename := "tmpl/form/form.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("tForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (tForm) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetctForm(args ...interface{}) (d fForm) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = fForm{}
	}
	return
}

func ctForm(args ...interface{}) (d fForm) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = fForm{}
	}
	return
}

func BtForm(intstr interface{}) string {
	return NettForm(intstr)
}

func NetatForm(args ...interface{}) string {

	var d aForm
	filename := "tmpl/form/aform.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atForm) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = aForm{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func batForm(d aForm) string {
	return NetbatForm(d)
}

func NetbatForm(d aForm) string {

	filename := "tmpl/form/aform.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("atForm")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (atForm) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func NetcatForm(args ...interface{}) (d aForm) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = aForm{}
	}
	return
}

func catForm(args ...interface{}) (d aForm) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = aForm{}
	}
	return
}

func BatForm(intstr interface{}) string {
	return NetatForm(intstr)
}

func Netang(args ...interface{}) string {

	var d NoStruct
	filename := "tmpl/momentum/ang.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (ang) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = NoStruct{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("ang")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func bang(d NoStruct) string {
	return Netbang(d)
}

func Netbang(d NoStruct) string {

	filename := "tmpl/momentum/ang.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("ang")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (ang) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func Netcang(args ...interface{}) (d NoStruct) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = NoStruct{}
	}
	return
}

func cang(args ...interface{}) (d NoStruct) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = NoStruct{}
	}
	return
}

func Bang(intstr interface{}) string {
	return Netang(intstr)
}

func Netserver(args ...interface{}) string {

	var d NoStruct
	filename := "tmpl/momentum/server.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (server) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = NoStruct{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("server")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func bserver(d NoStruct) string {
	return Netbserver(d)
}

func Netbserver(d NoStruct) string {

	filename := "tmpl/momentum/server.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("server")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (server) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func Netcserver(args ...interface{}) (d NoStruct) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = NoStruct{}
	}
	return
}

func cserver(args ...interface{}) (d NoStruct) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = NoStruct{}
	}
	return
}

func Bserver(intstr interface{}) string {
	return Netserver(intstr)
}

func Netjquery(args ...interface{}) string {

	var d NoStruct
	filename := "tmpl/momentum/jquery.tmpl"
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (jquery) : %s", filename))
			// log.Println(n)
			DebugTemplatePath(filename, &d)
			//http.Redirect(w,r,"/your-500-page",307)
		}
	}()
	if len(args) > 0 {
		jso := args[0].(string)
		var jsonBlob = []byte(jso)
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			return err.Error()
		}
	} else {
		d = NoStruct{}
	}

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("jquery")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))

	erro := t.Execute(output, &d)
	if erro != nil {
		color.Red(fmt.Sprintf("Error processing template %s", filename))
		DebugTemplatePath(filename, &d)
	}
	return html.UnescapeString(output.String())

}
func bjquery(d NoStruct) string {
	return Netbjquery(d)
}

func Netbjquery(d NoStruct) string {

	filename := "tmpl/momentum/jquery.tmpl"

	body, er := Asset(filename)
	if er != nil {
		return ""
	}
	output := new(bytes.Buffer)
	t := template.New("jquery")
	t = t.Funcs(template.FuncMap{"a": Netadd, "s": Netsubs, "m": Netmultiply, "d": Netdivided, "js": Netimportjs, "css": Netimportcss, "sd": NetsessionDelete, "sr": NetsessionRemove, "sc": NetsessionKey, "ss": NetsessionSet, "sso": NetsessionSetInt, "sgo": NetsessionGetInt, "sg": NetsessionGet, "form": formval, "eq": equalz, "neq": nequalz, "lte": netlt, "LoadWebAsset": NetLoadWebAsset, "aC": NetaC, "aO": NetaO, "IsIn": NetIsIn, "HasBody": NetHasBody, "Form": NetForm, "TokenizeForm": NetTokenizeForm, "TokenizeFormAng": NetTokenizeFormAng, "Build": NetBuild, "AngularForm": NetAngularForm, "GenerateToken": NetGenerateToken, "Jobs": NetJobs, "AddJob": NetAddJob, "tInput": NettInput, "btInput": NetbtInput, "ctInput": NetctInput, "atInput": NetatInput, "batInput": NetbatInput, "catInput": NetcatInput, "tForm": NettForm, "btForm": NetbtForm, "ctForm": NetctForm, "atForm": NetatForm, "batForm": NetbatForm, "catForm": NetcatForm, "ang": Netang, "bang": Netbang, "cang": Netcang, "server": Netserver, "bserver": Netbserver, "cserver": Netcserver, "jquery": Netjquery, "bjquery": Netbjquery, "cjquery": Netcjquery, "SampleForm": NetstructSampleForm, "isSampleForm": NetcastSampleForm, "fInput": NetstructfInput, "isfInput": NetcastfInput, "afInput": NetstructafInput, "isafInput": NetcastafInput, "fForm": NetstructfForm, "isfForm": NetcastfForm, "aForm": NetstructaForm, "isaForm": NetcastaForm, "Job": NetstructJob, "isJob": NetcastJob})
	t, _ = t.Parse(ReadyTemplate(body))
	defer func() {
		if n := recover(); n != nil {
			color.Red(fmt.Sprintf("Error loading template in path (jquery) : %s", filename))
			DebugTemplatePath(filename, &d)
		}
	}()
	erro := t.Execute(output, &d)
	if erro != nil {
		log.Println(erro)
	}
	return html.UnescapeString(output.String())
}
func Netcjquery(args ...interface{}) (d NoStruct) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = NoStruct{}
	}
	return
}

func cjquery(args ...interface{}) (d NoStruct) {
	if len(args) > 0 {
		var jsonBlob = []byte(args[0].(string))
		err := json.Unmarshal(jsonBlob, &d)
		if err != nil {
			log.Println("error:", err)
			return
		}
	} else {
		d = NoStruct{}
	}
	return
}

func Bjquery(intstr interface{}) string {
	return Netjquery(intstr)
}

func dummy_timer() {
	dg := time.Second * 5
	log.Println(dg)
}
func main() {
	fmt.Fprintf(os.Stdout, "%v\n", os.Getpid())

	//psss go code here : func main()

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0700)
	}

	//psss go code here : func main()
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashPort := 8700
	appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
	}
	color.Red("✅ Important!")
	log.Println("To see your traces, go to ", appdashURL)

	// Start the web UI in a separate goroutine.
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatal(err)
	}
	tapp.Store = store
	tapp.Queryer = store
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
	}()

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(fmt.Sprintf(":%d", collectorPort)))
	opentracing.InitGlobalTracer(tracer)

	port := ":8080"
	if envport := os.ExpandEnv("$PORT"); envport != "" {
		port = fmt.Sprintf(":%s", envport)
	}
	log.Printf("Listenning on Port %v\n", port)
	http.HandleFunc("/", MakeHandler(Handler))

	http.HandleFunc("/momentum/templates", func(w http.ResponseWriter, r *http.Request) {

		if r.FormValue("name") == "reset" {
			return
		} else if r.FormValue("name") == "tInput" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := NettInput(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))
		} else if r.FormValue("name") == "atInput" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := NetatInput(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))
		} else if r.FormValue("name") == "tForm" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := NettForm(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))
		} else if r.FormValue("name") == "atForm" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := NetatForm(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))
		} else if r.FormValue("name") == "ang" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := Netang(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))
		} else if r.FormValue("name") == "server" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := Netserver(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))
		} else if r.FormValue("name") == "jquery" {
			w.Header().Set("Content-Type", "text/html")
			tmplRendered := Netjquery(r.FormValue("payload"))
			w.Write([]byte(tmplRendered))

		}
	})

	http.HandleFunc("/funcfactory.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(`function tInput(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "tInput", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function atInput(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "atInput", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function tForm(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "tForm", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function atForm(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "atForm", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function ang(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "ang", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function server(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "server", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function jquery(dataOfInterface, cb){ jsrequestmomentum("/momentum/templates", {name: "jquery", payload: JSON.stringify(dataOfInterface)},"POST",  cb) }
function Jobs(  cb){
	var t = {}
	
	jsrequestmomentum("/momentum/funcs?name=Jobs", t, "POSTJSON", cb)
}
function AddJob(Job , cb){
	var t = {}
	
	t.Job = Job
	jsrequestmomentum("/momentum/funcs?name=AddJob", t, "POSTJSON", cb)
}
`))
	})

	http.HandleFunc("/momentum/funcs", func(w http.ResponseWriter, r *http.Request) {

		if r.FormValue("name") == "reset" {
			return
		} else if r.FormValue("name") == "Jobs" {
			w.Header().Set("Content-Type", "application/json")
			type PayloadJobs struct {
			}
			decoder := json.NewDecoder(r.Body)
			var tmvv PayloadJobs
			err := decoder.Decode(&tmvv)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", err.Error())))
				return
			}
			resp := db.O{}
			respreturnjobs0 := NetJobs()

			resp["returnjobs"] = respreturnjobs0
			w.Write([]byte(mResponse(resp)))
		} else if r.FormValue("name") == "AddJob" {
			w.Header().Set("Content-Type", "application/json")
			type PayloadAddJob struct {
				Job Job
			}
			decoder := json.NewDecoder(r.Body)
			var tmvv PayloadAddJob
			err := decoder.Decode(&tmvv)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", err.Error())))
				return
			}
			resp := db.O{}
			respreturnjobs0 := NetAddJob(tmvv.Job)

			resp["returnjobs"] = respreturnjobs0
			w.Write([]byte(mResponse(resp)))

		}
	})
	//+++extendgxmlmain+++
	http.Handle("/dist/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "web"}))

	errgos := http.ListenAndServe(port, nil)
	if errgos != nil {
		log.Fatal(errgos)
	}

}

//+++extendgxmlroot+++
