variable "function" {
    type = object({
      configuration = list(object({
        channel_id = string
        name = string
        secret_sa = string
        secret_calendar = string
        description = string
        memory = number
      }))
    })
}

variable "sa_roles" {
    type = list(string)
}

variable "pubsub_name" {
    type = string
}

variable "scheduler" {
    type = object({
      name = string
      description = string
      cron = string
    })
}
