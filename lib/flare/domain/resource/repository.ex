defmodule Flare.Domain.Resource.Repository do
  alias Flare.Domain.Resource

  @callback all() :: {:ok, [Resource]} | {:error, String.t()}
  @callback one(String.t(), keyword) :: {:ok, Resource} | {:error, Atom.t()}
  @callback create(Resource) :: {:ok}
end
