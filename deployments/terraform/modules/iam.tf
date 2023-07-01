resource "google_service_account" "live_to_calendar_sa" {
    account_id = "live-to-calendar"
    display_name = "apply youtube stream schedules to google calendar"
}

resource "google_project_iam_member" "live_to_calendar_sa_member" {
  project = data.google_project.project.project_id
  for_each = {for s,v in var.sa_roles : s => v}
  role = each.value
  member = "serviceAccount:${google_service_account.live_to_calendar_sa.email}"
}