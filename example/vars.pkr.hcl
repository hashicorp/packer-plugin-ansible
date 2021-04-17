variable "account_file" {
    default = env("GOOGLE_APPLICATION_CREDENTIALS")
}

variable "project_id" {
    default = env("GOOGLE_PROJECT_ID")
}

