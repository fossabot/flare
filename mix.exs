defmodule Flare.Mixfile do
  use Mix.Project

  def project do
    [
      app: :flare,
      version: "0.0.1",
      elixir: "~> 1.4",
      elixirc_paths: elixirc_paths(Mix.env()),
      compilers: [:phoenix] ++ Mix.compilers(),
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  def application do
    [
      mod: {Flare.Application, []},
      extra_applications: [:logger, :runtime_tools]
    ]
  end

  defp elixirc_paths(:test), do: ["lib", "test/support"]
  defp elixirc_paths(_), do: ["lib"]

  defp deps do
    [
      # Concurrency
      {:gen_stage, "~> 0.12.2"},

      # Phoenix
      {:phoenix, "~> 1.3.0"},
      {:cowboy, "~> 1.0"},

      # UUID
      {:uuid, "~> 1.1"},

      # AWS
      {:ex_aws, "~> 2.0"},
      {:ex_aws_sqs, "~> 2.0"},
      {:hackney, "~> 1.9"},
      {:sweet_xml, "~> 0.6"},

      # MongoDB
      {:mongodb, ">= 0.0.0"},
      {:poolboy, ">= 0.0.0"}
    ]
  end
end
