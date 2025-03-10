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
  dev = "${INSPECTION_DATABASE_URL}"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
