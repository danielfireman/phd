# Tarefas

[Link](https://docs.google.com/document/d/1d5MMw7AEDn28mY4ncQim5k3zYDUsXt565JglzoXrnuI/edit#heading=h.iix63l1day8w)

_Nota:  Assumindo que os dados e a explicação sobre os campos já estão locais_

---

Primeira coisa que fizemos:

* Curso de R: http://tryr.codeschool.com

## Resumo da tarefa

1. Criar um breve relatório que mostre um panorama dos dados considerados abaixo:
. Escolha um subconjunto das variáveis que você considera interessante para investigar como nossos deputados gastam sua verba
. comente qual a distribuição dos dados dessas variáveis, em termos de centralidade, extremos, concentração e simetria
. como você sugere que um analista de dados lide com a parte surpreendente/estranha?

1. Responda:
. Em que tipo de despesas nossos parlamentares gastam mais recursos de sua cota?
. Quais tipos de despesas têm despesas que mais variam, e que são mais desiguais?

## Primeira olhada na base de dados

* Campos interessantes
  * **txNomeParlamentar**: Nome adotado pelo Parlamentar ao tomar posse do seu mandato.
  * **ideCadastro**: Número que identifica unicamente um deputado federal na CD.
  * **sgUF**: No contexto da cota CEAP, representa a unidade da federação pela qual o deputado foi eleito e é utilizada para definir o valor da cota a que o deputado tem.
  * **sgPartido**:O seu conteúdo representa a sigla de um partido
  * **txtFornecedor**: O conteúdo deste dado representa o nome do fornecedor do produto ou serviço presente no documento fiscal
  * **txtCNPJCPF**:O conteúdo deste dado representa o CNPJ ou o CPF do emitente do documento fiscal, quando se tratar do uso da cota em razão do reembolso despesas comprovadas pela emissão de documentos fiscais.
  * **indTipoDocumento**:Este dado representa o tipo de documento do fiscal – 0 (Zero), para Nota Fiscal; 1 (um), para Recibo; e 2, para Despesa no Exterior.
  * **vlrDocumento**: O seu conteúdo é o valor de face do documento fiscal ou o valor
do documento que deu causa à despesa. Quando se tratar de
bilhete aéreo, esse valor poderá ser negativo, significando que o referido bilhete é um bilhete de compensação, pois compensa um outro bilhete emitido e não utilizado pelo deputado (idem para o dado vlrLiquido abaixo).
  * **vlrGlosa**: O seu conteúdo representa o valor da glosa do documento fiscal que incidirá sobre o Valor do Documento, ou o valor da glosa do documento que deu causa à despesa.
  * **vlrLiquido**: O seu conteúdo representa o valor líquido do documento fiscal ou do documento que deu causa à despesa e será calculado pela diferença entre o Valor do Documento e o Valor da Glosa. É este valor que será debitado da cota do deputado
  * **numMes**:O seu conteúdo representa o Mês da competência financeira do documento fiscal ou do documento que deu causa à despesa.
  * **numAno**: O seu conteúdo representa o Ano da competência financeira do documento fiscal ou do documento que deu causa à despesa
  * **txtPassageiro**:O conteúdo deste dado representa o nome do passageiro, quando o documento que deu causa à despesa se tratar de emissão de bilhete aéreo.
  * **txtTrecho**:O conteúdo deste dado representa o trecho da viagem, quando o documento que deu causa à despesa se tratar de emissão de bilhete aéreo.

## Perguntas interessantes
* Passagens aéreas:
  * Qual o deputado que mais viajou
  * Qual deputado comprou mais passagens para terceiros
  * Qual foi o trecho mais visitado pelos deputados
* Qual foi o fornecedor que mais prestou serviços a CD?
  * por partido
  * por deputado
* Houve reembolso de forma parcelada?

## Exemplos de aprendizagem

# Quais deputados realizaram compras no mês de janeiro?
distinct(select(filter(d, d$numMes == 1), txNomeParlamentar))

# Quantos os deputados viajaram
d %>% select(txtPassageiro) %>% distinct() %>% nrow()

# Quantas entradas não são de voos

d %>% filter(d$txtPassageiro == "-") %>% nrow()

# Passageiros ordenados por numero de voos

passageiros <- d %>% filter(d$txtPassageiro != "-") %>% count(txtPassageiro) %>% arrange(-n)

# Deputados que não colocaram os nomes dele como nomes do passageiro na compra de passagens

d %>% filter(d$txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar))

# Deputados que mais compraram passagens para terceiros

d %>% filter(d$txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar)) %>% count(txNomeParlamentar) %>% arrange(desc(n))

# Quem foram as pessoas que o deputado paulo freire comprou passagens

d  %>% filter(txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar), txNomeParlamentar == "PAULO FREIRE") %>% select(txNomeParlamentar, txtPassageiro) %>% group_by(txtPassageiro) %>% count(txtPassageiro) %>% arrange(desc(n))

# Quem foram as pessoas que o deputado marco feliciano comprou passagens

d  %>% filter(txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar), txNomeParlamentar == "PR. MARCO FELICIANO") %>% select(txNomeParlamentar, txtPassageiro) %>% group_by(txtPassageiro) %>% count(txtPassageiro) %>% arrange(desc(n))

# Quais partidos mais compraram passagens para terceiros

d %>% filter(d$txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar)) %>% count(sgPartido) %>% arrange(desc(n))

# Listar as bancadas ordenadas por maior tamanho

d %>% distinct(txNomeParlamentar) %>% count(sgPartido) %>% arrange(desc(n))

# Quanto cada partido gastou ordenado pelo maior gastou

d %>% group_by(sgPartido) %>% summarise(total=sum(vlrLiquido)) %>% arrange(desc(total))

# Quanto cada partido gastou com passagens para terceiros

d %>% filter(d$txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar)) %>% group_by(sgPartido) %>% summarise(total=sum(vlrLiquido)) %>% arrange(desc(total))

# Quanto cada deputado gastou em passagens para terceiros ordenados pelo maior valor
d %>% filter(d$txtPassageiro != "-", as.character(txtPassageiro) != as.character(txNomeParlamentar)) %>% group_by(txNomeParlamentar) %>% summarise(total=sum(vlrLiquido)) %>% arrange(desc(total))
