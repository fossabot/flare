defmodule Flare.Plugin do
  def repository(resource) do
    provider = Application.get_env(:flare, :repository)[:provider]
    module = Module.concat(Flare.Plugin, provider)

    case resource do
      :resource -> Module.concat(module, Resource)
      :subscription -> Module.concat(module, Subscription)
    end
  end
end
