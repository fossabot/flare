defmodule Flare.Plugin.MongoDB.Resource do
  @behaviour Flare.Domain.Resource.Repository
  alias Flare.Domain.Resource

  @bucket_size Application.get_env(:flare, :repository)[:resource][:bucket_size]
  @options Application.get_env(:flare, :repository)[:options]
  @pid :mongo
  @collection "resources"

  def all(opts \\ []) do
    opts = Keyword.merge(opts, @options)

    t1 =
      Task.async(fn ->
        %{
          step: :query,
          result: Mongo.find(@pid, @collection, %{}, opts) |> Enum.map(&transform/1)
        }
      end)

    t2 =
      Task.async(fn ->
        %{step: :count, result: Mongo.count(@pid, @collection, %{}, opts)}
      end)

    result =
      [t1, t2]
      |> Enum.map(&Task.await/1)
      |> Enum.reduce(%{}, fn t, result ->
        case Map.get(t, :step) do
          :query -> Map.put(result, :resources, Map.get(t, :result))
          :count -> Map.put(result, :total, elem(Map.get(t, :result), 1))
        end
      end)

    {:ok, result}
  end

  def one(id, opts \\ []) do
    opts = Keyword.merge(opts, @options)

    case Mongo.find_one(@pid, @collection, %{"id" => id}, opts) do
      doc when is_map(doc) -> {:ok, transform(doc)}
      nil -> {:error, :not_found}
    end
  end

  def bucketsCount(id, opts \\ []) do
    opts =
      opts
      |> Keyword.merge(@options)
      |> Keyword.put(:projection, %{"bucketsCount" => 1})

    Mongo.find_one(@pid, @collection, %{"id" => id}, opts) |> Map.get("bucketsCount")
  end

  def create(content) do
    content =
      content
      |> Map.put("id", UUID.uuid4())
      |> Map.put("createdAt", DateTime.utc_now())
      |> Map.put("buckets", %{"0" => 0})
      |> Map.put("bucketsCount", 0)

    case Mongo.insert_one(@pid, @collection, content, @options) do
      {:ok, _} -> {:ok, Map.get(content, "id")}
      {:error, %Mongo.Error{code: 11000}} -> {:error, "duplicate key on address or path"}
    end
  end

  def delete(id) do
    case Mongo.delete_one(@pid, @collection, %{"id" => id}, @options) do
      {:ok, %Mongo.DeleteResult{deleted_count: 0}} -> {:error, :not_found}
      {:ok, _} -> :ok
    end
  end

  # tem que ver os casos de erro agora.
  # Flare.Plugin.MongoDB.Resource.bucket_placement("9799aabe-81e9-49c2-a49b-99817751d036")
  # resource = Flare.Plugin.MongoDB.Resource.one("dcccf67b-f8c4-4c6c-8146-793ef29c856a")
  # key = Flare.Plugin.MongoDB.Resource.bucket_select(resource)
  # id = "dcccf67b-f8c4-4c6c-8146-793ef29c856a"
  def bucket_placement(id) do
    case one(id) do
      resource ->
        key = bucket_select(resource)

        case bucket_increment(id, key) do
          {:ok, _} -> {:ok, key}
        end
    end
  end

  def bucket_increment(id, key) do
    Mongo.update_one(
      @pid,
      @collection,
      %{"id" => id},
      %{"$inc" => %{"buckets.#{key}" => 1}},
      @options
    )
  end

  def bucket_select(resource) do
    buckets = Map.get(resource, "buckets")

    case Enum.find(buckets, fn x ->
           elem(x, 1) < @bucket_size
         end) do
      {key, _value} -> key
      nil -> Kernel.map_size(buckets) |> Integer.to_string()
    end
  end

  defp transform(doc) do
    %Resource{
      id: Map.get(doc, "id"),
      path: Map.get(doc, "path"),
      created_at: Map.get(doc, "createdAt"),
      addresses: Map.get(doc, "addresses"),
      change: Map.get(doc, "change")
    }
  end
end
