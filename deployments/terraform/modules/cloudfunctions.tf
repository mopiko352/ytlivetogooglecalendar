data "google_project" "project" {
}

resource "google_cloudfunctions_function" "live_to_calendar_function" {
    for_each = { for i in var.function.configuration : i.name => i } 
    name = each.value.name
    description = each.value.description
    available_memory_mb = each.value.memory
    runtime = "go119"
    source_archive_bucket = data.google_storage_bucket.live_to_calendar_function_src.name
    source_archive_object = google_storage_bucket_object.packages.name
    service_account_email = google_service_account.live_to_calendar_sa.email
    # 共通でいいや
    event_trigger {
        event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
        resource   = google_pubsub_topic.live_to_calendar_pubsub.name
    }
    environment_variables = {
        CHANNEL_ID = each.value.channel_id
        SA_SECRET_PATH = "projects/${data.google_project.project.number}/secrets/${each.value.secret_sa}/versions/latest"
        CALENDAR_ID_SECRET_PATH = "projects/${data.google_project.project.number}/secrets/${each.value.secret_calendar}/versions/latest"
    }
    entry_point           = "LiveToCalendar"
}

resource "google_storage_bucket_object" "packages" {
    name = "packages/functions.${data.archive_file.function_archive.output_md5}.zip"
    bucket = data.google_storage_bucket.live_to_calendar_function_src.name
    source = data.archive_file.function_archive.output_path
}

# tfstateとgcf sourceをおいておくばけっと　事前に手動で作っておく　terraform管理はしない
data "google_storage_bucket" "live_to_calendar_function_src" {
  name     = "${data.google_project.project.project_id}-live-to-calendar-srcs"
}

data "archive_file" "function_archive" {
    type = "zip"
    source_dir = "../../.."
    excludes = [
        "../../../Readme.md",
        "../../../Dockerfile",
        "../../../deployments",
        "../../../.gitignore"
    ]
    output_path = "index.zip"
}