defmodule Flare.Plugin do
  @provider Application.get_env(:flare, :repository)[:provider]

  def repository(resource) do
    module = Module.concat(Flare.Plugin, @provider)

    case resource do
      :resource -> Module.concat(module, Resource)
      :subscription -> Module.concat(module, Subscription)
    end
  end
end