package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackdanger/collectlinks"
)

const urlConsulta = "http://transparencia.uepb.edu.br/wp-content/themes/uepb2016-transparencia-final/controller/consulta.php"

var mes = flag.String("mes", "01/2016", "mm/yyyy de referência.")
var maxWorkers = flag.Int("max_workers", 20, "Número máximo de workers.")
var sep = flag.String("sep", ",", "CSV separator.")

func main() {
	startTime := time.Now()

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

	links := make(chan string, *maxWorkers)
	results := make(chan string, *maxWorkers)
	var wg sync.WaitGroup
	for i := 0; i < *maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doWork(links, results)
		}()
	}

	// Work is done, close the result channel.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fill up work queue.
	go func() {
		for _, link := range collectlinks.All(resp.Body) {
			links <- link
		}
		close(links)
	}()

	// Print results.
	for row := range results {
		fmt.Println(row)
	}

	fmt.Fprintf(os.Stderr, "\nFinished! Duration: %s\n", time.Now().Sub(startTime))
}

const maxRetries = 3

func doWork(links <-chan string, results chan<- string) {
	for link := range links {
		var doc *goquery.Document
		for i := 1; ; i++ {
			var err error
			doc, err = goquery.NewDocument(link)
			if err == nil {
				break
			}
			fmt.Fprintf(os.Stderr, "[Tentativa %d] Erro tentando processar página de servidor: %s. Erro: %q", i, link, err)
			if i == maxRetries {
				fmt.Fprintf(os.Stderr, "Página não processada: %s", link)
				return
			}
			time.Sleep(time.Duration(i) * time.Duration(rand.Intn(5)) * time.Second)
		}
		var row []string
		doc.Find("td.desc").Each(func(i int, s *goquery.Selection) {
			cell := strings.Replace(
				strings.Trim(s.Next().Text(), " \n"),
				",",
				".",
				1)
			row = append(row, cell)
		})
		if len(row) > 0 {
			results <- strings.Join(row, *sep)
		} else {
			fmt.Fprintf(os.Stderr, "Não achou td.desc: %s\n", link)
		}
	}
}
