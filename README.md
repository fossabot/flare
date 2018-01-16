-----
quero ligar o worker que vai ser um genstage, ele vai ligar os workers de acordo com a demanda.
dai vou ligar os producers que vao ler das filas e colocar no worker.
todo esse processo vai ter backpressure.

-----
o sqs limita em ateh 10 mensagens, como buscar mais no elixir way!?
vamos ter o consumidor que pede sempre 1 mensagem, isso eh meio que nao tem solucao...

-----
vamos ter o passo de producer, rate limit, buffer, consumer.
precisamos do rate limit para nao estourar as pontas, futuramente, pode ser distribuido.
precisamos do buffer para caso o cliente peca mais mensagens do que conseguimos pgar, podemos pegar de varios lugares e entao entregar para ele.

-----
qnd um processo genevent, genserver, etc.. roda em um outro processo e qnd roda no msm process!?

-----
usar o genstage pra implementar back pressure.
vamos ter o producer que vai ser o cara que vai ler da aws. e vamos ter o consumer que vai ser o worker.
nesse caso o worker nao vai precisar do poolboy.

ao inves de ficar carregando esses valores estranhos no env, vou carregar no meu padrao e entao expor isso por config para a aplicacao.

tem que criar um provider para a aws pra poder entao comecar a fazer o get la e entao chamar 
os workers com os dados.
na verdade, vai chamar o pool que vai fazer a distribuicao

renomear de plugin para provider, acho que vai fazer mais sentido.

vamos criar tb uma interface web para poder gerenciar o cluster.
acho que html mesmo com turbolinks do rails.

resource, subscription, document, e application?
como fazer o rate limit!?

1 - o documento chega, coloco no banco e gero a mensagem de atualizacao.
2 - pego o resource do documento e gero uma mensagem para cada bucket.
3 - cada bucket gera a mensagem para o seu subscription.
4 - a mensagem do subscription envia 1-n mensagens para os clientes.
(end)


como o mongodb nao eh seguro (nao tem transaction), temos que ter uma tarefa de limpeza do banco que roda de x em x tempos.
vamos varrer os tipos procurando por atualizacoes, delecoes, etc...
dependendo do banco(mongodb), talvez faca sentido desligar o write externo pra manter a consistencia.

podemos colocar pequenos locks.
tipo, podemos dar lock no resource e com isso nada nele entra, e conseguimos fazer tudo e dps desbloquear ele.


colocar o ecto para fazer as validacoes.

https://github.com/joakimk/toniq
https://github.com/ejholmes/exsidekiq
https://github.com/edgurgel/verk
https://github.com/akira/exq
acho que a gente tem que ir mais pro caminho do toniq e ter um job runner que chama o modulo ao inves de fiar me injetando la na parada.

de onde vai vir a queue!?


ler esse blog: https://medium.com/learn-elixir

precisa fazer tipo um init pra poder ir inicializando todas as partes..


nao esquecer a ideia de votacao de achados por local.
plataforma para achar roupas, as  pessoas  podem tirar fotos, vincular instagrams, etc...
gerar mapa, votacao de achados e classificacao.
colocar um algorimo de decay para ir descendo automaticamente com as pecas.

# Flare

To start your Phoenix server:

  * Install dependencies with `mix deps.get`
  * Start Phoenix endpoint with `mix phx.server`

Now you can visit [`localhost:8080`](http://localhost:8080) from your browser.

Ready to run in production? Please [check our deployment guides](http://www.phoenixframework.org/docs/deployment).

## Learn more

  * Official website: http://www.phoenixframework.org/
  * Guides: http://phoenixframework.org/docs/overview
  * Docs: https://hexdocs.pm/phoenix
  * Mailing list: http://groups.google.com/group/phoenix-talk
  * Source: https://github.com/phoenixframework/phoenix



































# --------------------------------------------------------------------------------------------------
# - http.server.addr
#   The address the HTTP server gonna listen. The value should have this format: "host:port".
#   The host can be omited, in this case, Flare gonna bind to all available interfaces. The port is
#   required. Default value: ":8080"
#
# - http.default-limit
#   The quantity of entities to be returned by default. Default value: 30.
#
# - http.timeout
#   The max time a request has to process. Default value: 1s.
#
[http]
addr          = ":8080"
default-limit = 30
timeout       = "1s"

# --------------------------------------------------------------------------------------------------
# - repository.engine
#   The location of the content. Default value: "memory". Possible values: "memory" and "mongodb".
#
[repository]
engine = "memory"

# --------------------------------------------------------------------------------------------------
# If "repository.engine" is "mongodb", the options bellow can be used:
#
#  - repository.addrs
#    List of "host:port" to MongoDB servers. Default value: ["localhost:27017"]
#
#  - repository.database
#    Name of the database on MongoDB. Default value: "flare"
#
#  - repository.username
#    Username used to connect to MongoDB. Default value is unset.
#
#  - repository.password
#    Password used to connect to MongoDB. Default value is unset.
#
#  - repository.replica-set
#    Replica set, if specified, will prevent the obtained session from communicating with any server
#    which is not part of a replica set with the given name. The default is to communicate with any
#    server specified or discovered via the servers contacted. Default value is unset.
#
#  - repository.pool-limit
#    Defines the per-server socket pool limit. Defaults value is 4096.
#
#  - repository.timeout
#    Timeout is the amount of time to wait for a server to respond when first connecting and on
#    follow up operations in the session. Default value is 1s.
#
[repository]
engine      = "mongodb"
addrs       = ["localhost:27017"]
database    = "flare"
username    = "flare"
password    = "flare"
replica-set = ""
pool-limit  = 4096
timeout     = "1s"

# --------------------------------------------------------------------------------------------------
# - task.engine
#   The engine used to enqueue jobs. If the 'sqs' is chosen, the 'aws' config block must be
#   configured. Possible values: "sqs" or "memory". Default value: "memory".
#
# - task.queue-subscription-bucket
#   If the SQS is used as engine, there is a option to set the queue name. Default value is
#   "flare-document-queue".
#
[task]
engine                      = "sqs"
queue-subscription-bucket   = "flare-subscription-bucket"
queue-subscription-spread   = "flare-subscription-spread"
queue-subscription-delivery = "flare-subscription-delivery"

# --------------------------------------------------------------------------------------------------
# - aws.key
#   Key used to connect to AWS. Default value is unset.
#
# - aws.secret
#   Secret used to connect to AWS. Default value is unset.
#
# - aws.region
#   Region used to connect to AWS. Default value is unset.
#
[aws]
key    = "key"
secret = "secret"
region = "us-east-1"

# --------------------------------------------------------------------------------------------------
# - log.level
#   The minimum log level to be displayed. Default value: "debug". Possible values: "debug", "info",
#   "warn" and "error".
#
# - log.output
#   Where the logs gonna be sent. Default value: "stdout". Possible values: "stdout" and "discard".
#
# - log.format
#   Format of the outputed log. Default value: "human". Possible values: "human" and "json".
#
[log]
level  = "debug"
output = "stdout"
format = "human"



