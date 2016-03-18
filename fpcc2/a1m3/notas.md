# Atividade 1 Milestone 3

Os dados podem ser baixados [aqui](https://drive.google.com/file/d/0B2rlaHwjOlZAZ2lidzUzRWVjeEE/view?usp=sharing).

Estes dados do Github Archive sobre atividade de repositórios no github entre o
início de 2012 e o início de 2015. Cada evento de criação de um repositório,
push, watch, criação de issue, etc. em cada repositório é registrado no archive.
Nos nossos dados, esses eventos foram agrupados por tipo do evento, linguagem de
programação, ano e trimestre. Assim, para cada observação nos dados, temos as
seguintes variáveis:
* year: ano da observação
* quarter: trimestre da observação
* type: tipo do evento (de atividade) observado em um conjunto de repositórios
* events: quantos eventos desse tipo foram registrados nesse trimestre
* active_repos_by_url: quantidade de repositórios onde os eventos aconteceram
* repository_language: linguagem dos repositórios

A linha 1 ("Ruby  ForkEvent               24540  73049 2014       3") diz que
houveram 73049 forks feitos em 24540 repositórios Ruby no 3o trimestre de 2014.

As atividades nesta 3a parte do lab são:

* Escolha e descreva 8 perguntas que na sua opinião são relevantes, não são
óbvias e que você gostaria de ver respondidas a partir dos dados. Para cada uma
escreva uma frase curta que documenta qual você acha que será a resposta.
(O resultado desta atividade é o checkpoint 3 do problema).
