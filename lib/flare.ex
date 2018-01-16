defmodule Flare do
  use Application

  def start(_type, _args) do
    Flare.Application.start(nil, nil)
  end
end
