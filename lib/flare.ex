defmodule Flare do
  use Application

  def start(_type, _args) do
    IO.puts("first thing...")
    Flare.Application.start(nil, nil)
  end
end
