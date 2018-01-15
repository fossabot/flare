defmodule Flare.Domain.Resource.Repository do
  alias Flare.Domain.Resource
  alias Flare.Domain.Resource.RepositoryError

  @callback all() :: {:ok, [Resource]} | {:error, RepositoryError}
  @callback one(String.t(), keyword) :: {:ok, Resource} | {:error, RepositoryError}
  @callback create(Resource) :: :ok | {:error, RepositoryError}
  @callback destroy(String.t()) :: :ok | {:error, RepositoryError}
end

