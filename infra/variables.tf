variable "namespace" {
  description = "Kubernetes namespace for SwipeLearn"
  type        = string
  default     = "swipelearn"
}

variable "app_name" {
  description = "Name of the application"
  type        = string
  default     = "swipelearn"
}
