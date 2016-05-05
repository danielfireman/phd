package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackdanger/collectlinks"
)

const urlConsulta = "http://transparencia.uepb.edu.br/wp-content/themes/uepb2016-transparencia-final/controller/consulta.php"

var mes = flag.String("mes", "01/2016", "mm/yyyy de referência.")
var maxWorkers = flag.Int("max_workers", 20, "Número máximo de workers.")

func main() {
	flag.Parse()

	if len(*mes) == 0 {
		log.Fatalf("Mês de referência inválido:%s. Formato esperado: mm/yyyy.", *mes)
	}

	param := fmt.Sprintf("parametros[]=%s&parametros[]=&parametros[]=&parametros[]=", *mes)
	resp, err := http.Post(
		urlConsulta,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(param))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	links := make(chan string)
	results := make(chan string)

	for i := 0; i < *maxWorkers; i++ {
		go worker(i, links, results)
	}

	go func() {
		for _, link := range collectlinks.All(resp.Body) {
			links <- link
		}
		close(links)
	}()

	for row := range results {
		fmt.Println(row)
	}

	close(results)
}

func worker(id int, links <-chan string, results chan<- string) {
	for link := range links {
		resp, err := http.Get(link)
		if err != nil {
			log.Printf("Erro tentando buscar informações: %s. Erro: %q", link, err)
			return
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocument(link)
		if err != nil {
			log.Printf("Erro tentando processar página de servidor: %s. Erro: %q", link, err)
			return
		}
		var rows []string
		doc.Find("td").Each(func(i int, s *goquery.Selection) {
			r := strings.Trim(s.Next().Text(), " \n")
			if len(r) > 0 {
				rows = append(rows, r)
			}
		})
		results <- strings.Join(rows, ";")
	}
}
