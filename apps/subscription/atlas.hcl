data "external_schema" "local" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader/main.go",
  ]
}

env "local" {
  src = data.external_schema.local.url
  dev = "postgres://v1nislpo@sphnet.com.sg:@localhost:5432/subs_temp?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
