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
}


job "__JOB_NAME__" {
  datacenters = var.dcs

  group "fe" {
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["172.17.0.1", "8.8.8.8", "1.1.1.1"]
      }
      
      port "app" { to = 8080 }
      port "http" { to = 80 }
    }

    service {
      name = "${JOB}-http"

      tags = [
        "public",
        "traefik.enable=true",
        "traefik.http.routers.${NOMAD_JOB_NAME}-http.rule=Host(`http-${var.fqdn}`)"
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

  } # END group FE

}
