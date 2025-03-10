{
  description = "Subscriptions Development Shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        postgres_port = toString 5433;
        data_dir = "./.data/postgres";
        postgres_user = "postgres";
        postgres_password = "postgres";
        db_name = "subscriptions";
        db_url = "postgresql://${postgres_user}:${postgres_password}@localhost:${postgres_port}/${db_name}?sslmode=disable";
        db_inspection_url = "postgresql://${postgres_user}:${postgres_password}@localhost:${postgres_port}/${db_name}.inspection?sslmode=disable";
      in
      {
        devShells.default =
          let
            packages = with pkgs; [
              postgresql
              go
              atlas
            ];
            scripts = [
              (pkgs.writeShellScriptBin "db" (''
                case $1 in
                  "prepare") ${db.prepare};;
                  "start") ${db.start};;
                  "stop") ${db.stop};;
                  "restart") ${db.restart};;
                  "migrate") ${db.migrate};;
                  "inspect") ${db.inspect};;
                  "generate-diff") ${db.diff};;
                  *) echo "Usage:
                  db prepare 
                    - prepare data directory. 
                  db start  
                    - start postgres. create necessary databases if not exist.
                  db stop
                    - stop postgres
                  db restart
                    - restart postgres
                  db migrate
                    - migrate database (atlas migrate apply)
                  db inspect
                    - inspect database (atlas schema inspect)
                  db generate-diff
                    - diff between database and migration files (atlas migrate diff)"
                  exit 1;;
                esac
              ''))
              (pkgs.writeShellScriptBin "app" (''
                case $1 in
                  "build") ${application.build};;
                  "run") ${application.run};;
                  "test") ${application.test};;
                  *) echo "Usage:
                  app build   
                    - build application
                  app run     
                    - run application
                  app test
                    - test application"
                  exit 1;;
                esac
              ''))
            ];
            application = {
              build = ''
                go build -o ./bin/subscriptions cmd/main.go
              '';
              run = ''
                go run cmd/main.go
              '';
              test = ''
                go test ./...
              '';
            };
            db = {
              prepare = ''
                mkdir -p ${data_dir}
                initdb -D ${data_dir} --auth=trust --username=${postgres_user}
              '';
              start = ''
                pg_ctl -D ${data_dir} -l ${data_dir}.log -o "-p ${postgres_port}" start
                pg_ctl -D ${data_dir} reload
                if ! psql -h localhost -p ${postgres_port} -U ${postgres_user} -lqt | cut -d \| -f 1 | grep -qw ${db_name}; then
                  createdb -h localhost -p ${postgres_port} -U ${postgres_user} ${db_name}
                else
                  echo "Database ${db_name} already exists"
                fi
                if ! psql -h localhost -p ${postgres_port} -U ${postgres_user} -lqt | cut -d \| -f 1 | grep -qw ${db_name}.inspection; then
                  createdb -h localhost -p ${postgres_port} -U ${postgres_user} ${db_name}.inspection
                else
                  echo "Database ${db_name}.inspection already exists"
                fi
              '';
              stop = ''
                pg_ctl -D ${data_dir} stop
              '';
              restart = ''
                pg_ctl -D ${data_dir} restart
              '';
              migrate = ''
                atlas migrate apply --url ${db_url}
              '';
              inspect = ''
                atlas schema inspect --url ${db_inspection_url}
              '';
              diff = ''
                atlas migrate diff --dev-url ${db_inspection_url} --dir file://migrations --to ${db_url}
              '';
            };
          in
          pkgs.mkShell {
            buildInputs = packages ++ scripts;
            shellHook = ''
              export DATABASE_URL="${db_url}"
              export INSPECTION_DATABASE_URL="${db_inspection_url}"

              echo "db  
                - manage postgres
              app 
                - manage application"
            '';
          };
      }
    );
}
