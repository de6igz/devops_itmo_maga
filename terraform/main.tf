terraform {
  required_providers {
    twc = {
      source = "tf.timeweb.cloud/timeweb-cloud/timeweb-cloud"
    }
  }
  required_version = ">= 1.4.4"
}

provider "twc" {
  token = var.twc_token
}

data "twc_configurator" "configurator" {
  location    = var.location
  preset_type = var.preset_type
}

data "twc_os" "os" {
  name    = var.os_name
  version = var.os_version
}

resource "twc_ssh_key" "lab_key" {
  name = "lab2-key"
  body = file(var.ssh_public_key_path)
}

resource "twc_floating_ip" "lab_ip" {
  availability_zone = var.availability_zone
}

resource "twc_server" "lab_server" {
  name              = var.server_name
  os_id             = data.twc_os.os.id
  availability_zone = var.availability_zone
  floating_ip_id    = twc_floating_ip.lab_ip.id
  ssh_keys_ids      = [twc_ssh_key.lab_key.id]

  configuration {
    configurator_id = data.twc_configurator.configurator.id
    disk            = var.disk
    cpu             = var.cpu
    ram             = var.ram
  }
}