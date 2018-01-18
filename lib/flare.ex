defmodule Flare do
  use Application

  def start(_type, _args) do
    Flare.Supervisor.start(nil, nil)
  end
end
