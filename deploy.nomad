variable "fqdn" {
  type    = string
  default = "chuvicka.check-this.link"
}

variable "dcs" {
  type    = list(string)
  default = ["dc1", "devel"]
}


variable "image" {
  type    = string
  default = "ghcr.io/theztd/chuvicka:24e12e6d344a6f3a453ecfb40e7f2a69f8158ec2"
}


//job "__JOB_NAME__" {

job "chuvicka-dev" {
  datacenters = var.dcs

  meta {
    fqdn    = var.fqdn
    git     = "github.com/theztd/chuvicka"
    managed = "github-pipeline"
    image   = var.image
  }

  group "fe" {
    count = 1

    network {
      dns {
        servers = ["172.17.0.1", "8.8.8.8", "1.1.1.1"]
      }

      port "server" { to = 8080 }
      // port "http" { to = 80 }
    }

    service {
      provider = "nomad"
      name     = "${JOB}-http"
      port     = "server"

      tags = [
        "public",
        "traefik.enable=true",
        "traefik.http.routers.${NOMAD_JOB_NAME}-http.rule=Host(`${var.fqdn}`)"
        //"traefik.http.routers.${NOMAD_JOB_NAME}-http.tls=true"
      ]
    }

    task "server" {

      driver = "docker"

      config {
        image      = var.image
        force_pull = true
        command    = ["server"]

        ports = ["server"]

        labels {
          group = "app"
        }
      }

      template {
        env         = true
        data        = <<EOH
          {{ range nomadVarList }}
            {{ . }}
          {{ end }}
          }
          EOH
        change_mode = "restart"
        destination = "${NOMAD_SECRETS_DIR}/.secret-env"
      }

      env {
        ADDRESS = ":8080"
      }

      resources {
        cpu        = 100
        memory     = 64
        memory_max = 96
      }

    } # END task app

  } # END group FE

  // group "backend" {
  //   count = 1

  //   network {
  //     dns {
  //       servers = ["172.17.0.1", "8.8.8.8", "1.1.1.1"]
  //     }

  //     port "influxdb" { 
  //       static = 18086 
  //       to = 8086 
  //     }


  //   }

  //   task "influx" {
  //     driver = "docker"

  //     config {
  //       image = "influxdb:2.6.1"

  //       ports = ["influxdb"]

  //       mount {
  //         type = "volume"
  //         source = "files-${local.prometheus.fqdn}"
  //         target = "/prometheus"
  //         readonly = false
  //         volume_options {
  //           labels {
  //             job = "${NOMAD_JOB_NAME}"
  //             domain = "${local.prometheus.fqdn}"
  //             backup = "true"
  //             type = "files"
  //           }
  //         }
  //       }
  //       // myInfluxVolume:/var/lib/influxdb2
  //     }

  //     env {
  //       DOCKER_INFLUXDB_INIT_MODE = "setup"
  //       DOCKER_INFLUXDB_INIT_USERNAME = "chuvicka"
  //       DOCKER_INFLUXDB_INIT_PASSWORD = "Heslicko"
  //       DOCKER_INFLUXDB_INIT_ORG = "chuvicka"
  //       DOCKER_INFLUXDB_INIT_BUCKET = "chuvicka"
  //       DOCKER_INFLUXDB_INIT_RETENTION = "1w"
  //     }

  //   } # END task influx
  // } # END group backend

}
