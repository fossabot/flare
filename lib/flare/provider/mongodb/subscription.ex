defmodule Flare.Provider.MongoDB.Subscription do
  @behaviour Flare.Domain.Subscription.Repository
  alias Flare.Provider.MongoDB.Resource

  @options Application.get_env(:flare, :repository)[:options]
  @pid :mongo
  @collection "subscriptions"

  def all do
    Mongo.find(@pid, @collection, %{}, @options) |> Enum.to_list()
  end

  def create(content) do
    IO.inspect(content)

    content =
      content
      |> Map.put("id", UUID.uuid4())
      |> Map.put("createdAt", DateTime.utc_now())

    # - validamos ante pra ver se ta tudo certo
    # - pede para o resource o bucket que vamos pertencer, em seguida colocamos no subscription, o count ja incrementa
    # - qnd quiermos o bucket, vamos  buscar logo os 1000 elementos
    {:ok, key} = Resource.bucket_placement(Map.get(content, "id"))
    content = Map.put(content, "bucket", key)

    Mongo.insert_one(@pid, @collection, content, @options)
  end
end
