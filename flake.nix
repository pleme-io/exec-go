# flake.nix — exec-go (GSDS Biblioteca) via substrate's go-library-flake.
# vendorHash OMITTED → spec-sourced (__from-spec__); clean nix build lands once
# errors-go is published. Pre-publish proof is `go test` (green).
{
  description = "exec-go — the fleet's one typed external-command wrapper";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    substrate = {
      url = "github:pleme-io/substrate";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ { self, nixpkgs, substrate, ... }:
    (import substrate.goLibraryFlakeBuilder { inherit nixpkgs; }) {
      name = "exec-go";
      version = "0.1.0";
      src = self;
      repo = "pleme-io/exec-go";
    };
}
