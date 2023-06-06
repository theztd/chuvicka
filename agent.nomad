variable "image" {
  type    = string
  default = "ghcr.io/theztd/chuvicka:main"
}

job "chuvicka-cron" {

  type = "batch"
  
  periodic {
    cron = "* * * * *"
    prohibit_overlap = true
  }

  group "agent" {
    task "frequent_tasks" {
      driver = "docker"
      config {
        image   = var.image
        args    = ["agent"]
      }
      resources {
        cpu    = 500
        memory = 64
      }

      template {
        env         = true
        data        = <<EOH
          # meta
          {{ with nomadVar "nomad/jobs/chuvicka" }}{{ .Parent.Items | sprig_toJson | parseJSON | toTOML }}{{end}}
        EOH
        change_mode = "restart"
        destination = "local/.secret-env"
      }

    }
  }
}