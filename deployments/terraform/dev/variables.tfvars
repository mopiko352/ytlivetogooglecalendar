function = {
    configuration = [
        {
            channel_id = "UC_BJOJtwijKGpyARVPUVw4w"
            secret_sa = "ytlive-to-calendar-sa-key"
            secret_calendar = "carol-calendar-id"
            name = "ytliveToCalendarForCarol"
            description = "apply youtube schedule to google calendar for caroline"
            memory = 128
        }
    ]
}

sa_roles = [
    "roles/cloudfunctions.invoker",
    "roles/secretmanager.viewer",
    "roles/secretmanager.secretAccessor"
]

pubsub_name = "live-to-calendar"

scheduler = {
    name = "ytlive-to-calendar",
    description = "schedule calendar batch",
    cron = "0 */3 * * *"
}