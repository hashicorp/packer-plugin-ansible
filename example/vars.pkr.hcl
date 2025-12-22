# Copyright IBM Corp. 2013, 2025
# SPDX-License-Identifier: MPL-2.0

variable "account_file" {
    default = env("GOOGLE_APPLICATION_CREDENTIALS")
}

variable "project_id" {
    default = env("GOOGLE_PROJECT_ID")
}

