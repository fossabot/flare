defmodule Flare.Plugin.Memory.Resource do
  @behaviour Flare.Domain.Resource.Repository
  use Agent

  def start_link do
    Agent.start_link(fn -> %{} end, name: __MODULE__)
  end

  def all() do
    Agent.get(__MODULE__, fn set ->
      set
    end)

    # [
    #   %{
    #     id: "diego",
    #     addresses: ["http://app.com"],
    #     path: "/users/{id}",
    #     change: %{
    #       field: "updatedAt",
    #       format: "2016"
    #     }
    #   }
    # ]
  end

  def create(_content) do
    # content = Map.put(content, :id, UUID.uuid4())

    # found = false

    # Agent.get(__MODULE__, fn resources ->
    #   Enum.each(fn resource ->
    #     nil
    #   end)
    # end)
  end
end
