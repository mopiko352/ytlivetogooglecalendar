provider "google" {
    region = "asia-northeast1"
    zone = "asia-northeast1"
}


module "live_to_calendar" {
    source = "../modules"
    function = var.function
    sa_roles = var.sa_roles
    pubsub_name = var.pubsub_name
    scheduler = var.scheduler
}