resource "kubernetes_namespace" "swipelearn" {
  metadata {
    name = var.namespace
    labels = {
      name = var.namespace
      app  = var.app_name
    }
  }
}
