package main

import (
	"errors"
	"fmt"
	"github.com/DanielFillol/goSpider"
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	users := []goSpider.Requests{
		{SearchString: "1017927-35.2023.8.26.0008"},
		{SearchString: "0002396-75.2013.8.26.0201"},
		{SearchString: "1551285-50.2021.8.26.0477"},
		{SearchString: "0015386-82.2013.8.26.0562"},
		{SearchString: "0007324-95.2015.8.26.0590"},
		{SearchString: "1545639-85.2023.8.26.0090"},
		{SearchString: "1557599-09.2021.8.26.0090"},
		{SearchString: "1045142-72.2021.8.26.0002"},
		{SearchString: "0208591-43.2009.8.26.0004"},
		{SearchString: "1024511-70.2022.8.26.0003"},
	}

	numberOfWorkers := 1
	duration := 0 * time.Millisecond

	results, err := goSpider.ParallelRequests(users, numberOfWorkers, duration, Crawler)
	if err != nil {
		log.Println("Expected %d results, but got %d, List results: %v", len(users), 0, len(results))
	}

	log.Println("Finish Parallel Requests!")

	type Lawsuit struct {
		Cover     Cover
		Persons   []Person
		Movements []Movement
	}
	var lawsuits []Lawsuit
	for _, result := range results {
		// Cover
		c, err := extractDataCover(result.Page, "//*[@id=\"numeroProcesso\"]", "//*[@id=\"labelSituacaoProcesso\"]", "//*[@id=\"classeProcesso\"]", "//*[@id=\"assuntoProcesso\"]", "//*[@id=\"foroProcesso\"]", "//*[@id=\"varaProcesso\"]", "//*[@id=\"juizProcesso\"]", "//*[@id=\"dataHoraDistribuicaoProcesso\"]", "//*[@id=\"numeroControleProcesso\"]", "//*[@id=\"areaProcesso\"]/span", "//*[@id=\"valorAcaoProcesso\"]")
		if err != nil {
			log.Printf("ExtractDataCover error: %v", err)
		}
		// Persons
		p, err := extractDataPerson(result.Page, "//*[@id=\"tableTodasPartes\"]/tbody/tr", "td[1]/span", "td[2]/text()", "\n")
		if err != nil {
			p, err = extractDataPerson(result.Page, "//*[@id=\"tablePartesPrincipais\"]/tbody/tr", "td[1]/text()", "td[2]/text()", "\n")
			if err != nil {
				log.Printf("Expected some person but got none: %v", err.Error())
			}
		}
		// Movements
		m, err := extractDataMovement(result.Page, "//*[@id=\"tabelaTodasMovimentacoes\"]/tr", "\n")
		if err != nil {
			log.Printf("Expected some movement but got none: %v", err.Error())
		}

		lawsuits = append(lawsuits, Lawsuit{
			Cover:     c,
			Persons:   p,
			Movements: m,
		})
	}

	if len(lawsuits) != len(users) {
		log.Printf("Expected %d lawsuits, but got %d", len(users), len(lawsuits))
	}

	fmt.Println(lawsuits)
}

func Crawler(d string) (*html.Node, error) {
	url := "https://esaj.tjsp.jus.br/cpopg/open.do"
	nav := goSpider.NewNavigator()

	err := nav.OpenURL(url)
	if err != nil {
		log.Printf("OpenURL error: %v", err)
		return nil, err
	}

	err = nav.CheckRadioButton("#interna_NUMPROC > div > fieldset > label:nth-child(5)")
	if err != nil {
		log.Printf("CheckRadioButton error: %v", err)
		return nil, err
	}

	err = nav.FillField("#nuProcessoAntigoFormatado", d)
	if err != nil {
		log.Printf("filling field error: %v", err)
		return nil, err
	}

	err = nav.ClickButton("#botaoConsultarProcessos")
	if err != nil {
		log.Printf("ClickButton error: %v", err)
		return nil, err
	}

	err = nav.WaitForElement("#tabelaUltimasMovimentacoes > tr:nth-child(1) > td.dataMovimentacao", 15*time.Second)
	if err != nil {
		log.Printf("WaitForElement error: %v", err)
		return nil, err
	}

	pageSource, err := nav.GetPageSource()
	if err != nil {
		log.Printf("GetPageSource error: %v", err)
		return nil, err
	}

	return pageSource, nil
}

type Cover struct {
	Title       string
	Tag         string
	Class       string
	Subject     string
	Location    string
	Unit        string
	Judge       string
	InitialDate string
	Control     string
	Field       string
	Value       string
	Error       string
}

func extractDataCover(pageSource *html.Node, xpathTitle string, xpathTag string, xpathClass string, xpathSubject string, xpathLocation string, xpathUnit string, xpathJudge string, xpathInitDate string, xpathControl string, xpathField string, xpathValue string) (Cover, error) {
	var i int //count errors
	title, err := goSpider.ExtractText(pageSource, xpathTitle, "                                                            ")
	if err != nil {
		log.Println("error extracting title")
	}

	tag, err := goSpider.ExtractText(pageSource, xpathTag, "")
	if err != nil {
		i++
		log.Println("error extracting tag")
	}

	class, err := goSpider.ExtractText(pageSource, xpathClass, "")
	if err != nil {
		i++
		log.Println("error extracting class")
	}

	subject, err := goSpider.ExtractText(pageSource, xpathSubject, "")
	if err != nil {
		i++
		log.Println("error extracting subject")
	}

	location, err := goSpider.ExtractText(pageSource, xpathLocation, "")
	if err != nil {
		i++
		log.Println("error extracting location")
	}

	unit, err := goSpider.ExtractText(pageSource, xpathUnit, "")
	if err != nil {
		i++
		log.Println("error extracting unit")
	}

	judge, err := goSpider.ExtractText(pageSource, xpathJudge, "")
	if err != nil {
		i++
		log.Println("error extracting existJudge")
	}

	initDate, err := goSpider.ExtractText(pageSource, xpathInitDate, "")
	if err != nil {
		i++
		log.Println("error extracting initDate")
	}

	control, err := goSpider.ExtractText(pageSource, xpathControl, "")
	if err != nil {
		i++
		log.Println("error extracting control")
	}

	field, err := goSpider.ExtractText(pageSource, xpathField, "")
	if err != nil {
		log.Println("error extracting field")
	}

	value, err := goSpider.ExtractText(pageSource, xpathValue, "R$         ")
	if err != nil {
		i++
		log.Println("error extracting field value")
	}

	var e string
	if err != nil {
		e = err.Error()
	}

	if i >= 5 {
		return Cover{}, fmt.Errorf("too many errors: %d", i)
	}

	return Cover{
		Title:       title,
		Tag:         tag,
		Class:       class,
		Subject:     subject,
		Location:    location,
		Unit:        unit,
		Judge:       judge,
		InitialDate: initDate,
		Control:     control,
		Field:       field,
		Value:       value,
		Error:       e,
	}, nil
}

type Person struct {
	Pole    string
	Name    string
	Lawyers []string
}

func extractDataPerson(pageSource *html.Node, xpathPeople string, xpathPole string, xpathLawyer string, dirt string) ([]Person, error) {
	Pole, err := goSpider.FindNodes(pageSource, xpathPeople)
	if err != nil {
		return nil, err
	}

	var personas []Person
	for i, person := range Pole {
		pole, err := goSpider.ExtractText(person, xpathPole, dirt)
		if err != nil {
			return nil, errors.New("error extract data person, pole not found: " + err.Error())
		}

		var name string
		_, err = goSpider.FindNodes(person, xpathPeople+"["+strconv.Itoa(i)+"]/td[2]")
		if err != nil {
			name, err = goSpider.ExtractText(person, "td[2]/text()", dirt)
			if err != nil {
				return nil, errors.New("error extract data person, name not found: " + err.Error())
			}
		} else {
			name, err = goSpider.ExtractText(person, "td[2]/text()["+strconv.Itoa(1)+"]", dirt)
			if err != nil {
				return nil, errors.New("error extract data person, name not found: " + err.Error())
			}
		}

		var lawyers []string
		ll, err := goSpider.FindNodes(person, xpathLawyer)
		if err != nil {
			lawyers = append(lawyers, "no lawyer found")
		}
		for j, _ := range ll {
			n, err := goSpider.ExtractText(person, "td[2]/text()["+strconv.Itoa(j+1)+"]", dirt)
			if err != nil {
				return nil, errors.New("error extract data person, lawyer not  found: " + err.Error())
			}
			lawyers = append(lawyers, n)
		}

		p := Person{
			Pole:    pole,
			Name:    name,
			Lawyers: lawyers,
		}

		personas = append(personas, p)
	}

	return personas, nil
}

type Movement struct {
	Date  string
	Title string
	Text  string
}

func extractDataMovement(pageSource *html.Node, node string, dirt string) ([]Movement, error) {
	xpathTable := node

	tableRows, err := goSpider.ExtractTable(pageSource, xpathTable)
	if err != nil {
		return nil, err
	}

	if len(tableRows) > 0 {
		var allMovements []Movement
		for _, row := range tableRows {
			date, err := goSpider.ExtractText(row, "td[1]", dirt)
			if err != nil {
				return nil, errors.New("error extracting table date: " + err.Error())
			}
			title, err := goSpider.ExtractText(row, "td[3]", dirt)
			if err != nil {
				return nil, errors.New("error extracting table title: " + err.Error())
			}
			text, err := goSpider.ExtractText(row, "td[3]/span", dirt)
			if err != nil {
				return nil, errors.New("error extracting table text: " + err.Error())
			}

			mv := Movement{
				Date:  strings.ReplaceAll(date, "\t", ""),
				Title: strings.ReplaceAll(strings.ReplaceAll(title, text, ""), dirt, ""),
				Text:  strings.TrimSpace(strings.ReplaceAll(text, "\t", "")),
			}

			allMovements = append(allMovements, mv)
		}
		return allMovements, nil
	}

	return nil, errors.New("error table: could not find any movements")
}
