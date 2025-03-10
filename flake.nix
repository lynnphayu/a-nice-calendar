{
  description = "Development environment shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nodejs-18_x
            postgresql
            mongodb-ce
            zsh
            go
            pnpm
          ];
          shellHook = ''
            # Create local data directories
            mkdir -p ./.data/mongodb
            mkdir -p ./.data/postgres

            # kill process on port 5433 and 27018
            mongopid=$(lsof -t -i:27018)
            postgpid=$(lsof -t -i:5433)
            if [ ! -z "$mongopid" ]; then
              kill -9 $mongopid
            fi
            if [! -z "$postgpid" ]; then
              kill -9 $postgpid
            fi

            # Initialize and start PostgreSQL
            initdb -D ./.data/postgres --auth=trust --username=postgres 
            pg_ctl -D ./.data/postgres -l ./.data/postgres.log -o "-p 5433" start

            pg_ctl -D ./.data/postgres reload
            psql -h localhost -p 5433 -U postgres -w -c "ALTER USER postgres WITH PASSWORD 'postgres';"
            createdb -h localhost -p 5433 -U postgres subscriptions
        
            # Configure MongoDB permissions
            chmod -R 777 ./.data/mongodb

            # Start MongoDB with local data directory 
            mongod --dbpath ./.data/mongodb --bind_ip 127.0.0.1 --port 27018 --fork --logpath ./.data/mongod.log

            export POSTGRESQL_URL="postgresql://postgres:postgres@localhost:5433/subs"
            export MONGODB_URL="mongodb://localhost:27018"

            # Add shutdown trap
            trap '\
            kill -9 $(lsof -t -i:27018) && \
            kill -9 $(lsof -t -i:5433) \
            ' EXIT
          '';
        };
      }
    );
}
