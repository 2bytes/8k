{
    inputs = {
        nixpkgs.url = "github:nixos/nixpkgs";
        flake-utils.url = "github:numtide/flake-utils";
    };

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachDefaultSystem (system:
            let
                pkgs = nixpkgs.legacyPackages.${system};
                packageName = "eightk";

                app = pkgs.buildGoModule rec {
                    pname = "8k";
                    version = "1.0.5";

                    # This src will pull from Github
                    # src = pkgs.fetchFromGitHub {
                    #     owner = "2bytes";
                    #     repo = "8k";
                    #     rev = "v${version}";
                    #     # Obtained by calling: test-build nix-prefetch-url --unpack https://github.com/2bytes/8k/archive/refs/tags/v1.0.5.zip
                    #     sha256 = "1nj9kfyn891bsysz4dgdqbb6cyavlj3pbif40xagzmc0anxp7w1z";
                    # };

                    # This src builds from the current dir (if this flake is in the root of the 8k project repo)
                    src = ./.;

                    runVend = true;
                    vendorSha256 = null;

                    buildFlagsArray = [
                        "-ldflags=-X github.com/2bytes/8k/internal/config.Version=v${version}"
                    ];

                    meta = with pkgs.lib; {
                        description = " https://8k.fyi - Text storage / link maker with a default byte max of 8192 (configurable).";
                        homepage = "https://github.com/2bytes/8k";
                        licence = licences.gpl3Only;
                        maintainers = with maintainers; [ hamid-elaosta ];
                        platforms = platforms.linux;
                    };
                };

            in {
                packages.${packageName} = app;
                defaultPackage = self.packages.${system}.${packageName};
            }
        );
}
