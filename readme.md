- verificar se os timeouts sao maior do que zero.
- add support to tsl and add tsl prefer (https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)
- falta corrigir os nomes 
- como suportar update? devemos? patch? update?
- unificar no infra o pagination
- um documento nunca eh deletado pq mandamos a trigger pro subscrition, mas efetivamente nao removemos do banco, como vamos resolver isso? soft delete? vaccum?
- mover os default value da inicializcao para o local quie o programa inicia, la eh mt melhor pq la fica centralizado.
- mover o ensureindex do mongodb para o setup!? acho que vale a pena...
- esta faltando metricas, muitas metricas... acho que temos que fazer isso antes da query.
- rodar o fluxo de tests em um docker de teste (tipo para testar versao beta do go e no tip)
- acho que o subscription nao precisa ter o partition, podemos abstrair isso do objeto.
- acho que temos que fzer tipo um reference count nos documents, tipo, toda vez que colocarmos na fila 
	o document, vamos incrementar um valor em 1, depois que processar, vamos tirar 1, dai quando um processo de
	limpeza rodar, e ver que o documento eiste a mais de 1 hr e o refrence count ta zero, significa que podemos
	deletar sem problema.
- precisamos ter um docker de build, docker de ci e docker para rodar.
- renomear os arquivos de teste json para 'estrutura.funcao.acao.qtd.json'
vamos botar um limit reader no body pra prevenir de alguem enviar um conteudo gitante de grande a gente ficar lendo!? acho valido sim...
https://golang.org/pkg/io/#LimitReader
- acho que podemos remover aqueles metodos de join/leave partition da interface, eles sao chamados dentro do msm pacote.... nao tem problema.
  na verdade nao sei, pq se nao a gente nao consegue testar....

- tirar aquele join/leave partition, fazer isso por interfaces internas e nao deixar na interface publica.
- renomear o findone para bindbyid e o findall para find
- fetchpartitions, vira partitions

- mover a criacao de filas e indices para o setup.

criar uma goroutine no logger com formato para o stdout e com formato humano de que, se nada for printado por 1 segundo, pular 1 linha.
se ja tiver pulado, nao pular mais ateh algo ser printado. (nice!!)

tem qu e fazer a mesma coisa do init do mongodb no repositorio da memoria...

temos que criar um conceito de transaction, pra qnd a gente usar um banco com transacao funcionar.
 +tem que ver como poderiamos fazer isso....
 +passar no contexto!? talvez faca mais sentido.

# <img src="misc/doc/logo.png" border="0" alt="flare" height="45">
<a href="https://travis-ci.org/diegobernardes/flare"><img src="https://img.shields.io/travis/diegobernardes/flare/master.svg?style=flat-square" alt="Build Status"></a>
<a href="https://coveralls.io/github/diegobernardes/flare"><img src="https://img.shields.io/coveralls/diegobernardes/flare/master.svg?style=flat-square" alt="Coveralls"></a>
<a href="https://godoc.org/github.com/diegobernardes/flare"><img src="https://img.shields.io/badge/api-reference-blue.svg?style=flat-square" alt="GoDoc"></a>

Flare is a service that listen to changes on HTTP endpoints and notify subscripted clients about the changes. It help reduce the pressure on APIs by avoiding the clients to do pooling requests to search for new/changed content and the need of the APIs to develop workers to notify the clients about document changes.

There is no need to the the service know anything about who is consuming it's updates, this is abstracted and lead to a simpler design on APIs. Problems like scaling the workers to notify the changes if the number of subscriptions increase, need to control the delivery success of the messages, include/update/delete the clients on your subscription list and so on are just solved.

## How to run
```bash
go get github.com/diegobernardes/flare/services/flare/cmd
cd github.com/diegobernardes/flare/services/flare/cmd
go run flare.go start
```

Currently the `memory` provider is broken, use `mongodb` and `aws` instead. Soon this will be solved.

## How it works
Flare has 3 basic entities: `Resource`, `Subscription` and `Document`. The origin of content is responsible for `Resource` and `Document` entities and the clients are responsible for `Subscription`.

### Resource
<p align="center">
	<img src="misc/doc/resource.jpg">
<p>

Resource represents the entity you want to track. It cannot be updated, only deleted, and to delete, first you need to remove all the associated subscriptions.


| Field  | Description |
| ------------- | ------------- |
| `addresses` | All the addresses a resource can have. Very useful when the same API is exposed to the internet and intranet with different hosts. |
| `path` | Is the actual document that gonna be tracked. `wildcards` are required to track the collection and they can be later used at subscriptions. |
| `change.field` | The field that is used to track changes on a document. It can be a string containing a date or a auto incremented integer. |
| `change.format` | If the field is a date, this fields has the format to parse the document date. More info about the format [here](https://golang.org/pkg/time/#pkg-constants). |


### Subscription
<p align="center">
	<img src="misc/doc/subscription.jpg">
<p>


Subscription is the responsible to notify the clients when a document from a resource changes.

| Field  | Description |
| ------------- | ------------- |
| `endpoint.url` | The address of the client that gonna receive the notification. |
| `endpoint.method` | The method used on the notification request. |
| `endpoint.headers` | A list of headers to sent within the request. |
| `delivery.success` | List of success status code. This is used to mark the notification as delivered for the respective client. |
| `delivery.discard` | List of status code to discard the notification. |
| `sendDocument` | Send the whole document body as it was received. Default value is false. |
| `skipEnvelope` | The trigger send a envelope with some useful informations. Default value is false. |
| `data` | Only sent if the `skipEnvelope` is false. Can be used to provide aditional information to the client that gonna receive the notification. It also can interpolate wildcards used at resource path definition. |

### Document
<p align="center">
	<img src="misc/doc/document-create.jpg">
<p>

To update a document, a `PUT` should be done at `http://flare/documents/{endpoint}`, where the `{endpoint}` is the real document endpoint and it should match the information inserted at the resource creation. The body should contain the document.
If the origin send the same document or older documents more then one time, the service don't gonna notify the clients again because it know the document version each client has. The notification only happens when is really needed.  

<p align="center">
	<img src="misc/doc/document-update.jpg">
<p>

The difference from previous request is that the email has changed from `@gmail.com` to `@outlook.com` and the `updatedAt` from `08:30` to `08:35`. After this, the client receive a notifications of update.

<p align="center">
	<img src="misc/doc/document-delete.jpg">
<p>

The delete should be sent with the delete method and no body.
