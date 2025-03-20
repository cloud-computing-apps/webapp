packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  region        = var.region
  profile       = var.profile
  instance_type = var.instance_type
  ssh_username  = var.ssh_username

  source_ami_filter {
    filters = {
      name                = var.aws_ami_name
      root-device-type    = var.root_device_type
      virtualization-type = var.virt_type
    }
    owners      = [var.aws_owner]
    most_recent = true
  }

  ami_name   = "${var.image_name}-{{timestamp}}"
  ami_groups = []

  launch_block_device_mappings {
    device_name           = var.device_name
    delete_on_termination = true
  }
}

build {
  sources = ["source.amazon-ebs.ubuntu"]

  provisioner "shell" {
    script = "./scripts/user_group.sh"
  }

  provisioner "shell" {
    script = "./scripts/install_golang.sh"
  }

  provisioner "file" {
    source      = "webapp"
    destination = "/tmp/webapp"
  }

  provisioner "shell" {
    script = "./scripts/build_webapp.sh"
  }

  provisioner "file" {
    source      = "webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    script = "./scripts/systemd_conf.sh"
  }
}
