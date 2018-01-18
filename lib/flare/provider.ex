defmodule Flare.Provider do
  def repository(resource) do
    provider = Application.get_env(:flare, :repository)[:provider]
    module = Module.concat(Flare.Provider, provider)

    case resource do
      :resource -> Module.concat(module, Resource)
      :subscription -> Module.concat(module, Subscription)
    end
  end
end
