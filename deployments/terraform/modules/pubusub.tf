resource "google_cloud_scheduler_job" "invoke_function" {
    name = var.scheduler.name
    description = var.scheduler.description
    schedule = var.scheduler.cron
    pubsub_target {
        topic_name = google_pubsub_topic.live_to_calendar_pubsub.id
        data =  base64encode("ok")
    }
}

resource "google_pubsub_topic" "live_to_calendar_pubsub" {
  name = var.pubsub_name
  message_retention_duration = "86600s"
}