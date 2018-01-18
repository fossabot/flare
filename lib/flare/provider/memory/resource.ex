defmodule Flare.Provider.Memory.Resource do
  @behaviour Flare.Domain.Resource.Repository
  use Agent
  alias Flare.Domain.Resource

  def start_link do
    Agent.start_link(fn -> [%Resource{}] end, name: __MODULE__)
  end

  def all() do
    {:ok, %{resources: Agent.get(__MODULE__, fn set -> set end)}}
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
