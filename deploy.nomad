variable "fqdn" {
  type    = string
  default = "chuvicka.check-this.link"
}

variable "dcs" {
  type    = list(string)
  default = ["dc1", "devel"]
}


variable "image" {
  type = string
  default = "ghcr.io/theztd/chuvicka:24e12e6d344a6f3a453ecfb40e7f2a69f8158ec2"
}


//job "__JOB_NAME__" {

job "chuvicka-dev" {
  datacenters = var.dcs

  meta {
    fqdn = var.fqdn
    git = "github.com/theztd/chuvicka"
    managed = "github-pipeline"
    image = var.image
  }

  group "fe" {
    count = 1

    network {
      dns {
        servers = ["172.17.0.1", "8.8.8.8", "1.1.1.1"]
      }
      
      port "app" { to = 8080 }
      port "http" { to = 80 }
    }

    service {
      provider = "nomad"
      name = "${JOB}-http"

      tags = [
        "public",
        "traefik.enable=true",
        "traefik.http.routers.${NOMAD_JOB_NAME}-http.rule=Host(`${var.fqdn}`)"
        //"traefik.http.routers.${NOMAD_JOB_NAME}-http.tls=true"
      ]

      port = "http"
    }

    task "nginx" {
      driver = "docker"

      config {
        image = "nginx:1.21"

        volumes = [
          "local:/etc/nginx/conf.d",
        ]

        ports = ["http"]
      }

      template {
        destination = "local/default.conf"
        perms       = "644"
        data        = file("nginx.conf")
      }

      # Resources:    https://www.nomadproject.io/docs/job-specification/resources
      resources {
        cpu        = 100 # MHz
        memory     = 16  # MB
        memory_max = 64  #MB
      }


      kill_timeout = "10s"
    }
    # END NGinx task


    task "app" {

      driver = "docker"

      config {
        image      = var.image
        force_pull = true

        ports = ["app"]

        labels {
          group = "app"
        }
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

  group "backend" {
    count = 1

    network {
      dns {
        servers = ["172.17.0.1", "8.8.8.8", "1.1.1.1"]
      }
      
      port "influxdb" { 
        static = 18086 
        to = 8086 
      }

      
    }

    task "influx" {
      driver = "docker"

      config {
        image = "influxdb:2.6.1"

        ports = ["influxdb"]
        
        // myInfluxVolume:/var/lib/influxdb2
      }

      env {
        DOCKER_INFLUXDB_INIT_MODE = "setup"
        DOCKER_INFLUXDB_INIT_USERNAME = "chuvicka"
        DOCKER_INFLUXDB_INIT_PASSWORD = "Heslicko"
        DOCKER_INFLUXDB_INIT_ORG = "chuvicka"
        DOCKER_INFLUXDB_INIT_BUCKET = "chuvicka"
        DOCKER_INFLUXDB_INIT_RETENTION = "1w"
      }

    } # END task influx
  } # END group backend

}
